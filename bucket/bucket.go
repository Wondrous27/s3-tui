package bucket

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/charmbracelet/bubbles/list"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)

type S3Repository struct {
	Client *s3.Client
}

type Bucket struct {
	Name         string
	CreationDate time.Time
}

// Implement the `Item` interface
func (b Bucket) Title() string { return b.Name }

func (b Bucket) Description() string {
	return fmt.Sprintf("CreationDate: %v", b.CreationDate.Format(DDMMYYYYhhmmss))
}

func (b Bucket) FilterValue() string { return b.Name }

// TODO: I'm not sure whether I should return []list.Item or []bucket.Bucket
func (s S3Repository) GetAllBuckets() ([]list.Item, error) {
	s3buckets, err := s.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("could not list buckets: %v", err)
	}

	var buckets []list.Item
	for _, b := range s3buckets.Buckets {
		buckets = append(buckets, Bucket{
			Name:         *b.Name,
			CreationDate: *b.CreationDate,
		})
	}
	return buckets, nil
}

func (s S3Repository) CreateBucket(bucketName string) error {
	_, err := s.Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: &bucketName,
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(os.Getenv("AWS_REGION")),
		},
	})
	if err != nil {
		log.Printf("Failed to create bucket %s %v", bucketName, err)
		return fmt.Errorf("could not create bucket %v", err)
	}
	return nil
}

func (s S3Repository) DeleteBucket(bucketName string) error {
	_, err := s.Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{Bucket: &bucketName})
	if err != nil {
		log.Printf("Failed to delete bucket %s %v", bucketName, err)
		return fmt.Errorf("failed to delete %s: %v", bucketName, err)
	}
	return nil
}
