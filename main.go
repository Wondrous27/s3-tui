package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Wondrous27/s3-tui/bucket"
	"github.com/Wondrous27/s3-tui/object"
	"github.com/Wondrous27/s3-tui/tui"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	region, ok := os.LookupEnv("AWS_REGION")
	if !ok {
		fmt.Println("AWS_REGION environment variable not set")
		os.Exit(1)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		fmt.Println(err)
	}

	client := s3.NewFromConfig(cfg)
	br := &bucket.S3Repository{Client: client}
	or := &object.S3Repository{Client: client}
	tui.StartTea(br, or)
}
