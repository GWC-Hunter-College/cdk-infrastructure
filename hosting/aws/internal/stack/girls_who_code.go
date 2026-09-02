package stack

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfront"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscloudfrontorigins"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
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

	girlsWhoCodeAccount := awscdk.Stack_Of(girlsWhoCodeHostingStack).Account()
	girlsWhoCodeStagingDistributionARN := awscdk.Arn_Format(&awscdk.ArnComponents{
		Service:      jsii.String("cloudfront"),
		Account:      girlsWhoCodeAccount,
		Resource:     jsii.String("distribution"),
		ResourceName: girlsWhoCodeStagingDistribution.DistributionId(),
		Region:       jsii.String(""),
	}, girlsWhoCodeHostingStack)
	girlsWhoCodeProductionDistributionARN := awscdk.Arn_Format(&awscdk.ArnComponents{
		Service:      jsii.String("cloudfront"),
		Account:      girlsWhoCodeAccount,
		Resource:     jsii.String("distribution"),
		ResourceName: girlsWhoCodeProductionDistribution.DistributionId(),
		Region:       jsii.String(""),
	}, girlsWhoCodeHostingStack)

	girlsWhoCodeS3ObjectsStatement := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("s3:PutObject", "s3:DeleteObject"),
		Resources: jsii.Strings(*girlsWhoCodeWebsiteBucket.ArnForObjects(jsii.String("*"))),
	})
	girlsWhoCodeS3ListStatement := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("s3:ListBucket"),
		Resources: jsii.Strings(*girlsWhoCodeWebsiteBucket.BucketArn()),
	})
	girlsWhoCodeCloudFrontInvalidateStatement := awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Effect:    awsiam.Effect_ALLOW,
		Actions:   jsii.Strings("cloudfront:CreateInvalidation"),
		Resources: jsii.Strings(*girlsWhoCodeStagingDistributionARN, *girlsWhoCodeProductionDistributionARN),
	})

	girlsWhoCodeCIPolicy := awsiam.NewPolicy(girlsWhoCodeHostingStack, jsii.String("GirlsWhoCodeCiPolicy"), &awsiam.PolicyProps{
		PolicyName: jsii.String("frontend-gwc-ci-policy"),
		Statements: &[]awsiam.PolicyStatement{
			girlsWhoCodeS3ObjectsStatement,
			girlsWhoCodeS3ListStatement,
			girlsWhoCodeCloudFrontInvalidateStatement,
		},
	})

	girlsWhoCodeCIUser := awsiam.NewUser(girlsWhoCodeHostingStack, jsii.String("GirlsWhoCodeCiUser"), &awsiam.UserProps{
		UserName: jsii.String("gwc-website-ci-deployer"),
	})
	girlsWhoCodeCIPolicy.AttachToUser(girlsWhoCodeCIUser)

	awscdk.NewCfnOutput(girlsWhoCodeHostingStack, jsii.String("GirlsWhoCodeS3StagingDestination"), &awscdk.CfnOutputProps{
		Description: jsii.String("Girls Who Code staging S3 deployment destination"),
		Value:       jsii.String("s3://" + *girlsWhoCodeWebsiteBucket.BucketName() + "/staging"),
	})
	awscdk.NewCfnOutput(girlsWhoCodeHostingStack, jsii.String("GirlsWhoCodeS3ProductionDestination"), &awscdk.CfnOutputProps{
		Description: jsii.String("Girls Who Code production S3 deployment destination"),
		Value:       jsii.String("s3://" + *girlsWhoCodeWebsiteBucket.BucketName() + "/production"),
	})

	return girlsWhoCodeHostingStack
}
