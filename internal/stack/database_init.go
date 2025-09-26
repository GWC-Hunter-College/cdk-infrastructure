package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/customresources"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type Props struct {
	awscdk.StackProps

	Vpc                               awsec2.Vpc
	DbInstance                        awsrds.DatabaseInstance
	LambdaSecretsManagerSecurityGroup awsec2.SecurityGroup
	LambdaSecurityGroup               awsec2.SecurityGroup
	Proxy                             awsrds.DatabaseProxy
}

func NewDatabaseInitStack(scope constructs.Construct, id string, props *Props) {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	lambdaSecretsManagerSecurityGroup := props.LambdaSecretsManagerSecurityGroup
	lambdaSecurityGroup := props.LambdaSecurityGroup
	dbInstance := props.DbInstance
	proxy := props.Proxy
	vpc := props.Vpc
	proxyEndpoint := proxy.Endpoint()

	initRDSFunc := awslambda.NewDockerImageFunction(stack, jsii.String("RDS Init Function"),
		&awslambda.DockerImageFunctionProps{
			FunctionName: jsii.String("InitRDS"),
			Description:  jsii.String("Lambda function to initialize RDS database"),
			Code:         awslambda.DockerImageCode_FromImageAsset(jsii.String("lambda/internal/database/init"), nil),
			Timeout:      awscdk.Duration_Minutes(jsii.Number(1)),
			MemorySize:   jsii.Number(256),
			Architecture: awslambda.Architecture_X86_64(),
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
		},
	)

	dbInstance.Secret().GrantRead(initRDSFunc, nil)
	dbInstance.GrantConnect(initRDSFunc, nil)

	// Create a custom resource provider to invoke the RDS initialization function on deployment
	provider := customresources.NewProvider(stack, jsii.String("RdsInitProvider"), &customresources.ProviderProps{
		OnEventHandler: initRDSFunc,
	})

	awscdk.NewCustomResource(stack, jsii.String("RdsInitializer"), &awscdk.CustomResourceProps{
		ServiceToken: provider.ServiceToken(),
	})
}
