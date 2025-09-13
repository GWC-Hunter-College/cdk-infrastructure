package gateway_helpers

import (
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
)

func GrantRdsAccessToLambda(
	lambdaFn awslambda.Function,
	dbInstance awsrds.IDatabaseInstance,
	dbSecret awssecretsmanager.ISecret,
) {
	dbInstance.GrantConnect(lambdaFn, nil)
	dbSecret.GrantRead(lambdaFn, nil)
}
