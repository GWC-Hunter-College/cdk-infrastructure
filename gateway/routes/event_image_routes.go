package gateway_routes

import (
	gateway_helpers "cdk-infrastructure/gateway/helpers"
	"cdk-infrastructure/gateway/integrations"
	gateway_parameters "cdk-infrastructure/gateway/parameters"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

// Helper function to add event image routes to the API.
//
// Should be noted that this also includes the event images and thumbnail endpoints even though
// they are routed under clubs because they are event-specific.
//
// Routes:
//
/// Event Images:
// POST /clubs/{clubId}/events/{eventId}/images
// GET /clubs/{clubId}/events/{eventId}/images
// POST /clubs/{clubId}/events/{eventId}/images/confirm

// POST /clubs/{clubId}/events/{eventId}/thumbnails
func EventImageRoutes(
	httpApi awsapigatewayv2.HttpApi,
	stack awscdk.Stack,
	vpc awsec2.Vpc,
	s3Params gateway_parameters.S3PermissionsParameters,
	dbParams gateway_parameters.DatabaseConnectionParameters,
) {
	dbInstance := dbParams.DbInstance
	dbSecret := dbParams.Secret

	/// Event Images
	// GET /clubs/{clubId}/events/{eventId}/images
	getEventImagesInt, getEventImagesFn := integrations.GetEventImagesIntegration(stack, vpc, s3Params, dbParams)
	gateway_helpers.GrantRdsAccessToLambda(getEventImagesFn, dbInstance, dbSecret)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events/{eventId}/images"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
		Integration: getEventImagesInt,
	})

	// POST /clubs/{clubId}/events/{eventId}/images
	postEventImagesInt, _ := integrations.PostEventImagesIntegration(stack, vpc, s3Params)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events/{eventId}/images"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
			awsapigatewayv2.HttpMethod_OPTIONS,
		},
		Integration: postEventImagesInt,
	})

	// POST /clubs/{clubId}/events/{eventId}/images/confirm
	confirmEventImagesInt, confirmEventImagesFn := integrations.ConfirmEventImagesIntegration(stack, vpc, dbParams)
	gateway_helpers.GrantRdsAccessToLambda(confirmEventImagesFn, dbInstance, dbSecret)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events/{eventId}/images/confirm"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
			awsapigatewayv2.HttpMethod_OPTIONS,
		},
		Integration: confirmEventImagesInt,
	})

	// Event Thumbnails
	// POST /clubs/{clubId}/events/{eventId}/thumbnails
	postEventThumbnailsInt, _ := integrations.PostEventThumbnailsIntegration(stack, vpc, s3Params)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events/{eventId}/thumbnails"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
			awsapigatewayv2.HttpMethod_OPTIONS,
		},
		Integration: postEventThumbnailsInt,
	})
}
