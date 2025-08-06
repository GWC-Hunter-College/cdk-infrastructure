package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const CIDR = "10.1.0.0/16"

type NetworkStackProps struct {
	Props awscdk.StackProps
}

type NetworkStack struct {
	Stack                  awscdk.Stack
	Vpc                    awsec2.Vpc
	LambdaSecretsManagerSg awsec2.SecurityGroup
	// DatabaseSecurityGroup  awsec2.SecurityGroup
	// LambdaSecurityGroup    awsec2.SecurityGroup
}

func NewNetworkStack(scope constructs.Construct, id string, props *NetworkStackProps) *NetworkStack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	vpc := awsec2.NewVpc(stack, jsii.String("vpc"), &awsec2.VpcProps{
		IpAddresses:                  awsec2.IpAddresses_Cidr(jsii.String(CIDR)),
		MaxAzs:                       jsii.Number(2),
		NatGateways:                  jsii.Number(0),
		RestrictDefaultSecurityGroup: jsii.Bool(true),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("private-subnet-isolated"),
				SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
				CidrMask:   jsii.Number(26),
			},
			{
				Name:       jsii.String("public-subnet"),
				SubnetType: awsec2.SubnetType_PUBLIC,
				CidrMask:   jsii.Number(26),
			},
		},
	})

	secretsManagerVpcEndpointSg := createSecurityGroup(stack, vpc, "secrets-manager-vpc-endpoint")
	lambdaSecretsManagerSg := createSecurityGroup(stack, vpc, "lambda-secrets-manager")

	// took this code from:
	// https://github.com/ProductDock/lambda-vpc-rds-example/blob/main/stacks/network.go
	// We need to configure SecretsManagerInterfaceVPCEndpoint security group to allow port 443 inbound traffic
	// from the Lambda security group and, the Lambda security group to allow port 443 outbound traffic
	// to SecretsManagerInterfaceVPCEndpoint security group.
	secretsManagerVpcEndpointSg.AddIngressRule(
		lambdaSecretsManagerSg,
		awsec2.Port_Tcp(jsii.Number(443)),
		jsii.String("Allow connections from lambda."),
		jsii.Bool(false))

	lambdaSecretsManagerSg.AddEgressRule(
		secretsManagerVpcEndpointSg,
		awsec2.Port_Tcp(jsii.Number(443)),
		jsii.String("Allow connections to SecretsManager VPC endpoint."),
		jsii.Bool(false))

	vpc.AddInterfaceEndpoint(jsii.String("secrets-manager-endpoint"), &awsec2.InterfaceVpcEndpointOptions{
		Service:           awsec2.InterfaceVpcEndpointAwsService_SECRETS_MANAGER(),
		PrivateDnsEnabled: jsii.Bool(true),
		Open:              jsii.Bool(false),
		SecurityGroups:    &[]awsec2.ISecurityGroup{secretsManagerVpcEndpointSg},
	})

	return &NetworkStack{
		Stack: stack,
		Vpc:   vpc,

		LambdaSecretsManagerSg: lambdaSecretsManagerSg,
		// DatabaseSecurityGroup:  dbSecurityGroup,
		// LambdaSecurityGroup:    lambdaSecurityGroup,
	}
}

func createSecurityGroup(stack awscdk.Stack, vpc awsec2.Vpc, name string) awsec2.SecurityGroup {
	return awsec2.NewSecurityGroup(stack, jsii.String(name+"SecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		AllowAllOutbound:  jsii.Bool(false),
		SecurityGroupName: jsii.String(name),
	})
}
