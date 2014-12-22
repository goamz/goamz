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
