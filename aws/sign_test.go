package aws_test

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goamz/goamz/aws"
	. "gopkg.in/check.v1"
)

var _ = Suite(&V4SignerSuite{})

type V4SignerSuite struct {
	auth   aws.Auth
	region aws.Region
	cases  []V4SignerSuiteCase
}

type V4SignerSuiteCase struct {
	label            string
	request          V4SignerSuiteCaseRequest
	canonicalRequest string
	stringToSign     string
	signature        string
	authorization    string
}

type V4SignerSuiteCaseRequest struct {
	method  string
	host    string
	url     string
	headers []string
	body    string
}

func (s *V4SignerSuite) SetUpSuite(c *C) {
	s.auth = aws.Auth{AccessKey: "AKIDEXAMPLE", SecretKey: "wJalrXUtnFEMI/K7MDENG+bPxRfiCYEXAMPLEKEY"}
	s.region = aws.USEast

	// Test cases from the Signature Version 4 Test Suite (http://goo.gl/nguvs0)
	s.cases = append(s.cases,

		// get-header-key-duplicate
		V4SignerSuiteCase{
			label: "get-header-key-duplicate",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"DATE:Mon, 09 Sep 2011 23:36:00 GMT", "ZOO:zoobar", "zoo:foobar", "zoo:zoobar"},
			},
			canonicalRequest: "POST\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\nzoo:foobar,zoobar,zoobar\n\ndate;host;x-amz-content-sha256;zoo\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nb47a355a1f4aed01ddaae2a0992cf11758413038797b8768c7668c04292994f0",
			signature:        "b48d915443fc61838f63037d445e630c55e7d36efa07ce417e28313b06162bb2",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256;zoo, Signature=b48d915443fc61838f63037d445e630c55e7d36efa07ce417e28313b06162bb2",
		},

		// get-header-value-order
		V4SignerSuiteCase{
			label: "get-header-value-order",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"DATE:Mon, 09 Sep 2011 23:36:00 GMT", "p:z", "p:a", "p:p", "p:a"},
			},
			canonicalRequest: "POST\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\np:a,a,p,z\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;p;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n76823ce9cd3e250ee87dbc7142653ecabf6d80a9a91f3199dc7c2ed22eeb668b",
			signature:        "0a5450064fe28a63773c01290500cc67b2e52ffb8efd29af6d15df191e313e8b",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;p;x-amz-content-sha256, Signature=0a5450064fe28a63773c01290500cc67b2e52ffb8efd29af6d15df191e313e8b",
		},

		// get-header-value-trim
		V4SignerSuiteCase{
			label: "get-header-value-trim",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"DATE:Mon, 09 Sep 2011 23:36:00 GMT", "p: phfft "},
			},
			canonicalRequest: "POST\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\np:phfft\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;p;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n5fc8bd66e043705a5414321c921183106f96b3ab394b106383ea62f7de105cb5",
			signature:        "c7876fcc7fa5fdfd0d087d47c110c945d1bd1528828507b9c2c7cff6f3e8d838",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;p;x-amz-content-sha256, Signature=c7876fcc7fa5fdfd0d087d47c110c945d1bd1528828507b9c2c7cff6f3e8d838",
		},

		// get-relative-relative
		V4SignerSuiteCase{
			label: "get-relative-relative",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/foo/bar/../..",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nedf3a954ef9b74cf366ad17b25bc502050eb85bcd3d86b24431cb9aea8761c84",
			signature:        "6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
		},

		// get-relative
		V4SignerSuiteCase{
			label: "get-relative",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/foo/..",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nedf3a954ef9b74cf366ad17b25bc502050eb85bcd3d86b24431cb9aea8761c84",
			signature:        "6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
		},

		// get-slash-dot-slash
		V4SignerSuiteCase{
			label: "get-slash-dot-slash",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/./",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nedf3a954ef9b74cf366ad17b25bc502050eb85bcd3d86b24431cb9aea8761c84",
			signature:        "6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
		},

		// get-slash-pointless-dot
		V4SignerSuiteCase{
			label: "get-slash-pointless-dot",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/./foo",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/foo\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nc0dd19b115fdcd7b17e21b9845d31c5caead32a194b4dd2d720280770c182db3",
			signature:        "b79af6d6802186bdd3c7a22ec726b88ee0ff8bd15cb3ae7fa81d3e1decf2a3f2",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=b79af6d6802186bdd3c7a22ec726b88ee0ff8bd15cb3ae7fa81d3e1decf2a3f2",
		},

		// get-slash
		V4SignerSuiteCase{
			label: "get-slash",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "//",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nedf3a954ef9b74cf366ad17b25bc502050eb85bcd3d86b24431cb9aea8761c84",
			signature:        "6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
		},

		// get-slashes
		V4SignerSuiteCase{
			label: "get-slashes",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "//foo//",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/foo/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n4a7592c23ce09aa58e45d8cb1dc827287cdcf621e98e784f3d86715bb4c4e24e",
			signature:        "cac8d6fc96a93d075848795a33138fa0dfb81a9b9a1387486b67c5803659107b",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=cac8d6fc96a93d075848795a33138fa0dfb81a9b9a1387486b67c5803659107b",
		},

		// get-space
		V4SignerSuiteCase{
			label: "get-space",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/%20/foo",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/%20/foo\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nf8efe2a424f379ce780bcf0b1bbd93a2377fac1cd9b34283473c5d287fb98daa",
			signature:        "280de7094bfa4f41dee69c9af17de9c4652699c284495c6769426170fe9b0407",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=280de7094bfa4f41dee69c9af17de9c4652699c284495c6769426170fe9b0407",
		},

		// get-unreserved
		V4SignerSuiteCase{
			label: "get-unreserved",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n28e6e45f8b295710940c557546b115694d7f5bfe7bc65a3922c1d980f8dcb14a",
			signature:        "fc33535e9ffe760da308afc673033bb2f89254ef7abe2faaf5bf7c37a64cdf91",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=fc33535e9ffe760da308afc673033bb2f89254ef7abe2faaf5bf7c37a64cdf91",
		},

		// get-utf8
		V4SignerSuiteCase{
			label: "get-utf8",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/%E1%88%B4",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/%E1%88%B4\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n02fde5f405aa0ca5dc7851735b4ec75d0dfb6082b92f12c13abcd8f268eea660",
			signature:        "b366e667085aa65f8b3d7be4826dce4507b2d5edd78dde12105fc44014cf65e4",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=b366e667085aa65f8b3d7be4826dce4507b2d5edd78dde12105fc44014cf65e4",
		},

		// get-vanilla-empty-query-key
		V4SignerSuiteCase{
			label: "get-vanilla-empty-query-key",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/?foo=bar",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\nfoo=bar\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n763e50808f6a9e1ba0bf89728d3f34aafb9d3719350e474dbe810bf8e7b270a3",
			signature:        "c3af9ab0d58331f70981239ab27cab9169ce2a87e5b5157b829648b084987e73",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=c3af9ab0d58331f70981239ab27cab9169ce2a87e5b5157b829648b084987e73",
		},

		// get-vanilla-space-query-parameters
		V4SignerSuiteCase{
			label: "get-vanilla-space-query-parameters",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/?foo foo=bar bar",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\nfoo%20foo=bar%20bar\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nf190cec7b75aed4693c771b29c4cfa619a28c2c54c8f7897c4744d1534b23cfa",
			signature:        "d31134b01ce7e5e2356464e64256f9dcfb652301bdd15e171d2b23c7f5e3f626",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=d31134b01ce7e5e2356464e64256f9dcfb652301bdd15e171d2b23c7f5e3f626",
		},

		// get-vanilla-query-order-key-case
		V4SignerSuiteCase{
			label: "get-vanilla-query-order-key-case",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/?foo=Zoo&foo=aha",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\nfoo=Zoo&foo=aha\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nd651e96c251a3375a9248bacb97e05013ed9c36ebe8bb39f24be8dd522f58948",
			signature:        "9f4456ed08b128fe09d5490f93591f1c4eaf3882a04b163920908f72c1ed7243",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=9f4456ed08b128fe09d5490f93591f1c4eaf3882a04b163920908f72c1ed7243",
		},

		// get-vanilla-query-order-key
		V4SignerSuiteCase{
			label: "get-vanilla-query-order-key",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/?a=foo&b=foo",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\na=foo&b=foo\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n1da40574a1309ebbd7b2c4523059d1b6acd99ac198fc6ae5e0f5f90e85726e75",
			signature:        "800aed7d844fd17d0b98bbfcabfb0f4fc79fc9d2aae345baf75d696df36241d7",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=800aed7d844fd17d0b98bbfcabfb0f4fc79fc9d2aae345baf75d696df36241d7",
		},

		// get-vanilla-query-order-value
		V4SignerSuiteCase{
			label: "get-vanilla-query-order-value",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/?foo=b&foo=a",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\nfoo=a&foo=b\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n781593a5fd3c02af91b2956df4c1f56305bc2289cb88179df74f905e9108854b",
			signature:        "0784708ecdb06ffca4f0e2cffb2825849f1fa5c947bd91da312403d3ac060ba8",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=0784708ecdb06ffca4f0e2cffb2825849f1fa5c947bd91da312403d3ac060ba8",
		},

		// get-vanilla-query-unreserved
		V4SignerSuiteCase{
			label: "get-vanilla-query-unreserved",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/?-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz=-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz=-._~0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\neddf477a78e0d7fd1617d20021167184e11634e9b3cccb472fae4881555c7585",
			signature:        "032038308853317f8898d90c3678a1e269e68e30e50fb97c6654eaeff9480799",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=032038308853317f8898d90c3678a1e269e68e30e50fb97c6654eaeff9480799",
		},

		// get-vanilla-query
		V4SignerSuiteCase{
			label: "get-vanilla-query",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nedf3a954ef9b74cf366ad17b25bc502050eb85bcd3d86b24431cb9aea8761c84",
			signature:        "6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
		},

		// get-vanilla-ut8-query
		V4SignerSuiteCase{
			label: "get-vanilla-ut8-query",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/?ሴ=bar",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n%E1%88%B4=bar\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\ne06a77d8fca5b57869e3a5b84b6147b8a38bb8caecbe7e14c2b5d2a0f98df7ec",
			signature:        "05991e997a8d3bab48124e73396fe4826fb47b714072423d765a2a3bdd4a2a7c",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=05991e997a8d3bab48124e73396fe4826fb47b714072423d765a2a3bdd4a2a7c",
		},

		// get-vanilla
		V4SignerSuiteCase{
			label: "get-vanilla",
			request: V4SignerSuiteCaseRequest{
				method:  "GET",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "GET\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nedf3a954ef9b74cf366ad17b25bc502050eb85bcd3d86b24431cb9aea8761c84",
			signature:        "6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=6ad5d9fd6d93c5855f08c4675163da70623fcca6f49f9c96a3890c1c32877a2e",
		},

		// post-header-key-case
		V4SignerSuiteCase{
			label: "post-header-key-case",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"DATE:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "POST\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\ne99db283d41c8d76b65f2bde6003746997018685efaf38028979f3b4c01b9fdd",
			signature:        "1f2d4161c5845a97c13548fd1f6421c1a5930730a177e00c3684579c68a99d1e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=1f2d4161c5845a97c13548fd1f6421c1a5930730a177e00c3684579c68a99d1e",
		},

		// post-header-key-sort
		V4SignerSuiteCase{
			label: "post-header-key-sort",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"DATE:Mon, 09 Sep 2011 23:36:00 GMT", "ZOO:zoobar"},
			},
			canonicalRequest: "POST\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\nzoo:zoobar\n\ndate;host;x-amz-content-sha256;zoo\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nfec201022f7e94b0c824b3485230ca29fd026618dc8d759923aa48bc581762ed",
			signature:        "7daf56768cfc0d366c47044fbf00fda5e7b80043cf55c316317ef29a83e2555c",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256;zoo, Signature=7daf56768cfc0d366c47044fbf00fda5e7b80043cf55c316317ef29a83e2555c",
		},

		// post-header-value-case
		V4SignerSuiteCase{
			label: "post-header-value-case",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"DATE:Mon, 09 Sep 2011 23:36:00 GMT", "zoo:ZOOBAR"},
			},
			canonicalRequest: "POST\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\nzoo:ZOOBAR\n\ndate;host;x-amz-content-sha256;zoo\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\nfe62c1b0b49383f0a25e528d62a5fba9fce82a01db9a9e37dc5f2fd70fe55e08",
			signature:        "9fd3fe58fa020bd61f0899abcb638d335bd91fe85a0bd26a8ee272164a1a122e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256;zoo, Signature=9fd3fe58fa020bd61f0899abcb638d335bd91fe85a0bd26a8ee272164a1a122e",
		},

		// post-vanilla-empty-query-value
		V4SignerSuiteCase{
			label: "post-vanilla-empty-query-value",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/?foo=bar",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "POST\n/\nfoo=bar\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n48d51a068e0fe8aff12e1f1278aad9d750ce3633480760373c9d2d1556ce2ac3",
			signature:        "8e7cde9089c5c8675dc1306869fd1f2089d84f5008c86332bcb4def9245b15a4",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=8e7cde9089c5c8675dc1306869fd1f2089d84f5008c86332bcb4def9245b15a4",
		},

		// post-vanilla-query
		V4SignerSuiteCase{
			label: "post-vanilla-query",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/?foo=bar",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "POST\n/\nfoo=bar\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n48d51a068e0fe8aff12e1f1278aad9d750ce3633480760373c9d2d1556ce2ac3",
			signature:        "8e7cde9089c5c8675dc1306869fd1f2089d84f5008c86332bcb4def9245b15a4",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=8e7cde9089c5c8675dc1306869fd1f2089d84f5008c86332bcb4def9245b15a4",
		},

		// post-vanilla
		V4SignerSuiteCase{
			label: "post-vanilla",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"Date:Mon, 09 Sep 2011 23:36:00 GMT"},
			},
			canonicalRequest: "POST\n/\n\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855\n\ndate;host;x-amz-content-sha256\ne3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\ne99db283d41c8d76b65f2bde6003746997018685efaf38028979f3b4c01b9fdd",
			signature:        "1f2d4161c5845a97c13548fd1f6421c1a5930730a177e00c3684579c68a99d1e",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=date;host;x-amz-content-sha256, Signature=1f2d4161c5845a97c13548fd1f6421c1a5930730a177e00c3684579c68a99d1e",
		},

		// post-x-www-form-urlencoded-parameters
		V4SignerSuiteCase{
			label: "post-x-www-form-urlencoded-parameters",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"Content-Type:application/x-www-form-urlencoded; charset=utf8", "Date:Mon, 09 Sep 2011 23:36:00 GMT"},
				body:    "foo=bar",
			},
			canonicalRequest: "POST\n/\n\ncontent-type:application/x-www-form-urlencoded; charset=utf8\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:3ba8907e7a252327488df390ed517c45b96dead033600219bdca7107d1d3f88a\n\ncontent-type;date;host;x-amz-content-sha256\n3ba8907e7a252327488df390ed517c45b96dead033600219bdca7107d1d3f88a",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\ndfae45948f3fdc5780d014f391042d4e220f577961bb96b84f7bbf8e1a91fe5e",
			signature:        "3bfa5e222990f3f1f14e746a68568a647a210c0705a869e9c19b40d8c86ca309",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=content-type;date;host;x-amz-content-sha256, Signature=3bfa5e222990f3f1f14e746a68568a647a210c0705a869e9c19b40d8c86ca309",
		},

		// post-x-www-form-urlencoded
		V4SignerSuiteCase{
			label: "post-x-www-form-urlencoded",
			request: V4SignerSuiteCaseRequest{
				method:  "POST",
				host:    "host.foo.com",
				url:     "/",
				headers: []string{"Content-Type:application/x-www-form-urlencoded", "Date:Mon, 09 Sep 2011 23:36:00 GMT"},
				body:    "foo=bar",
			},
			canonicalRequest: "POST\n/\n\ncontent-type:application/x-www-form-urlencoded\ndate:Mon, 09 Sep 2011 23:36:00 GMT\nhost:host.foo.com\nx-amz-content-sha256:3ba8907e7a252327488df390ed517c45b96dead033600219bdca7107d1d3f88a\n\ncontent-type;date;host;x-amz-content-sha256\n3ba8907e7a252327488df390ed517c45b96dead033600219bdca7107d1d3f88a",
			stringToSign:     "AWS4-HMAC-SHA256\n20110909T233600Z\n20110909/us-east-1/host/aws4_request\n8434678b03c6f956b391ec4423570bc1f0be17aa53b81118dbba82873a64b700",
			signature:        "771ad8937d4a60e1b926fd12b91c01d904afeada0acf79c019b4c2827c902670",
			authorization:    "AWS4-HMAC-SHA256 Credential=AKIDEXAMPLE/20110909/us-east-1/host/aws4_request, SignedHeaders=content-type;date;host;x-amz-content-sha256, Signature=771ad8937d4a60e1b926fd12b91c01d904afeada0acf79c019b4c2827c902670",
		},
	)
}

func (s *V4SignerSuite) TestCases(c *C) {
	signer := aws.NewV4Signer(s.auth, "host", s.region.Name)

	for _, testCase := range s.cases {

		req, err := http.NewRequest(testCase.request.method, "http://"+testCase.request.host+testCase.request.url, strings.NewReader(testCase.request.body))
		c.Assert(err, IsNil, Commentf("Testcase: %s", testCase.label))
		for _, v := range testCase.request.headers {
			h := strings.SplitN(v, ":", 2)
			req.Header.Add(h[0], h[1])
		}
		req.Header.Set("host", req.Host)

		t := signer.RequestTime(req)

		canonicalRequest, _ := signer.CanonicalRequest(req)
		c.Check(canonicalRequest, Equals, testCase.canonicalRequest, Commentf("Testcase: %s", testCase.label))

		stringToSign := signer.StringToSign(t, canonicalRequest)
		c.Check(stringToSign, Equals, testCase.stringToSign, Commentf("Testcase: %s", testCase.label))

		signature := signer.Signature(t, stringToSign)
		c.Check(signature, Equals, testCase.signature, Commentf("Testcase: %s", testCase.label))

		authorization := signer.Authorization(req.Header, t, signature)
		c.Check(authorization, Equals, testCase.authorization, Commentf("Testcase: %s", testCase.label))

		signer.Sign(req)
		c.Check(req.Header.Get("Authorization"), Equals, testCase.authorization, Commentf("Testcase: %s", testCase.label))
	}
}

func ExampleV4Signer() {
	// Get auth from env vars
	auth, err := aws.EnvAuth()
	if err != nil {
		fmt.Println(err)
	}

	// Create a signer with the auth, name of the service, and aws region
	signer := aws.NewV4Signer(auth, "dynamodb", aws.USEast.Name)

	// Create a request
	req, err := http.NewRequest("POST", aws.USEast.DynamoDBEndpoint, strings.NewReader("sample_request"))
	if err != nil {
		fmt.Println(err)
	}

	// Date or x-amz-date header is required to sign a request
	req.Header.Add("Date", time.Now().UTC().Format(http.TimeFormat))

	// Sign the request
	signer.Sign(req)

	// Issue signed request
	http.DefaultClient.Do(req)
}
