//
// ecs: This package provides types and functions to interact with the AWS EC2 Container Service API
//
// Depends on https://github.com/goamz/goamz
//
// Author boyann@gmail.com

package ecs

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/goamz/goamz/aws"
)

const debug = false

var timeNow = time.Now

// ECS contains the details of the AWS region to perform operations against.
type ECS struct {
	aws.Auth
	aws.Region
}

// New creates a new ECS Client.
func New(auth aws.Auth, region aws.Region) *ECS {
	return &ECS{auth, region}
}

// ----------------------------------------------------------------------------
// Request dispatching logic.

// Error encapsulates an error returned by the AWS ECS API.
//
// See http://goo.gl/VZGuC for more details.
type Error struct {
	// HTTP status code (200, 403, ...)
	StatusCode int
	// ECS error code ("UnsupportedOperation", ...)
	Code string
	// The error type
	Type string
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

func (e *ECS) query(params map[string]string, resp interface{}) error {
	params["Version"] = "2014-11-13"
	data := strings.NewReader(multimap(params).Encode())

	hreq, err := http.NewRequest("POST", e.Region.ECSEndpoint+"/", data)
	if err != nil {
		return err
	}

	hreq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	token := e.Auth.Token()
	if token != "" {
		hreq.Header.Set("X-Amz-Security-Token", token)
	}

	signer := aws.NewV4Signer(e.Auth, "ecs", e.Region)
	signer.Sign(hreq)

	if debug {
		log.Printf("%v -> {\n", hreq)
	}
	r, err := http.DefaultClient.Do(hreq)

	if err != nil {
		log.Printf("Error calling Amazon %v", err)
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

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}

func makeParams(action string) map[string]string {
	params := make(map[string]string)
	params["Action"] = action
	return params
}

func addParamsList(params map[string]string, label string, ids []string) {
	for i, id := range ids {
		params[label+"."+strconv.Itoa(i+1)] = id
	}
}

// ----------------------------------------------------------------------------
// Filtering helper.

// Filter builds filtering parameters to be used in an e query which supports
// filtering.  For example:
//
//     filter := NewFilter()
//     filter.Add("architecture", "i386")
//     filter.Add("launch-index", "0")
//     resp, err := e.DescribeTags(filter,nil,nil)
//
type Filter struct {
	m map[string][]string
}

// NewFilter creates a new Filter.
func NewFilter() *Filter {
	return &Filter{make(map[string][]string)}
}

// Add appends a filtering parameter with the given name and value(s).
func (f *Filter) Add(name string, value ...string) {
	f.m[name] = append(f.m[name], value...)
}

func (f *Filter) addParams(params map[string]string) {
	if f != nil {
		a := make([]string, len(f.m))
		i := 0
		for k := range f.m {
			a[i] = k
			i++
		}
		sort.StringSlice(a).Sort()
		for i, k := range a {
			prefix := "Filters.member." + strconv.Itoa(i+1)
			params[prefix+".Name"] = k
			for j, v := range f.m[k] {
				params[prefix+".Values.member."+strconv.Itoa(j+1)] = v
			}
		}
	}
}

// ----------------------------------------------------------------------------
// ECS types and related functions.

// SimpleResp is the beic response from most actions.
type SimpleResp struct {
	XMLName   xml.Name
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

// Cluster encapsulates the cluster datatype
//
// See
type Cluster struct {
	ClusterArn  string `xml:"clusterArn"`
	ClusterName string `xml:"clusterName"`
	Status      string `xml:"status"`
}

// CreateClusterReq encapsulates the createcluster req params
type CreateClusterReq struct {
	ClusterName string
}

// CreateClusterResp encapsulates the createcluster response
type CreateClusterResp struct {
	Cluster   Cluster `xml:"CreateClusterResult>cluster"`
	RequestId string  `xml:"ResponseMetadata>RequestId"`
}

// CreateCluster creates a new Amazon ECS cluster. By default, your account
// will receive a default cluster when you launch your first container instance
func (e *ECS) CreateCluster(req *CreateClusterReq) (resp *CreateClusterResp, err error) {
	if req == nil {
		return nil, fmt.Errorf("The req params cannot be nil")
	}

	params := makeParams("CreateCluster")
	params["clusterName"] = req.ClusterName

	resp = new(CreateClusterResp)
	if err := e.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// Resource describes the resources available for a container instance.
type Resource struct {
	DoubleValue    float64  `xml:"doubleValue"`
	IntegerValue   int64    `xml:"integerValue"`
	LongValue      int64    `xml:"longValue"`
	Name           string   `xml:"name"`
	StringSetValue []string `xml:"stringSetValue>member"`
	Type           string   `xml:"type"`
}

// ContainerInstance represents n Amazon EC2 instance that is running
// the Amazon ECS agent and has been registered with a cluster
type ContainerInstance struct {
	AgentConnected       bool       `xml:"agentConnected"`
	ContainerInstanceArn string     `xml:"containerInstanceArn"`
	Ec2InstanceId        string     `xml:"ec2InstanceId"`
	RegisteredResources  []Resource `xml:"registeredResources>member"`
	RemainingResources   []Resource `xml:"remainingResources>member"`
	Status               string     `xml:"status"`
}

// DeregisterContainerInstanceReq encapsulates DeregisterContainerInstance request params
type DeregisterContainerInstanceReq struct {
	Cluster string
	// arn:aws:ecs:region:aws_account_id:container-instance/container_instance_UUID.
	ContainerInstance string
	Force             bool
}

// DeregisterContainerInstanceResp encapsulates DeregisterContainerInstance response
type DeregisterContainerInstanceResp struct {
	ContainerInstance ContainerInstance `xml:"DeregisterContainerInstanceResult>containerInstance"`
	RequestId         string            `xml:"ResponseMetadata>RequestId"`
}

// DeregisterContainerInstance deregisters an Amazon ECS container instance from the specified cluster
func (e *ECS) DeregisterContainerInstance(req *DeregisterContainerInstanceReq) (
	resp *DeregisterContainerInstanceResp, err error) {
	if req == nil {
		return nil, fmt.Errorf("The req params cannot be nil")
	}

	params := makeParams("DeregisterContainerInstance")
	params["containerInstance"] = req.ContainerInstance
	params["force"] = strconv.FormatBool(req.Force)

	if req.Cluster != "" {
		params["cluster"] = req.Cluster
	}

	resp = new(DeregisterContainerInstanceResp)
	if err := e.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// PortMapping encapsulates the PortMapping data type
type PortMapping struct {
	ContainerPort int64 `xml:containerPort`
	HostPort      int64 `xml:hostPort`
}

// KeyValuePair encapsulates the KeyValuePair data type
type KeyValuePair struct {
	Name  string `xml:"name"`
	Value string `xml:"value"`
}

// ContainerDefinition encapsulates the container definition type
// Container definitions are used in task definitions to describe
// the different containers that are launched as part of a task
type ContainerDefinition struct {
	Command      []string       `xml:"command>member"`
	Cpu          int64          `xml:"cpu"`
	EntryPoint   []string       `xml:"entryPoint>member"`
	Environment  []KeyValuePair `xml:"environment>member"`
	Essential    bool           `xml:"essential"`
	Image        string         `xml:"image"`
	Links        []string       `xml:"links>member"`
	Memory       int64          `xml:"memory"`
	Name         string         `xml:"name"`
	PortMappings []PortMapping  `xml:"portMappings>member"`
}

// TaskDefinition encapsulates the task definition type
type TaskDefinition struct {
	ContainerDefinitions []ContainerDefinition `xml:"containerDefinitions>member"`
	Family               string                `xml:"family"`
	Revision             int64                 `xml:"revision"`
	TaskDefinitionArn    string                `xml:"taskDefinitionArn"`
}

// DeregisterTaskDefinitionReq encapsulates DeregisterTaskDefinitionReq req params
type DeregisterTaskDefinitionReq struct {
	TaskDefinition string
}

type DeregisterTaskDefinitionResp struct {
	TaskDefinition TaskDefinition `xml:"DeregisterTaskDefinitionResult>taskDefinition"`
	RequestId      string         `xml:"ResponseMetadata>RequestId"`
}

func (e *ECS) DeregisterTaskDefinition(req *DeregisterTaskDefinitionReq) (
	*DeregisterTaskDefinitionResp, error) {
	if req == nil {
		return nil, fmt.Errorf("The req params cannot be nil")
	}

	params := makeParams("DeregisterTaskDefinition")
	params["taskDefinition"] = req.TaskDefinition

	resp := new(DeregisterTaskDefinitionResp)
	if err := e.query(params, resp); err != nil {
		return nil, err
	}
	return resp, nil
}
