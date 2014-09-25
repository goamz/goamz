package aws

import (
	"bytes"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"time"
)

type RetryableFunc func(*http.Request, *http.Response, error) bool
type WaitFunc func(try int)
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

	// MaxTries, if non-zero, specifies the number of times we will retry on
	// failure. Retries are only attempted for temporary network errors or known
	// safe failures.
	MaxTries    int
	Deadline    DeadlineFunc
	ShouldRetry RetryableFunc
	Wait        WaitFunc
	transport   *http.Transport
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
		Proxy: http.ProxyFromEnvironment,
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
	DialTimeout: 10 * time.Second,
	MaxTries:    3,
	ShouldRetry: awsRetry,
	Wait:        ExpBackoff,
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

	for try := 0; try < t.MaxTries; try += 1 {
		// Each retry should use a copy of the body.
		// This fixes a bug where subsequent retries would be using a reader
		// that was already read.
		if resetBody {
			req.Body = ioutil.NopCloser(bytes.NewReader(buf.Bytes()))
		}
		res, err = t.transport.RoundTrip(req)

		if !t.ShouldRetry(req, res, err) {
			break
		}
		if res != nil {
			res.Body.Close()
		}
		if t.Wait != nil {
			t.Wait(try)
		}
	}

	return
}

func ExpBackoff(try int) {
	time.Sleep(100 * time.Millisecond *
		time.Duration(math.Exp2(float64(try))))
}

func LinearBackoff(try int) {
	time.Sleep(time.Duration(try*100) * time.Millisecond)
}

// Decide if we should retry a request.
// In general, the criteria for retrying a request is described here
// http://docs.aws.amazon.com/general/latest/gr/api-retries.html
func awsRetry(req *http.Request, res *http.Response, err error) bool {
	// Retry if there's a temporary network error.
	if neterr, ok := err.(net.Error); ok {
		if neterr.Temporary() {
			return true
		}
	}

	// Retry if we get a 5xx series error.
	if res != nil {
		if res.StatusCode >= 500 && res.StatusCode < 600 {
			return true
		}
	}

	// Check the body to see if it matches ContentLength
	if res.ContentLength > 0 && res.Body != nil {
		body, _ := ioutil.ReadAll(res.Body)                // Read the body
		res.Body = ioutil.NopCloser(bytes.NewReader(body)) // Restore the reader
		if int64(len(body)) != res.ContentLength {
			return true
		}
	}

	return false
}
