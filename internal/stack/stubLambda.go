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
}

func NewStubLambdaStack(scope constructs.Construct, id string, props *StubLambdaStackProps) *StubLambdaStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	//  =======================================
	//  Test ping and s3 image storage test
	//  =======================================
	// create ping lambda function
	pingFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Ping Function"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("PingTest"),
		Entry:        jsii.String("./lambda/ping/main.go"),
	})

	return &StubLambdaStack{
		Stack:        stack,
		PingFunction: pingFunction,
	}
}
