package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core

	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StubApiStackProps struct {
	Props awscdk.StackProps

	PingFunction awslambda.IFunction

	StubStudentBundle StubStudentFunctions
}

func NewStubApiStack(scope constructs.Construct, id string, props *StubApiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	//  =======================================
	// Api Creation
	//  =======================================
	// create HTTP API
	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("StubbedClubEventApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("StubbedClubEventApi"),
	})

	//  =======================================
	// import health function through props
	//  =======================================
	pingFunction := props.PingFunction

	// add route to HTTP API
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/health"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("StubHealthIntegration"),
			pingFunction,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	//  =======================================
	// import student endpoint functions through props
	//  =======================================
	studentsBundle := props.StubStudentBundle

	studentsMeClubsFunction := studentsBundle.StubStudentsMeClubsGetEndpoint
	// add route to HTTP API
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/me/clubs"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("tubStudentsMeClubsGetIntegration"),
			studentsMeClubsFunction,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	//  =======================================
	// 	prints
	//  =======================================
	// log HTTP API endpoint
	awscdk.NewCfnOutput(stack, jsii.String("myHttpApiEndpoint"), &awscdk.CfnOutputProps{
		Value:       httpApi.ApiEndpoint(),
		Description: jsii.String("HTTP API Endpoint"),
	})
	return stack
}
