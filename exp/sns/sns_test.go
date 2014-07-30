package sns_test

import (
	"github.com/czgo/goamz/aws"
	"github.com/czgo/goamz/exp/sns"
	"github.com/czgo/goamz/testutil"
	"github.com/motain/gocheck"
	"testing"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

var _ = gocheck.Suite(&S{})

type S struct {
	sns *sns.SNS
}

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.sns = sns.New(auth, aws.Region{SNSEndpoint: testServer.URL})
}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

func (s *S) TestListTopicsOK(c *gocheck.C) {
	testServer.Response(200, nil, TestListTopicsXmlOK)

	resp, err := s.sns.ListTopics(nil)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "bd10b26c-e30e-11e0-ba29-93c3aca2f103")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestCreateTopic(c *gocheck.C) {
	testServer.Response(200, nil, TestCreateTopicXmlOK)

	resp, err := s.sns.CreateTopic("My-Topic")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.Topic.TopicArn, gocheck.Equals, "arn:aws:sns:us-east-1:123456789012:My-Topic")
	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "a8dec8b3-33a4-11df-8963-01868b7c937a")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestDeleteTopic(c *gocheck.C) {
	testServer.Response(200, nil, TestDeleteTopicXmlOK)

	t := sns.Topic{nil, "arn:aws:sns:us-east-1:123456789012:My-Topic"}
	resp, err := s.sns.DeleteTopic(t)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "f3aa9ac9-3c3d-11df-8235-9dab105e9c32")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestListSubscriptions(c *gocheck.C) {
	testServer.Response(200, nil, TestListSubscriptionsXmlOK)

	resp, err := s.sns.ListSubscriptions(nil)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(len(resp.Subscriptions), gocheck.Not(gocheck.Equals), 0)
	c.Assert(resp.Subscriptions[0].Protocol, gocheck.Equals, "email")
	c.Assert(resp.Subscriptions[0].Endpoint, gocheck.Equals, "example@amazon.com")
	c.Assert(resp.Subscriptions[0].SubscriptionArn, gocheck.Equals, "arn:aws:sns:us-east-1:123456789012:My-Topic:80289ba6-0fd4-4079-afb4-ce8c8260f0ca")
	c.Assert(resp.Subscriptions[0].TopicArn, gocheck.Equals, "arn:aws:sns:us-east-1:698519295917:My-Topic")
	c.Assert(resp.Subscriptions[0].Owner, gocheck.Equals, "123456789012")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestGetTopicAttributes(c *gocheck.C) {
	testServer.Response(200, nil, TestGetTopicAttributesXmlOK)

	resp, err := s.sns.GetTopicAttributes("arn:aws:sns:us-east-1:123456789012:My-Topic")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(len(resp.Attributes), gocheck.Not(gocheck.Equals), 0)
	c.Assert(resp.Attributes[0].Key, gocheck.Equals, "Owner")
	c.Assert(resp.Attributes[0].Value, gocheck.Equals, "123456789012")
	c.Assert(resp.Attributes[1].Key, gocheck.Equals, "Policy")
	c.Assert(resp.Attributes[1].Value, gocheck.Equals, `{"Version":"2008-10-17","Id":"us-east-1/698519295917/test__default_policy_ID","Statement" : [{"Effect":"Allow","Sid":"us-east-1/698519295917/test__default_statement_ID","Principal" : {"AWS": "*"},"Action":["SNS:GetTopicAttributes","SNS:SetTopicAttributes","SNS:AddPermission","SNS:RemovePermission","SNS:DeleteTopic","SNS:Subscribe","SNS:ListSubscriptionsByTopic","SNS:Publish","SNS:Receive"],"Resource":"arn:aws:sns:us-east-1:698519295917:test","Condition" : {"StringLike" : {"AWS:SourceArn": "arn:aws:*:*:698519295917:*"}}}]}`)
	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "057f074c-33a7-11df-9540-99d0768312d3")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestPublish(c *gocheck.C) {
	testServer.Response(200, nil, TestPublishXmlOK)

	pubOpt := &sns.PublishOpt{"foobar", "", "subject", "arn:aws:sns:us-east-1:123456789012:My-Topic"}
	resp, err := s.sns.Publish(pubOpt)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.MessageId, gocheck.Equals, "94f20ce6-13c5-43a0-9a9e-ca52d816e90b")
	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "f187a3c1-376f-11df-8963-01868b7c937a")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestSetTopicAttributes(c *gocheck.C) {
	testServer.Response(200, nil, TestSetTopicAttributesXmlOK)

	resp, err := s.sns.SetTopicAttributes("DisplayName", "MyTopicName", "arn:aws:sns:us-east-1:123456789012:My-Topic")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "a8763b99-33a7-11df-a9b7-05d48da6f042")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestSubscribe(c *gocheck.C) {
	testServer.Response(200, nil, TestSubscribeXmlOK)

	resp, err := s.sns.Subscribe("example@amazon.com", "email", "arn:aws:sns:us-east-1:123456789012:My-Topic")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.SubscriptionArn, gocheck.Equals, "pending confirmation")
	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "a169c740-3766-11df-8963-01868b7c937a")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestUnsubscribe(c *gocheck.C) {
	testServer.Response(200, nil, TestUnsubscribeXmlOK)

	resp, err := s.sns.Unsubscribe("arn:aws:sns:us-east-1:123456789012:My-Topic:a169c740-3766-11df-8963-01868b7c937a")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "18e0ac39-3776-11df-84c0-b93cc1666b84")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestConfirmSubscription(c *gocheck.C) {
	testServer.Response(200, nil, TestConfirmSubscriptionXmlOK)

	opt := &sns.ConfirmSubscriptionOpt{"", "51b2ff3edb475b7d91550e0ab6edf0c1de2a34e6ebaf6", "arn:aws:sns:us-east-1:123456789012:My-Topic"}
	resp, err := s.sns.ConfirmSubscription(opt)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.SubscriptionArn, gocheck.Equals, "arn:aws:sns:us-east-1:123456789012:My-Topic:80289ba6-0fd4-4079-afb4-ce8c8260f0ca")
	c.Assert(resp.ResponseMetadata.RequestId, gocheck.Equals, "7a50221f-3774-11df-a9b7-05d48da6f042")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestAddPermission(c *gocheck.C) {
	testServer.Response(200, nil, TestAddPermissionXmlOK)
	perm := make([]sns.Permission, 2)
	perm[0].ActionName = "Publish"
	perm[1].ActionName = "GetTopicAttributes"
	perm[0].AccountId = "987654321000"
	perm[1].AccountId = "876543210000"

	resp, err := s.sns.AddPermission(perm, "NewPermission", "arn:aws:sns:us-east-1:123456789012:My-Topic")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.RequestId, gocheck.Equals, "6a213e4e-33a8-11df-9540-99d0768312d3")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestRemovePermission(c *gocheck.C) {
	testServer.Response(200, nil, TestRemovePermissionXmlOK)

	resp, err := s.sns.RemovePermission("NewPermission", "arn:aws:sns:us-east-1:123456789012:My-Topic")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.RequestId, gocheck.Equals, "d170b150-33a8-11df-995a-2d6fbe836cc1")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestListSubscriptionByTopic(c *gocheck.C) {
	testServer.Response(200, nil, TestListSubscriptionsByTopicXmlOK)

	opt := &sns.ListSubscriptionByTopicOpt{"", "arn:aws:sns:us-east-1:123456789012:My-Topic"}
	resp, err := s.sns.ListSubscriptionByTopic(opt)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(len(resp.Subscriptions), gocheck.Not(gocheck.Equals), 0)
	c.Assert(resp.Subscriptions[0].TopicArn, gocheck.Equals, "arn:aws:sns:us-east-1:123456789012:My-Topic")
	c.Assert(resp.Subscriptions[0].SubscriptionArn, gocheck.Equals, "arn:aws:sns:us-east-1:123456789012:My-Topic:80289ba6-0fd4-4079-afb4-ce8c8260f0ca")
	c.Assert(resp.Subscriptions[0].Owner, gocheck.Equals, "123456789012")
	c.Assert(resp.Subscriptions[0].Endpoint, gocheck.Equals, "example@amazon.com")
	c.Assert(resp.Subscriptions[0].Protocol, gocheck.Equals, "email")
	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestCreatePlatformApplication(c *gocheck.C) {
	testServer.Response(200, nil, TestCreatePlatformApplicationXmlOK)

	attrs := []sns.AttributeEntry{
		sns.AttributeEntry{Key: "PlatformCredential", Value: "AIzaSyClE2lcV2zEKTLYYo645zfk2jhQPFeyxDo"},
		sns.AttributeEntry{Key: "PlatformPrincipal", Value: "There is no principal for GCM"},
	}
	opt := &sns.PlatformApplicationOpt{Name: "gcmpushapp", Platform: "GCM", Attributes: attrs}
	resp, err := s.sns.CreatePlatformApplication(opt)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.PlatformApplicationArn, gocheck.Equals, "arn:aws:sns:us-west-2:123456789012:app/GCM/gcmpushapp")
	c.Assert(resp.RequestId, gocheck.Equals, "b6f0e78b-e9d4-5a0e-b973-adc04e8a4ff9")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestCreatePlatformEndpoint(c *gocheck.C) {
	testServer.Response(200, nil, TestCreatePlatformEndpointXmlOK)

	opt := &sns.PlatformEndpointOpt{PlatformApplicationArn: "arn:aws:sns:us-west-2:123456789012:app/GCM/gcmpushapp",
		CustomUserData: "UserId=27576823", Token: "APA91bGi7fFachkC1xjlqT66VYEucGHochmf1VQAr9k...jsM0PKPxKhddCzx6paEsyay9Zn3D4wNUJb8m6HZrBEXAMPLE"}

	resp, err := s.sns.CreatePlatformEndpoint(opt)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.EndpointArn, gocheck.Equals, "arn:aws:sns:us-west-2:123456789012:endpoint/GCM/gcmpushapp/5e3e9847-3183-3f18-a7e8-671c3a57d4b3")
	c.Assert(resp.RequestId, gocheck.Equals, "6613341d-3e15-53f7-bf3c-7e56994ba278")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestDeleteEndpoint(c *gocheck.C) {
	testServer.Response(200, nil, DeleteEndpointXmlOK)

	resp, err := s.sns.DeleteEndpoint("arn:aws:sns:us-west-2:123456789012:endpoint/GCM%/gcmpushapp/5e3e9847-3183-3f18-a7e8-671c3a57d4b3")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.RequestId, gocheck.Equals, "c1d2b191-353c-5a5f-8969-fbdd3900afa8")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestDeletePlatformApplication(c *gocheck.C) {
	testServer.Response(200, nil, TestDeletePlatformApplicationXmlOK)

	resp, err := s.sns.DeletePlatformApplication("arn:aws:sns:us-west-2:123456789012:app/GCM/gcmpushapp")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.RequestId, gocheck.Equals, "097dac18-7a77-5823-a8dd-e65476dcb037")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestGetEndpointAttributes(c *gocheck.C) {
	testServer.Response(200, nil, TestGetEndpointAttributesXmlOK)

	resp, err := s.sns.GetEndpointAttributes("arn:aws:sns:us-west-2:123456789012:endpoint/GCM/gcmpushapp/5e3e9847-3183-3f18-a7e8-671c3a57d4b3")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(len(resp.Attributes), gocheck.Equals, 3)
	c.Assert(resp.Attributes[0].Key, gocheck.Equals, "Enabled")
	c.Assert(resp.Attributes[0].Value, gocheck.Equals, "true")

	c.Assert(resp.Attributes[1].Key, gocheck.Equals, "CustomUserData")
	c.Assert(resp.Attributes[1].Value, gocheck.Equals, "UserId=01234567")

	c.Assert(resp.Attributes[2].Key, gocheck.Equals, "Token")
	c.Assert(resp.Attributes[2].Value, gocheck.Equals, "APA91bGi7fFachkC1xjlqT66VYEucGHochmf1VQAr9k...jsM0PKPxKhddCzx6paEsyay9Zn3D4wNUJb8m6HZrBEXAMPLE")

	c.Assert(resp.RequestId, gocheck.Equals, "6c725a19-a142-5b77-94f9-1055a9ea04e7")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestGetPlatformApplicationAttributes(c *gocheck.C) {
	testServer.Response(200, nil, TestGetPlatformApplicationAttributesXmlOK)

	resp, err := s.sns.GetPlatformApplicationAttributes("arn:aws:sns:us-west-2:123456789012:app/GCM/gcmpushapp", "")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(len(resp.Attributes), gocheck.Not(gocheck.Equals), 0)
	c.Assert(resp.Attributes[0].Key, gocheck.Equals, "AllowEndpointPolicies")
	c.Assert(resp.Attributes[0].Value, gocheck.Equals, "false")
	c.Assert(resp.RequestId, gocheck.Equals, "74848df2-87f6-55ed-890c-c7be80442462")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestListEndpointsByPlatformApplication(c *gocheck.C) {
	testServer.Response(200, nil, TestListEndpointsByPlatformApplicationXmlOK)

	resp, err := s.sns.ListEndpointsByPlatformApplication("arn:aws:sns:us-west-2:123456789012:app/GCM/gcmpushapp", "")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(len(resp.Endpoints), gocheck.Not(gocheck.Equals), 0)
	c.Assert(resp.Endpoints[0].EndpointArn, gocheck.Equals, "arn:aws:sns:us-west-2:123456789012:endpoint/GCM/gcmpushapp/5e3e9847-3183-3f18-a7e8-671c3a57d4b3")
	c.Assert(len(resp.Endpoints[0].Attributes), gocheck.Equals, 3)
	c.Assert(resp.Endpoints[0].Attributes[0].Key, gocheck.Equals, "Enabled")
	c.Assert(resp.Endpoints[0].Attributes[0].Value, gocheck.Equals, "true")
	c.Assert(resp.Endpoints[0].Attributes[1].Key, gocheck.Equals, "CustomUserData")
	c.Assert(resp.Endpoints[0].Attributes[1].Value, gocheck.Equals, "UserId=27576823")
	c.Assert(resp.Endpoints[0].Attributes[2].Key, gocheck.Equals, "Token")
	c.Assert(resp.Endpoints[0].Attributes[2].Value, gocheck.Equals, "APA91bGi7fFachkC1xjlqT66VYEucGHochmf1VQAr9k...jsM0PKPxKhddCzx6paEsyay9Zn3D4wNUJb8m6HZrBEXAMPLE")

	c.Assert(resp.RequestId, gocheck.Equals, "9a48768c-dac8-5a60-aec0-3cc27ea08d96")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestListPlatformApplications(c *gocheck.C) {
	testServer.Response(200, nil, TestListPlatformApplicationsXmlOK)

	resp, err := s.sns.ListPlatformApplications("")
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(len(resp.PlatformApplications), gocheck.Not(gocheck.Equals), 0)

	c.Assert(resp.PlatformApplications[0].PlatformApplicationArn, gocheck.Equals, "arn:aws:sns:us-west-2:123456789012:app/APNS_SANDBOX/apnspushapp")
	c.Assert(len(resp.PlatformApplications[0].Attributes), gocheck.Equals, 1)
	c.Assert(resp.PlatformApplications[0].Attributes[0].Key, gocheck.Equals, "AllowEndpointPolicies")
	c.Assert(resp.PlatformApplications[0].Attributes[0].Value, gocheck.Equals, "false")

	c.Assert(resp.PlatformApplications[1].PlatformApplicationArn, gocheck.Equals, "arn:aws:sns:us-west-2:123456789012:app/GCM/gcmpushapp")
	c.Assert(len(resp.PlatformApplications[1].Attributes), gocheck.Equals, 1)
	c.Assert(resp.PlatformApplications[1].Attributes[0].Key, gocheck.Equals, "AllowEndpointPolicies")
	c.Assert(resp.PlatformApplications[1].Attributes[0].Value, gocheck.Equals, "false")

	c.Assert(resp.RequestId, gocheck.Equals, "315a335e-85d8-52df-9349-791283cbb529")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestSetEndpointAttributes(c *gocheck.C) {
	testServer.Response(200, nil, TestSetEndpointAttributesXmlOK)

	attrs := []sns.AttributeEntry{
		sns.AttributeEntry{Key: "CustomUserData", Value: "My custom userdata"},
	}

	opts := &sns.SetEndpointAttributesOpt{
		EndpointArn: "arn:aws:sns:us-west-2:123456789012:endpoint/GCM/gcmpushapp/5e3e9847-3183-3f18-a7e8-671c3a57d4b3",
		Attributes:  attrs}

	resp, err := s.sns.SetEndpointAttributes(opts)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.RequestId, gocheck.Equals, "2fe0bfc7-3e85-5ee5-a9e2-f58b35e85f6a")

	c.Assert(err, gocheck.IsNil)
}

func (s *S) TestSetPlatformApplicationAttributes(c *gocheck.C) {
	testServer.Response(200, nil, TestSetPlatformApplicationAttributesXmlOK)

	attrs := []sns.AttributeEntry{
		sns.AttributeEntry{Key: "EventEndpointCreated", Value: "arn:aws:sns:us-west-2:123456789012:topicarn"},
	}

	opts := &sns.SetPlatformApplicationAttributesOpt{
		PlatformApplicationArn: "arn:aws:sns:us-west-2:123456789012:app/GCM/gcmpushapp",
		Attributes:             attrs}

	resp, err := s.sns.SetPlatformApplicationAttributes(opts)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "GET")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	c.Assert(resp.RequestId, gocheck.Equals, "cf577bcc-b3dc-5463-88f1-3180b9412395")

	c.Assert(err, gocheck.IsNil)
}
