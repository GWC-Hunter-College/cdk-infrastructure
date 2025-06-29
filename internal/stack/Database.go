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
}

type DatabaseStack struct {
	Stack awscdk.Stack

	Vpc                 awsec2.IVpc
	DbSecurityGroup     awsec2.SecurityGroup
	DatabaseInformation DatabaseAttributes
}

type DatabaseAttributes struct {
	DbEndpoint *string
	DbPort     *string
	DbSecret   awssecretsmanager.ISecret
}

func NewDatabaseStack(scope constructs.Construct, id string, props *DatabaseStackProps) *DatabaseStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here
	vpc := awsec2.Vpc_FromLookup(stack, jsii.String("DefaultVPC"), &awsec2.VpcLookupOptions{
		IsDefault: jsii.Bool(true),
	})

	dbSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String("DBSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:              vpc,
		AllowAllOutbound: jsii.Bool(true),
	})

	dbSecurityGroup.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow public mysql (only for dev)"),
		nil,
	)

	dbInstance := awsrds.NewDatabaseInstance(stack, jsii.String("ClubEventDb"), &awsrds.DatabaseInstanceProps{
		DatabaseName: jsii.String("ClubEventDb"),
		Engine: awsrds.DatabaseInstanceEngine_Mysql(&awsrds.MySqlInstanceEngineProps{
			Version: awsrds.MysqlEngineVersion_VER_8_0_36(),
		}),
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
		Vpc:          vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		PubliclyAccessible:  jsii.Bool(true), // make this false in production?
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
		Stack:           stack,
		Vpc:             vpc,
		DbSecurityGroup: dbSecurityGroup,
		DatabaseInformation: DatabaseAttributes{
			DbEndpoint: dbInstance.DbInstanceEndpointAddress(),
			DbPort:     dbInstance.DbInstanceEndpointPort(),
			DbSecret:   dbInstance.Secret(),
		},
	}
}

/*
package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/jsii-runtime-go"
	"time"
)

type DataStackProps struct {
	awscdk.StackProps
	HttpAPI awsapigatewayv2.IHttpApi
}

func NewDataStack(scope constructs.Construct, id string, props *DataStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	vpc := awsec2.Vpc_FromLookup(stack, jsii.String("DefaultVPC"), &awsec2.VpcLookupOptions{
		IsDefault: jsii.Bool(true),
	})
	dbSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String("DBSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:             vpc,
		AllowAllOutbound: jsii.Bool(true),
	})
	dbSecurityGroup.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
I		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow public MySql (only for development)"),
	)
	dbInstance := awsrds.NewDatabaseInstance(stack, jsii.String("mageGalleryDB"), &awsrds.DatabaseInstanceProps{
		Engine: awsrds.DatabaseInstanceEngine_Mysql(&awsrds.MySqlInstanceEngineProps{
			Version: awsrds.MysqlEngineVersion_VER_8_0_36(),
		}),
		InstanceType: awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
		Vpc:          vpc,
		VpcSubnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PUBLIC,
		},
		PubliclyAccessible: jsii.Bool(true),
		SecurityGroups:     &[]awsec2.ISecurityGroup{dbSecurityGroup},
		Credentials:        awsrds.Credentials_FromGeneratedSecret(jsii.String("dbadmin"), nil),
		AllocatedStorage:   jsii.Number(20),
		MaxAllocatedStorage: jsii.Number(100),
		BackupRetention:    awscdk.Duration_Days(jsii.Number(7)),
		MultiAz:            jsii.Bool(false),
		RemovalPolicy:      awscdk.RemovalPolicy_DESTROY,
		DeletionProtection: jsii.Bool(false),
	})

	awscdk.NewCfnOutput(stack, jsii.String("ImageGalleryDBEndpoint"), &awscdk.CfnOutputProps{
		Value: dbInstance.DbInstanceEndpointAddress(),
	})

	awscdk.NewCfnOutput(stack, jsii.String("ImageGalleryDBSecretArn"), &awscdk.CfnOutputProps{
		Value: dbInstance.Secret().SecretArn(),
	})

	dbTestFunction := awslambda.NewFunction(stack, jsii.String("DBTestFunction"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_18_X(),
		Handler: jsii.String("index.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("lambdas/db-test"), nil),
		MemorySize: jsii.Number(256),
		Timeout:    awscdk.Duration_Seconds(jsii.Number(10)),
		Environment: &map[string]*string{
			"DB_HOST":     dbInstance.DbInstanceEndpointAddress(),
			"DB_USER":     dbInstance.Secret().SecretValueFromJson(jsii.String("username")).UnsafeUnwrap().ToString(),
			"DB_PASSWORD": dbInstance.Secret().SecretValueFromJson(jsii.String("password")).UnsafeUnwrap().ToString(),
			"DB_NAME":     jsii.String("ImageGalleryDB"),
		},
		Vpc:            vpc,
		SecurityGroups: &[]awsec2.ISecurityGroup{dbSecurityGroup},
	})

	dbTestIntegration := awsapigatewayv2integrations.NewHttpLambdaIntegration(
		jsii.String("DBTestIntegration"),
		dbTestFunction,
		&awsapigatewayv2integrations.HttpLambdaIntegrationProps{},
	)

	props.HttpAPI.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.String("/dbtest"),
		Methods:     &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: dbTestIntegration,
	})

	awscdk.NewCfnOutput(stack, jsii.String("DBTestEndpoint"), &awscdk.CfnOutputProps{
		Value:       jsii.String(*props.HttpAPI.Url() + "dbtest"),
		Description: jsii.String("Invoke URL for /dbtest Lambda -> RDS test"),
	})

	return stack
}
*/
