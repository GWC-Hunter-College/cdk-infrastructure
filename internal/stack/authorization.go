package stack

import (
	"net/url"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AuthorizationStackProps struct {
	Props awscdk.StackProps
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

	userPool := awscognito.NewUserPool(stack, jsii.String("UserPool"), &awscognito.UserPoolProps{
		SelfSignUpEnabled: jsii.Bool(true),
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
		StandardAttributes: &awscognito.StandardAttributes{
			Email: &awscognito.StandardAttribute{
				Required: jsii.Bool(true),
				Mutable:  jsii.Bool(false),
			},
		},
		AccountRecovery: awscognito.AccountRecovery_EMAIL_ONLY,
		PasswordPolicy: &awscognito.PasswordPolicy{
			MinLength:            jsii.Number(7),
			RequireLowercase:     jsii.Bool(false),
			RequireUppercase:     jsii.Bool(false),
			RequireDigits:        jsii.Bool(false),
			RequireSymbols:       jsii.Bool(false),
			TempPasswordValidity: awscdk.Duration_Days(jsii.Number(7)),
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY, // <- forces deletion
	})

	awscdk.NewCfnOutput(stack, jsii.String("UserPoolId"), &awscdk.CfnOutputProps{
		Value: userPool.UserPoolId(),
	})

	//

	domain := userPool.AddDomain(jsii.String("CognitoDomain"), &awscognito.UserPoolDomainOptions{
		CognitoDomain: &awscognito.CognitoDomainOptions{
			DomainPrefix: jsii.String("event-manager-auth"),
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("HostedUIDomainBaseUrl"), &awscdk.CfnOutputProps{
		// e.g., https://gwc-auth-yourprefix.auth.us-east-1.amazoncognito.com
		Value: domain.BaseUrl(&awscognito.BaseUrlOptions{}),
	})

	//

	webClient := userPool.AddClient(jsii.String("WebClient"), &awscognito.UserPoolClientOptions{
		GenerateSecret: jsii.Bool(false), // public client (browser)
		AuthFlows: &awscognito.AuthFlow{
			UserPassword: jsii.Bool(true),
		},
		PreventUserExistenceErrors: jsii.Bool(true),
		OAuth: &awscognito.OAuthSettings{
			Flows: &awscognito.OAuthFlows{
				AuthorizationCodeGrant: jsii.Bool(true),
			},
			Scopes: &[]awscognito.OAuthScope{
				awscognito.OAuthScope_OPENID(),
				awscognito.OAuthScope_EMAIL(),
				awscognito.OAuthScope_PROFILE(),
			},
			CallbackUrls: &[]*string{
				jsii.String("http://localhost:5173/callback"),
			},
			LogoutUrls: &[]*string{
				jsii.String("http://localhost:5173/"),
			},
		},
		SupportedIdentityProviders: &[]awscognito.UserPoolClientIdentityProvider{
			awscognito.UserPoolClientIdentityProvider_COGNITO(),
			// Google will be added in the next step
		},
		// Recommended for browsers
		// DisableOAuthScopesRequired: jsii.Bool(false),
		EnableTokenRevocation: jsii.Bool(true),
		AccessTokenValidity:   awscdk.Duration_Hours(jsii.Number(1)),
		IdTokenValidity:       awscdk.Duration_Hours(jsii.Number(1)),
		RefreshTokenValidity:  awscdk.Duration_Days(jsii.Number(30)),
	})

	awscdk.NewCfnOutput(stack, jsii.String("UserPoolClientId"), &awscdk.CfnOutputProps{
		Value: webClient.UserPoolClientId(),
	})

	//

	awscdk.NewCfnOutput(stack, jsii.String("AuthorizeUrlTemplate"), &awscdk.CfnOutputProps{
		// Helpful to test the redirect immediately
		Value: jsii.String("{{HostedUIDomainBaseUrl}}/oauth2/authorize?client_id={{UserPoolClientId}}&response_type=code&scope=openid+email+profile&redirect_uri=http://localhost:5173/callback"),
	})

	redirectUri := "http://localhost:5173/callback"
	awscdk.NewCfnOutput(stack, jsii.String("AuthorizeUrl"), &awscdk.CfnOutputProps{
		Value: jsii.String(
			*domain.BaseUrl(&awscognito.BaseUrlOptions{}) +
				"/oauth2/authorize" +
				"?client_id=" + *webClient.UserPoolClientId() +
				"&response_type=code" +
				"&scope=openid+email+profile" +
				"&redirect_uri=" + url.QueryEscape(redirectUri),
		),
	})

	return &AuthorizationStack{
		Stack: stack,
	}
}
