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
