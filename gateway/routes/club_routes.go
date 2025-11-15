package gateway_routes

import (
	"cdk-infrastructure/gateway/integrations"
	gateway_parameters "cdk-infrastructure/gateway/parameters"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/jsii-runtime-go"
)

// Helper function to add club routes to the API.
// Routes:
//
// POST /clubs/{clubId}/events
func ClubRoutes(
	httpApi awsapigatewayv2.HttpApi,
	stack awscdk.Stack,
	vpc awsec2.Vpc,
	dbParams gateway_parameters.DatabaseConnectionParameters,
	authorizer awsapigatewayv2.IHttpRouteAuthorizer,
) {
	// GET /clubs
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
		Integration: integrations.GetClubsByVerification(stack, vpc, dbParams),
	})

	// Get /clubs/{clubId}
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
		Integration: integrations.GetClubById(stack, vpc, dbParams),
	})

	// GET /clubs/{clubId}/events
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_GET,
		},
		Integration: integrations.GetClubEventsByPeriod(stack, vpc, dbParams),
	})

	// POST /clubs/{clubId}/events
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs/{clubId}/events"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
		},
		Integration: integrations.PostNewClubEvent(stack, vpc, dbParams),
		Authorizer:  authorizer,
	})

	// POST /clubs
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path: jsii.String("/clubs"),
		Methods: &[]awsapigatewayv2.HttpMethod{
			awsapigatewayv2.HttpMethod_POST,
		},
		Integration: integrations.PostNewClub(stack, vpc, dbParams),
		Authorizer:  authorizer,
	})
}
