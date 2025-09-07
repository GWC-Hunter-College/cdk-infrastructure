package stack

import (
	gateway_routes "cdk-infrastructure/gateway/routes"

	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApiStackProps struct {
	Props awscdk.StackProps

	Vpc                               awsec2.Vpc
	LambdaSecretsManagerSecurityGroup awsec2.SecurityGroup
	DbInstance                        awsrds.DatabaseInstance
	ProxyEndpoint                     *string
	LambdaSecurityGroup               awsec2.SecurityGroup

	ImagesBucket awss3.IBucket
}

func NewApiStack(scope constructs.Construct, id string, props *ApiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("ClubEventApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("ClubEventApi"),
		CorsPreflight: &awsapigatewayv2.CorsPreflightOptions{
			AllowHeaders: &[]*string{
				jsii.String("*"),
			},
			AllowMethods: &[]awsapigatewayv2.CorsHttpMethod{
				awsapigatewayv2.CorsHttpMethod_GET,
				awsapigatewayv2.CorsHttpMethod_POST,
				awsapigatewayv2.CorsHttpMethod_OPTIONS,
				awsapigatewayv2.CorsHttpMethod_PATCH,
			},
			AllowOrigins: &[]*string{
				jsii.String("*"), // allowing from all origins atm, should be locked down later
			},
		},
	})

	gateway_routes.TestRoutes(httpApi, stack)
	gateway_routes.ClubRoutes(httpApi, stack, props.ImagesBucket)
	gateway_routes.EventRoutes(httpApi, stack, props.ImagesBucket)
	gateway_routes.DatabaseRoutes(httpApi, stack, gateway_routes.DatabaseRouteProps{
		Vpc:                               props.Vpc,
		LambdaSecretsManagerSecurityGroup: props.LambdaSecretsManagerSecurityGroup,
		DbInstance:                        props.DbInstance,
		ProxyEndpoint:                     props.ProxyEndpoint,
		LambdaSecurityGroup:               props.LambdaSecurityGroup,
	})

	// log HTTP API endpoint
	awscdk.NewCfnOutput(stack, jsii.String("myHttpApiEndpoint"), &awscdk.CfnOutputProps{
		Value:       httpApi.ApiEndpoint(),
		Description: jsii.String("HTTP API Endpoint"),
	})

	return stack
}
