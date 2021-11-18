package aws

var USGovWest = Region{
	"us-gov-west-1",
	"https://ec2.us-gov-west-1.amazonaws.com",
	"https://s3-fips-us-gov-west-1.amazonaws.com",
	"",
	true,
	true,
	"",
	"",
	"https://sns.us-gov-west-1.amazonaws.com",
	"https://sqs.us-gov-west-1.amazonaws.com",
	"https://iam.us-gov.amazonaws.com",
	"https://elasticloadbalancing.us-gov-west-1.amazonaws.com",
	"https://dynamodb.us-gov-west-1.amazonaws.com",
	ServiceInfo{"https://monitoring.us-gov-west-1.amazonaws.com", V2Signature},
	"https://autoscaling.us-gov-west-1.amazonaws.com",
	ServiceInfo{"https://rds.us-gov-west-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.us-gov-west-1.amazonaws.com",
	"https://ecs.us-gov-west-1.amazonaws.com",
	"https://streams.dynamodb.us-gov-west-1.amazonaws.com",
}

var USEast = Region{
	"us-east-1",
	"https://ec2.us-east-1.amazonaws.com",
	"https://s3.amazonaws.com",
	"",
	false,
	false,
	"https://sdb.amazonaws.com",
	"https://email.us-east-1.amazonaws.com",
	"https://sns.us-east-1.amazonaws.com",
	"https://sqs.us-east-1.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.us-east-1.amazonaws.com",
	"https://dynamodb.us-east-1.amazonaws.com",
	ServiceInfo{"https://monitoring.us-east-1.amazonaws.com", V2Signature},
	"https://autoscaling.us-east-1.amazonaws.com",
	ServiceInfo{"https://rds.us-east-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.us-east-1.amazonaws.com",
	"https://ecs.us-east-1.amazonaws.com",
	"https://streams.dynamodb.us-east-1.amazonaws.com",
}

var USEast2 = Region{
	"us-east-2",
	"https://ec2.us-east-2.amazonaws.com",
	"https://s3.amazonaws.com",
	"",
	true,
	true,
	"",
	"",
	"https://sns.us-east-2.amazonaws.com",
	"https://sqs.us-east-2.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.us-east-2.amazonaws.com",
	"https://dynamodb.us-east-2.amazonaws.com",
	ServiceInfo{"https://monitoring.us-east-2.amazonaws.com", V2Signature},
	"https://autoscaling.us-east-1.amazonaws.com",
	ServiceInfo{"https://rds.us-east-2.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.us-east-2.amazonaws.com",
	"https://ecs.us-east-2.amazonaws.com",
	"https://streams.dynamodb.us-east-2.amazonaws.com",
}

var USWest = Region{
	"us-west-1",
	"https://ec2.us-west-1.amazonaws.com",
	"https://s3-us-west-1.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.us-west-1.amazonaws.com",
	"",
	"https://sns.us-west-1.amazonaws.com",
	"https://sqs.us-west-1.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.us-west-1.amazonaws.com",
	"https://dynamodb.us-west-1.amazonaws.com",
	ServiceInfo{"https://monitoring.us-west-1.amazonaws.com", V2Signature},
	"https://autoscaling.us-west-1.amazonaws.com",
	ServiceInfo{"https://rds.us-west-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.us-west-1.amazonaws.com",
	"https://ecs.us-west-1.amazonaws.com",
	"https://streams.dynamodb.us-west-1.amazonaws.com",
}

var USWest2 = Region{
	"us-west-2",
	"https://ec2.us-west-2.amazonaws.com",
	"https://s3-us-west-2.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.us-west-2.amazonaws.com",
	"https://email.us-west-2.amazonaws.com",
	"https://sns.us-west-2.amazonaws.com",
	"https://sqs.us-west-2.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.us-west-2.amazonaws.com",
	"https://dynamodb.us-west-2.amazonaws.com",
	ServiceInfo{"https://monitoring.us-west-2.amazonaws.com", V2Signature},
	"https://autoscaling.us-west-2.amazonaws.com",
	ServiceInfo{"https://rds.us-west-2.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.us-west-2.amazonaws.com",
	"https://ecs.us-west-2.amazonaws.com",
	"https://streams.dynamodb.us-west-2.amazonaws.com",
}

var EUWest = Region{
	"eu-west-1",
	"https://ec2.eu-west-1.amazonaws.com",
	"https://s3-eu-west-1.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.eu-west-1.amazonaws.com",
	"https://email.eu-west-1.amazonaws.com",
	"https://sns.eu-west-1.amazonaws.com",
	"https://sqs.eu-west-1.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.eu-west-1.amazonaws.com",
	"https://dynamodb.eu-west-1.amazonaws.com",
	ServiceInfo{"https://monitoring.eu-west-1.amazonaws.com", V2Signature},
	"https://autoscaling.eu-west-1.amazonaws.com",
	ServiceInfo{"https://rds.eu-west-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.eu-west-1.amazonaws.com",
	"https://ecs.eu-west-1.amazonaws.com",
	"https://streams.dynamodb.eu-west-1.amazonaws.com",
}

var EUWest2 = Region{
	"eu-west-2",
	"https://ec2.eu-west-2.amazonaws.com",
	"https://s3-eu-west-2.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.eu-west-2.amazonaws.com",
	"https://email.eu-west-2.amazonaws.com",
	"https://sns.eu-west-2.amazonaws.com",
	"https://sqs.eu-west-2.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.eu-west-2.amazonaws.com",
	"https://dynamodb.eu-west-2.amazonaws.com",
	"https://monitoring.eu-west-2.amazonaws.com",
	"https://autoscaling.eu-west-2.amazonaws.com",
	"https://rds.eu-west-2.amazonaws.com",
	"https://sts.amazonaws.com",
	"https://cloudformation.eu-west-2.amazonaws.com",
	"https://ecs.eu-west-2.amazonaws.com",
	"https://streams.dynamodb.eu-west-2.amazonaws.com",
	SignV4Region("eu-west-2"),
}

var EUCentral = Region{
	"eu-central-1",
	"https://ec2.eu-central-1.amazonaws.com",
	"https://s3-eu-central-1.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.eu-central-1.amazonaws.com",
	"https://email.eu-central-1.amazonaws.com",
	"https://sns.eu-central-1.amazonaws.com",
	"https://sqs.eu-central-1.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.eu-central-1.amazonaws.com",
	"https://dynamodb.eu-central-1.amazonaws.com",
	ServiceInfo{"https://monitoring.eu-central-1.amazonaws.com", V2Signature},
	"https://autoscaling.eu-central-1.amazonaws.com",
	ServiceInfo{"https://rds.eu-central-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.eu-central-1.amazonaws.com",
	"https://ecs.eu-central-1.amazonaws.com",
	"https://streams.dynamodb.eu-central-1.amazonaws.com",
}

var APSoutheast = Region{
	"ap-southeast-1",
	"https://ec2.ap-southeast-1.amazonaws.com",
	"https://s3-ap-southeast-1.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.ap-southeast-1.amazonaws.com",
	"",
	"https://sns.ap-southeast-1.amazonaws.com",
	"https://sqs.ap-southeast-1.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.ap-southeast-1.amazonaws.com",
	"https://dynamodb.ap-southeast-1.amazonaws.com",
	ServiceInfo{"https://monitoring.ap-southeast-1.amazonaws.com", V2Signature},
	"https://autoscaling.ap-southeast-1.amazonaws.com",
	ServiceInfo{"https://rds.ap-southeast-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.ap-southeast-1.amazonaws.com",
	"https://ecs.ap-southeast-1.amazonaws.com",
	"https://streams.dynamodb.ap-southeast-1.amazonaws.com",
}

var APSoutheast2 = Region{
	"ap-southeast-2",
	"https://ec2.ap-southeast-2.amazonaws.com",
	"https://s3-ap-southeast-2.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.ap-southeast-2.amazonaws.com",
	"",
	"https://sns.ap-southeast-2.amazonaws.com",
	"https://sqs.ap-southeast-2.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.ap-southeast-2.amazonaws.com",
	"https://dynamodb.ap-southeast-2.amazonaws.com",
	ServiceInfo{"https://monitoring.ap-southeast-2.amazonaws.com", V2Signature},
	"https://autoscaling.ap-southeast-2.amazonaws.com",
	ServiceInfo{"https://rds.ap-southeast-2.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.ap-southeast-2.amazonaws.com",
	"https://ecs.ap-southeast-2.amazonaws.com",
	"https://streams.dynamodb.ap-southeast-2.amazonaws.com",
}

var APNortheast = Region{
	"ap-northeast-1",
	"https://ec2.ap-northeast-1.amazonaws.com",
	"https://s3-ap-northeast-1.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.ap-northeast-1.amazonaws.com",
	"",
	"https://sns.ap-northeast-1.amazonaws.com",
	"https://sqs.ap-northeast-1.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.ap-northeast-1.amazonaws.com",
	"https://dynamodb.ap-northeast-1.amazonaws.com",
	ServiceInfo{"https://monitoring.ap-northeast-1.amazonaws.com", V2Signature},
	"https://autoscaling.ap-northeast-1.amazonaws.com",
	ServiceInfo{"https://rds.ap-northeast-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.ap-northeast-1.amazonaws.com",
	"https://ecs.ap-northeast-1.amazonaws.com",
	"https://streams.dynamodb.ap-northeast-1.amazonaws.com",
}

var APNortheast2 = Region{
	"ap-northeast-2",
	"https://ec2.ap-northeast-2.amazonaws.com",
	"https://s3-ap-northeast-2.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.ap-northeast-2.amazonaws.com",
	"",
	"https://sns.ap-northeast-2.amazonaws.com",
	"https://sqs.ap-northeast-2.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.ap-northeast-2.amazonaws.com",
	"https://dynamodb.ap-northeast-2.amazonaws.com",
	ServiceInfo{"https://monitoring.ap-northeast-2.amazonaws.com", V2Signature},
	"https://autoscaling.ap-northeast-2.amazonaws.com",
	ServiceInfo{"https://rds.ap-northeast-2.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.ap-northeast-2.amazonaws.com",
	"https://ecs.ap-northeast-2.amazonaws.com",
	"https://streams.dynamodb.ap-northeast-2.amazonaws.com",
}

var SAEast = Region{
	"sa-east-1",
	"https://ec2.sa-east-1.amazonaws.com",
	"https://s3-sa-east-1.amazonaws.com",
	"",
	true,
	true,
	"https://sdb.sa-east-1.amazonaws.com",
	"",
	"https://sns.sa-east-1.amazonaws.com",
	"https://sqs.sa-east-1.amazonaws.com",
	"https://iam.amazonaws.com",
	"https://elasticloadbalancing.sa-east-1.amazonaws.com",
	"https://dynamodb.sa-east-1.amazonaws.com",
	ServiceInfo{"https://monitoring.sa-east-1.amazonaws.com", V2Signature},
	"https://autoscaling.sa-east-1.amazonaws.com",
	ServiceInfo{"https://rds.sa-east-1.amazonaws.com", V2Signature},
	"https://sts.amazonaws.com",
	"https://cloudformation.sa-east-1.amazonaws.com",
	"https://ecs.sa-east-1.amazonaws.com",
	"https://streams.dynamodb.sa-east-1.amazonaws.com",
}

var CNNorth = Region{
	"cn-north-1",
	"https://ec2.cn-north-1.amazonaws.com.cn",
	"https://s3.cn-north-1.amazonaws.com.cn",
	"",
	true,
	true,
	"https://sdb.cn-north-1.amazonaws.com.cn",
	"",
	"https://sns.cn-north-1.amazonaws.com.cn",
	"https://sqs.cn-north-1.amazonaws.com.cn",
	"https://iam.cn-north-1.amazonaws.com.cn",
	"https://elasticloadbalancing.cn-north-1.amazonaws.com.cn",
	"https://dynamodb.cn-north-1.amazonaws.com.cn",
	ServiceInfo{"https://monitoring.cn-north-1.amazonaws.com.cn", V4Signature},
	"https://autoscaling.cn-north-1.amazonaws.com.cn",
	ServiceInfo{"https://rds.cn-north-1.amazonaws.com.cn", V4Signature},
	"https://sts.cn-north-1.amazonaws.com.cn",
	"https://cloudformation.cn-north-1.amazonaws.com.cn",
	"https://ecs.cn-north-1.amazonaws.com.cn",
	"https://streams.dynamodb.cn-north-1.amazonaws.com.cn",
}
