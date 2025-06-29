package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StorageStackProps struct {
	awscdk.StackProps
}

type StorageStack struct {
	Stack  awscdk.Stack
	Bucket awss3.Bucket
}

func NewStorageStack(scope constructs.Construct, id string, props *StorageStackProps) *StorageStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	imageBucket := awss3.NewBucket(stack,
		jsii.String("ImageBucket"), // logical ID
		&awss3.BucketProps{
			BucketName:        jsii.String("gwc-image-storage"),
			PublicReadAccess:  jsii.Bool(false),
			RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
			AutoDeleteObjects: jsii.Bool(true),
		})

	// MAKE SURE TO TIGHTEN FOR PRODUCTION,
	// ONLY HAVE PUBLIC ACCESS DURING DEVELOPMENT
	imageBucket.AddCorsRule(&awss3.CorsRule{
		AllowedOrigins: &[]*string{
			jsii.String("*"),
		},
		AllowedMethods: &[]awss3.HttpMethods{
			awss3.HttpMethods_PUT,
		},
		AllowedHeaders: &[]*string{
			jsii.String("*"), // Allow all headers
		},
	})

	// Output S3 bucket name
	awscdk.NewCfnOutput(stack, jsii.String("websiteBucketName"), &awscdk.CfnOutputProps{
		Value: imageBucket.BucketName(),
	})

	return &StorageStack{Stack: stack, Bucket: imageBucket}
}
