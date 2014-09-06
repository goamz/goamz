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
	"strconv"
	"strings"
	"time"

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
	// Error type
	Type string `xml:"Type"`
	// CloudFormation error code
	Code string `xml:"Code"`
	// The human-oriented error message
	Message   string `xml:"Message"`
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

// -----------------------------------------------------------------------
// API Supported Types and Methods

// SimpleResp is the basic response from most actions.
type SimpleResp struct {
	XMLName   xml.Name
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

// CancelUpdateStack cancels an update on the specified stack.
// If the call completes successfully, the stack will roll back the update and revert
// to the previous stack configuration.
//
// See http://goo.gl/ZE6fOa for more details
func (c *CloudFormation) CancelUpdateStack(stackName string) (resp *SimpleResp, err error) {
	params := makeParams("CancelUpdateStack")

	params["StackName"] = stackName

	resp = new(SimpleResp)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Parameter encapsulates the cloudstack paramter data type
//
// See http://goo.gl/2rg9eG for more details
type Parameter struct {
	ParameterKey     string `xml:"ParameterKey"`
	ParameterValue   string `xml:"ParameterValue"`
	UsePreviousValue bool   `xml:"UsePreviousValue"`
}

type Tag struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

// CreateStackParams wraps CreateStack request options
//
// See http://goo.gl/yDZYuV for more information
type CreateStackParams struct {
	Capabilities     []string
	DisableRollback  bool
	NotificationARNs []string
	OnFailure        string
	Parameters       []Parameter
	StackName        string
	StackPolicyBody  string
	StackPolicyURL   string
	Tags             []Tag
	TemplateBody     string
	TemplateURL      string
	TimeoutInMinutes int
}

// CreateStackResponse wraps a CreateStack call response
//
// See http://goo.gl/yDZYuV for more details
type CreateStackResponse struct {
	StackId   string `xml:"CreateStackResult>StackId"`
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

// CreateStack creates a stack as specified in the template. After the call completes successfully,
// the stack creation starts.
//
// Required params: StackName
//
// See http://goo.gl/yDZYuV for more details
func (c *CloudFormation) CreateStack(options *CreateStackParams) (
	resp *CreateStackResponse, err error) {
	params := makeParams("CreateStack")

	params["StackName"] = options.StackName

	if options.DisableRollback {
		params["DisableRollback"] = strconv.FormatBool(options.DisableRollback)
	}
	if options.OnFailure != "" {
		params["OnFailure"] = options.OnFailure
	}
	if options.StackPolicyBody != "" {
		params["StackPolicyBody"] = options.StackPolicyBody
	}
	if options.StackPolicyURL != "" {
		params["StackPolicyURL"] = options.StackPolicyURL
	}
	if options.TemplateBody != "" {
		params["TemplateBody"] = options.TemplateBody
	}
	if options.TemplateURL != "" {
		params["TemplateURL"] = options.TemplateURL
	}
	if options.TimeoutInMinutes != 0 {
		params["TimeoutInMinutes"] = strconv.Itoa(options.TimeoutInMinutes)
	}
	if len(options.Capabilities) > 0 {
		addParamsList(params, "Capabilities.member", options.Capabilities)
	}
	if len(options.NotificationARNs) > 0 {
		addParamsList(params, "NotificationARNs.member", options.NotificationARNs)
	}
	// Add any parameters
	for i, t := range options.Parameters {
		key := "Parameters.member.%d.%s"
		index := i + 1
		params[fmt.Sprintf(key, index, "ParameterKey")] = t.ParameterKey
		params[fmt.Sprintf(key, index, "ParameterValue")] = t.ParameterValue
		params[fmt.Sprintf(key, index, "UsePreviousValue")] = strconv.FormatBool(t.UsePreviousValue)
	}
	// Add any tags
	for i, t := range options.Tags {
		key := "Tags.member.%d.%s"
		index := i + 1
		params[fmt.Sprintf(key, index, "Key")] = t.Key
		params[fmt.Sprintf(key, index, "Value")] = t.Value
	}

	resp = new(CreateStackResponse)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteStack deletes a specified stack.
// Once the call completes successfully, stack deletion starts.
//
// See http://goo.gl/CVMpxC for more details
func (c *CloudFormation) DeleteStack(stackName string) (resp *SimpleResp, err error) {
	params := makeParams("DeleteStack")

	params["StackName"] = stackName

	resp = new(SimpleResp)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// StackEvent encapsulates the StackEvent data type
//
// See http://goo.gl/EHwiMf for more details
type StackEvent struct {
	EventId              string    `xml:"EventId"`
	LogicalResourceId    string    `xml:"LogicalResourceId"`
	PhysicalResourceId   string    `xml:"PhysicalResourceId"`
	ResourceProperties   string    `xml:"ResourceProperties"`
	ResourceStatus       string    `xml:"ResourceStatus"`
	ResourceStatusReason string    `xml:"ResourceStatusReason"`
	ResourceType         string    `xml:"ResourceType"`
	StackId              string    `xml:"StackId"`
	StackName            string    `xml:"StackName"`
	Timestamp            time.Time `xml:"Timestamp"`
}

// DescribeStackEventsResponse wraps a response returned by DescribeStackEvents request
//
// See http://goo.gl/zqj4Bz for more details
type DescribeStackEventsResponse struct {
	NextToken   string       `xml:"DescribeStackEventsResult>NextToken"`
	StackEvents []StackEvent `xml:"DescribeStackEventsResult>StackEvents>member"`
	RequestId   string       `xml:"ResponseMetadata>RequestId"`
}

// DescribeStackEvents returns all stack related events for a specified stack.
//
// See http://goo.gl/zqj4Bz for more details
func (c *CloudFormation) DescribeStackEvents(stackName string, nextToken string) (
	resp *DescribeStackEventsResponse, err error) {
	params := makeParams("DescribeStackEvents")

	if stackName != "" {
		params["StackName"] = stackName
	}
	if nextToken != "" {
		params["NextToken"] = nextToken
	}

	resp = new(DescribeStackEventsResponse)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// StackResourceDetail encapsulates the StackResourceDetail data type
//
// See http://goo.gl/flce6I for more details
type StackResourceDetail struct {
	Description          string    `xml:"Description"`
	LastUpdatedTimestamp time.Time `xml:"LastUpdatedTimestamp"`
	LogicalResourceId    string    `xml:"LogicalResourceId"`
	Metadata             string    `xml:"Metadata"`
	PhysicalResourceId   string    `xml:"PhysicalResourceId"`
	ResourceStatus       string    `xml:"ResourceStatus"`
	ResourceStatusReason string    `xml:"ResourceStatusReason"`
	ResourceType         string    `xml:"ResourceType"`
	StackId              string    `xml:"StackId"`
	StackName            string    `xml:"StackName"`
}

// DescribeStackResourceResponse wraps a response returned by DescribeStackResource request
//
// See http://goo.gl/6pfPFs for more details
type DescribeStackResourceResponse struct {
	StackResourceDetail StackResourceDetail `xml:"DescribeStackResourceResult>StackResourceDetail"`
	RequestId           string              `xml:"ResponseMetadata>RequestId"`
}

// DescribeStackResource returns a description of the specified resource in the specified stack.
// For deleted stacks, DescribeStackResource returns resource information
// for up to 90 days after the stack has been deleted.
//
// Required params: stackName, logicalResourceId
//
// See http://goo.gl/6pfPFs for more details
func (c *CloudFormation) DescribeStackResource(stackName string, logicalResourceId string) (
	resp *DescribeStackResourceResponse, err error) {
	params := makeParams("DescribeStackResource")

	params["StackName"] = stackName
	params["LogicalResourceId"] = logicalResourceId

	resp = new(DescribeStackResourceResponse)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// StackResource encapsulates the StackResource data type
//
// See http://goo.gl/j4eli5 for more details
type StackResource struct {
	Description          string    `xml:"Description"`
	LogicalResourceId    string    `xml:"LogicalResourceId"`
	PhysicalResourceId   string    `xml:"PhysicalResourceId"`
	ResourceStatus       string    `xml:"ResourceStatus"`
	ResourceStatusReason string    `xml:"ResourceStatusReason"`
	ResourceType         string    `xml:"ResourceType"`
	StackId              string    `xml:"StackId"`
	StackName            string    `xml:"StackName"`
	Timestamp            time.Time `xml:"Timestamp"`
}

// DescribeStackResourcesResponse wraps a response returned by DescribeStackResources request
//
// See http://goo.gl/YnY5rs for more details
type DescribeStackResourcesResponse struct {
	StackResources []StackResource `xml:"DescribeStackResourcesResult>StackResources>member"`
	RequestId      string          `xml:"ResponseMetadata>RequestId"`
}

// DescribeStackResources returns AWS resource descriptions for running and deleted stacks.
// If stackName is specified, all the associated resources that are part of the stack are returned.
// If physicalResourceId is specified, the associated resources of the stack that the resource
// belongs to are returned.
//
// Only the first 100 resources will be returned. If your stack has more resources than this,
// you should use ListStackResources instead.
//
// See http://goo.gl/YnY5rs for more details
func (c *CloudFormation) DescribeStackResources(stackName, physicalResourceId, logicalResourceId string) (
	resp *DescribeStackResourcesResponse, err error) {
	params := makeParams("DescribeStackResources")

	if stackName != "" {
		params["StackName"] = stackName
	}
	if physicalResourceId != "" {
		params["PhysicalResourceId"] = physicalResourceId
	}
	if logicalResourceId != "" {
		params["LogicalResourceId"] = logicalResourceId
	}

	resp = new(DescribeStackResourcesResponse)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Output encapsulates the Output AWS data type
//
// See http://goo.gl/UOn7q6 for more information
type Output struct {
	Description string `xml:"Description"`
	OutputKey   string `xml:"OutputKey"`
	OutputValue string `xml:"OutputValue"`
}

// Stack encapsulates the Stack AWS data type
//
// See http://goo.gl/yDZYuV for more information
type Stack struct {
	Capabilities      []string    `xml:"Capabilities>member"`
	CreationTime      time.Time   `xml:"CreationTime"`
	Description       string      `xml:"Description"`
	DisableRollback   bool        `xml:"DisableRollback"`
	LastUpdatedTime   time.Time   `xml:"LastUpdatedTime"`
	NotificationARNs  []string    `xml:"NotificationARNs>member"`
	Outputs           []Output    `xml:"Outputs>member"`
	Parameters        []Parameter `xml:"Parameters>member"`
	StackId           string      `xml:"StackId"`
	StackName         string      `xml:"StackName"`
	StackStatus       string      `xml:"StackStatus"`
	StackStatusReason string      `xml:"StackStatusReason"`
	Tags              []Tag       `xml:"Tags>member"`
	TimeoutInMinutes  int         `xml:"TimeoutInMinutes"`
}

// DescribeStacksResponse wraps a response returned by DescribeStacks request
//
// See http://goo.gl/UOLsXD for more information
type DescribeStacksResponse struct {
	NextToken string  `xml:"DescribeStacksResult>NextToken"`
	Stacks    []Stack `xml:"DescribeStacksResult>Stacks>member"`
	RequestId string  `xml:"ResponseMetadata>RequestId"`
}

// DescribeStacks returns the description for the specified stack;
// If no stack name was specified, then it returns the description for all the stacks created.
//
// See http://goo.gl/UOLsXD for more information
func (c *CloudFormation) DescribeStacks(stackName string, nextToken string) (
	resp *DescribeStacksResponse, err error) {
	params := makeParams("DescribeStacks")

	if stackName != "" {
		params["StackName"] = stackName
	}
	if nextToken != "" {
		params["NextToken"] = nextToken
	}

	resp = new(DescribeStacksResponse)
	if err := c.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
