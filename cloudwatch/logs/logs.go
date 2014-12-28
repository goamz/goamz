//
// cloudwatch/logs: This package provides types and functions to interact with
// the AWS CloudWatch Logs API
//
// Depends on https://github.com/goamz/goamz/aws
//

package logs

import (
	"encoding/json"
	"fmt"
	"github.com/goamz/goamz/aws"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// The CloudWatchLogs type encapsulates operations within a specific EC2 region.
type CloudWatchLogs struct {
	aws.Auth
	aws.Region
}

// New creates a new CloudWatchLogs Client.
func New(auth aws.Auth, region aws.Region) *CloudWatchLogs {
	return &CloudWatchLogs{auth, region}
}

// Returns a human-readable string version of the supplied Error
func (err *Error) Error() string {
	if err.Code == "" {
		return err.Message
	}
	return fmt.Sprintf("[HTTP %d] %s : %s\n", err.StatusCode, err.Code, err.Message)
}

const debug = false

// ----------------------------------------------------------------------------
// Request dispatching logic.

func (l *CloudWatchLogs) doQuery(target string, query *Query) ([]byte, error) {
	// construct HTTP request
	data := strings.NewReader(query.String())
	hreq, err := http.NewRequest("POST", l.Region.CloudWatchLogsEndpoint+"/", data)
	if err != nil {
		return nil, err
	}
	hreq.Header.Set("Content-Type", "application/x-amz-json-1.1")
	hreq.Header.Set("X-Amz-Date", time.Now().UTC().Format(aws.ISO8601BasicFormat))
	hreq.Header.Set("X-Amz-Target", target)
	signer := aws.NewV4Signer(l.Auth, "logs", l.Region)
	signer.Sign(hreq)
	// perform HTTP request
	if debug {
		log.Printf("cloudwatchLogs: request body: %v", query.String())
	}
	resp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		log.Printf("cloudwatchLogs: Error calling Amazon\n: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	// parse the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("cloudwatchLogs: Could not read response body\n")
		return nil, err
	}
	if debug {
		log.Printf("cloudwatchLogs: response code: %v\n", resp.StatusCode)
		log.Printf("cloudwatchLogs: response headers: %v\n", resp.Header)
		log.Printf("cloudwatchLogs: response body: %s\n", string(body))
	}
	// return the body if the response status code was 200, otherwise nil
	if resp.StatusCode != 200 {
		err = buildError(resp, body)
		return nil, err
	}
	return body, nil
}

func buildError(r *http.Response, jsonBody []byte) error {
	cloudwatchLogsError := &Error{
		StatusCode: r.StatusCode,
		Status:     r.Status,
	}

	err := json.Unmarshal(jsonBody, cloudwatchLogsError)
	if err != nil {
		log.Printf("cloudwatchLogs: Failed to parse body as JSON")
		return err
	}

	return cloudwatchLogsError
}
