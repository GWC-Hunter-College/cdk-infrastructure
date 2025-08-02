package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"

	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// type ApiStackProps struct {
// 	awscdk.StackProps
// 	ImagesBucket awss3.IBucket

//		Vpc                 awsec2.IVpc
//		DbSecurityGroup     awsec2.SecurityGroup
//		DatabaseInformation DatabaseAttributes
//	}
type ApiStackProps struct {
	awscdk.StackProps
	ImagesBucket awss3.IBucket

	NetworkStackData  NetworkStack
	DatabaseStackData DatabaseStack
}

func NewApiStack(scope constructs.Construct, id string, props *ApiStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
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
	pingFunc := awslambda.NewFunction(stack, jsii.String("Ping Function"), &awslambda.FunctionProps{
		FunctionName: jsii.String("APIPing"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./lambda/ping/main.go"), nil),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_GO_1_X(),
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
	presignFunc := awslambda.NewFunction(stack, jsii.String("Presign Function"), &awslambda.FunctionProps{
		FunctionName: jsii.String("S3Presign"),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("./lambda/presign/main.go"), nil),
		Handler:      jsii.String("main"),
		Runtime:      awslambda.Runtime_GO_1_X(),
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
	//  Lambda for RDS database initialization
	//  =======================================)

	initRDSFunc := awslambda.NewDockerImageFunction(stack, jsii.String("RDS Init Function"),
		&awslambda.DockerImageFunctionProps{
			FunctionName: jsii.String("InitRDS"),
			Description:  jsii.String("Lambda function to initialize RDS database"),
			Code:         awslambda.DockerImageCode_FromImageAsset(jsii.String("lambda/database/init"), nil),
			MemorySize:   jsii.Number(256),
			Architecture: awslambda.Architecture_X86_64(),
			Environment: &map[string]*string{
				"DB_SECRET_ARN": props.DatabaseStackData.DbInstance.Secret().SecretArn(),
			},
			Vpc: props.NetworkStackData.Vpc,
			SecurityGroups: &[]awsec2.ISecurityGroup{
				props.NetworkStackData.LambdaSecretsManagerSg,
				props.NetworkStackData.LambdaSecurityGroup,
			},
			AllowPublicSubnet: jsii.Bool(true),
		},
	)

	props.DatabaseStackData.DbInstance.Secret().
		GrantRead(initRDSFunc, nil)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/database/init"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("DBInitFuncIntegration"),
			initRDSFunc,
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
