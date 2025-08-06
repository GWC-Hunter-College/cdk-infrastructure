package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2" // core

	"github.com/aws/jsii-runtime-go"

	stack "cdk-infrastructure/internal/stack"
)

// type FrontendStackProps struct {
// 	awscdk.StackProps
// }

// type TutorialStackProps struct {
// 	awscdk.StackProps
// }

// type ApiStackProps struct {
// 	awscdk.StackProps
// 	ImagesBucket awss3.IBucket
// }

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	stack.NewFrontendStack(app, "FrontendStack", &stack.FrontendStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	images := stack.NewStorageStack(app, "StorageStack", &stack.StorageStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	network := stack.NewNetworkStack(app, "NetworkStack", &stack.NetworkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	database := stack.NewDatabaseStack(app, "DatabaseStack", &stack.DatabaseStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		NetworkStackData: *network,
	})

	stack.NewApiStack(app, "ApiStack", &stack.ApiStackProps{
		StackProps: awscdk.StackProps{
			Env: env(),
		},
		ImagesBucket: images.Bucket,

		// NetworkStackData:  *network,
		DatabaseStackData: *database,
	})

	stack.NewBastionStack(app, "BastionStack", &stack.BastionStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
