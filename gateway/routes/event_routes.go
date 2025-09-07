package gateway_routes

import (
	"cdk-infrastructure/gateway/integrations"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/jsii-runtime-go"
)

// Helper function to add event routes to the API.
//
// Should be noted that this also includes the event images and thumbnail endpoints even though
// they are routed under clubs because they are event-specific.
//
// Routes:
//
//	POST /clubs/{clubId}/events/{eventId}/images
//	POST /clubs/{clubId}/events/{eventId}/thumbnails
func EventRoutes(httpApi awsapigatewayv2.HttpApi, stack awscdk.Stack, eventImagesBucket awss3.IBucket) {
	// Event Images
	// POST /clubs/{clubId}/events/{eventId}/images
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events/{eventId}/images"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
			// awsapigatewayv2.HttpMethod_OPTIONS,
		},
		Integration: integrations.EventImagesIntegration(stack, eventImagesBucket),
	})

	// Event Thumbnails
	// POST /clubs/{clubId}/events/{eventId}/thumbnails
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events/{eventId}/thumbnails"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
			awsapigatewayv2.HttpMethod_OPTIONS,
		},
		Integration: integrations.EventThumbnailsIntegration(stack, eventImagesBucket),
	})
}
