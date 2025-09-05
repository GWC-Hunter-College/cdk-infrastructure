package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type StubLambdaStackProps struct {
	Props awscdk.StackProps
}

type StubLambdaStack struct {
	Stack awscdk.Stack

	PingFunction awslambda.IFunction

	StubStudentBundle StubStudentFunctions
}

type StubStudentFunctions struct {
	StubStudentsMeClubsGetEndpoint       awslambda.IFunction
	StubStudentsMeClubsEventsGetEndpoint awslambda.IFunction
	StubStudentsMeClubsEboardGetEndpoint awslambda.IFunction
}

func NewStubLambdaStack(scope constructs.Construct, id string, props *StubLambdaStackProps) *StubLambdaStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	//  =======================================
	//  health
	//  =======================================
	// create health check lambda function
	pingFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Stub Health Function"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("StubHealthTest"),
		Entry:        jsii.String("./stub/lambda/health/main.go"),
	})

	//  =======================================
	//  students
	//  =======================================
	studentsMeClubsFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Stub /students/me/club Get EndPoint Funtion"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("StubStudentsMeClubsGetEndpoint"),
		Entry:        jsii.String("./stub/lambda/me/clubs/get.go"),
	})

	studentsMeClubsEventsFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Stub /students/me/club/Events Get EndPoint Funtion"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("StubStudentsMeClubsEventsGetEndpoint"),
		Entry:        jsii.String("./stub/lambda/me/clubs/events/get.go"),
	})

	studentsMeClubsEboardFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Stub /students/me/club/Eboard Get EndPoint Funtion"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("StubStudentsMeClubsEBoardGetEndpoint"),
		Entry:        jsii.String("./stub/lambda/me/clubs/eboard/get.go"),
	})

	studentsBundle := StubStudentFunctions{
		StubStudentsMeClubsGetEndpoint:       studentsMeClubsFunction,
		StubStudentsMeClubsEventsGetEndpoint: studentsMeClubsEventsFunction,
		StubStudentsMeClubsEboardGetEndpoint: studentsMeClubsEboardFunction,
	}

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
	// add routes to HTTP API
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/me/clubs"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("StubStudentsMeClubsGetIntegration"),
			studentsMeClubsFunction,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/me/clubs/events"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("StubStudentsMeClubsGetIntegration"),
			studentsMeClubsEventsFunction,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/me/clubs/eboard"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("StubStudentsMeClubsGetIntegration"),
			studentsMeClubsEboardFunction,
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

	return &StubLambdaStack{
		Stack:        stack,
		PingFunction: pingFunction,

		StubStudentBundle: studentsBundle,
	}
}
