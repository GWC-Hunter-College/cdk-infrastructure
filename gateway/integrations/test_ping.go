package integrations

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

func PingTestIntegration(stack awscdk.Stack) awsapigatewayv2integrations.HttpLambdaIntegration {
	pingFunc := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Ping Function"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("PingTest"),
		Entry:        jsii.String("./lambda/api/test/ping/main.go"),
	})

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("PingLambdaIntegration"),
		pingFunc,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration
}
