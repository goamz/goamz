package logs_test

import (
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/cloudwatch/logs"
	"github.com/goamz/goamz/testutil"
	"github.com/motain/gocheck"
	"testing"
)

func Test(t *testing.T) {
	gocheck.TestingT(t)
}

type S struct {
	logs *logs.CloudWatchLogs
}

var _ = gocheck.Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *gocheck.C) {
	testServer.Start()
	auth := aws.Auth{AccessKey: "abc", SecretKey: "123"}
	localRegion := aws.Region{}
	localRegion.CloudWatchLogsEndpoint = testServer.URL
	s.logs = logs.New(auth, localRegion)
}

func (s *S) TearDownTest(c *gocheck.C) {
	testServer.Flush()
}

func getTestLogGroup() *logs.LogGroup {
	group := new(logs.LogGroup)
	group.LogGroupName = "testGroup"
	return group
}

func (s *S) TestCreateLogGroup(c *gocheck.C) {
	testServer.Response(200, nil, "")
	group := getTestLogGroup()

	err := s.logs.CreateLogGroup(group.LogGroupName)
	c.Assert(err, gocheck.IsNil)

	req := testServer.WaitRequest()
	c.Assert(req.Method, gocheck.Equals, "POST")
	c.Assert(req.URL.Path, gocheck.Equals, "/")
	c.Assert(req.Header.Get("X-Amz-Target"),
		gocheck.Equals, "Logs_20140328.CreateLogGroup")
	// c.Assert(body, gocheck.DeepEquals, `{"logGroupName":"testGroup"}`)
}
