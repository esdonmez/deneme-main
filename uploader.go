package pgo

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

type Uploader interface {
	Create(ctx context.Context, bucketName string, location string)
	Upload(ctx context.Context, bucketName string, objectName string, buf []byte)
	UploadFile(ctx context.Context, bucketName string, objectName string, filePath string)
	Download(ctx context.Context, bucketName string, objectName string, filePath string)
	DownloadAll(ctx context.Context, bucketName string) []string
}

type uploader struct {
	client *minio.Client
}

func NewUploader(endpoint string, username string, password string, token string) Uploader {
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, token),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &uploader{client: minioClient}
}

func (c *uploader) Create(ctx context.Context, bucketName string, location string) {
	err := c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := c.client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func (c *uploader) UploadFile(ctx context.Context, bucketName string, objectName string, filePath string) {
	info, err := c.client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

func (c *uploader) Upload(ctx context.Context, bucketName string, objectName string, buf []byte) {
	info, err := c.client.PutObject(ctx, bucketName, objectName, bytes.NewReader(buf), int64(len(buf)), minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
}

func (c *uploader) Download(ctx context.Context, bucketName string, objectName string, filePath string) {
	err := c.client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
}

func (c *uploader) DownloadAll(ctx context.Context, bucketName string) []string {
	objectsChan := c.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})
	objectNames := make([]string, 0, len(objectsChan))
	for object := range objectsChan {
		objectName := object.Key
		objectNames = append(objectNames, objectName)
		filePath := "./" + objectName + ".pprof"
		err := c.client.FGetObject(ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
		if err != nil {
			log.Fatalln(err)
		}
	}

	return objectNames
}
