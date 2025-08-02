package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

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
	getHandler := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("PingLambda"), &awscdklambdagoalpha.GoFunctionProps{
		Entry: jsii.String("./lambda/ping/main.go"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
		},
	})

	// add route to HTTP API
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/pingTest"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("PingLambdaIntegration"),
			getHandler,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	// create presign lambda function
	presign := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("presign"), &awscdklambdagoalpha.GoFunctionProps{
		Entry: jsii.String("./lambda/presign/main.go"),
		Bundling: &awscdklambdagoalpha.BundlingOptions{
			GoBuildFlags: jsii.Strings(`-ldflags "-s -w"`),
		},
	})

	// add route to HTTP API
	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/presign"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("PresignOptionsIntegration"),
			presign,
			&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
		),
	})

	//  =======================================
	//  Lambda for RDS database initialization
	//  =======================================)

	initRDSFunction := awslambda.NewDockerImageFunction(stack, jsii.String("Database Initialization Function"),
		&awslambda.DockerImageFunctionProps{
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
		GrantRead(initRDSFunction, nil)

	httpApi.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:    jsii.String("/database/init"),
		Methods: &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: awsapigatewayv2integrations.NewHttpLambdaIntegration(
			jsii.String("DBInitFuncIntegration"),
			initRDSFunction,
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
