package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type FrontendStackProps struct {
	Props awscdk.StackProps
}

func NewFrontendStack(scope constructs.Construct, id string, props *FrontendStackProps) awscdk.Stack {
	var stackProps awscdk.StackProps
	if props != nil {
		stackProps = props.Props
	}
	stack := awscdk.NewStack(scope, &id, &stackProps)

	websiteBucket := awss3.NewBucket(stack, jsii.String("GwcWebsiteBucket"), &awss3.BucketProps{
		BucketName:        jsii.String("gwc-club-site"),
		PublicReadAccess:  jsii.Bool(false),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})

	awscdk.NewCfnOutput(stack, jsii.String("websiteBucketName"), &awscdk.CfnOutputProps{
		Value: websiteBucket.BucketName(),
	})

	cloudfrontOAI := awscloudfront.NewOriginAccessIdentity(stack, jsii.String("FrontendOAI"), &awscloudfront.OriginAccessIdentityProps{})
	websiteBucket.GrantRead(cloudfrontOAI, nil)

	cloudfrontStagingBehavior := &awscloudfront.BehaviorOptions{
		Origin: awscloudfrontorigins.NewS3Origin(websiteBucket, &awscloudfrontorigins.S3OriginProps{
			OriginAccessIdentity: cloudfrontOAI,
			OriginPath:           jsii.String("/staging"),
		}),
		ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
	}

	// Keep the reference construct ID stable while serving the staging branch.
	frontendStaging := awscloudfront.NewDistribution(stack, jsii.String("FrontendMain"), &awscloudfront.DistributionProps{
		DefaultRootObject: jsii.String("index.html"),
		DefaultBehavior:   cloudfrontStagingBehavior,
		ErrorResponses: &[]*awscloudfront.ErrorResponse{
			{
				HttpStatus:         jsii.Number(404),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(0)),
			},
			{
				HttpStatus:         jsii.Number(403),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(0)),
			},
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("CloudFront_Main_Info"), &awscdk.CfnOutputProps{
		Description: jsii.String("Staging Branch CloudFront Info"),
		Value: jsii.String("Staging URL: https://" + *frontendStaging.DomainName() +
			" | ID: " + *frontendStaging.DistributionId()),
	})

	cloudfrontProductionBehavior := &awscloudfront.BehaviorOptions{
		Origin: awscloudfrontorigins.NewS3Origin(websiteBucket, &awscloudfrontorigins.S3OriginProps{
			OriginAccessIdentity: cloudfrontOAI,
			OriginPath:           jsii.String("/production"),
		}),
		ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
	}

	frontendProduction := awscloudfront.NewDistribution(stack, jsii.String("FrontendProduction"), &awscloudfront.DistributionProps{
		DefaultRootObject: jsii.String("index.html"),
		DefaultBehavior:   cloudfrontProductionBehavior,
		ErrorResponses: &[]*awscloudfront.ErrorResponse{
			{
				HttpStatus:         jsii.Number(404),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(0)),
			},
			{
				HttpStatus:         jsii.Number(403),
				ResponseHttpStatus: jsii.Number(200),
				ResponsePagePath:   jsii.String("/index.html"),
				Ttl:                awscdk.Duration_Seconds(jsii.Number(0)),
			},
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("CloudFront_Production_Info"), &awscdk.CfnOutputProps{
		Description: jsii.String("Production Branch CloudFront Info"),
		Value: jsii.String("Production URL: https://" + *frontendProduction.DomainName() +
			" | ID: " + *frontendProduction.DistributionId()),
	})

	return stack
}
