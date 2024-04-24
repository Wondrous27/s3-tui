package object

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	DDMMYYYYhhmmss = "2006-01-02 15:04:05"
)

type S3Repository struct {
	Client *s3.Client
}

type Object struct {
	Key          string
	LastModified time.Time
	Size         int64
	ETag         string
	StorageClass types.ObjectStorageClass
	Content      string
}

func (o Object) FilterValue() string { return o.Key }

func (o Object) Description() string {
	return fmt.Sprintf("LastModified: %v", o.LastModified.Format(DDMMYYYYhhmmss))
}

func (o Object) Title() string { return o.Key }

func (s S3Repository) ListObjects(bucketName string) ([]string, error) {
	out, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: &bucketName})
	if err != nil {
		return nil, fmt.Errorf("could not get objects: %v", err)
	}

	var keys []string
	for _, obj := range out.Contents {
		keys = append(keys, *obj.Key)
	}

	return keys, nil
}

func (s S3Repository) GetObject(bucket, key string) (*Object, error) {
	result, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	log.Println("last modified: ", *result.LastModified)
	if err != nil {
		return nil, fmt.Errorf("could not get object: %v", err)
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read object: %v", err)
	}
	return &Object{
		Key:          key,
		LastModified: *result.LastModified,
		Size:         *result.ContentLength,
		ETag:         *result.ETag,
		Content:      string(body),
	}, nil
}

func (s S3Repository) PutObject(r io.Reader, bucket string, key string) error {
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   r,
	})
	if err != nil {
		return fmt.Errorf("could not put object %v", err)
	}
	return nil
}
