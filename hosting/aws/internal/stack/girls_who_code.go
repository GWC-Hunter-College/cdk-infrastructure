package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type GirlsWhoCodeHostingStackProps struct {
	Props awscdk.StackProps
}

func NewGirlsWhoCodeHostingStack(scope constructs.Construct, id string, props *GirlsWhoCodeHostingStackProps) awscdk.Stack {
	var stackProps awscdk.StackProps
	if props != nil {
		stackProps = props.Props
	}
	girlsWhoCodeHostingStack := awscdk.NewStack(scope, &id, &stackProps)

	girlsWhoCodeWebsiteBucket := awss3.NewBucket(girlsWhoCodeHostingStack, jsii.String("GwcWebsiteBucket"), &awss3.BucketProps{
		BucketName:        jsii.String("gwc-club-site"),
		PublicReadAccess:  jsii.Bool(false),
		RemovalPolicy:     awscdk.RemovalPolicy_DESTROY,
		AutoDeleteObjects: jsii.Bool(true),
	})

	awscdk.NewCfnOutput(girlsWhoCodeHostingStack, jsii.String("websiteBucketName"), &awscdk.CfnOutputProps{
		Value: girlsWhoCodeWebsiteBucket.BucketName(),
	})

	girlsWhoCodeOAI := awscloudfront.NewOriginAccessIdentity(girlsWhoCodeHostingStack, jsii.String("FrontendOAI"), &awscloudfront.OriginAccessIdentityProps{})
	girlsWhoCodeWebsiteBucket.GrantRead(girlsWhoCodeOAI, nil)

	girlsWhoCodeStagingBehavior := &awscloudfront.BehaviorOptions{
		Origin: awscloudfrontorigins.NewS3Origin(girlsWhoCodeWebsiteBucket, &awscloudfrontorigins.S3OriginProps{
			OriginAccessIdentity: girlsWhoCodeOAI,
			OriginPath:           jsii.String("/staging"),
		}),
		ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
	}

	// Keep the historical FrontendMain construct ID stable while serving the staging branch.
	girlsWhoCodeStagingDistribution := awscloudfront.NewDistribution(girlsWhoCodeHostingStack, jsii.String("FrontendMain"), &awscloudfront.DistributionProps{
		DefaultRootObject: jsii.String("index.html"),
		DefaultBehavior:   girlsWhoCodeStagingBehavior,
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

	awscdk.NewCfnOutput(girlsWhoCodeHostingStack, jsii.String("CloudFront_Main_Info"), &awscdk.CfnOutputProps{
		Description: jsii.String("Staging Branch CloudFront Info"),
		Value: jsii.String("Staging URL: https://" + *girlsWhoCodeStagingDistribution.DomainName() +
			" | ID: " + *girlsWhoCodeStagingDistribution.DistributionId()),
	})

	girlsWhoCodeProductionBehavior := &awscloudfront.BehaviorOptions{
		Origin: awscloudfrontorigins.NewS3Origin(girlsWhoCodeWebsiteBucket, &awscloudfrontorigins.S3OriginProps{
			OriginAccessIdentity: girlsWhoCodeOAI,
			OriginPath:           jsii.String("/production"),
		}),
		ViewerProtocolPolicy: awscloudfront.ViewerProtocolPolicy_REDIRECT_TO_HTTPS,
	}

	girlsWhoCodeProductionDistribution := awscloudfront.NewDistribution(girlsWhoCodeHostingStack, jsii.String("FrontendProduction"), &awscloudfront.DistributionProps{
		DefaultRootObject: jsii.String("index.html"),
		DefaultBehavior:   girlsWhoCodeProductionBehavior,
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

	awscdk.NewCfnOutput(girlsWhoCodeHostingStack, jsii.String("CloudFront_Production_Info"), &awscdk.CfnOutputProps{
		Description: jsii.String("Production Branch CloudFront Info"),
		Value: jsii.String("Production URL: https://" + *girlsWhoCodeProductionDistribution.DomainName() +
			" | ID: " + *girlsWhoCodeProductionDistribution.DistributionId()),
	})

	return girlsWhoCodeHostingStack
}
