package integrations

import (
	gateway_parameters "cdk-infrastructure/gateway/parameters"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/jsii-runtime-go"
)

// Integrations for event image presign endpoint
func GetEventImagesIntegration(
	stack awscdk.Stack,
	vpc awsec2.IVpc,
	s3Params gateway_parameters.S3PermissionsParameters,
	dbParams gateway_parameters.DatabaseConnectionParameters,
) (awsapigatewayv2integrations.HttpLambdaIntegration, awscdklambdagoalpha.GoFunction) {
	bucket := s3Params.Bucket

	dbSecretArn := dbParams.Secret.SecretArn()
	dbHost := jsii.String(dbParams.DbHost)
	dbName := jsii.String(dbParams.DbName)

	lambdaToProxySG := dbParams.LambdaToProxySG
	lambdaToSecretsSG := dbParams.LambdaToSecretsManagerSG

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("GetEventImagesPresignFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("GetClubEventImagesPresign"),
		Entry:        jsii.String("lambda/api/clubs/events/images/get/get.go"),
		Environment: &map[string]*string{
			"S3_BUCKET":     bucket.BucketName(),
			"DB_SECRET_ARN": dbSecretArn,
			"DB_HOST":       dbHost,
			"DB_NAME":       dbName,
		},
		Vpc: vpc,
		SecurityGroups: &[]awsec2.ISecurityGroup{
			lambdaToProxySG,
			lambdaToSecretsSG,
		},
		Timeout: awscdk.Duration_Minutes(jsii.Number(1)),
	})

	bucket.GrantRead(function, "events/*")

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("GetEventImagePresignIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration, function
}

func PostEventImagesIntegration(stack awscdk.Stack, vpc awsec2.IVpc, s3Params gateway_parameters.S3PermissionsParameters) (awsapigatewayv2integrations.HttpLambdaIntegration, awscdklambdagoalpha.GoFunction) {
	bucket := s3Params.Bucket

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("PostEventImagesPresignFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("PostClubEventImagesPresign"),
		Entry:        jsii.String("lambda/api/clubs/events/images/post/post.go"),
		Environment: &map[string]*string{
			"S3_BUCKET": bucket.BucketName(),
		},
		Vpc: vpc,
	})

	bucket.GrantPut(function, "events/*")

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("PostEventImagePresignIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration, function
}

func ConfirmEventImagesIntegration(
	stack awscdk.Stack,
	vpc awsec2.IVpc,
	dbParams gateway_parameters.DatabaseConnectionParameters,
) (awsapigatewayv2integrations.HttpLambdaIntegration, awscdklambdagoalpha.GoFunction) {
	dbSecretArn := dbParams.Secret.SecretArn()
	dbHost := jsii.String(dbParams.DbHost)
	dbName := jsii.String(dbParams.DbName)

	lambdaToProxySG := dbParams.LambdaToProxySG
	lambdaToSecretsSG := dbParams.LambdaToSecretsManagerSG

	function := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("ConfirmEventImagesFunction"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("ConfirmClubEventImages"),
		Entry:        jsii.String("lambda/api/clubs/events/images/confirm/post.go"),
		Environment: &map[string]*string{
			"DB_SECRET_ARN": dbSecretArn,
			"DB_HOST":       dbHost,
			"DB_NAME":       dbName,
		},
		Vpc: vpc,
		SecurityGroups: &[]awsec2.ISecurityGroup{
			lambdaToProxySG,
			lambdaToSecretsSG,
		},
		Timeout: awscdk.Duration_Minutes(jsii.Number(1)),
	})

	integration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("ConfirmEventImageIntegration"),
		function,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	return integration, function
}
