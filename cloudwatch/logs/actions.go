package logs

import "encoding/json"

// Actions for AWS CloudWatch Logs service - see http://goo.gl/aGKDYr

// CreateLogGroup - see http://goo.gl/IXnbI6
func (l *CloudWatchLogs) CreateLogGroup(name string) error {
	query := NewEmptyQuery().AddLogGroupName(name)
	_, err := l.doQuery(target("CreateLogGroup"), query)
	return err
}

// CreateLogStream - see http://goo.gl/1IEVYN
func (l *CloudWatchLogs) CreateLogStream(groupName, streamName string) error {
	query := NewEmptyQuery().AddLogGroupName(groupName).AddLogStreamName(streamName)
	_, err := l.doQuery(target("CreateLogStream"), query)
	return err
}

// DeleteLogGroup - see http://goo.gl/jz4l6O
func (l *CloudWatchLogs) DeleteLogGroup(name string) error {
	query := NewEmptyQuery().AddLogGroupName(name)
	_, err := l.doQuery(target("DeleteLogGroup"), query)
	return err
}

// DeleteLogStream - see http://goo.gl/izxZRW
func (l *CloudWatchLogs) DeleteLogStream(groupName, streamName string) error {
	query := NewEmptyQuery().AddLogGroupName(groupName).AddLogStreamName(streamName)
	_, err := l.doQuery(target("DeleteLogStream"), query)
	return err
}

// DescribeLogGroups - see http://goo.gl/ChjhjZ
func (l *CloudWatchLogs) DescribeLogGroups(
	prefix string, limit int, token string) (*DescribeLogGroupsResult, error) {

	// define query
	query := NewEmptyQuery()
	query.AddLimit(limit).AddLogGroupNamePrefix(prefix).AddNextToken(token)
	// perform query
	body, err := l.doQuery(target("DescribeLogGroups"), query)
	if err != nil {
		return nil, err
	}
	// parse and return results
	result := &DescribeLogGroupsResult{}
	err = json.Unmarshal(body, result)
	return result, err
}

// DescribeLogStreams - see http://goo.gl/t95xWC
func (l *CloudWatchLogs) DescribeLogStreams(
	groupName, prefix string, limit int, token string) (
	*DescribeLogStreamsResult, error) {

	// define query
	query := NewEmptyQuery().AddLogGroupName(groupName)
	query.AddLimit(limit).AddLogGroupNamePrefix(prefix).AddNextToken(token)
	// perform query
	body, err := l.doQuery(target("DescribeLogStreams"), query)
	if err != nil {
		return nil, err
	}
	// parse and return results
	result := &DescribeLogStreamsResult{}
	err = json.Unmarshal(body, result)
	return result, err
}

// PutLogEvents - see http://goo.gl/cOKlDV
func (l *CloudWatchLogs) PutLogEvents(
	events []InputLogEvent, groupName, streamName, nextToken string) (
	string, error) {
	// define query
	query := NewEmptyQuery()
	query.AddLogGroupName(groupName).AddLogStreamName(streamName)
	query.AddLogEvents(events).AddSequenceToken(nextToken)
	// perform query
	body, err := l.doQuery(target("PutLogEvents"), query)
	if err != nil {
		return "", err
	}
	// parse and return results
	result := &PutLogEventsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return "", err
	}
	return result.NextSequenceToken, err
}

// GetLogEvents - see http://goo.gl/eRiUsW
func (l *CloudWatchLogs) GetLogEvents(groupName, streamName string,
	startTime, endTime int64, nextToken string, limit int, startFromHead bool) (
	*GetLogEventsResult, error) {
	// define query
	query := NewEmptyQuery().AddStartFromHead(startFromHead)
	query.AddLogGroupName(groupName).AddLogStreamName(streamName)
	query.AddLimit(limit).AddStartTime(startTime).AddEndTime(endTime)
	// perform query
	body, err := l.doQuery(target("GetLogEvents"), query)
	if err != nil {
		return nil, err
	}
	// parse and return results
	result := &GetLogEventsResult{}
	err = json.Unmarshal(body, result)
	return result, err
}

// PutRetentionPolicy - see http://goo.gl/ktuPXT
func (l *CloudWatchLogs) PutRetentionPolicy(
	groupName string, retentionInDays int32) error {
	query := NewEmptyQuery()
	query.AddLogGroupName(groupName).AddRetentionInDays(retentionInDays)
	_, err := l.doQuery(target("PutRetentionPolicy"), query)
	return err
}

// DeleteRetentionPolicy - see http://goo.gl/rmH58u
func (l *CloudWatchLogs) DeleteRetentionPolicy(name string) error {
	query := NewEmptyQuery().AddLogGroupName(name)
	_, err := l.doQuery(target("DeleteRetentionPolicy"), query)
	return err
}

// the value set here is used in the "X-Amz-Target" header,
// and must track the API version
func target(name string) string {
	return "Logs_20140328." + name
}
