package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/customresources"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DatabaseStackProps struct {
	Props awscdk.StackProps

	NetworkStackData NetworkStack
}

type DatabaseStack struct {
	Stack awscdk.Stack
}

func NewDatabaseStack(scope constructs.Construct, id string, props *DatabaseStackProps) *DatabaseStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	stack.AddDependency(props.NetworkStackData.Stack, jsii.String("Required Network stack"))

	vpc := props.NetworkStackData.Vpc
	dbSecurityGroup := props.NetworkStackData.DatabaseSecurityGroup

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

	initRDSFunc := awslambda.NewDockerImageFunction(stack, jsii.String("RDS Init Function"),
		&awslambda.DockerImageFunctionProps{
			FunctionName: jsii.String("InitRDS"),
			Description:  jsii.String("Lambda function to initialize RDS database"),
			Code:         awslambda.DockerImageCode_FromImageAsset(jsii.String("lambda/database/init"), nil),
			MemorySize:   jsii.Number(256),
			Architecture: awslambda.Architecture_X86_64(),
			Environment: &map[string]*string{
				"DB_SECRET_ARN": dbInstance.Secret().SecretArn(),
			},
			Vpc: props.NetworkStackData.Vpc,
			SecurityGroups: &[]awsec2.ISecurityGroup{
				props.NetworkStackData.LambdaSecretsManagerSg,
				props.NetworkStackData.LambdaSecurityGroup,
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
	}
}
