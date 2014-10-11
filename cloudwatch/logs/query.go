package logs

import (
	"encoding/json"
)

type msi map[string]interface{}
type Query struct {
	buffer msi
}

func NewEmptyQuery() *Query {
	return &Query{msi{}}
}

func (q *Query) AddLogGroupName(logGroupName string) *Query {
	if logGroupName != "" {
		q.buffer["logGroupName"] = logGroupName
	}
	return q
}

func (q *Query) AddLogEvents(logEvents []InputLogEvent) *Query {
	q.buffer["logEvents"] = logEvents
	return q
}

func (q *Query) AddLogStreamName(logStreamName string) *Query {
	if logStreamName != "" {
		q.buffer["logStreamName"] = logStreamName
	}
	return q
}

func (q *Query) AddLimit(limit int) *Query {
	if limit > 0 {
		q.buffer["limit"] = limit
	}
	return q
}

func (q *Query) AddStartFromHead(startFromHead bool) *Query {
	if startFromHead { // default is false - see http://goo.gl/jGh7vu
		q.buffer["startFromHead"] = startFromHead
	}
	return q
}

func (q *Query) AddNextToken(token string) *Query {
	if token != "" {
		q.buffer["nextToken"] = token
	}
	return q
}

func (q *Query) AddStartTime(startTime int64) *Query {
	if startTime > 0 {
		q.buffer["startTime"] = startTime
	}
	return q
}

func (q *Query) AddEndTime(endTime int64) *Query {
	if endTime > 0 {
		q.buffer["endTime"] = endTime
	}
	return q
}

func (q *Query) AddSequenceToken(token string) *Query {
	if token != "" {
		q.buffer["sequenceToken"] = token
	}
	return q
}

func (q *Query) AddLogGroupNamePrefix(prefix string) *Query {
	if prefix != "" {
		q.buffer["logGroupNamePrefix"] = prefix
	}
	return q
}

func (q *Query) AddRetentionInDays(retentionInDays int32) *Query {
	q.buffer["retentionInDays"] = retentionInDays
	return q
}

func (q *Query) String() string {
	bytes, err := json.Marshal(q.buffer)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
