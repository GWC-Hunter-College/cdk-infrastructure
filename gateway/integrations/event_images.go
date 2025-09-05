package integrations

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

// Integration for event image presign endpoint
func EventImagesIntegration(stack awscdk.Stack, bucket awss3.IBucket) awsapigatewayv2integrations.HttpLambdaIntegration {
	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("EventImagePresignFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("ClubEventImagePresign"),
		Entry:        jsii.String("lambda/api/clubs/events/images"),
	})

	bucket.GrantPut(function, "events/*")

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("EventImagePresignIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration
}
