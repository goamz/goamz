package cloudformation_test

import (
	"testing"
	//"time"

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
	// Response test

	c.Assert(resp.CancelUpdateStackResult, gocheck.Equals, "")
}




type CreateStackParams
