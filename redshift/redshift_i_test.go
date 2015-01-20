package redshift_test

import (
	"github.com/czos/goamz/aws"
	"github.com/czos/goamz/redshift"
	"github.com/czos/goamz/testutil"
	"github.com/motain/gocheck"
)

// AmazonServer represents an Amazon server.
type AmazonServer struct {
	auth aws.Auth
}

func (s *AmazonServer) SetUp(c *gocheck.C) {
	auth, err := aws.EnvAuth()
	if err != nil {
		c.Fatal(err.Error())
	}
	s.auth = auth
}

// Suite cost per run: 0.02 USD
var _ = gocheck.Suite(&AmazonClientSuite{})

// AmazonClientSuite tests the client against a live EC2 server.
type AmazonClientSuite struct {
	srv AmazonServer
	ClientTests
}

func (s *AmazonClientSuite) SetUpSuite(c *gocheck.C) {
	if !testutil.Amazon {
		c.Skip("AmazonClientSuite tests not enabled")
	}
	s.srv.SetUp(c)
	s.redshift = redshift.NewWithClient(s.srv.auth, aws.USEast, testutil.DefaultClient)
}

// ClientTests defines integration tests designed to test the client.
// It is not used as a test suite in itself, but embedded within
// another type.
type ClientTests struct {
	redshift *redshift.Redshift
}

// Cost: 0.00 USD
func (s *ClientTests) TestDescribeClusters(c *gocheck.C) {
	resp, err := s.redshift.DescribeClusters("", []string{}, []string{}, "", 0)
	c.Assert(err, gocheck.IsNil)

	c.Assert(resp.RequestId, gocheck.Matches, ".+")
}
