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

func NewImageStack(scope constructs.Construct, id string, props *ImageStackProps) (awscdk.Stack, awss3.IBucket) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	imagesBucket := awss3.NewBucket(stack, jsii.String("hc-images-bucket"), &awss3.BucketProps{
		BucketName: jsii.String("hunter-event-sys-uploaded-images"),
		Cors: &[]*awss3.CorsRule{
			{
				AllowedMethods: &[]awss3.HttpMethods{
					awss3.HttpMethods_GET,
					awss3.HttpMethods_PUT,
					awss3.HttpMethods_POST,
				},
				AllowedOrigins: &[]*string{
					jsii.String("*"), // allowing from all origins atm, should be locked down later
				},
				AllowedHeaders: &[]*string{
					jsii.String("*"),
				},
			},
		},
		BlockPublicAccess: awss3.BlockPublicAccess_BLOCK_ALL(),
		EnforceSSL:        jsii.Bool(true),
	})

	// TODO: Add Event Notification to trigger a Lambda to add image details to RDS

	// Needs to be generic for each endpoint that adds images
	// insertImageInDbLambda := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Insert Images into RDS Function"), &awscdklambdagoalpha.GoFunctionProps{
	// 	FunctionName: jsii.String("InsertS3ImageDetails"),
	// 	Entry:        jsii.String("./lambda/images/main.go"),
	// })

	// imagesBucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3notifications.NewLambdaDestination(insertImageInDbLambda), nil)

	return stack, imagesBucket
}
