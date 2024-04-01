package config

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
)

func NewAws() *s3.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("S3_ID"),         //v.GetString("aws.id"),
			os.Getenv("S3_SECRET_KEY"), // v.GetString("aws.secret"),
			"",
		)),
		config.WithRegion(os.Getenv("S3_REGION")),
	)
	if err != nil {
		log.Fatalf("Failed connect aws: %v", err)
		return nil
	}

	client := s3.NewFromConfig(cfg)
	return client
}
