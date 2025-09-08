package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"

	cr "github.com/aws/aws-cdk-go/awscdk/v2/customresources"
)

type AuthorizationStackProps struct {
	Props awscdk.StackProps

	Vpc                               awsec2.Vpc
	LambdaSecretsManagerSecurityGroup awsec2.SecurityGroup
	DbInstance                        awsrds.DatabaseInstance
	ProxyEndpoint                     *string
	LambdaSecurityGroup               awsec2.SecurityGroup

	UserPool awscognito.IUserPool
}

type AuthorizationStack struct {
	Stack awscdk.Stack
}

func NewAuthorizationStack(scope constructs.Construct, id string, props *AuthorizationStackProps) *AuthorizationStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	//  =======================================
	//  read props
	//  =======================================
	vpc := props.Vpc
	dbInstance := props.DbInstance
	proxyEndpoint := props.ProxyEndpoint

	lambdaSecretsManagerSecurityGroup := props.LambdaSecretsManagerSecurityGroup
	lambdaSecurityGroup := props.LambdaSecurityGroup

	//  =======================================
	//  authenticaion lambda
	//  =======================================

	postConfirmFunction := awslambda.NewDockerImageFunction(stack, jsii.String("PostConfirmUserUpsertFunction"),
		&awslambda.DockerImageFunctionProps{
			FunctionName: jsii.String("PostConfirmUserUpsert"),
			Description:  jsii.String("Lambda function to initialize RDS database"),
			Code:         awslambda.DockerImageCode_FromImageAsset(jsii.String("lambda/database/init"), nil),
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

	dbInstance.Secret().
		GrantRead(postConfirmFunction, nil)

		// allow Cognito to invoke your Lambda
	postConfirmFunction.AddPermission(jsii.String("AllowCognitoInvoke"), &awslambda.Permission{
		Principal: awsiam.NewServicePrincipal(jsii.String("cognito-idp.amazonaws.com"), nil),
		SourceArn: props.UserPool.UserPoolArn(),
	})

	physID := "Wire-" + *props.UserPool.UserPoolId()

	onUp := &cr.AwsSdkCall{
		Service: jsii.String("CognitoIdentityServiceProvider"), // <<< v2 name
		Action:  jsii.String("updateUserPool"),
		Parameters: &map[string]interface{}{
			"UserPoolId": *props.UserPool.UserPoolId(),
			"LambdaConfig": map[string]interface{}{
				"PostConfirmation": postConfirmFunction.FunctionArn(),
				// include if you want every login too:
				"PostAuthentication": postConfirmFunction.FunctionArn(),
			},
		},
		PhysicalResourceId: cr.PhysicalResourceId_Of(jsii.String(physID)),
	}

	onDel := &cr.AwsSdkCall{
		Service: jsii.String("CognitoIdentityServiceProvider"), // <<< v2 name
		Action:  jsii.String("updateUserPool"),
		Parameters: &map[string]interface{}{
			"UserPoolId":   *props.UserPool.UserPoolId(),
			"LambdaConfig": map[string]interface{}{}, // clears triggers
		},
		PhysicalResourceId:       cr.PhysicalResourceId_Of(jsii.String(physID)),
		IgnoreErrorCodesMatching: jsii.String(".*ResourceNotFound.*"),
	}

	// Scope permissions if you want (instead of ANY_RESOURCE)
	policy := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   &[]*string{jsii.String("cognito-idp:UpdateUserPool")},
		Resources: &[]*string{props.UserPool.UserPoolArn()},
	})

	cr.NewAwsCustomResource(stack, jsii.String("WireCognitoTriggers"), &cr.AwsCustomResourceProps{
		Policy:              cr.AwsCustomResourcePolicy_FromStatements(&[]awsiam.PolicyStatement{policy}),
		OnCreate:            onUp,
		OnUpdate:            onUp,
		OnDelete:            onDel,
		InstallLatestAwsSdk: jsii.Bool(false), // use SDK v2
	})

	return &AuthorizationStack{
		Stack: stack,
	}
}
