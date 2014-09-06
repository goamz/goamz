package cloudformation_test

import (
	"testing"
	"time"

	"github.com/motain/gocheck"

	"github.com/goamz/goamz/aws"
	cf "github.com/goamz/goamz/cloudformation"
	"github.com/goamz/goamz/testutil"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

var _ = gocheck.Suite(&S{})

type S struct {
	cf *cf.CloudFormation
}

var testServer = testutil.NewHTTPServer()

var mockTest bool

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.cf = cf.New(auth, aws.Region{CloudFormationEndpoint: testServer.URL})
}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

func (s *S) TestCancelUpdateStack(c *gocheck.C) {
	testServer.Response(200, nil, CancelUpdateStackResponse)

	resp, err := s.cf.CancelUpdateStack("foo")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "CancelUpdateStack")
	c.Assert(values.Get("StackName"), gocheck.Equals, "foo")

	// Response test
	c.Assert(resp.RequestId, gocheck.Equals, "4af14eec-350e-11e4-b260-EXAMPLE")
}

func (s *S) TestCreateStack(c *gocheck.C) {
	testServer.Response(200, nil, CreateStackResponse)

	stackParams := &cf.CreateStackParams{
		NotificationARNs: []string{"arn:aws:sns:us-east-1:1234567890:my-topic"},
		Parameters: []cf.Parameter{
			{
				ParameterKey:   "AvailabilityZone",
				ParameterValue: "us-east-1a",
			},
		},
		StackName:    "MyStack",
		TemplateBody: "[Template Document]",
	}
	resp, err := s.cf.CreateStack(stackParams)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "CreateStack")
	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	c.Assert(values.Get("NotificationARNs.member.1"), gocheck.Equals, "arn:aws:sns:us-east-1:1234567890:my-topic")
	c.Assert(values.Get("TemplateBody"), gocheck.Equals, "[Template Document]")
	c.Assert(values.Get("Parameters.member.1.ParameterKey"), gocheck.Equals, "AvailabilityZone")
	c.Assert(values.Get("Parameters.member.1.ParameterValue"), gocheck.Equals, "us-east-1a")
	// Response test
	c.Assert(resp.RequestId, gocheck.Equals, "4af14eec-350e-11e4-b260-EXAMPLE")
	c.Assert(resp.StackId, gocheck.Equals, "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83")
}

func (s *S) TestCreateStackWithInvalidParams(c *gocheck.C) {
	testServer.Response(400, nil, CreateStackWithInvalidParamsResponse)
	//testServer.Response(200, nil, DeleteAutoScalingGroupResponse)

	cfTemplate := `
{
  "AWSTemplateFormatVersion" : "2010-09-09",
  "Description" : "Sample template",
  "Parameters" : {
    "KeyName" : {
      "Description" : "key pair",
      "Type" : "String"
    }
  },
  "Resources" : {
    "Ec2Instance" : {
      "Type" : "AWS::EC2::Instance",
      "Properties" : {
        "KeyName" : { "Ref" : "KeyName" },
        "ImageId" : "ami-7f418316",
        "UserData" : { "Fn::Base64" : "80" }
      }
    }
  },
  "Outputs" : {
    "InstanceId" : {
      "Description" : "InstanceId of the newly created EC2 instance",
      "Value" : { "Ref" : "Ec2Instance" }
    }
}`

	stackParams := &cf.CreateStackParams{
		Capabilities:    []string{"CAPABILITY_IAM"},
		DisableRollback: true,
		NotificationARNs: []string{
			"arn:aws:sns:us-east-1:1234567890:my-topic",
			"arn:aws:sns:us-east-1:1234567890:my-topic2",
		},
		OnFailure: "ROLLBACK",
		Parameters: []cf.Parameter{
			{
				ParameterKey:   "AvailabilityZone",
				ParameterValue: "us-east-1a",
			},
		},
		StackName:       "MyStack",
		StackPolicyBody: "{PolicyBody}",
		StackPolicyURL:  "http://stack-policy-url",
		Tags: []cf.Tag{
			{
				Key:   "TagKey",
				Value: "TagValue",
			},
		},
		TemplateBody:     cfTemplate,
		TemplateURL:      "http://url",
		TimeoutInMinutes: 20,
	}
	resp, err := s.cf.CreateStack(stackParams)
	c.Assert(err, gocheck.NotNil)
	c.Assert(resp, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm

	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "CreateStack")
	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	c.Assert(values.Get("NotificationARNs.member.1"), gocheck.Equals, "arn:aws:sns:us-east-1:1234567890:my-topic")
	c.Assert(values.Get("NotificationARNs.member.2"), gocheck.Equals, "arn:aws:sns:us-east-1:1234567890:my-topic2")
	c.Assert(values.Get("Capabilities.member.1"), gocheck.Equals, "CAPABILITY_IAM")
	c.Assert(values.Get("TemplateBody"), gocheck.Equals, cfTemplate)
	c.Assert(values.Get("TemplateURL"), gocheck.Equals, "http://url")
	c.Assert(values.Get("StackPolicyBody"), gocheck.Equals, "{PolicyBody}")
	c.Assert(values.Get("StackPolicyURL"), gocheck.Equals, "http://stack-policy-url")
	c.Assert(values.Get("OnFailure"), gocheck.Equals, "ROLLBACK")
	c.Assert(values.Get("DisableRollback"), gocheck.Equals, "true")
	c.Assert(values.Get("Tags.member.1.Key"), gocheck.Equals, "TagKey")
	c.Assert(values.Get("Tags.member.1.Value"), gocheck.Equals, "TagValue")
	c.Assert(values.Get("Parameters.member.1.ParameterKey"), gocheck.Equals, "AvailabilityZone")
	c.Assert(values.Get("Parameters.member.1.ParameterValue"), gocheck.Equals, "us-east-1a")
	c.Assert(values.Get("TimeoutInMinutes"), gocheck.Equals, "20")

	// Response test
	c.Assert(err.(*cf.Error).RequestId, gocheck.Equals, "70a76d42-9665-11e2-9fdf-211deEXAMPLE")
	c.Assert(err.(*cf.Error).Message, gocheck.Equals, "Either Template URL or Template Body must be specified.")
	c.Assert(err.(*cf.Error).Type, gocheck.Equals, "Sender")
	c.Assert(err.(*cf.Error).Code, gocheck.Equals, "ValidationError")
	c.Assert(err.(*cf.Error).StatusCode, gocheck.Equals, 400)

}

func (s *S) TestDeleteStack(c *gocheck.C) {
	testServer.Response(200, nil, DeleteStackResponse)

	resp, err := s.cf.DeleteStack("foo")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "DeleteStack")
	c.Assert(values.Get("StackName"), gocheck.Equals, "foo")
	// Response test
	c.Assert(resp.RequestId, gocheck.Equals, "4af14eec-350e-11e4-b260-EXAMPLE")
}

func (s *S) TestDescribeStackEvents(c *gocheck.C) {
	testServer.Response(200, nil, DescribeStackEventsResponse)

	resp, err := s.cf.DescribeStackEvents("MyStack", "")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm

	// Post request test
	t1, _ := time.Parse(time.RFC3339, "2010-07-27T22:26:28Z")
	t2, _ := time.Parse(time.RFC3339, "2010-07-27T22:27:28Z")
	t3, _ := time.Parse(time.RFC3339, "2010-07-27T22:28:28Z")
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeStackEvents")
	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	c.Assert(values.Get("NextToken"), gocheck.Equals, "")

	// Response test
	expected := &cf.DescribeStackEventsResponse{
		StackEvents: []cf.StackEvent{
			{
				EventId:              "Event-1-Id",
				StackId:              "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83",
				StackName:            "MyStack",
				LogicalResourceId:    "MyStack",
				PhysicalResourceId:   "MyStack_One",
				ResourceType:         "AWS::CloudFormation::Stack",
				Timestamp:            t1,
				ResourceStatus:       "CREATE_IN_PROGRESS",
				ResourceStatusReason: "User initiated",
			},
			{
				EventId:            "Event-2-Id",
				StackId:            "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83",
				StackName:          "MyStack",
				LogicalResourceId:  "MyDBInstance",
				PhysicalResourceId: "MyStack_DB1",
				ResourceType:       "AWS::SecurityGroup",
				Timestamp:          t2,
				ResourceStatus:     "CREATE_IN_PROGRESS",
				ResourceProperties: "{\"GroupDescription\":...}",
			},
			{
				EventId:            "Event-3-Id",
				StackId:            "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83",
				StackName:          "MyStack",
				LogicalResourceId:  "MySG1",
				PhysicalResourceId: "MyStack_SG1",
				ResourceType:       "AWS::SecurityGroup",
				Timestamp:          t3,
				ResourceStatus:     "CREATE_COMPLETE",
			},
		},
		NextToken: "",
		RequestId: "4af14eec-350e-11e4-b260-EXAMPLE",
	}
	c.Assert(resp, gocheck.DeepEquals, expected)
}

func (s *S) TestDescribeStackResource(c *gocheck.C) {
	testServer.Response(200, nil, DescribeStackResourceResponse)

	resp, err := s.cf.DescribeStackResource("MyStack", "MyDBInstance")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeStackResource")
	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	c.Assert(values.Get("LogicalResourceId"), gocheck.Equals, "MyDBInstance")
	t, _ := time.Parse(time.RFC3339, "2011-07-07T22:27:28Z")
	// Response test
	expected := &cf.DescribeStackResourceResponse{
		StackResourceDetail: cf.StackResourceDetail{
			StackId:              "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83",
			StackName:            "MyStack",
			LogicalResourceId:    "MyDBInstance",
			PhysicalResourceId:   "MyStack_DB1",
			ResourceType:         "AWS::RDS::DBInstance",
			LastUpdatedTimestamp: t,
			ResourceStatus:       "CREATE_COMPLETE",
		},
		RequestId: "4af14eec-350e-11e4-b260-EXAMPLE",
	}
	c.Assert(resp, gocheck.DeepEquals, expected)
}

func (s *S) TestDescribeStackResources(c *gocheck.C) {
	testServer.Response(200, nil, DescribeStackResourcesResponse)

	resp, err := s.cf.DescribeStackResources("MyStack", "", "")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm

	// Post request test
	t1, _ := time.Parse(time.RFC3339, "2010-07-27T22:27:28Z")
	t2, _ := time.Parse(time.RFC3339, "2010-07-27T22:28:28Z")
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeStackResources")
	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	c.Assert(values.Get("PhysicalResourceId"), gocheck.Equals, "")
	c.Assert(values.Get("LogicalResourceId"), gocheck.Equals, "")

	// Response test
	expected := &cf.DescribeStackResourcesResponse{
		StackResources: []cf.StackResource{
			{
				StackId:            "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83",
				StackName:          "MyStack",
				LogicalResourceId:  "MyDBInstance",
				PhysicalResourceId: "MyStack_DB1",
				ResourceType:       "AWS::DBInstance",
				Timestamp:          t1,
				ResourceStatus:     "CREATE_COMPLETE",
			},
			{
				StackId:            "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83",
				StackName:          "MyStack",
				LogicalResourceId:  "MyAutoScalingGroup",
				PhysicalResourceId: "MyStack_ASG1",
				ResourceType:       "AWS::AutoScalingGroup",
				Timestamp:          t2,
				ResourceStatus:     "CREATE_IN_PROGRESS",
			},
		},
		RequestId: "4af14eec-350e-11e4-b260-EXAMPLE",
	}
	c.Assert(resp, gocheck.DeepEquals, expected)
}

func (s *S) TestDescribeStacks(c *gocheck.C) {
	testServer.Response(200, nil, DescribeStacksResponse)

	resp, err := s.cf.DescribeStacks("MyStack", "")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm

	// Post request test
	t, _ := time.Parse(time.RFC3339, "2010-07-27T22:28:28Z")
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "DescribeStacks")
	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	c.Assert(values.Get("NextToken"), gocheck.Equals, "")

	// Response test
	expected := &cf.DescribeStacksResponse{
		Stacks: []cf.Stack{
			{
				StackId:          "arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83",
				StackName:        "MyStack",
				Description:      "My Description",
				Capabilities:     []string{"CAPABILITY_IAM"},
				NotificationARNs: []string{"arn:aws:sns:region-name:account-name:topic-name"},
				Parameters: []cf.Parameter{
					{
						ParameterKey:   "MyKey",
						ParameterValue: "MyValue",
					},
				},
				Tags: []cf.Tag{
					{
						Key:   "MyTagKey",
						Value: "MyTagValue",
					},
				},
				CreationTime:    t,
				StackStatus:     "CREATE_COMPLETE",
				DisableRollback: false,
				Outputs: []cf.Output{
					{
						Description: "ServerUrl",
						OutputKey:   "StartPage",
						OutputValue: "http://my-load-balancer.amazonaws.com:80/index.html",
					},
				},
			},
		},
		NextToken: "",
		RequestId: "4af14eec-350e-11e4-b260-EXAMPLE",
	}
	c.Assert(resp, gocheck.DeepEquals, expected)
}

func (s *S) TestEstimateTemplateCost(c *gocheck.C) {
	testServer.Response(200, nil, EstimateTemplateCostResponse)

	resp, err := s.cf.EstimateTemplateCost(nil, "", "https://s3.amazonaws.com/cloudformation-samples-us-east-1/Drupal_Simple.template")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "EstimateTemplateCost")
	c.Assert(values.Get("TemplateBody"), gocheck.Equals, "")
	c.Assert(values.Get("TemplateURL"), gocheck.Equals, "https://s3.amazonaws.com/cloudformation-samples-us-east-1/Drupal_Simple.template")
	// Response test
	c.Assert(resp.Url, gocheck.Equals, "http://calculator.s3.amazonaws.com/calc5.html?key=cf-2e351785-e821-450c-9d58-625e1e1ebfb6")
	c.Assert(resp.RequestId, gocheck.Equals, "4af14eec-350e-11e4-b260-EXAMPLE")
}

func (s *S) TestGetStackPolicy(c *gocheck.C) {
	testServer.Response(200, nil, GetStackPolicyResponse)

	resp, err := s.cf.GetStackPolicy("MyStack")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "GetStackPolicy")

	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	// Response test
	policy := `{
      "Statement" : [
        {
          "Effect" : "Deny",
          "Action" : "Update:*",
          "Principal" : "*",
          "Resource" : "LogicalResourceId/ProductionDatabase"
        },
        {
          "Effect" : "Allow",
          "Action" : "Update:*",
          "Principal" : "*",
          "Resource" : "*"
        }
      ]
    }`
	c.Assert(resp.StackPolicyBody, gocheck.Equals, policy)
	c.Assert(resp.RequestId, gocheck.Equals, "4af14eec-350e-11e4-b260-EXAMPLE")
}

func (s *S) TestGetTemplate(c *gocheck.C) {
	testServer.Response(200, nil, GetTemplateResponse)

	resp, err := s.cf.GetTemplate("MyStack")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "GetTemplate")

	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	// Response test
	templateBody := `{
      "AWSTemplateFormatVersion" : "2010-09-09",
      "Description" : "Simple example",
      "Resources" : {
        "MySQS" : {
           "Type" : "AWS::SQS::Queue",
           "Properties" : {
            }
         }
        }
      }`
	c.Assert(resp.TemplateBody, gocheck.Equals, templateBody)
	c.Assert(resp.RequestId, gocheck.Equals, "4af14eec-350e-11e4-b260-EXAMPLE")
}

func (s *S) TestListStackResources(c *gocheck.C) {
	testServer.Response(200, nil, ListStackResourcesResponse)

	resp, err := s.cf.ListStackResources("MyStack", "4dad1-32131da-d-31")
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm

	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2010-05-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "ListStackResources")
	c.Assert(values.Get("StackName"), gocheck.Equals, "MyStack")
	c.Assert(values.Get("NextToken"), gocheck.Equals, "4dad1-32131da-d-31")

	// Response test
	t1, _ := time.Parse(time.RFC3339, "2011-06-21T20:15:58Z")
	t2, _ := time.Parse(time.RFC3339, "2011-06-21T20:25:57Z")
	t3, _ := time.Parse(time.RFC3339, "2011-06-21T20:26:12Z")
	t4, _ := time.Parse(time.RFC3339, "2011-06-21T20:28:48Z")
	t5, _ := time.Parse(time.RFC3339, "2011-06-21T20:29:06Z")
	t6, _ := time.Parse(time.RFC3339, "2011-06-21T20:29:23Z")

	expected := &cf.ListStackResourcesResponse{
		StackResourceSummaries: []cf.StackResourceSummary{
			{
				LogicalResourceId:    "DBSecurityGroup",
				PhysicalResourceId:   "gmarcteststack-dbsecuritygroup-1s5m0ez5lkk6w",
				ResourceType:         "AWS::RDS::DBSecurityGroup",
				LastUpdatedTimestamp: t1,
				ResourceStatus:       "CREATE_COMPLETE",
			},
			{
				LogicalResourceId:    "SampleDB",
				PhysicalResourceId:   "MyStack-sampledb-ycwhk1v830lx",
				ResourceType:         "AWS::RDS::DBInstance",
				LastUpdatedTimestamp: t2,
				ResourceStatus:       "CREATE_COMPLETE",
			},
			{
				LogicalResourceId:    "SampleApplication",
				PhysicalResourceId:   "MyStack-SampleApplication-1MKNASYR3RBQL",
				ResourceType:         "AWS::ElasticBeanstalk::Application",
				LastUpdatedTimestamp: t3,
				ResourceStatus:       "CREATE_COMPLETE",
			},
			{
				LogicalResourceId:    "SampleEnvironment",
				PhysicalResourceId:   "myst-Samp-1AGU6ERZX6M3Q",
				ResourceType:         "AWS::ElasticBeanstalk::Environment",
				LastUpdatedTimestamp: t4,
				ResourceStatus:       "CREATE_COMPLETE",
			},
			{
				LogicalResourceId:    "AlarmTopic",
				PhysicalResourceId:   "arn:aws:sns:us-east-1:803981987763:MyStack-AlarmTopic-SW4IQELG7RPJ",
				ResourceType:         "AWS::SNS::Topic",
				LastUpdatedTimestamp: t5,
				ResourceStatus:       "CREATE_COMPLETE",
			},
			{
				LogicalResourceId:    "CPUAlarmHigh",
				PhysicalResourceId:   "MyStack-CPUAlarmHigh-POBWQPDJA81F",
				ResourceType:         "AWS::CloudWatch::Alarm",
				LastUpdatedTimestamp: t6,
				ResourceStatus:       "CREATE_COMPLETE",
			},
		},
		NextToken: "",
		RequestId: "2d06e36c-ac1d-11e0-a958-f9382b6eb86b",
	}
	c.Assert(resp, gocheck.DeepEquals, expected)
}
