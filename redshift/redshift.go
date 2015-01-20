package redshift

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/czos/goamz/aws"
)

// The Redshift type encapsulates redshfit operations within a specific region.
type Redshift struct {
	aws.Auth
	aws.Region
	httpClient *http.Client
}

// New creates a new IAM instance.
func New(auth aws.Auth, region aws.Region) *Redshift {
	return NewWithClient(auth, region, aws.RetryingClient)
}

// NewWithClient creates a new EC2 with a custom http client
func NewWithClient(
	auth aws.Auth, region aws.Region, httpClient *http.Client) *Redshift {
	return &Redshift{auth, region, httpClient}
}

// ----------------------------------------------------------------------------
// Request dispatching logic.

// Error encapsulates an error returned by Redshift
//
// See http://goo.gl/VZGuC for more details.
type Error struct {
	// HTTP status code (200, 403, ...)
	StatusCode int
	// Redshift error code ("UnsupportedOperation", ...)
	Code string
	// The human-oriented error message
	Message   string
	RequestId string
	Type      string
}

func (err *Error) Error() string {
	if err.Code == "" {
		return err.Message
	}

	return fmt.Sprintf("Type: %s, Code: %s, Message: %s",
		err.Type, err.Code, err.Message)
}

type ErrorResponse struct {
	Errors    Error  `xml:"Error"`
	RequestId string `xml:"RequestID"`
}

var timeNow = time.Now

func (rs *Redshift) query(params map[string]string, resp interface{}) error {
	params["Version"] = "2012-12-01"

	// Create the request
	endpoint, err := url.Parse(rs.Region.RedshiftEndpoint)
	if err != nil {
		return err
	}
	if endpoint.Path == "" {
		endpoint.Path = "/"
	}
	endpoint.RawQuery = multimap(params).Encode()
	req, err := http.NewRequest("GET", endpoint.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-Amz-Date", time.Now().UTC().Format(aws.ISO8601BasicFormat))

	// sign the request
	signer := aws.NewV4Signer(rs.Auth, "redshift", rs.Region)
	signer.Sign(req)

	log.SetLevel(log.DebugLevel)

	// make the request
	//log.Debugf("GET %v", endpoint.String())
	d, _ := httputil.DumpRequest(req, true)
	log.Debugf("Request\n%v", string(d))
	r, err := rs.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	dump := []byte{}
	if tracingEnabled() {
		dump, _ = httputil.DumpResponse(r, true)
	} else {
		dump, _ = httputil.DumpResponse(r, false)
	}
	log.Debugf("%v\n", string(dump))

	if r.StatusCode != 200 {
		return buildError(r)
	}
	err = xml.NewDecoder(r.Body).Decode(resp)
	return err
}

func multimap(p map[string]string) url.Values {
	q := make(url.Values, len(p))
	for k, v := range p {
		q[k] = []string{v}
	}
	return q
}

func buildError(r *http.Response) error {
	errors := ErrorResponse{}
	xml.NewDecoder(r.Body).Decode(&errors)
	var err Error
	err = errors.Errors
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

func addParamsList(params map[string]string, label string, ids []string) {
	for i, id := range ids {
		params[label+"."+strconv.Itoa(i+1)] = id
	}
}

func tracingEnabled() bool {
	t, err := strconv.ParseBool(os.Getenv("TRACE"))
	if err != nil {
		return false
	}
	return t
}

// ----------------------------------------------------------------------------
// Redshift cluster management functions and types

// CreateClusterOptions encapsulates options for creating a redshift cluster
//
// See http://goo.gl/mnsUTI for more details.
type CreateClusterOptions struct {
	ClusterIdentifier                string
	MasterUserPassword               string
	MasterUsername                   string
	NodeType                         string
	AllowVersionUpgrade              bool
	AutomatedSnapshotRetentionPeriod int
	AvailabilityZone                 string
	ClusterParameterGroupName        string
	ClusterSecurityGroups            []string
	ClusterSubnetGroupName           string
	ClusterType                      string
	ClusterVersion                   string
	DBName                           string
	ElasticIp                        string
	Encrypted                        bool
	HsmClientCertificateIdentifier   string
	HsmConfigurationIdentifier       string
	KmsKeyId                         string
	NumberOfNodes                    int
	Port                             int
	PreferredMaintenanceWindow       string
	PubliclyAccessible               bool
	Tags                             []Tag
	VpcSecurityGroupIds              []string
}

// DefaultedCreateClusterOptions creates an instance of
// CreateClusterOptions with default values
//
// See http://goo.gl/mnsUTI for more details.
func NewDefaultedCreateClusterOptions() *CreateClusterOptions {
	return &CreateClusterOptions{
		AllowVersionUpgrade:              true,
		AutomatedSnapshotRetentionPeriod: 1,
		Encrypted:                        false,
		NumberOfNodes:                    1,
		Port:                             5439,
		PubliclyAccessible:               false,
	}
}

// CreateClusterResp encapsulates a response to CreateCluster request
//
// See http://goo.gl/mnsUTI for more details.
type CreateClusterResp struct {
	RequestId string  `xml:"ResponseMetadata>RequestId"`
	Cluster   Cluster `xml:"CreateClusterResult>Cluster"`
}

// DeleteClusterResp encapsulates a response to DeleteCluster request
//
// See http://goo.gl/mnsUTI for more details.
type DeleteClusterResp struct {
	RequestId string  `xml:"ResponseMetadata>RequestId"`
	Cluster   Cluster `xml:"DeleteClusterResult>Cluster"`
}

// CreateClusterSubnetGroupResp encapsulates a
// response to CreateClusterSubnetGroup request
//
// See http://goo.gl/j1BoQG for more details.
type CreateClusterSubnetGroupResp struct {
	RequestId          string             `xml:"ResponseMetadata>RequestId"`
	ClusterSubnetGroup ClusterSubnetGroup `xml:"CreateClusterSubnetGroupResult>ClusterSubnetGroup"`
}

// DeleteClusterSubnetGroupResp encapsulates a response to
// a DeleteClusterSubnetGroup request
//
// See http://goo.gl/CDYQnX for more details.
type DeleteClusterSubnetGroupResp struct {
	RequestId string `xml:"ResponseMetadata>RequestId"`
}

// Cluster resource
//
// See http://goo.gl/mnsUTI for more details.
type Cluster struct {
	AllowVersionUpgrade              bool                         `xml:"AllowVersionUpgrade"`
	AutomatedSnapshotRetentionPeriod int                          `xml:"AutomatedSnapshotRetentionPeriod"`
	AvailabilityZone                 string                       `xml:"AvailabilityZone"`
	ClusterCreateTime                string                       `xml:"ClusterCreateTime"`
	ClusterIdentifier                string                       `xml:"ClusterIdentifier"`
	ClusterNodes                     []ClusterNode                `xml:"ClusterNodes>ClusterNode"`
	ClusterParameterGroups           []ClusterParameterGroup      `xml:"ClusterParameterGroups>ClusterParameterGroup"`
	ClusterPublicKey                 string                       `xml:"ClusterPublicKey"`
	ClusterRevisionNumber            string                       `xml:"ClusterRevisionNumber"`
	ClusterSecurityGroups            []ClusterSecurityGroup       `xml:"ClusterSecurityGroups>ClusterSecurityGroup"`
	ClusterSnapshotCopyStatus        ClusterSnapshotCopyStatus    `xml:"ClusterSnapshotCopyStatus"`
	ClusterStatus                    string                       `xml:"ClusterStatus"`
	ClusterSubnetGroupName           string                       `xml:"ClusterSubnetGroupName"`
	ClusterVersion                   string                       `xml:"ClusterVersion"`
	DBName                           string                       `xml:"DBName"`
	ElasticIpStatus                  ElasticIpStatus              `xml:"ElasticIpStatus"`
	Encrypted                        bool                         `xml:"Encrypted"`
	Endpoint                         Endpoint                     `xml:"Endpoint"`
	HsmStatus                        HsmStatus                    `xml:"HsmStatus"`
	KmsKeyId                         string                       `xml:"KmsKeyId"`
	MasterUsername                   string                       `xml:"MasterUsername"`
	ModifyStatus                     string                       `xml:"ModifyStatus"`
	NodeType                         string                       `xml:"NodeType"`
	NumberOfNodes                    int                          `xml:"NumberOfNodes"`
	PendingModifiedValues            PendingModifiedValues        `xml:"PendingModifiedValues"`
	PreferredMaintenanceWindow       string                       `xml:"PreferredMaintenanceWindow"`
	PubliclyAccessible               bool                         `xml:"PubliclyAccessible"`
	RestoreStatus                    RestoreStatus                `xml:"RestoreStatus"`
	Tags                             []Tag                        `xml:"Tags"`
	VpcId                            string                       `xml:"VpcId"`
	VpcSecurityGroups                []VpcSecurityGroupMembership `xml:"VpcSecurityGroups"`
}

type ClusterNode struct {
	NodeRole         string `xml:"NodeRole"`
	PrivateIPAddress string `xml:"PrivateIPAddress"`
	PublicIPAddress  string `xml:"PublicIPAddress"`
}

type ClusterParameterGroup struct {
	ParameterApplyStatus string `xml:"ParameterApplyStatus"`
	ParameterGroupName   string `xml:"ParameterGroupName"`
}

type ClusterSecurityGroup struct {
	Status                   string `xml:"Status"`
	ClusterSecurityGroupName string `xml:"ClusterSecurityGroupName"`
}

type ClusterSnapshotCopyStatus struct {
	DestinationRegion string `xml:"DestinationRegion"`
	RetentionPeriod   int64  `xml:"RetentionPeriod"`
}

type ElasticIpStatus struct {
	ElasticIp string `xml:"ElasticIp"`
	Status    string `xml:"Status"`
}

type Endpoint struct {
	Address string `xml:"Address"`
	Port    int    `xml:"Port"`
}

type HsmStatus struct {
	HsmClientCertificateIdentifier string `xml:"HsmClientCertificateIdentifier"`
	HsmConfigurationIdentifier     string `xml:"HsmConfigurationIdentifier"`
	Status                         string `xml:"Status"`
}

type PendingModifiedValues struct {
	AutomatedSnapshotRetentionPeriod int    `xml:"AutomatedSnapshotRetentionPeriod"`
	ClusterIdentifier                string `xml:"ClusterIdentifier"`
	ClusterType                      string `xml:"ClusterType"`
	ClusterVersion                   string `xml:"ClusterVersion"`
	MasterUserPassword               string `xml:"MasterUserPassword"`
	NodeType                         string `xml:"NodeType"`
	NumberOfNodes                    int    `xml:"NumberOfNodes"`
}

type RestoreStatus struct {
	CurrentRestoreRateInMegaBytesPerSecond float64 `xml:"CurrentRestoreRateInMegaBytesPerSecond"`
	ElapsedTimeInSeconds                   int64   `xml:"ElapsedTimeInSeconds"`
	EstimatedTimeToCompletionInSeconds     int64   `xml:"EstimatedTimeToCompletionInSeconds"`
	ProgressInMegaBytes                    int64   `xml:"ProgressInMegaBytes"`
	SnapshotSizeInMegaBytes                int64   `xml:"SnapshotSizeInMegaBytes"`
	Status                                 string  `xml:"Status"`
}

type VpcSecurityGroupMembership struct {
	Status             string `xml:"Status"`
	VpcSecurityGroupId string `xml:"VpcSecurityGroupId"`
}

type Tag struct {
	Key   string `xml:"Key"`
	Value string `xml:"Value"`
}

type ClusterSubnetGroup struct {
	ClusterSubnetGroupName string
	Description            string
	SubnetGroupStatus      string
	Subnets                []Subnet `xml:"Subnets>Subnet"`
	Tags                   []Tag    `xml:"Tags>Tag"`
	VpcId                  string
}

type Subnet struct {
	SubnetAvailabilityZone AvailabilityZone
	SubnetIdentifier       string
	SubnetStatus           string
}

type AvailabilityZone struct {
	Name string
}

// DescribeClustersResp encapsulates a response to DescribeClusters request
//
// See http://goo.gl/iNXhTr for more details.
type DescribeClustersResp struct {
	RequestId string    `xml:"ResponseMetadata>RequestId"`
	Clusters  []Cluster `xml:"DescribeClustersResult>Clusters>Cluster"`
}

// DescribeClusters returns details about redshift clusters.  All parameters
// are optional, and if provided will limit the clusters returned to those
// matching parameters
//
// See http://goo.gl/iNXhTr or more details.
func (rs *Redshift) DescribeClusters(
	clusterIdentifier string,
	tagKeys []string,
	tagValues []string,
	marker string,
	maxRecords int) (*DescribeClustersResp, error) {

	params := makeParams("DescribeClusters")
	if clusterIdentifier != "" {
		params["ClusterIdentifier"] = clusterIdentifier
	}

	if maxRecords >= 20 && maxRecords <= 100 {
		params["MaxRecords"] = strconv.Itoa(maxRecords)
	}

	if marker != "" {
		params["Marker"] = marker
	}

	if len(tagKeys) > 0 {
		addParamsList(params, "TagKeys.TagKey", tagKeys)
	}

	if len(tagValues) > 0 {
		addParamsList(params, "TagValues.TagValue", tagValues)
	}

	resp := &DescribeClustersResp{}
	err := rs.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateCluster creates a new redshift cluster.
//
// See http://goo.gl/iNXhTr or more details.
func (rs *Redshift) CreateCluster(
	options *CreateClusterOptions) (*CreateClusterResp, error) {

	// Required params
	params := makeParams("CreateCluster")
	params["ClusterIdentifier"] = options.ClusterIdentifier
	params["MasterUserPassword"] = options.MasterUserPassword
	params["MasterUsername"] = options.MasterUsername
	params["NodeType"] = options.NodeType

	// Optional params
	params["AllowVersionUpgrade"] = strconv.FormatBool(options.AllowVersionUpgrade)
	params["AutomatedSnapshotRetentionPeriod"] = strconv.Itoa(options.AutomatedSnapshotRetentionPeriod)

	if options.AvailabilityZone != "" {
		params["AvailabilityZone"] = options.AvailabilityZone
	}

	if options.ClusterParameterGroupName != "" {
		params["ClusterParameterGroupName"] = options.ClusterParameterGroupName
	}

	if len(options.ClusterSecurityGroups) > 0 {
		addParamsList(params, "ClusterSecurityGroups.ClusterSecurityGroupName", options.ClusterSecurityGroups)
	}

	if options.ClusterSubnetGroupName != "" {
		params["ClusterSubnetGroupName"] = options.ClusterSubnetGroupName
	}

	if options.ClusterType != "" {
		params["ClusterType"] = options.ClusterType
	}

	if options.ClusterVersion != "" {
		params["ClusterVersion"] = options.ClusterVersion
	}

	if options.DBName != "" {
		params["DBName"] = options.DBName
	}

	if options.ElasticIp != "" {
		params["ElasticIp"] = options.ElasticIp
	}
	params["Encrypted"] = strconv.FormatBool(options.Encrypted)

	if options.HsmClientCertificateIdentifier != "" {
		params["HsmClientCertificateIdentifier"] = options.HsmClientCertificateIdentifier
	}

	if options.HsmConfigurationIdentifier != "" {
		params["HsmConfigurationIdentifier"] = options.HsmConfigurationIdentifier
	}

	if options.KmsKeyId != "" {
		params["KmsKeyId"] = options.KmsKeyId
	}

	if options.KmsKeyId != "" {
		params["KmsKeyId"] = options.KmsKeyId
	}

	params["NumberOfNodes"] = strconv.Itoa(options.NumberOfNodes)
	params["Port"] = strconv.Itoa(options.Port)

	if options.PreferredMaintenanceWindow != "" {
		params["PreferredMaintenanceWindow"] = options.PreferredMaintenanceWindow
	}

	params["Encrypted"] = strconv.FormatBool(options.Encrypted)
	params["PubliclyAccessible"] = strconv.FormatBool(options.PubliclyAccessible)

	for j, tag := range options.Tags {
		params["Tags.Tag."+strconv.Itoa(j+1)+".Key"] = tag.Key
		params["Tags.Tag."+strconv.Itoa(j+1)+".Value"] = tag.Value
	}

	if len(options.VpcSecurityGroupIds) > 0 {
		addParamsList(
			params, "VpcSecurityGroupIds.VpcSecurityGroupId", options.VpcSecurityGroupIds)
	}

	resp := &CreateClusterResp{}
	err := rs.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteCluster deletes an existing redshift cluster
//
// See http://goo.gl/iNXhTr or more details.
func (rs *Redshift) DeleteCluster(
	clusterIdentifier string,
	finalClusterSnapshotIdentifier string,
	skipFinalClusterSnapshot bool) (*DeleteClusterResp, error) {

	params := makeParams("DeleteCluster")
	if clusterIdentifier != "" {
		params["ClusterIdentifier"] = clusterIdentifier
	}

	if finalClusterSnapshotIdentifier != "" {
		params["FinalClusterSnapshotIdentifier"] = finalClusterSnapshotIdentifier
	}

	params["SkipFinalClusterSnapshot"] = strconv.FormatBool(skipFinalClusterSnapshot)

	resp := &DeleteClusterResp{}
	err := rs.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateClusterSubnetGroup creates a CreateClusterSubnetGroup
//
// See http://goo.gl/j1BoQG for more details.
func (rs *Redshift) CreateClusterSubnetGroup(
	clusterSubnetGroupName string,
	description string,
	subnetIds []string,
	tags []Tag) (*CreateClusterSubnetGroupResp, error) {

	params := makeParams("CreateClusterSubnetGroup")
	if clusterSubnetGroupName != "" {
		params["ClusterSubnetGroupName"] = clusterSubnetGroupName
	}

	if description != "" {
		params["Description"] = description
	}

	if len(subnetIds) > 0 {
		addParamsList(params, "SubnetIds.SubnetIdentifier", subnetIds)
	}

	for j, tag := range tags {
		params["Tags.Tag."+strconv.Itoa(j+1)+".Key"] = tag.Key
		params["Tags.Tag."+strconv.Itoa(j+1)+".Value"] = tag.Value
	}

	resp := &CreateClusterSubnetGroupResp{}
	err := rs.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteClusterSubnetGroup deletes an existing cluster subnet group
//
// See http://goo.gl/CDYQnX or more details.
func (rs *Redshift) DeleteClusterSubnetGroup(
	clusterSubnetGroupName string) (*DeleteClusterSubnetGroupResp, error) {

	params := makeParams("DeleteClusterSubnetGroup")
	if clusterSubnetGroupName != "" {
		params["ClusterSubnetGroupName"] = clusterSubnetGroupName
	}

	resp := &DeleteClusterSubnetGroupResp{}
	err := rs.query(params, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
