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
	StubStudentsMeClubsGetEndpoint awslambda.IFunction
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
	studentsMeClubsFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Stub Students Me Club Get EndPoint Funtion"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("StubStudentsMeClubsGetEndpoint"),
		Entry:        jsii.String("./stub/lambda/me/clubs/get.go"),
	})

	studentsBundle := StubStudentFunctions{
		StubStudentsMeClubsGetEndpoint: studentsMeClubsFunction,
	}

	return &StubLambdaStack{
		Stack:        stack,
		PingFunction: pingFunction,

		StubStudentBundle: studentsBundle,
	}
}
