package redshift_test

import (
	"testing"

	"github.com/czos/goamz/aws"
	"github.com/czos/goamz/redshift"
	"github.com/czos/goamz/testutil"
	"github.com/motain/gocheck"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

var _ = gocheck.Suite(&S{})

type S struct {
	redshift *redshift.Redshift
}

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.redshift = redshift.NewWithClient(
		auth,
		aws.Region{RedshiftEndpoint: testServer.URL},
		testutil.DefaultClient,
	)
}

func (s *S) TearDownSuite(c *gocheck.C) {
	testServer.Stop()
}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

func (s *S) TestRunInstancesErrorDump(c *gocheck.C) {
	testServer.Response(400, nil, ErrorDump)
	options := redshift.NewDefaultedCreateClusterOptions()

	msg := `something failed`

	resp, err := s.redshift.CreateCluster(options)

	testServer.WaitRequest()

	c.Assert(resp, gocheck.IsNil)

	rs2err, ok := err.(*redshift.Error)
	c.Assert(ok, gocheck.Equals, true)
	c.Assert(rs2err.StatusCode, gocheck.Equals, 400)
	c.Assert(rs2err.Code, gocheck.Equals, "UnsupportedOperation")
	c.Assert(rs2err.Message, gocheck.Matches, msg)
	c.Assert(rs2err.RequestId, gocheck.Equals, "0503f4e9-bbd6-483c-b54f-c4ae9f3b30f4")
}

func (s *S) TestDescribeClustersExample(c *gocheck.C) {
	testServer.Response(200, nil, DescribeClustersExample1)
	resp, err := s.redshift.DescribeClusters("abc", []string{}, []string{}, "def", 20)

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DescribeClusters"})
	c.Assert(req.Form["ClusterIdentifier"], gocheck.DeepEquals, []string{"abc"})
	c.Assert(req.Form["Marker"], gocheck.DeepEquals, []string{"def"})
	c.Assert(req.Form["MaxRecords"], gocheck.DeepEquals, []string{"20"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "837d45d6-64f0-11e2-b07c-f7fbdd006c67")
	c.Assert(resp.Clusters, gocheck.HasLen, 1)

	c0 := resp.Clusters[0]
	c.Assert(c0.PendingModifiedValues, gocheck.DeepEquals,
		redshift.PendingModifiedValues{MasterUserPassword: "****"})
	c.Assert(c0.ClusterVersion, gocheck.Equals, "1.0")
	c.Assert(c0.ClusterStatus, gocheck.Equals, "creating")
	c.Assert(c0.NumberOfNodes, gocheck.Equals, 2)
	c.Assert(c0.AutomatedSnapshotRetentionPeriod, gocheck.Equals, 1)
	c.Assert(c0.PubliclyAccessible, gocheck.Equals, true)
	c.Assert(c0.Encrypted, gocheck.Equals, false)
	c.Assert(c0.DBName, gocheck.Equals, "dev")
	c.Assert(c0.PreferredMaintenanceWindow, gocheck.Equals, "sun:10:30-sun:11:00")
	c.Assert(c0.ClusterParameterGroups, gocheck.DeepEquals,
		[]redshift.ClusterParameterGroup{
			{ParameterApplyStatus: "in-sync", ParameterGroupName: "default.redshift-1.0"},
		})
	c.Assert(c0.ClusterSecurityGroups, gocheck.DeepEquals,
		[]redshift.ClusterSecurityGroup{
			{Status: "active", ClusterSecurityGroupName: "default"},
		})
	c.Assert(c0.AvailabilityZone, gocheck.Equals, "us-east-1a")
	c.Assert(c0.NodeType, gocheck.Equals, "dw1.xlarge")
	c.Assert(c0.ClusterIdentifier, gocheck.Equals, "examplecluster")
	c.Assert(c0.AllowVersionUpgrade, gocheck.Equals, true)
	c.Assert(c0.MasterUsername, gocheck.Equals, "masteruser")
}

func (s *S) TestCreateClusterExample(c *gocheck.C) {
	testServer.Response(200, nil, CreateClusterExample1)
	cco := redshift.NewDefaultedCreateClusterOptions()
	cco.ClusterIdentifier = "abc"
	cco.MasterUserPassword = "def"
	cco.MasterUsername = "ghi"
	cco.NodeType = "dw1.xlarge"
	cco.AllowVersionUpgrade = false
	cco.AutomatedSnapshotRetentionPeriod = 5
	cco.AvailabilityZone = "us-east-1a"
	cco.ClusterParameterGroupName = "abc"
	cco.ClusterSecurityGroups = []string{"sec1", "sec2"}
	cco.ClusterSubnetGroupName = "subnet123"
	cco.ClusterType = "single-node"
	cco.ClusterVersion = "1.0"
	cco.DBName = "testdb"
	cco.ElasticIp = "10.100.10.55"
	cco.Encrypted = true
	cco.HsmClientCertificateIdentifier = "abcid"
	cco.HsmConfigurationIdentifier = "hsmid"
	cco.KmsKeyId = "ksmkeyid"
	cco.NumberOfNodes = 5
	cco.Port = 8000
	cco.PreferredMaintenanceWindow = "abc"
	cco.PubliclyAccessible = true
	cco.Tags = []redshift.Tag{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"}}
	cco.VpcSecurityGroupIds = []string{"vpcsec1", "vpcsec2"}

	resp, err := s.redshift.CreateCluster(cco)
	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"CreateCluster"})
	c.Assert(req.Form["ClusterIdentifier"], gocheck.DeepEquals, []string{"abc"})
	c.Assert(req.Form["MasterUserPassword"], gocheck.DeepEquals, []string{"def"})
	c.Assert(req.Form["MasterUsername"], gocheck.DeepEquals, []string{"ghi"})
	c.Assert(req.Form["NodeType"], gocheck.DeepEquals, []string{"dw1.xlarge"})
	c.Assert(req.Form["AllowVersionUpgrade"], gocheck.DeepEquals, []string{"false"})
	c.Assert(req.Form["AutomatedSnapshotRetentionPeriod"],
		gocheck.DeepEquals, []string{"5"})
	c.Assert(req.Form["AvailabilityZone"], gocheck.DeepEquals, []string{"us-east-1a"})
	c.Assert(req.Form["ClusterParameterGroupName"], gocheck.DeepEquals, []string{"abc"})
	c.Assert(req.Form["ClusterSecurityGroups.ClusterSecurityGroupName.1"],
		gocheck.DeepEquals, []string{"sec1"})
	c.Assert(req.Form["ClusterSecurityGroups.ClusterSecurityGroupName.2"],
		gocheck.DeepEquals, []string{"sec2"})
	c.Assert(req.Form["ClusterSubnetGroupName"],
		gocheck.DeepEquals, []string{"subnet123"})
	c.Assert(req.Form["ClusterType"], gocheck.DeepEquals, []string{"single-node"})
	c.Assert(req.Form["ClusterVersion"], gocheck.DeepEquals, []string{"1.0"})
	c.Assert(req.Form["DBName"], gocheck.DeepEquals, []string{"testdb"})
	c.Assert(req.Form["ElasticIp"], gocheck.DeepEquals, []string{"10.100.10.55"})
	c.Assert(req.Form["Encrypted"], gocheck.DeepEquals, []string{"true"})
	c.Assert(req.Form["HsmClientCertificateIdentifier"],
		gocheck.DeepEquals, []string{"abcid"})
	c.Assert(req.Form["HsmConfigurationIdentifier"],
		gocheck.DeepEquals, []string{"hsmid"})
	c.Assert(req.Form["KmsKeyId"], gocheck.DeepEquals, []string{"ksmkeyid"})
	c.Assert(req.Form["NumberOfNodes"], gocheck.DeepEquals, []string{"5"})
	c.Assert(req.Form["Port"], gocheck.DeepEquals, []string{"8000"})
	c.Assert(req.Form["PreferredMaintenanceWindow"],
		gocheck.DeepEquals, []string{"abc"})
	c.Assert(req.Form["PubliclyAccessible"], gocheck.DeepEquals, []string{"true"})
	c.Assert(req.Form["Tags.Tag.1.Key"], gocheck.DeepEquals, []string{"key1"})
	c.Assert(req.Form["Tags.Tag.2.Key"], gocheck.DeepEquals, []string{"key2"})
	c.Assert(req.Form["Tags.Tag.1.Value"], gocheck.DeepEquals, []string{"value1"})
	c.Assert(req.Form["Tags.Tag.2.Value"], gocheck.DeepEquals, []string{"value2"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "e69b1294-64ef-11e2-b07c-f7fbdd006c67")

	c0 := resp.Cluster
	c.Assert(c0.PendingModifiedValues, gocheck.DeepEquals,
		redshift.PendingModifiedValues{MasterUserPassword: "****"})
	c.Assert(c0.ClusterVersion, gocheck.Equals, "1.0")
	c.Assert(c0.ClusterStatus, gocheck.Equals, "creating")
	c.Assert(c0.NumberOfNodes, gocheck.Equals, 2)
	c.Assert(c0.AutomatedSnapshotRetentionPeriod, gocheck.Equals, 1)
	c.Assert(c0.PubliclyAccessible, gocheck.Equals, true)
	c.Assert(c0.Encrypted, gocheck.Equals, false)
	c.Assert(c0.DBName, gocheck.Equals, "dev")
	c.Assert(c0.PreferredMaintenanceWindow, gocheck.Equals, "sun:10:30-sun:11:00")
	c.Assert(c0.ClusterParameterGroups, gocheck.DeepEquals,
		[]redshift.ClusterParameterGroup{
			{ParameterApplyStatus: "in-sync", ParameterGroupName: "default.redshift-1.0"},
		})
	c.Assert(c0.ClusterSecurityGroups, gocheck.DeepEquals,
		[]redshift.ClusterSecurityGroup{
			{Status: "active", ClusterSecurityGroupName: "default"},
		})
	c.Assert(c0.NodeType, gocheck.Equals, "dw1.xlarge")
	c.Assert(c0.ClusterIdentifier, gocheck.Equals, "examplecluster")
	c.Assert(c0.AllowVersionUpgrade, gocheck.Equals, true)
	c.Assert(c0.MasterUsername, gocheck.Equals, "masteruser")
}

func (s *S) TestDeleteClusterExample(c *gocheck.C) {
	testServer.Response(200, nil, DeleteClusterExample1)
	resp, err := s.redshift.DeleteCluster("abc", "def", true)

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DeleteCluster"})
	c.Assert(req.Form["ClusterIdentifier"], gocheck.DeepEquals, []string{"abc"})
	c.Assert(req.Form["FinalClusterSnapshotIdentifier"], gocheck.DeepEquals, []string{"def"})
	c.Assert(req.Form["SkipFinalClusterSnapshot"], gocheck.DeepEquals, []string{"true"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "f2e6b87e-6503-11e2-b343-393adc3f0a21")

	c0 := resp.Cluster
	c.Assert(c0.ClusterVersion, gocheck.Equals, "1.0")
	c.Assert(c0.ClusterStatus, gocheck.Equals, "deleting")
	c.Assert(c0.NumberOfNodes, gocheck.Equals, 2)
	c.Assert(c0.AutomatedSnapshotRetentionPeriod, gocheck.Equals, 1)
	c.Assert(c0.PubliclyAccessible, gocheck.Equals, true)
	c.Assert(c0.Encrypted, gocheck.Equals, true)
	c.Assert(c0.DBName, gocheck.Equals, "dev")
	c.Assert(c0.PreferredMaintenanceWindow, gocheck.Equals, "sun:10:30-sun:11:00")
	c.Assert(c0.ClusterParameterGroups, gocheck.DeepEquals,
		[]redshift.ClusterParameterGroup{
			{ParameterApplyStatus: "in-sync", ParameterGroupName: "default.redshift-1.0"},
		})
	c.Assert(c0.ClusterSecurityGroups, gocheck.DeepEquals,
		[]redshift.ClusterSecurityGroup{
			{Status: "active", ClusterSecurityGroupName: "default"},
		})
	c.Assert(c0.AvailabilityZone, gocheck.Equals, "us-east-1a")
	c.Assert(c0.NodeType, gocheck.Equals, "dw1.xlarge")
	c.Assert(c0.ClusterIdentifier, gocheck.Equals, "examplecluster2")
	c.Assert(c0.AllowVersionUpgrade, gocheck.Equals, true)
	c.Assert(c0.MasterUsername, gocheck.Equals, "masteruser")
}

func (s *S) TestCreateClusterSubnetGroupExample1(c *gocheck.C) {
	testServer.Response(200, nil, CreateClusterSubnetGroupExample1)
	tags := []redshift.Tag{
		{Key: "key1", Value: "value1"},
		{Key: "key2", Value: "value2"},
	}

	resp, err := s.redshift.CreateClusterSubnetGroup(
		"abc", "def", []string{"sub1", "sub2"}, tags)

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"],
		gocheck.DeepEquals, []string{"CreateClusterSubnetGroup"})
	c.Assert(req.Form["ClusterSubnetGroupName"], gocheck.DeepEquals, []string{"abc"})
	c.Assert(req.Form["Description"], gocheck.DeepEquals, []string{"def"})
	c.Assert(
		req.Form["SubnetIds.SubnetIdentifier.1"], gocheck.DeepEquals, []string{"sub1"})
	c.Assert(
		req.Form["SubnetIds.SubnetIdentifier.2"], gocheck.DeepEquals, []string{"sub2"})
	c.Assert(req.Form["Tags.Tag.1.Key"], gocheck.DeepEquals, []string{"key1"})
	c.Assert(req.Form["Tags.Tag.2.Key"], gocheck.DeepEquals, []string{"key2"})
	c.Assert(req.Form["Tags.Tag.1.Value"], gocheck.DeepEquals, []string{"value1"})
	c.Assert(req.Form["Tags.Tag.2.Value"], gocheck.DeepEquals, []string{"value2"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "0a60660f-6a4a-11e2-aad2-71d00c36728e")

	c0 := resp.ClusterSubnetGroup
	c.Assert(c0.ClusterSubnetGroupName, gocheck.Equals, "mysubnetgroup1")
	c.Assert(c0.VpcId, gocheck.Equals, "vpc-796a5913")
	c.Assert(c0.Description, gocheck.Equals, "My subnet group 1")
	c.Assert(c0.SubnetGroupStatus, gocheck.Equals, "Complete")
	c.Assert(c0.Subnets[0].SubnetStatus, gocheck.Equals, "Active")
	c.Assert(c0.Subnets[0].SubnetIdentifier, gocheck.Equals, "subnet-756a591f")
	c.Assert(c0.Subnets[0].SubnetAvailabilityZone.Name, gocheck.Equals, "us-east-1c")
}

func (s *S) TestDeleteClusterSubnetGroupExample1(c *gocheck.C) {
	testServer.Response(200, nil, DeleteClusterSubnetGroupExample1)
	resp, err := s.redshift.DeleteClusterSubnetGroup("abc")

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DeleteClusterSubnetGroup"})
	c.Assert(req.Form["ClusterSubnetGroupName"], gocheck.DeepEquals, []string{"abc"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "3a63806b-6af4-11e2-b27b-4d850b1c672d")
}
