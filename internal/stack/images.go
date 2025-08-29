package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ImageStackProps struct {
	Props awscdk.StackProps
}

func NewImageStack(scope constructs.Construct, id string, props *ImageStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	_ = awss3.NewBucket(stack, jsii.String("club-thumbnails-bucket"), &awss3.BucketProps{
		BucketName: jsii.String("hc-club-thumbnails"),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedMethods: &[]awss3.HttpMethods{
					awss3.HttpMethods_GET,
					awss3.HttpMethods_PUT,
				},
				AllowedOrigins: &[]*string{
					jsii.String("*"), // allowing from all origins atm, should be locked down later
				},
			},
		},
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		EnforceSSL:        jsii.Bool(true),
	})

	_ = awss3.NewBucket(stack, jsii.String("event-thumbnails-bucket"), &awss3.BucketProps{
		BucketName: jsii.String("hc-event-thumbnails"),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedMethods: &[]awss3.HttpMethods{
					awss3.HttpMethods_GET,
					awss3.HttpMethods_PUT,
				},
				AllowedOrigins: &[]*string{
					jsii.String("*"), // allowing from all origins atm, should be locked down later
				},
			},
		},
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		EnforceSSL:        jsii.Bool(true),
	})

	_ = awss3.NewBucket(stack, jsii.String("event-images-bucket"), &awss3.BucketProps{
		BucketName: jsii.String("hc-event-images"),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedMethods: &[]awss3.HttpMethods{
					awss3.HttpMethods_GET,
					awss3.HttpMethods_POST,
				},
				AllowedOrigins: &[]*string{
					jsii.String("*"), // allowing from all origins atm, should be locked down later
				},
			},
		},
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		EnforceSSL:        jsii.Bool(true),
	})

	return stack
}
