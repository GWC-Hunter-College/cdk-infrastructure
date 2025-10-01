package integrations

import (
	gateway_helpers "cdk-infrastructure/gateway/helpers"
	gateway_parameters "cdk-infrastructure/gateway/parameters"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

// Integration for GET /clubs?verified=true
func GetStudentMeBySub(
	stack awscdk.Stack,
	vpc awsec2.IVpc,
	dbParams gateway_parameters.DatabaseConnectionParameters,
) awsapigatewayv2integrations.HttpLambdaIntegration {
	host := dbParams.DbHost
	dbName := dbParams.DbName
	arn := dbParams.Secret.SecretArn()
	lambdaToProxySG := dbParams.LambdaToProxySG
	lambdaToSecretsManagerSG := dbParams.LambdaToSecretsManagerSG

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("GetStudentMeFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("GetStudentMe"),
		Description:  jsii.String("Returns the student data associated with the sub from jwt"),
		Entry:        jsii.String("./lambda/api/me/get.go"),
		Environment: &map[string]*string{
			"DB_SECRET_ARN": arn,
			"DB_HOST":       jsii.String(host),
			"DB_NAME":       jsii.String(dbName),
		},
		Vpc:     vpc,
		Timeout: awscdk.Duration_Minutes(jsii.Number(1)),
		SecurityGroups: &[]awsec2.ISecurityGroup{
			lambdaToProxySG,
			lambdaToSecretsManagerSG,
		},
	})

	gateway_helpers.GrantRdsAccessToLambda(function, dbParams.DbInstance, dbParams.Secret)

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("GetStudentMeIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration
}

// Integration for GET /clubs?verified=true
func GetMyClubs(
	stack awscdk.Stack,
	vpc awsec2.IVpc,
	dbParams gateway_parameters.DatabaseConnectionParameters,
) awsapigatewayv2integrations.HttpLambdaIntegration {
	host := dbParams.DbHost
	dbName := dbParams.DbName
	arn := dbParams.Secret.SecretArn()
	lambdaToProxySG := dbParams.LambdaToProxySG
	lambdaToSecretsManagerSG := dbParams.LambdaToSecretsManagerSG

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("GetMyClubsFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("GetMyClubs"),
		Description:  jsii.String("Returns clubs of the student of asoociated jwt"),
		Entry:        jsii.String("lambda/api/me/clubs/get.go"),
		Environment: &map[string]*string{
			"DB_SECRET_ARN": arn,
			"DB_HOST":       jsii.String(host),
			"DB_NAME":       jsii.String(dbName),
		},
		Vpc:     vpc,
		Timeout: awscdk.Duration_Minutes(jsii.Number(1)),
		SecurityGroups: &[]awsec2.ISecurityGroup{
			lambdaToProxySG,
			lambdaToSecretsManagerSG,
		},
	})

	gateway_helpers.GrantRdsAccessToLambda(function, dbParams.DbInstance, dbParams.Secret)

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("GetMyClubsIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration
}
