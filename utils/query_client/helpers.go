package query_client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func loadSecrets(ctx context.Context, arn string) (*databaseCredentials, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	sm := secretsmanager.NewFromConfig(cfg)

	out, err := sm.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &arn,
	})

	if err != nil {
		return nil, err
	}

	var creds struct {
		User string `json:"username"`
		Pass string `json:"password"`
		Host string `json:"host"`
		Port int    `json:"port"`
	}

	if err = json.Unmarshal([]byte(*out.SecretString), &creds); err != nil {
		return nil, err
	}

	return &databaseCredentials{
		user:     creds.User,
		password: creds.Pass,
		host:     creds.Host,
		port:     fmt.Sprintf("%d", creds.Port),
		address:  fmt.Sprintf("%s:%d", creds.Host, creds.Port),
	}, nil
}

func loadSQLFromFile(filepath string) (string, error) {
	file, err := queriesDir.Open(filepath)
	defer file.Close()

	if err != nil {
		return "", err
	}

	fileBytes := make([]byte, 4096)
	file.Read(fileBytes)

	return string(fileBytes), nil
}

func NewQuery(filepath string, args ...any) Query {
	return Query{
		Filepath: filepath,
		Args:     args,
	}
}
