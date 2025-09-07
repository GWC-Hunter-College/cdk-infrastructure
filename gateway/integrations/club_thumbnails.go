package integrations

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

// Integration for club thumbnail presign endpoint
func ClubThumbnailsIntegration(stack awscdk.Stack, bucket awss3.IBucket) awsapigatewayv2integrations.HttpLambdaIntegration {
	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("ClubThumbnailPresignFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("ClubThumbnailPresign"),
		Entry:        jsii.String("lambda/api/clubs/thumbnails/post/post.go"),
		Environment: &map[string]*string{
			"S3_BUCKET": bucket.BucketName(),
		},
	})

	bucket.GrantPut(function, "clubs/*/thumbnails/*")

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("ClubThumbnailPresignIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration
}
