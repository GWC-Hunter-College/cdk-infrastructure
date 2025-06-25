package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	s3Client *s3.Client
	bucket   string
)

func init() {
	cfg, _ := config.LoadDefaultConfig(context.Background())
	s3Client = s3.NewFromConfig(cfg)
}

type res struct {
	UploadURL string `json:"uploadUrl"`
	Key       string `json:"key"`
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fileName := request.QueryStringParameters["fileName"]
	fileType := request.QueryStringParameters["fileType"]
	if fileName == "" || fileType == "" {
		return events.APIGatewayProxyResponse{Body: string(`message: "Missing fileName or fileType"`), StatusCode: 400}, nil
	}

	// dont know why i have this since in my ts i basically didnt use this
	key := fmt.Sprintf("uploads/%d-%s", time.Now().UnixMilli(), fileName)

	command := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(fileType),
	}

	presignClient := s3.NewPresignClient(s3Client)

	url, err := presignClient.PresignPutObject(ctx, command, func(o *s3.PresignOptions) {
		o.Expires = time.Minute
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	out, _ := json.Marshal(res{UploadURL: url.URL, Key: key})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       string(out),
	}, nil
}

func main() {
	lambda.Start(handleRequest)
}
