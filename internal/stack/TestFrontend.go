package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"       // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3" // ← FIXED: v2 path!

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type FrontendStackProps struct {
	awscdk.StackProps
}

func NewFrontendStack(scope constructs.Construct, id string, props *FrontendStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	websiteBucket := awss3.NewBucket(stack,
		jsii.String("testBucket"), // logical ID
		&awss3.BucketProps{
			BucketName:        jsii.String("test-site-for-hunter"),
			PublicReadAccess:  jsii.Bool(false),
			RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
			AutoDeleteObjects: jsii.Bool(true),
		})

	// Output S3 bucket name
	awscdk.NewCfnOutput(stack, jsii.String("websiteBucketName"), &awscdk.CfnOutputProps{
		Value: websiteBucket.BucketName(),
	})

	// var cloudfrontDistribution cloudfront.Distribution

	return stack
}
