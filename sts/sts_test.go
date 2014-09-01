package sts_test

import (
	"testing"
	"time"

	"github.com/motain/gocheck"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/sts"
	"github.com/goamz/goamz/testutil"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

var _ = gocheck.Suite(&S{})

type S struct {
	sts *sts.STS
}

var testServer = testutil.NewHTTPServer()

var mockTest bool

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.sts = sts.New(auth, aws.Region{STSEndpoint: testServer.URL})
}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

func (s *S) TestAssumeRole(c *gocheck.C) {
	testServer.Response(200, nil, AssumeRoleResponse)
	request := &sts.AssumeRoleParams{
		DurationSeconds: 3600,
		ExternalId:      "123ABC",
		Policy:          `{"Version":"2012-10-17","Statement":[{"Sid":"Stmt1","Effect":"Allow","Action":"s3:*","Resource":"*"}]}`,
		RoleArn:         "arn:aws:iam::123456789012:role/demo",
		RoleSessionName: "Bob",
	}
	resp, err := s.sts.AssumeRole(request)
	c.Assert(err, gocheck.IsNil)
	values := testServer.WaitRequest().PostForm
	// Post request test
	c.Assert(values.Get("Version"), gocheck.Equals, "2011-06-15")
	c.Assert(values.Get("Action"), gocheck.Equals, "AssumeRole")
	c.Assert(values.Get("DurationSeconds"), gocheck.Equals, "3600")
	c.Assert(values.Get("ExternalId"), gocheck.Equals, "123ABC")
	c.Assert(values.Get("Policy"), gocheck.Equals, `{"Version":"2012-10-17","Statement":[{"Sid":"Stmt1","Effect":"Allow","Action":"s3:*","Resource":"*"}]}`)
	c.Assert(values.Get("RoleArn"), gocheck.Equals, "arn:aws:iam::123456789012:role/demo")
	c.Assert(values.Get("RoleSessionName"), gocheck.Equals, "Bob")
	// Response test
	exp, _ := time.Parse(time.RFC3339, "2011-07-15T23:28:33.359Z")
	c.Assert(resp.RequestId, gocheck.Equals, "c6104cbe-af31-11e0-8154-cbc7ccf896c7")
	c.Assert(resp.PackedPolicySize, gocheck.Equals, 6)
	c.Assert(resp.AssumedRoleUser, gocheck.DeepEquals, sts.AssumedRoleUser{
		Arn:           "arn:aws:sts::123456789012:assumed-role/demo/Bob",
		AssumedRoleId: "ARO123EXAMPLE123:Bob",
	})
	c.Assert(resp.Credentials, gocheck.DeepEquals, sts.Credentials{
		SessionToken: `
       AQoDYXdzEPT//////////wEXAMPLEtc764bNrC9SAPBSM22wDOk4x4HIZ8j4FZTwdQW
       LWsKWHGBuFqwAeMicRXmxfpSPfIeoIYRqTflfKD8YUuwthAx7mSEI/qkPpKPi/kMcGd
       QrmGdeehM4IC1NtBmUpp2wUE8phUZampKsburEDy0KPkyQDYwT7WZ0wq5VSXDvp75YU
       9HFvlRd8Tx6q6fE8YQcHNVXAkiY9q6d+xo0rKwT38xVqr7ZD0u0iPPkUL64lIZbqBAz
       +scqKmlzm8FDrypNC9Yjc8fPOLn9FX9KSYvKTr4rvx3iSIlTJabIQwj2ICCR/oLxBA==
      `,
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYzEXAMPLEKEY",
		AccessKeyId:     "AKIAIOSFODNN7EXAMPLE",
		Expiration:      exp,
	})

}
