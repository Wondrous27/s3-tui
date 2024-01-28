package bucket

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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

func (s S3Repository) GetAllBuckets() ([]Bucket, error) {
	s3buckets, err := s.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("could not list buckets: %v", err)
	}

	buckets := make([]Bucket, 0, len(s3buckets.Buckets))
	for _, b := range s3buckets.Buckets {
		buckets = append(buckets, Bucket{
			Name:         *b.Name,
			CreationDate: *b.CreationDate,
		})
	}
	return buckets, nil
}

