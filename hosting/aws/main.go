package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"

	"github.com/GWC-Hunter-College/cdk-infrastructure/hosting/aws/internal/stack"
)

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	stack.NewFrontendStack(app, "FrontendStack", &stack.FrontendStackProps{
		Props: awscdk.StackProps{
			Env:         environment(),
			Description: jsii.String("Stack for the GWC website deployment"),
		},
	})

	stack.NewFrontendHccStack(app, "FrontendHccStack", &stack.FrontendHccStackProps{
		Props: awscdk.StackProps{
			Env:         environment(),
			Description: jsii.String("Stack for the EMS website deployment"),
		},
	})

	app.Synth(nil)
}

// environment leaves the stacks environment-agnostic, matching the reference
// hosting application. Account and region are resolved when the stacks deploy.
func environment() *awscdk.Environment {
	return nil
}
