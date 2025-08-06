package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2" // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type BastionStackProps struct {
	awscdk.StackProps

	// NetworkStackData NetworkStack
}

func NewBastionStack(scope constructs.Construct, id string, props *BastionStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	// stack.AddDependency(props.NetworkStackData.Stack, jsii.String("Required Network stack"))

	// // The code that defines your stack goes here
	// vpc := props.NetworkStackData.Vpc
	// // dbSecurityGroup := props.NetworkStackData.DatabaseSecurityGroup
	// endpointSecurityGroup := props.NetworkStackData.LambdaSecretsManagerSg
	// bastionSecurityGroup := props.NetworkStackData.BastionSecurityGroup

	vpc := awsec2.NewVpc(stack, jsii.String("BastionVpc"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("Private"),
				SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
				CidrMask:   jsii.Number(24),
			},
		},
	})

	endpointSecurityGroup := createSecurityGroup(stack, vpc, "endpointSecurityGroup")
	bastionSecurityGroup := createSecurityGroup(stack, vpc, "bastionSecurityGroup")

	bastionSecurityGroup.AddEgressRule(
		endpointSecurityGroup,
		awsec2.Port_Tcp(jsii.Number(443)),
		jsii.String("Allow https to smm endpoints."),
		jsii.Bool(false))

	endpointSecurityGroup.AddIngressRule(
		bastionSecurityGroup,
		awsec2.Port_Tcp(jsii.Number(443)),
		jsii.String("Allow HTTPS from Bastion SG"),
		jsii.Bool(false))
	// ===========================
	// create vpc endpoints for ssm
	// ===========================
	addEndpoint(
		vpc,
		endpointSecurityGroup,
		*jsii.String("ssm"),
		awsec2.InterfaceVpcEndpointAwsService_SSM(),
	)
	addEndpoint(
		vpc,
		endpointSecurityGroup,
		*jsii.String("ssm-messages"),
		awsec2.InterfaceVpcEndpointAwsService_SSM_MESSAGES(),
	)
	addEndpoint(
		vpc,
		endpointSecurityGroup,
		*jsii.String("ec2-messages"),
		awsec2.InterfaceVpcEndpointAwsService_EC2_MESSAGES(),
	)

	// ===========================
	// create vpc endpoints for ssm
	// ===========================
	ssmPolicy := awsiam.ManagedPolicy_FromAwsManagedPolicyName(aws.String("AmazonSSMManagedInstanceCore"))

	bastionRole := awsiam.NewRole(stack, aws.String("webinstancerole"),
		&awsiam.RoleProps{
			AssumedBy:       awsiam.NewServicePrincipal(aws.String("ec2.amazonaws.com"), nil),
			Description:     aws.String("Bastion Role"),
			ManagedPolicies: &[]awsiam.IManagedPolicy{ssmPolicy},
		},
	)

	linuxImage := awsec2.MachineImage_LatestAmazonLinux(&awsec2.AmazonLinuxImageProps{
		Generation: awsec2.AmazonLinuxGeneration_AMAZON_LINUX_2,
	})

	instance := awsec2.NewInstance(stack, jsii.String("BastionHost"),
		&awsec2.InstanceProps{
			InstanceType:  awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
			MachineImage:  linuxImage,
			Vpc:           vpc,
			InstanceName:  aws.String("monolith"),
			Role:          bastionRole,
			SecurityGroup: bastionSecurityGroup,
			VpcSubnets:    &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED},
		})

	// awsec2.NewBastionHostLinux(stack, jsii.String("BastionHost"), &awsec2.BastionHostLinuxProps{
	// 	Vpc:             vpc,
	// 	SubnetSelection: &awsec2.SubnetSelection{SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED},
	// 	SecurityGroup:   bastionSecurityGroup,
	// 	InstanceType:    awsec2.InstanceType_Of(awsec2.InstanceClass_T3, awsec2.InstanceSize_MICRO),
	// 	Role:            bastionRole,
	// })
	awscdk.NewCfnOutput(stack, jsii.String("InstanceId"), &awscdk.CfnOutputProps{
		Value: instance.InstanceId(),
	})

	return stack
}

func addEndpoint(vpc awsec2.Vpc, securityGroup awsec2.SecurityGroup, id string, svc awsec2.IInterfaceVpcEndpointService) {
	vpc.AddInterfaceEndpoint(jsii.String(id), &awsec2.InterfaceVpcEndpointOptions{
		Service:           svc,
		PrivateDnsEnabled: jsii.Bool(true),
		Open:              jsii.Bool(false),
		SecurityGroups:    &[]awsec2.ISecurityGroup{securityGroup},
		Subnets: &awsec2.SubnetSelection{
			SubnetType: awsec2.SubnetType_PRIVATE_ISOLATED,
		},
	})
}
