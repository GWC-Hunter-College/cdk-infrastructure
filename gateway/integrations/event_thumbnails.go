package integrations

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

// Integration for event image presign endpoint
func EventThumbnailsIntegration(stack awscdk.Stack, bucket awss3.IBucket) awsapigatewayv2integrations.HttpLambdaIntegration {
	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("EventThumbnailPresignFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("ClubEventThumbnailPresign"),
		Entry:        jsii.String("lambda/api/clubs/events/thumbnails/post.go"),
	})

	bucket.GrantPut(function, "events/*/thumbnails/*")

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("EventThumbnailPresignIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration
}
