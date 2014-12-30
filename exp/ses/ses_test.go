// ses_test
package ses_test

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/exp/ses"
	"github.com/goamz/goamz/testutil"
	gocheck "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

type S struct {
	ses *ses.SES
}

var _ = gocheck.Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	s.ses = ses.NewSES(auth, aws.Region{Name: "faux-region-1", S3Endpoint: testServer.URL})
}

func (s *S) TearDownStrategy(c *gocheck.C) {

}

func (s *S) SetUpTest(c *gocheck.C) {

}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

func (s *S) TestSendEmail(c *gocheck.C) {
	testServer.Response(200, nil, "")

	email := ses.NewEmail()
	email.AddTo("test@test.com")
	email.AddSource("test@test.com")
	email.SetSubject("test")
	email.SetBodyHtml("test")

	s.ses.SendEmail(email)
	req := testServer.WaitRequest()

	c.Assert(req.Method, gocheck.Equals, "POST")
	c.Assert(req.Header["Date"], gocheck.Not(gocheck.Equals), "")

	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	body, _ := url.ParseQuery(buf.String())

	c.Assert(body, gocheck.Not(gocheck.IsNil))
	c.Assert(body["Destination.ToAddresses.member.1"], gocheck.Equals, "test@test.com")
	c.Assert(body["Source"], gocheck.Equals, "test@test.com")
	c.Assert(body["Message.Subject.Data"], gocheck.Equals, "test")
	c.Assert(body["Message.Body.Html.Data"], gocheck.Equals, "test")
}
