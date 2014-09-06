package cloudformation_test

var CancelUpdateStackResponse = `
<CancelUpdateStackResponse xmlns="http://cloudformation.amazonaws.com/doc/2010-05-15/">
  <CancelUpdateStackResult/>
  <ResponseMetadata>
    <RequestId>4af14eec-350e-11e4-b260-EXAMPLE</RequestId>
  </ResponseMetadata>
</CancelUpdateStackResponse>
`

var CreateStackResponse = `
<CreateStackResponse xmlns="http://cloudformation.amazonaws.com/doc/2010-05-15/">
<CreateStackResult>
  <StackId>arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83</StackId>
</CreateStackResult>
 <ResponseMetadata>
    <RequestId>4af14eec-350e-11e4-b260-EXAMPLE</RequestId>
  </ResponseMetadata>
</CreateStackResponse>
`

var CreateStackWithInvalidParamsResponse = `
<ErrorResponse xmlns="http://cloudformation.amazonaws.com/doc/2010-05-15/">
  <Error>
    <Type>Sender</Type>
    <Code>ValidationError</Code>
    <Message>Either Template URL or Template Body must be specified.</Message>
  </Error>
  <RequestId>70a76d42-9665-11e2-9fdf-211deEXAMPLE</RequestId>
</ErrorResponse>
`

var DeleteStackResponse = `
<DeleteStackResponse xmlns="http://cloudformation.amazonaws.com/doc/2010-05-15/">
  <DeleteStackResult/>
  <ResponseMetadata>
    <RequestId>4af14eec-350e-11e4-b260-EXAMPLE</RequestId>
  </ResponseMetadata>
</DeleteStackResponse>
`
var DescribeStackEventsResponse = `
<DescribeStackEventsResponse xmlns="http://cloudformation.amazonaws.com/doc/2010-05-15/">
  <DescribeStackEventsResult>
    <StackEvents>
      <member>
        <EventId>Event-1-Id</EventId>
        <StackId>arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83</StackId>
        <StackName>MyStack</StackName>
        <LogicalResourceId>MyStack</LogicalResourceId>
        <PhysicalResourceId>MyStack_One</PhysicalResourceId>
        <ResourceType>AWS::CloudFormation::Stack</ResourceType>
        <Timestamp>2010-07-27T22:26:28Z</Timestamp>
        <ResourceStatus>CREATE_IN_PROGRESS</ResourceStatus>
        <ResourceStatusReason>User initiated</ResourceStatusReason>
      </member>
      <member>
        <EventId>Event-2-Id</EventId>
        <StackId>arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83</StackId>
        <StackName>MyStack</StackName>
        <LogicalResourceId>MyDBInstance</LogicalResourceId>
        <PhysicalResourceId>MyStack_DB1</PhysicalResourceId>
        <ResourceType>AWS::SecurityGroup</ResourceType>
        <Timestamp>2010-07-27T22:27:28Z</Timestamp>
        <ResourceStatus>CREATE_IN_PROGRESS</ResourceStatus>
        <ResourceProperties>{"GroupDescription":...}</ResourceProperties>
      </member>
      <member>
        <EventId>Event-3-Id</EventId>
        <StackId>arn:aws:cloudformation:us-east-1:123456789:stack/MyStack/aaf549a0-a413-11df-adb3-5081b3858e83</StackId>
        <StackName>MyStack</StackName>
        <LogicalResourceId>MySG1</LogicalResourceId>
        <PhysicalResourceId>MyStack_SG1</PhysicalResourceId>
        <ResourceType>AWS::SecurityGroup</ResourceType>
        <Timestamp>2010-07-27T22:28:28Z</Timestamp>
        <ResourceStatus>CREATE_COMPLETE</ResourceStatus>
      </member>
    </StackEvents>
    <NextToken/>
  </DescribeStackEventsResult>
  <ResponseMetadata>
    <RequestId>4af14eec-350e-11e4-b260-EXAMPLE</RequestId>
  </ResponseMetadata>
</DescribeStackEventsResponse>
`
