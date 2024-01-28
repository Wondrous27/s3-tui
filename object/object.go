package object

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// TODO: group this and bucket.go format together
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

func (s S3Repository) ListObjects(bucketName string) ([]Object, error) {
	out, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{Bucket: &bucketName})
	if err != nil {
		return nil, fmt.Errorf("could not get objects: %v", err)
	}

	objCh := make(chan Object)
	// errCh := make(chan error)

	// Iterate through objects and spawn goroutines for each
	var wg sync.WaitGroup
	for _, obj := range out.Contents {
		wg.Add(1)
		go func(obj types.Object) {
			defer wg.Done()
			// TODO: err handling
			content, _ := s.GetObjectContent(bucketName, *obj.Key)
			object := Object{
				Key:          *obj.Key,
				LastModified: *obj.LastModified,
				Size:         *obj.Size,
				ETag:         *obj.ETag,
				StorageClass: obj.StorageClass,
				Content:      string(content),
			}
			objCh <- object
		}(obj)
	}

	// Close the channel once all goroutines are done
	go func() {
		wg.Wait()
		close(objCh)
		// close(errCh)
	}()

	// Collect results from channels
	objects := make([]Object, 0, len(out.Contents))
	for obj := range objCh {
		objects = append(objects, obj)
	}

	fmt.Println("objects:", objects)
	return objects, nil
}

func (s S3Repository) GetObjectContent(bucket string, key string) ([]byte, error) {
	result, err := s.Client.GetObject(context.TODO(), &s3.GetObjectInput{Bucket: &bucket, Key: &key})
	if err != nil {
		return nil, fmt.Errorf("could not get object: %v", err)
	}
	defer result.Body.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read object: %v", err)
	}
	return body, nil
}

func (s S3Repository) PutObject(file *os.File, bucket string, key string) error {
	_, err := s.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("could not put object %v", err)
	}
	return nil
}
