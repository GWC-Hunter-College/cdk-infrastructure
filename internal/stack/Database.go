package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/aws-cdk-go/awscdk/v2/customresources"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DatabaseStackProps struct {
	Props awscdk.StackProps

	NetworkStackData NetworkStack
}

// type DatabaseStack struct {
// 	Stack awscdk.Stack
// }

type DatabaseStack struct {
	Stack awscdk.Stack

	DbInstance       awsrds.DatabaseInstance
	DbSecurityGroup  awsec2.SecurityGroup
	NetworkStackData NetworkStack

	LambdaSecurityGroup awsec2.SecurityGroup
	ProxySecurityGroup  awsec2.SecurityGroup

	ProxyEndpoint *string
}

func NewDatabaseStack(scope constructs.Construct, id string, props *DatabaseStackProps) *DatabaseStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	stack.AddDependency(props.NetworkStackData.Stack, jsii.String("Required Network stack"))

	// ====================================
	// infrasctructure security groups and rules
	// ====================================
	vpc := props.NetworkStackData.Vpc

	// labda to rds security groups
	proxySecurityGroup := createSecurityGroup(stack, vpc, "Proxy")
	lambdaSecurityGroup := createSecurityGroup(stack, vpc, "Lambda")
	dbSecurityGroup := createSecurityGroup(stack, vpc, "RdsDb")

	lambdaSecurityGroup.AddEgressRule(
		proxySecurityGroup,
		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow connections to the proxy"),
		jsii.Bool(false),
	)
	proxySecurityGroup.AddIngressRule(
		lambdaSecurityGroup,
		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow connections from lambda"),
		jsii.Bool(false),
	)

	proxySecurityGroup.AddEgressRule(
		dbSecurityGroup,
		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow connections to the database (RDS)."),
		jsii.Bool(false),
	)
	dbSecurityGroup.AddIngressRule(
		proxySecurityGroup,
		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow connections from the proxy"),
		jsii.Bool(false),
	)

	// ====================================
	// db instance and proxy inilitialization
	// ====================================

	dbInstance := awsrds.NewDatabaseInstance(stack, jsii.String("ClubEventDb"), &awsrds.DatabaseInstanceProps{
		Engine: awsrds.DatabaseInstanceEngine_Mysql(&awsrds.MySqlInstanceEngineProps{
			Version: awsrds.MysqlEngineVersion_VER_8_0_36(),
		}),
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
		Vpc:          vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		PubliclyAccessible:  jsii.Bool(true),
		SecurityGroups:      &[]awsec2.ISecurityGroup{dbSecurityGroup},
		Credentials:         awsrds.Credentials_FromGeneratedSecret(jsii.String("dbadmin"), nil),
		AllocatedStorage:    jsii.Number(20),
		MaxAllocatedStorage: jsii.Number(100),
		BackupRetention:     awscdk.Duration_Days(jsii.Number(7)),
		MultiAz:             jsii.Bool(false),
		RemovalPolicy:       awscdk.RemovalPolicy_DESTROY,
		DeletionProtection:  jsii.Bool(false),
	})

	proxy := awsrds.NewDatabaseProxy(stack, jsii.String("ClubEventProxy"), &awsrds.DatabaseProxyProps{
		ProxyTarget:       awsrds.ProxyTarget_FromInstance(dbInstance),
		Secrets:           &[]awssecretsmanager.ISecret{dbInstance.Secret()},
		Vpc:               vpc,
		RequireTLS:        jsii.Bool(true),
		SecurityGroups:    &[]awsec2.ISecurityGroup{proxySecurityGroup},
		IdleClientTimeout: awscdk.Duration_Minutes(jsii.Number(30)),
	})

	initRDSFunc := awslambda.NewDockerImageFunction(stack, jsii.String("RDS Init Function"),
		&awslambda.DockerImageFunctionProps{
			FunctionName: jsii.String("InitRDS"),
			Description:  jsii.String("Lambda function to initialize RDS database"),
			Code:         awslambda.DockerImageCode_FromImageAsset(jsii.String("lambda/database/init"), nil),
			Timeout:      awscdk.Duration_Minutes(jsii.Number(2)),
			MemorySize:   jsii.Number(256),
			Architecture: awslambda.Architecture_X86_64(),
			Environment: &map[string]*string{
				"DB_SECRET_ARN": dbInstance.Secret().SecretArn(),
				"DB_HOST":       proxy.Endpoint(),
			},
			Vpc: props.NetworkStackData.Vpc,
			SecurityGroups: &[]awsec2.ISecurityGroup{
				props.NetworkStackData.LambdaSecretsManagerSg,
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

	rdsInitializer := awscdk.NewCustomResource(stack, jsii.String("RdsInitializer"), &awscdk.CustomResourceProps{
		ServiceToken: provider.ServiceToken(),
	})

	// Ensure the database is ready before the Lambda runs
	rdsInitializer.Node().AddDependency(dbInstance)

	return &DatabaseStack{
		Stack: stack,

		DbInstance:      dbInstance,
		DbSecurityGroup: dbSecurityGroup,

		LambdaSecurityGroup: lambdaSecurityGroup,
		ProxySecurityGroup:  proxySecurityGroup,

		NetworkStackData: props.NetworkStackData,

		ProxyEndpoint: proxy.Endpoint(),
	}
}
