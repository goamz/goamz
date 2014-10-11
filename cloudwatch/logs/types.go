package logs

// Data Types for AWS CloudWatch Logs service - see http://goo.gl/e5gDtN

// JSON Response to DescribeLogGroups - see http://goo.gl/5btnIM
type DescribeLogGroupsResult struct {
	LogGroups []LogGroup `json:"logGroups"`
	NextToken string     `json:"nextToken"`
}

// JSON Response to DescribeLogStreams - see http://goo.gl/2gvQZN
type DescribeLogStreamsResult struct {
	LogStreams []LogStream `json:"logStreams"`
	NextToken  string      `json:"nextToken"`
}

// JSON Response to DescribeMetricFilters - see http://goo.gl/MwcLXV
type DescribeMetricFiltersResult struct {
	MetricFilters []MetricFilter `json:"metricFilters"`
	NextToken     string         `json:"nextToken"`
}

// JSON Response to GetLogEvents - see http://goo.gl/jSIZll
type GetLogEventsResult struct {
	Events            []OutputLogEvent `json:"events"`
	NextForwardToken  string           `json:"nextForwardToken"`
	NextBackwardToken string           `json:"nextBackwardToken"`
}

// InputLogEvent - a record of some activity that was recorded by the
// application or resource being monitored - see http://goo.gl/vCh2Hg
type InputLogEvent struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// LogGroup - see http://goo.gl/qeDqDW
type LogGroup struct {
	Arn               string `json:"arn"`
	CreationTime      int64  `json:"creationTime"`
	LogGroupName      string `json:"logGroupName"`
	MetricFilterCount int    `json:"metricFilterCount"`
	StoredBytes       int64  `json:"storedBytes"`
	RetentionInDays   int32  `json:"retentionInDays"`
}

// LogStream - see http://goo.gl/4b92WP
type LogStream struct {
	Arn                 string `json:"arn"`
	CreationTime        int64  `json:"creationTime"`
	FirstEventTimestamp int64  `json:"firstEventTimestamp"`
	LastEventTimestamp  int64  `json:"lastEventTimestamp"`
	LastIngestionTime   int64  `json:"lastIngestionTime"`
	LogStreamName       string `json:"logStreamName"`
	StoredBytes         int64  `json:"storedBytes"`
	UploadSequenceToken string `json:"uploadSequenceToken"`
}

// MetricFilter - see http://goo.gl/Ic2iux
type MetricFilter struct {
	CreationTime          int64                  `json:"creationTime"`
	FilterName            string                 `json:"filterName"`
	FilterPattern         string                 `json:"filterPattern"`
	MetricTransformations []MetricTransformation `json:"metricTransformations"`
}

// MetricFilterMatchRecord - see http://goo.gl/KBo1ZK
type MetricFilterMatchRecord struct {
	EventMessage    string            `json:"eventMessage"`
	EventNumber     int64             `json:"eventNumber"`
	ExtractedValues map[string]string `json:"extractedValues"`
}

// MetricTransformation - see http://goo.gl/dQhuCt
type MetricTransformation struct {
	MetricName      string `json:"metricName"`
	MetricNameSpace string `json:"metricNameSpace"`
	MetricValue     string `json:"metricValue"`
}

// OutputLogEvent - see http://goo.gl/6G91PI
type OutputLogEvent struct {
	IngestionTime int64  `json:"ingestionTime"`
	Message       string `json:"message"`
	Timestamp     int64  `json:"timestamp"`
}

// JSON response to PutLogEvents
type PutLogEventsResult struct {
	NextSequenceToken string `json:"nextSequenceToken"`
}

// JSON response to TestMetricFilter
type TestMetricFilterResult struct {
	Matches []MetricFilterMatchRecord `json:"matches"`
}

// Error represents an error in an operation with CloudWatch Logs
type Error struct {
	StatusCode int // HTTP status code (200, 403, ...)
	Status     string
	Code       string `json:"__type"`
	Message    string `json:"message"`
}
