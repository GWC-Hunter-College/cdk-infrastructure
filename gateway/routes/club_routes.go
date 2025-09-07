package gateway_routes

import (
	"cdk-infrastructure/gateway/integrations"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

// Helper function to add club routes to the API.
// Routes:
//
//	POST /clubs/{clubId}/thumbnails
func ClubRoutes(httpApi awsapigatewayv2.HttpApi, stack awscdk.Stack, clubImagesBucket awss3.IBucket) {
	// POST /clubs/{clubId}/thumbnails
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/thumbnails"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
			// awsapigatewayv2.HttpMethod_OPTIONS,
		},
		Integration: integrations.ClubThumbnailsIntegration(stack, clubImagesBucket),
	})
}
