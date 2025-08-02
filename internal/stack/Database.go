package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DatabaseStackProps struct {
	awscdk.StackProps

	NetworkStackData NetworkStack
}

// type DatabaseStack struct {
// 	Stack awscdk.Stack

//		Vpc                 awsec2.IVpc
//		DbSecurityGroup     awsec2.SecurityGroup
//		DatabaseInformation DatabaseAttributes
//	}
type DatabaseStack struct {
	Stack awscdk.Stack

	DbInstance      awsrds.DatabaseInstance
	DbSecurityGroup awsec2.SecurityGroup
}

type DatabaseAttributes struct {
	DbEndpoint *string
	DbPort     *string
	DbSecret   awssecretsmanager.ISecret

	// passing these for unsafe opening
	DbUser     *string
	DbPassword *string
}

func NewDatabaseStack(scope constructs.Construct, id string, props *DatabaseStackProps) *DatabaseStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	stack.AddDependency(props.NetworkStackData.Stack, jsii.String("Required Network stack"))

	// The code that defines your stack goes here
	vpc := props.NetworkStackData.Vpc
	dbSecurityGroup := props.NetworkStackData.DatabaseSecurityGroup

	dbInstance := awsrds.NewDatabaseInstance(stack, jsii.String("ClubEventDb"), &awsrds.DatabaseInstanceProps{
		// DatabaseName: jsii.String("ClubEventDb"),
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

	awscdk.NewCfnOutput(stack, jsii.String("ImageGalleryDBEndpoint"), &awscdk.CfnOutputProps{
		Value: dbInstance.DbInstanceEndpointAddress(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("ImageGalleryDBSecretArn"), &awscdk.CfnOutputProps{
		Value: dbInstance.Secret().SecretArn(),
	})

	_ = dbInstance

	return &DatabaseStack{
		Stack: stack,

		DbInstance:      dbInstance,
		DbSecurityGroup: dbSecurityGroup,
	}
}
