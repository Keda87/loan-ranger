package files

import (
	"bytes"
	"context"
	"errors"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type BucketInterface interface {
	Upload(ctx context.Context, key string, reader io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
	GetSignURL(ctx context.Context, key string) (string, error)
}

type S3Client struct {
	bucketName string
	client     *s3.S3
}

func NewS3Storage(access, secret, bucketName, endpoint, region string) S3Client {
	var (
		s3conf = &aws.Config{
			Credentials:      credentials.NewStaticCredentials(access, secret, ""),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String(region),
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: aws.Bool(true),
		}
	)

	s3session, err := session.NewSession(s3conf)
	if err != nil {
		panic(err)
	}
	s3Client := s3.New(s3session)

	return S3Client{
		client:     s3Client,
		bucketName: bucketName,
	}
}

func (s S3Client) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	buff, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	out, err := s.client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(buff),
		ACL:    aws.String("private"),
	})
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func (s S3Client) Delete(ctx context.Context, key string) error {
	return errors.New("not implemented")
}

func (s S3Client) GetSignURL(ctx context.Context, key string) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})

	signedURL, err := req.Presign(time.Minute * 10)
	if err != nil {
		return "", err
	}

	return signedURL, nil
}
