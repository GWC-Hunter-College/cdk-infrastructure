package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core

	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"

	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

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

	return &StubLambdaStack{
		Stack:        stack,
		PingFunction: pingFunction,

		StubStudentBundle: studentsBundle,
	}
}
