package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"        // core
	"github.com/aws/aws-cdk-go/awscdk/v2/awssqs" // SQS (demo queue)
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type TutorialStackProps struct {
	awscdk.StackProps
}

func NewTutorial(scope constructs.Construct, id string, props *TutorialStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// The code that defines your stack goes here

	// example resource
	queue := awssqs.NewQueue(stack, jsii.String("CdkInfrastructureQueue"), &awssqs.QueueProps{
		VisibilityTimeout: awscdk.Duration_Seconds(jsii.Number(300)),
	})
	_ = queue // or add Outputs / other logic

	// importedBucket := awss3.Bucket_FromBucketAttributes(stack, jsii.String("MyImportedBucket"), &awss3.BucketAttributes{
	// 	BucketArn:  jsii.String("arn:aws:s3:::cdkinfrastructurestack-myfirstbucketb8884501-9bupkiwgx00v"),
	// 	BucketName: jsii.String("cdkinfrastructurestack-myfirstbucketb8884501-9bupkiwgx00v"),
	// })
	// _ = importedBucket
	// // importedBucket.ApplyRemovalPolicy(awscdk.RemovalPolicy_DESTROY)

	return stack
}
