package blob

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// New creates a new blob store using the default configuration.
func New(ctx context.Context, bucketName string) (s Store, err error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return
	}
	return NewFromConfig(cfg, bucketName), nil
}

// NewFromConfig creates a blob store using the provided configuration.
func NewFromConfig(cfg aws.Config, bucketName string, optFns ...func(*s3.Options)) Store {
	return Store{
		client:     s3.NewFromConfig(cfg, optFns...),
		bucketName: bucketName,
	}
}

type Store struct {
	client     *s3.Client
	bucketName string
}

// Put into the bucket.
func (s Store) Put(ctx context.Context, key string, r io.Reader) (err error) {
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
		Body:   r,
	})
	return
}

// Get from the bucket.
func (s Store) Get(ctx context.Context, key string) (r io.ReadCloser, err error) {
	goo, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return goo.Body, err
}

// List keys in the bucket.
func (s Store) List(ctx context.Context, prefix string) (keys []string, err error) {
	pager := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: &s.bucketName,
		Prefix: aws.String(prefix),
	})
	for pager.HasMorePages() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return keys, err
		}
		for _, obj := range page.Contents {
			key := *obj.Key
			keys = append(keys, key)
		}
	}
	return
}
