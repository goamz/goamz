package sts

import (
	"encoding/xml"
	"fmt"
	"github.com/hailocab/goamz/aws"
	"log"
	"net/http"
	"net/http/httputil"
	"sort"
	"strconv"
	"strings"
	"time"
)

// The STS type encapsulates operations within a specific EC2 region.
type STS struct {
	aws.Auth
	aws.Region
	private byte // Reserve the right of using private data.
}

// New creates a new STS Client.
// We can only use us-east for region because AWS..
func New(auth aws.Auth, region aws.Region) *STS {
	return &STS{auth, aws.Regions["us-east-1"], 0}
}

const debug = false

// ----------------------------------------------------------------------------
// Request dispatching logic.

// Error encapsulates an error returned by the AWS STS API.
//
// See http://goo.gl/zDZbuQ  for more details.
type Error struct {
	// HTTP status code (200, 403, ...)
	StatusCode int
	// STS error code
	Code string
	// The human-oriented error message
	Message   string
	RequestId string `xml:"RequestID"`
}

func (err *Error) Error() string {
	if err.Code == "" {
		return err.Message
	}

	return fmt.Sprintf("%s (%s)", err.Message, err.Code)
}

type xmlErrors struct {
	RequestId string  `xml:"RequestId"`
	Errors    []Error `xml:"Error"`
}

func (sts *STS) query(params map[string]string, resp interface{}) error {
	params["Version"] = "2011-06-15"

	data := strings.NewReader(prepareParams(params))

	hreq, err := http.NewRequest("POST", sts.Region.STSEndpoint+"/", data)
	if err != nil {
		return err
	}

	hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	token := sts.Auth.Token()
	if token != "" {
		hreq.Header.Set("X-Amz-Security-Token", token)
	}

	signer := aws.NewV4Signer(sts.Auth, "sts", sts.Region)
	signer.Sign(hreq)

	if debug {
		log.Printf("%v -> {\n", hreq)
	}
	r, err := http.DefaultClient.Do(hreq)

	if err != nil {
		log.Printf("Error calling Amazon")
		return err
	}

	defer r.Body.Close()

	if debug {
		dump, _ := httputil.DumpResponse(r, true)
		log.Printf("response:\n")
		log.Printf("%v\n}\n", string(dump))
	}
	if r.StatusCode != 200 {
		return buildError(r)
	}
	err = xml.NewDecoder(r.Body).Decode(resp)
	return err
}

func buildError(r *http.Response) error {
	var (
		err    Error
		errors xmlErrors
	)
	xml.NewDecoder(r.Body).Decode(&errors)
	if len(errors.Errors) > 0 {
		err = errors.Errors[0]
	}

	err.RequestId = errors.RequestId
	err.StatusCode = r.StatusCode
	if err.Message == "" {
		err.Message = r.Status
	}
	return &err
}

func makeParams(action string) map[string]string {
	params := make(map[string]string)
	params["Action"] = action
	return params
}

func prepareParams(params map[string]string) string {
	var keys, sarray []string

	for k, _ := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		sarray = append(sarray, aws.Encode(k)+"="+aws.Encode(params[k]))
	}

	return strings.Join(sarray, "&")
}

// options for the AssumeRole function
//
// See http://goo.gl/Ld6Dbk for details
type AssumeRoleOptions struct {
	DurationSeconds int
	ExternalId      string
	Policy          string
	RoleArn         string
	RoleSessionName string
}

type AssumedRoleUser struct {
	Arn           string
	AssumedRoleId string
}

type Credentials struct {
	AccessKeyId     string
	Expiration      time.Time
	SecretAccessKey string
	SessionToken    string
}

type AssumeRoleResp struct {
	AssumedRoleUser  AssumedRoleUser `xml:"AssumeRoleResult>AssumedRoleUser"`
	Credentials      Credentials     `xml:"AssumeRoleResult>Credentials"`
	PackedPolicySize int             `xml:"AssumeRoleResult>PackedPolicySize"`
}

// AssumeRole assumes the specified role
//
// See http://goo.gl/zDZbuQ for more details.
func (sts *STS) AssumeRole(options *AssumeRoleOptions) (resp *AssumeRoleResp, err error) {
	params := makeParams("AssumeRole")

	params["RoleArn"] = options.RoleArn
	params["RoleSessionName"] = options.RoleSessionName

	if options.DurationSeconds != 0 {
		params["DurationSeconds"] = strconv.Itoa(options.DurationSeconds)
	}

	if options.ExternalId != "" {
		params["ExternalId"] = options.ExternalId
	}

	if options.Policy != "" {
		params["Policy"] = options.Policy
	}

	resp = new(AssumeRoleResp)
	if err := sts.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
