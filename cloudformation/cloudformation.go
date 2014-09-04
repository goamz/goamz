//
// cloudformation: This package provides types and functions to interact with the AWS CloudFormation API
//
// Depends on https://github.com/goamz/goamz
//

package cloudformation

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	//"strconv"
	"strings"
	//"time"

	"github.com/goamz/goamz/aws"
)

// The CloudFormation type encapsulates operations within a specific EC2 region.
type CloudFormation struct {
	aws.Auth
	aws.Region
}

// New creates a new CloudFormation Client.
func New(auth aws.Auth, region aws.Region) *CloudFormation {

	return &CloudFormation{auth, region}

}

const debug = false

// ----------------------------------------------------------------------------
// Request dispatching logic.

// Error encapsulates an error returned by the AWS CloudFormation API.
//
// See http://goo.gl/zDZbuQ  for more details.
type Error struct {
	// HTTP status code (200, 403, ...)
	StatusCode int
	// CloudFormation error code
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

func (c *CloudFormation) query(params map[string]string, resp interface{}) error {
	params["Version"] = "2010-05-15"

	data := strings.NewReader(multimap(params).Encode())

	hreq, err := http.NewRequest("POST", c.Region.CloudFormationEndpoint+"/", data)
	if err != nil {
		return err
	}

	hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	token := c.Auth.Token()
	if token != "" {
		hreq.Header.Set("X-Amz-Security-Token", token)
	}

	signer := aws.NewV4Signer(c.Auth, "c", c.Region)
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

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}

// addParamsList adds params in the form of param.member.N to the params map
func addParamsList(params map[string]string, label string, ids []string) {
	for i, id := range ids {
		params[label+"."+strconv.Itoa(i+1)] = id
	}
}

// addComplexParamsList adds params in the form of param.member.key.value
func addComplexParamsList(params map[string]string, label, key string, ids []string) {
	for i, id := range ids {
		params[label+"."+strconv.Itoa(i+1)] = id
	}
}

// -----------------------------------------------------------------------
// API Supported Types and Methods

type CancelUpdateStackResponse struct {
	CancelUpdateStackResult string `xml:"CancelUpdateStackResult"`
	RequestId               string `xml:"ResponseMetadata>RequestId"`
}

// CancelUpdateStack cancels an update on the specified stack.
// If the call completes successfully, the stack will roll back the update and revert
// to the previous stack configuration.
//
// See http://goo.gl/ZE6fOa for more details
func (c *CloudFormation) CancelUpdateStack(stackName string) (resp *CancelUpdateStackResponse, err error) {
	params := makeParams("CancelUpdateStack")

	params["StackName"] = stackName

	resp = new(CancelUpdateStackResponse)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
