package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"

	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApiStackProps struct {
	Props        awscdk.StackProps
	ImagesBucket awss3.IBucket

	// DatabaseStackData DatabaseStack
	Vpc                               awsec2.Vpc
	LambdaSecretsManagerSecurityGroup awsec2.SecurityGroup
	DbInstance                        awsrds.DatabaseInstance
	ProxyEndpoint                     *string
	LambdaSecurityGroup               awsec2.SecurityGroup
}

func NewApiStack(scope constructs.Construct, id string, props *ApiStackProps) awscdk.Stack {
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
	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("ClubEventApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("ClubEventApi"),
	})

	//  =======================================
	//  Test ping and s3 image storage test
	//  =======================================
	// create ping lambda function
	pingFunc := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Ping Function"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("PingTest"),
		Entry:        jsii.String("./lambda/ping/main.go"),
	})

	// add route to HTTP API
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/pingTest"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("PingLambdaIntegration"),
			pingFunc,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	// create presign lambda function
	presignFunc := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Presign Function"), &awscdklambdagoalpha.GoFunctionProps{
		FunctionName: jsii.String("S3Presign"),
		Entry:        jsii.String("./lambda/presign/main.go"),
	})

	// add route to HTTP API
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/presign"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("PresignOptionsIntegration"),
			presignFunc,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	//  =======================================
	//  Lamnds to rds
	//  =======================================
	// networkStackData := props.DatabaseStackData.NetworkStackData
	vpc := props.Vpc
	lambdaSecretsManagerSecurityGroup := props.LambdaSecretsManagerSecurityGroup

	dbInstance := props.DbInstance
	proxyEndpoint := props.ProxyEndpoint
	lambdaSecurityGroup := props.LambdaSecurityGroup

	dbTestFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("DBTestFunction"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:      jsii.String("lambda/database/test/main.go"), // path to folder with main.go
		MemorySize: jsii.Number(256),
		Timeout:    awscdk.Duration_Seconds(jsii.Number(10)),
		Environment: &map[string]*string{
			"DB_SECRET_ARN": dbInstance.Secret().SecretArn(),
			"DB_HOST":       proxyEndpoint,
		},
		Vpc: vpc,
		SecurityGroups: &[]awsec2.ISecurityGroup{
			lambdaSecretsManagerSecurityGroup,
			lambdaSecurityGroup,
		},
		AllowPublicSubnet: jsii.Bool(true),
	})
	dbInstance.Secret().
		GrantRead(dbTestFunction, nil)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/database/test"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("DBTestIntegration"),
			dbTestFunction,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	// log HTTP API endpoint
	awscdk.NewCfnOutput(stack, jsii.String("myHttpApiEndpoint"), &awscdk.CfnOutputProps{
		Value:       httpApi.ApiEndpoint(),
		Description: jsii.String("HTTP API Endpoint"),
	})

	return stack
}
