package handlers

import (
	"context"
	"database/sql"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Handler struct {
	DbConn     *sql.DB
	S3Client   *s3.Client
	BucketName string
	Region     string
}

func (h Handler) CreateBucket(ctx context.Context) error {
	if h.bucketExists(ctx) {
		return nil
	}
	input := &s3.CreateBucketInput{
		Bucket: aws.String(h.BucketName),
	}
	input.CreateBucketConfiguration = &types.CreateBucketConfiguration{
		LocationConstraint: types.BucketLocationConstraint(h.Region),
	}
	_, err := h.S3Client.CreateBucket(ctx, input)

	return err
}

func (h Handler) DeleteBucket(ctx context.Context) error {
	_, err := h.S3Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
		Bucket: aws.String(h.BucketName),
	})

	return err
}

func (h Handler) bucketExists(ctx context.Context) bool {
	_, err := h.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(h.BucketName),
	})
	if err != nil {
		return false
	}
	return true
}
