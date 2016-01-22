package aws

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

type RetryableFunc func(*http.Request, *http.Response, error) bool
type WaitFunc func(try int, minWait time.Duration, maxWait time.Duration)
type DeadlineFunc func() time.Time

type ResilientTransport struct {
	// Timeout is the maximum amount of time a dial will wait for
	// a connect to complete.
	//
	// The default is no timeout.
	//
	// With or without a timeout, the operating system may impose
	// its own earlier timeout. For instance, TCP timeouts are
	// often around 3 minutes.
	DialTimeout time.Duration
	Deadline    DeadlineFunc
	ShouldRetry RetryableFunc
	Wait        WaitFunc
	transport   *http.Transport

	// Retries are only attempted for temporary network errors or known
	// safe failures.

	// MinRetryWait is the minimum time to wait between retries
	MinRetryWait time.Duration

	// MaxRetryWait is the maximum time to wait between retries
	MaxRetryWait time.Duration

	// Total cumulative time before giving up on retries. Note
	// this is only an estimate. Actual timeout
	// may be upto RetryingTimeout + MaxRetryWait
	// The value for this should usually be less that 15 minutes
	// as the signatures in the request becomes invalid
	// after 15 minutes. From AWS doc: http://goo.gl/TNqMCr
	// A request must reach AWS within 15 minutes of the time stamp
	// in the request. Otherwise, AWS denies the request.
	RetryingTimeout time.Duration
}

// Convenience method for creating an http client
func NewClient(rt *ResilientTransport) *http.Client {
	rt.transport = &http.Transport{
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, rt.DialTimeout)
			if err != nil {
				return nil, err
			}
			err = c.SetDeadline(rt.Deadline())
			return c, err
		},
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: -1,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	// TODO: Would be nice is ResilientTransport allowed clients to initialize
	// with http.Transport attributes.
	return &http.Client{
		Transport: rt,
	}
}

var retryingTransport = &ResilientTransport{
	Deadline: func() time.Time {
		return time.Now().Add(10 * time.Second)
	},
	DialTimeout:     10 * time.Second,
	ShouldRetry:     AwsRetry,
	Wait:            ExpBackoff,
	MaxRetryWait:    10 * time.Second,
	MinRetryWait:    100 * time.Millisecond,
	RetryingTimeout: 60 * time.Second,
}

// Exported default client
var RetryingClient = NewClient(retryingTransport)

func (t *ResilientTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.tries(req)
}

// Retry a request a maximum of t.MaxTries times.
// We'll only retry if the proper criteria are met.
// If a wait function is specified, wait that amount of time
// In between requests.
func (t *ResilientTransport) tries(req *http.Request) (res *http.Response, err error) {
	// Save a copy of the body.
	buf := new(bytes.Buffer)
	resetBody := false
	if req.Body != nil {
		resetBody = true
		buf.ReadFrom(req.Body)
	}

	// deadline for retrying
	retryingDeadline := time.Now().Add(t.RetryingTimeout)
	retries := 0
	for { // Watch the infinite loop here
		// Each retry should use a copy of the body.
		// This fixes a bug where subsequent retries would be using a reader
		// that was already read.
		if resetBody {
			req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		}
		logRequest(req)
		res, err = t.transport.RoundTrip(req)

		if !t.ShouldRetry(req, res, err) {
			if res != nil {
				logResponse(res)
			}
			break
		}

		if res != nil {
			res.Body.Close()
		}

		if retryingDeadline.Sub(time.Now()) > 0 && t.Wait != nil {
			t.Wait(retries, t.MinRetryWait, t.MaxRetryWait)
			retries += 1
			continue
		}
		break
	}
	return
}

func ExpBackoff(try int, minWait time.Duration, maxWait time.Duration) {
	wait := time.Duration(math.Pow(2, float64(try))) * minWait
	if wait < minWait || wait > maxWait { // check < minWait to deal with overflow
		wait = maxWait
	}

	log.Warnf("Waiting %v before retry #%d\n", wait, try+1)
	time.Sleep(wait)
}

func LinearBackoff(try int) {
	time.Sleep(time.Duration(try*100) * time.Millisecond)
}

// Decide if we should retry a request.
// In general, the criteria for retrying a request is described here
// http://docs.aws.amazon.com/general/latest/gr/api-retries.html
func AwsRetry(req *http.Request, res *http.Response, err error) bool {
	noSuchHostErr := "no such host"
	if err != nil {
		if err == io.EOF {
			log.Warnf(
				"Retryable network IO error on (%s %s)\n%s",
				req.Method,
				req.URL.String(),
				err)
			return true
		} else if neterr, ok := err.(net.Error); ok && neterr.Temporary() {
			log.Warnf(
				"Retryable network error on (%s %s)\n%s",
				req.Method,
				req.URL.String(),
				err)
			return true
		} else if operr, ok := neterr.(*net.OpError); ok {
			dnsErrStr := ""
			if dnsErr, ok := operr.Err.(*net.DNSError); ok {
				// extract error string from DNSErr. Examples of DNSError
				// "lookup login.microsoftonline.com on 10.0.2.4:53: no such host"
				// "lookup management.azure.com on 10.0.2.4:53: no such host",
				// breaks down to DNSErr:{Err:"no such host", Name:"management.azure.com",
				// Server:"10.0.2.4:53", IsTimeout:false}
				dnsErrStr = strings.TrimSpace(strings.ToLower(dnsErr.Err))
			}
			e := operr.Err.Error()
			if e == syscall.ECONNRESET.Error() || e == syscall.ECONNABORTED.Error() ||
				dnsErrStr == noSuchHostErr {
				log.Warnf(
					"Retryable network error on (%s %s)\n%s",
					req.Method,
					req.URL.String(),
					err)
				return true
			}
		}
		return false // non-retryable error
	}

	if res == nil {
		return false // no error, no response - retryable?
	}

	// Retry if we get a 5xx series error.
	if res.StatusCode >= 500 && res.StatusCode < 600 {
		dump, _ := httputil.DumpResponse(res, false)
		log.Warnf(
			"Retryable error on (%s %s)\n%v", req.Method, req.URL.String(), string(dump))
		return true
	}

	// Check the body to see if it matches ContentLength
	// XXX This is a temporary hack to work around an issue where
	// it appears deadline is triggering a timeout without returning a
	// timeout error.
	if res.ContentLength > 0 && res.Body != nil {
		body, _ := ioutil.ReadAll(res.Body)                // Read the body
		res.Body = ioutil.NopCloser(bytes.NewReader(body)) // Restore the reader
		bodyLen := int64(len(body))
		if bodyLen != res.ContentLength {
			dump, _ := httputil.DumpResponse(res, true)
			log.Warnf("Retryable error. Content length mismatch (%d vs %d).\n%s %s\n%v",
				res.ContentLength, bodyLen, req.Method, req.URL.String(), string(dump))
			return true
		}
	}
	return false // everything is fine. no need to retry
}

func logRequest(req *http.Request) {
	level := traceLevel()
	if level >= 1 {
		log.Debugf("%s %v", req.Method, req.URL.String())
	}
}

func logResponse(res *http.Response) {
	level := traceLevel()
	if level == 1 {
		dump, _ := httputil.DumpResponse(res, false)
		log.Debugf("%v", string(dump))
	} else if level == 2 {
		dump, _ := httputil.DumpResponse(res, true)
		log.Debugf("%v", string(dump))
	}
}

func traceLevel() int64 { // TODO: replace with glog
	t, err := strconv.ParseInt(os.Getenv("TRACE"), 10, 0)
	if err != nil {
		return 0
	}
	return t
}
