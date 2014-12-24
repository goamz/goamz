package ecs

import (
	"testing"

	"github.com/motain/gocheck"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/testutil"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

var _ = gocheck.Suite(&S{})

type S struct {
	ecs *ECS
}

var testServer = testutil.NewHTTPServer()

var mockTest bool

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.ecs = New(auth, aws.Region{ECSEndpoint: testServer.URL})
}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

// --------------------------------------------------------------------------
// Detailed Unit Tests

func (s *S) TestCreateCluster(c *gocheck.C) {
	testServer.Response(200, nil, CreateClusterResponse)
	req := &CreateClusterReq{
		ClusterName: "default",
	}
	resp, err := s.ecs.CreateCluster(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "CreateCluster")
	c.Assert(values.Get("clusterName"), gocheck.Equals, "default")

	c.Assert(resp.Cluster.ClusterArn, gocheck.Equals, "arn:aws:ecs:region:aws_account_id:cluster/default")
	c.Assert(resp.Cluster.ClusterName, gocheck.Equals, "default")
	c.Assert(resp.Cluster.Status, gocheck.Equals, "ACTIVE")
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestDeregisterContainerInstance(c *gocheck.C) {
	testServer.Response(200, nil, DeregisterContainerInstanceResponse)
	req := &DeregisterContainerInstanceReq{
		Cluster:           "default",
		ContainerInstance: "uuid",
		Force:             true,
	}
	resp, err := s.ecs.DeregisterContainerInstance(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "DeregisterContainerInstance")
	c.Assert(values.Get("cluster"), gocheck.Equals, "default")
	c.Assert(values.Get("containerInstance"), gocheck.Equals, "uuid")
	c.Assert(values.Get("force"), gocheck.Equals, "true")

	expectedResource := []Resource{
		{
			DoubleValue:  0.0,
			IntegerValue: 2048,
			LongValue:    0,
			Name:         "CPU",
			Type:         "INTEGER",
		},
		{
			DoubleValue:  0.0,
			IntegerValue: 3955,
			LongValue:    0,
			Name:         "MEMORY",
			Type:         "INTEGER",
		},
		{
			DoubleValue:    0.0,
			IntegerValue:   0,
			LongValue:      0,
			Name:           "PORTS",
			StringSetValue: []string{"2376", "22", "51678", "2375"},
			Type:           "STRINGSET",
		},
	}

	c.Assert(resp.ContainerInstance.AgentConnected, gocheck.Equals, false)
	c.Assert(resp.ContainerInstance.ContainerInstanceArn, gocheck.Equals, "arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID")
	c.Assert(resp.ContainerInstance.Status, gocheck.Equals, "INACTIVE")
	c.Assert(resp.ContainerInstance.Ec2InstanceId, gocheck.Equals, "instance_id")
	c.Assert(resp.ContainerInstance.RegisteredResources, gocheck.DeepEquals, expectedResource)
	c.Assert(resp.ContainerInstance.RemainingResources, gocheck.DeepEquals, expectedResource)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestDeregisterTaskDefinition(c *gocheck.C) {
	testServer.Response(200, nil, DeregisterTaskDefinitionResponse)
	req := &DeregisterTaskDefinitionReq{
		TaskDefinition: "sleep360:2",
	}
	resp, err := s.ecs.DeregisterTaskDefinition(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "DeregisterTaskDefinition")
	c.Assert(values.Get("taskDefinition"), gocheck.Equals, "sleep360:2")

	expected := TaskDefinition{
		Family:            "sleep360",
		Revision:          2,
		TaskDefinitionArn: "arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2",
		ContainerDefinitions: []ContainerDefinition{
			{
				Command:    []string{"sleep", "360"},
				Cpu:        10,
				EntryPoint: []string{"/bin/sh"},
				Environment: []KeyValuePair{
					{
						Name:  "envVar",
						Value: "foo",
					},
				},
				Essential: true,
				Image:     "busybox",
				Memory:    10,
				Name:      "sleep",
			},
		},
	}

	c.Assert(resp.TaskDefinition, gocheck.DeepEquals, expected)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestDescribeClusters(c *gocheck.C) {
	testServer.Response(200, nil, DescribeClustersResponse)
	req := &DescribeClustersReq{
		Clusters: []string{"test", "default"},
	}
	resp, err := s.ecs.DescribeClusters(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeClusters")
	c.Assert(values.Get("clusters.member.1"), gocheck.Equals, "test")
	c.Assert(values.Get("clusters.member.2"), gocheck.Equals, "default")

	expected := []Cluster{
		{
			ClusterName: "test",
			ClusterArn:  "arn:aws:ecs:us-east-1:aws_account_id:cluster/test",
			Status:      "ACTIVE",
		},
		{
			ClusterName: "default",
			ClusterArn:  "arn:aws:ecs:us-east-1:aws_account_id:cluster/default",
			Status:      "ACTIVE",
		},
	}

	c.Assert(resp.Clusters, gocheck.DeepEquals, expected)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestDescribeContainerInstances(c *gocheck.C) {
	testServer.Response(200, nil, DescribeContainerInstancesResponse)
	req := &DescribeContainerInstancesReq{
		Cluster:            "test",
		ContainerInstances: []string{"arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID"},
	}
	resp, err := s.ecs.DescribeContainerInstances(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeContainerInstances")
	c.Assert(values.Get("cluster"), gocheck.Equals, "test")
	c.Assert(values.Get("containerInstances.member.1"),
		gocheck.Equals, "arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID")

	expected := []ContainerInstance{
		ContainerInstance{
			AgentConnected:       true,
			ContainerInstanceArn: "arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID",
			Status:               "ACTIVE",
			Ec2InstanceId:        "instance_id",
			RegisteredResources: []Resource{
				{
					DoubleValue:  0.0,
					IntegerValue: 2048,
					LongValue:    0,
					Name:         "CPU",
					Type:         "INTEGER",
				},
				{
					DoubleValue:  0.0,
					IntegerValue: 3955,
					LongValue:    0,
					Name:         "MEMORY",
					Type:         "INTEGER",
				},
				{
					DoubleValue:    0.0,
					IntegerValue:   0,
					LongValue:      0,
					Name:           "PORTS",
					StringSetValue: []string{"2376", "22", "51678", "2375"},
					Type:           "STRINGSET",
				},
			},
			RemainingResources: []Resource{
				{
					DoubleValue:  0.0,
					IntegerValue: 2048,
					LongValue:    0,
					Name:         "CPU",
					Type:         "INTEGER",
				},
				{
					DoubleValue:  0.0,
					IntegerValue: 3955,
					LongValue:    0,
					Name:         "MEMORY",
					Type:         "INTEGER",
				},
				{
					DoubleValue:    0.0,
					IntegerValue:   0,
					LongValue:      0,
					Name:           "PORTS",
					StringSetValue: []string{"2376", "22", "51678", "2375"},
					Type:           "STRINGSET",
				},
			},
		},
	}

	c.Assert(resp.ContainerInstances, gocheck.DeepEquals, expected)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestDescribeTaskDefinition(c *gocheck.C) {
	testServer.Response(200, nil, DescribeTaskDefinitionResponse)
	req := &DescribeTaskDefinitionReq{
		TaskDefinition: "sleep360:2",
	}
	resp, err := s.ecs.DescribeTaskDefinition(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeTaskDefinition")
	c.Assert(values.Get("taskDefinition"), gocheck.Equals, "sleep360:2")

	expected := TaskDefinition{
		Family:            "sleep360",
		Revision:          2,
		TaskDefinitionArn: "arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2",
		ContainerDefinitions: []ContainerDefinition{
			{
				Command:    []string{"sleep", "360"},
				Cpu:        10,
				EntryPoint: []string{"/bin/sh"},
				Environment: []KeyValuePair{
					{
						Name:  "envVar",
						Value: "foo",
					},
				},
				Essential: true,
				Image:     "busybox",
				Memory:    10,
				Name:      "sleep",
			},
		},
	}

	c.Assert(resp.TaskDefinition, gocheck.DeepEquals, expected)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestDescribeTasks(c *gocheck.C) {
	testServer.Response(200, nil, DescribeTasksResponse)
	req := &DescribeTasksReq{
		Cluster: "test",
		Tasks:   []string{"arn:aws:ecs:us-east-1:aws_account_id:task/UUID"},
	}
	resp, err := s.ecs.DescribeTasks(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeTasks")
	c.Assert(values.Get("cluster"), gocheck.Equals, "test")
	c.Assert(values.Get("tasks.member.1"),
		gocheck.Equals, "arn:aws:ecs:us-east-1:aws_account_id:task/UUID")

	expected := []Task{
		Task{
			Containers: []Container{
				{
					TaskArn:      "arn:aws:ecs:us-east-1:aws_account_id:task/UUID",
					Name:         "sleep",
					ContainerArn: "arn:aws:ecs:us-east-1:aws_account_id:container/UUID",
					LastStatus:   "RUNNING",
				},
			},
			Overrides: TaskOverride{
				ContainerOverrides: []ContainerOverride{
					{
						Name: "sleep",
					},
				},
			},
			DesiredStatus:        "RUNNING",
			TaskArn:              "arn:aws:ecs:us-east-1:aws_account_id:task/UUID",
			ContainerInstanceArn: "arn:aws:ecs:us-east-1:aws_account_id:container-instance/UUID",
			LastStatus:           "RUNNING",
			TaskDefinitionArn:    "arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2",
		},
	}

	c.Assert(resp.Tasks, gocheck.DeepEquals, expected)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestDiscoverPollEndpoint(c *gocheck.C) {
	testServer.Response(200, nil, DiscoverPollEndpointResponse)
	req := &DiscoverPollEndpointReq{
		ContainerInstance: "arn:aws:ecs:us-east-1:aws_account_id:container-instance/UUID",
	}
	resp, err := s.ecs.DiscoverPollEndpoint(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "DiscoverPollEndpoint")
	c.Assert(values.Get("containerInstance"),
		gocheck.Equals, "arn:aws:ecs:us-east-1:aws_account_id:container-instance/UUID")

	c.Assert(resp.Endpoint, gocheck.Equals, "https://ecs-x-1.us-east-1.amazonaws.com/")
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestListClusters(c *gocheck.C) {
	testServer.Response(200, nil, ListClustersResponse)
	req := &ListClustersReq{
		MaxResults: 2,
		NextToken:  "Token_UUID",
	}
	resp, err := s.ecs.ListClusters(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "ListClusters")
	c.Assert(values.Get("maxResults"), gocheck.Equals, "2")
	c.Assert(values.Get("nextToken"), gocheck.Equals, "Token_UUID")

	c.Assert(resp.ClusterArns, gocheck.DeepEquals, []string{"arn:aws:ecs:us-east-1:aws_account_id:cluster/default",
		"arn:aws:ecs:us-east-1:aws_account_id:cluster/test"})
	c.Assert(resp.NextToken, gocheck.Equals, "token_UUID")
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestListContainerInstances(c *gocheck.C) {
	testServer.Response(200, nil, ListContainerInstancesResponse)
	req := &ListContainerInstancesReq{
		MaxResults: 2,
		NextToken:  "Token_UUID",
		Cluster:    "test",
	}
	resp, err := s.ecs.ListContainerInstances(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "ListContainerInstances")
	c.Assert(values.Get("maxResults"), gocheck.Equals, "2")
	c.Assert(values.Get("cluster"), gocheck.Equals, "test")
	c.Assert(values.Get("nextToken"), gocheck.Equals, "Token_UUID")

	c.Assert(resp.ContainerInstanceArns, gocheck.DeepEquals, []string{
		"arn:aws:ecs:us-east-1:aws_account_id:container-instance/uuid-1",
		"arn:aws:ecs:us-east-1:aws_account_id:container-instance/uuid-2"})
	c.Assert(resp.NextToken, gocheck.Equals, "token_UUID")
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestListTaskDefinitions(c *gocheck.C) {
	testServer.Response(200, nil, ListTaskDefinitionsResponse)
	req := &ListTaskDefinitionsReq{
		MaxResults:   2,
		NextToken:    "Token_UUID",
		FamilyPrefix: "sleep360",
	}
	resp, err := s.ecs.ListTaskDefinitions(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "ListTaskDefinitions")
	c.Assert(values.Get("maxResults"), gocheck.Equals, "2")
	c.Assert(values.Get("familyPrefix"), gocheck.Equals, "sleep360")
	c.Assert(values.Get("nextToken"), gocheck.Equals, "Token_UUID")

	c.Assert(resp.TaskDefinitionArns, gocheck.DeepEquals, []string{
		"arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:1",
		"arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2"})
	c.Assert(resp.NextToken, gocheck.Equals, "token_UUID")
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestListTasks(c *gocheck.C) {
	testServer.Response(200, nil, ListTasksResponse)
	req := &ListTasksReq{
		MaxResults:        2,
		NextToken:         "Token_UUID",
		Family:            "sleep360",
		Cluster:           "test",
		ContainerInstance: "container_uuid",
	}
	resp, err := s.ecs.ListTasks(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "ListTasks")
	c.Assert(values.Get("maxResults"), gocheck.Equals, "2")
	c.Assert(values.Get("family"), gocheck.Equals, "sleep360")
	c.Assert(values.Get("containerInstance"), gocheck.Equals, "container_uuid")
	c.Assert(values.Get("cluster"), gocheck.Equals, "test")
	c.Assert(values.Get("nextToken"), gocheck.Equals, "Token_UUID")

	c.Assert(resp.TaskArns, gocheck.DeepEquals, []string{
		"arn:aws:ecs:us-east-1:aws_account_id:task/uuid_1",
		"arn:aws:ecs:us-east-1:aws_account_id:task/uuid_2"})
	c.Assert(resp.NextToken, gocheck.Equals, "token_UUID")
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestRegisterContainerInstance(c *gocheck.C) {
	testServer.Response(200, nil, RegisterContainerInstanceResponse)

	resources := []Resource{
		{
			DoubleValue:  0.0,
			IntegerValue: 2048,
			LongValue:    0,
			Name:         "CPU",
			Type:         "INTEGER",
		},
		{
			DoubleValue:  0.0,
			IntegerValue: 3955,
			LongValue:    0,
			Name:         "MEMORY",
			Type:         "INTEGER",
		},
		{
			DoubleValue:    0.0,
			IntegerValue:   0,
			LongValue:      0,
			Name:           "PORTS",
			StringSetValue: []string{"2376", "22", "51678", "2375"},
			Type:           "STRINGSET",
		},
	}

	req := &RegisterContainerInstanceReq{
		Cluster:                           "default",
		InstanceIdentityDocument:          "foo",
		InstanceIdentityDocumentSignature: "baz",
		TotalResources:                    resources,
	}

	resp, err := s.ecs.RegisterContainerInstance(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "RegisterContainerInstance")
	c.Assert(values.Get("cluster"), gocheck.Equals, "default")
	c.Assert(values.Get("instanceIdentityDocument"), gocheck.Equals, "foo")
	c.Assert(values.Get("instanceIdentityDocumentSignature"), gocheck.Equals, "baz")
	c.Assert(values.Get("totalResources.member.1.doubleValue"), gocheck.Equals, "0.0")
	c.Assert(values.Get("totalResources.member.1.integerValue"), gocheck.Equals, "2048")
	c.Assert(values.Get("totalResources.member.1.longValue"), gocheck.Equals, "0")
	c.Assert(values.Get("totalResources.member.1.name"), gocheck.Equals, "CPU")
	c.Assert(values.Get("totalResources.member.1.type"), gocheck.Equals, "INTEGER")
	c.Assert(values.Get("totalResources.member.2.doubleValue"), gocheck.Equals, "0.0")
	c.Assert(values.Get("totalResources.member.2.integerValue"), gocheck.Equals, "3955")
	c.Assert(values.Get("totalResources.member.2.longValue"), gocheck.Equals, "0")
	c.Assert(values.Get("totalResources.member.2.name"), gocheck.Equals, "MEMORY")
	c.Assert(values.Get("totalResources.member.2.type"), gocheck.Equals, "INTEGER")
	c.Assert(values.Get("totalResources.member.3.doubleValue"), gocheck.Equals, "0.0")
	c.Assert(values.Get("totalResources.member.3.integerValue"), gocheck.Equals, "0")
	c.Assert(values.Get("totalResources.member.3.longValue"), gocheck.Equals, "0")
	c.Assert(values.Get("totalResources.member.3.name"), gocheck.Equals, "PORTS")
	c.Assert(values.Get("totalResources.member.3.stringSetValue.member.1"), gocheck.Equals, "2376")
	c.Assert(values.Get("totalResources.member.3.stringSetValue.member.2"), gocheck.Equals, "22")
	c.Assert(values.Get("totalResources.member.3.stringSetValue.member.3"), gocheck.Equals, "51678")
	c.Assert(values.Get("totalResources.member.3.stringSetValue.member.4"), gocheck.Equals, "2375")
	c.Assert(values.Get("totalResources.member.3.type"), gocheck.Equals, "STRINGSET")

	c.Assert(resp.ContainerInstance.AgentConnected, gocheck.Equals, true)
	c.Assert(resp.ContainerInstance.ContainerInstanceArn, gocheck.Equals, "arn:aws:ecs:us-east-1:aws_account_id:container-instance/container_instance_UUID")
	c.Assert(resp.ContainerInstance.Status, gocheck.Equals, "ACTIVE")
	c.Assert(resp.ContainerInstance.Ec2InstanceId, gocheck.Equals, "instance_id")
	c.Assert(resp.ContainerInstance.RegisteredResources, gocheck.DeepEquals, resources)
	c.Assert(resp.ContainerInstance.RemainingResources, gocheck.DeepEquals, resources)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}

func (s *S) TestRegisterTaskDefinition(c *gocheck.C) {
	testServer.Response(200, nil, RegisterTaskDefinitionResponse)

	CDefinitions := []ContainerDefinition{
		{
			Command:    []string{"sleep", "360"},
			Cpu:        10,
			EntryPoint: []string{"/bin/sh"},
			Environment: []KeyValuePair{
				{
					Name:  "envVar",
					Value: "foo",
				},
			},
			Essential: true,
			Image:     "busybox",
			Memory:    10,
			Name:      "sleep",
		},
	}

	req := &RegisterTaskDefinitionReq{
		Family:               "sleep360",
		ContainerDefinitions: CDefinitions,
	}
	resp, err := s.ecs.RegisterTaskDefinition(req)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	c.Assert(values.Get("Version"), gocheck.Equals, "2014-11-13")
	c.Assert(values.Get("Action"), gocheck.Equals, "RegisterTaskDefinition")
	c.Assert(values.Get("containerDefinitions.member.1.command.member.1"), gocheck.Equals, "sleep")
	c.Assert(values.Get("containerDefinitions.member.1.command.member.2"), gocheck.Equals, "360")
	c.Assert(values.Get("containerDefinitions.member.1.cpu"), gocheck.Equals, "10")
	c.Assert(values.Get("containerDefinitions.member.1.memory"), gocheck.Equals, "10")
	c.Assert(values.Get("containerDefinitions.member.1.entryPoint.member.1"), gocheck.Equals, "/bin/sh")
	c.Assert(values.Get("containerDefinitions.member.1.environment.member.1.name"), gocheck.Equals, "envVar")
	c.Assert(values.Get("containerDefinitions.member.1.environment.member.1.value"), gocheck.Equals, "foo")
	c.Assert(values.Get("containerDefinitions.member.1.essential"), gocheck.Equals, "true")
	c.Assert(values.Get("containerDefinitions.member.1.image"), gocheck.Equals, "busybox")
	c.Assert(values.Get("containerDefinitions.member.1.memory"), gocheck.Equals, "10")
	c.Assert(values.Get("containerDefinitions.member.1.name"), gocheck.Equals, "sleep")
	c.Assert(values.Get("family"), gocheck.Equals, "sleep360")

	expected := TaskDefinition{
		Family:               "sleep360",
		Revision:             2,
		TaskDefinitionArn:    "arn:aws:ecs:us-east-1:aws_account_id:task-definition/sleep360:2",
		ContainerDefinitions: CDefinitions,
	}

	c.Assert(resp.TaskDefinition, gocheck.DeepEquals, expected)
	c.Assert(resp.RequestId, gocheck.Equals, "8d798a29-f083-11e1-bdfb-cb223EXAMPLE")
}
