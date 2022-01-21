package minio

import (
	"content-management/pkg/log"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.elastic.co/apm"
)

type Client struct {
	minioClient *minio.Client
	config      *Config
}

type Config struct {
	Endpoint        string
	SecretAccessKey string
	AccessKey       string
	BucketName      string
}

func New(cfg *Config) *Client {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretAccessKey, ""),
		Secure: true,
	})
	if err != nil {
		log.Fatal(err, nil, nil)
	}

	return &Client{
		minioClient: minioClient,
		config:      cfg,
	}
}

func (c *Client) UploadFile(ctx context.Context, args *UploadObjectArgs) (*UploadObjectResponse, error) {
	trx := apm.TransactionFromContext(ctx)
	minioServiceSpan := trx.StartSpan(fmt.Sprintf("Upload file to Minio server"), "External service.Minio", nil)
	defer minioServiceSpan.End()
	var (
		err                   error
		contentType, fileName string
		fileSize              int64
	)

	if len(args.FileHeaders) > 0 && args.FileHeaders[0] != nil {
		for _, v := range args.FileHeaders {
			fileName = v.Filename
			fileSize = v.Size
			contentTypes := v.Header["Content-Type"]
			if len(contentTypes) > 0 && contentTypes[0] != "" {
				for _, t := range contentTypes {
					contentType = t
				}
			}
		}
	}
	mediaObj := args.MediaStorage
	mediaObj.SelectedType = contentType

	fileURL := c.makeFileURL(fileName)
	mediaObj.Url = fileURL
	file, err := mediaObj.FileHeader.Open()
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	// Upload file to CMC Cloud (Minio)
	_, err = c.minioClient.PutObject(ctx, c.config.BucketName, fileName, file.(io.Reader), fileSize, minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: args.UserMetaData,
	})
	if err != nil {
		return nil, nil
	}
	return &UploadObjectResponse{
		URL: fileURL,
	}, nil
}

func (c *Client) RemoveObject(ctx context.Context, fileName string) error {
	trx := apm.TransactionFromContext(ctx)
	minioServiceSpan := trx.StartSpan(fmt.Sprintf("Remove file in Minio server"), "External service.Minio", nil)
	defer minioServiceSpan.End()
	if err := c.minioClient.RemoveObject(ctx, c.config.BucketName, fileName, minio.RemoveObjectOptions{}); err != nil {
		return nil
	}
	return nil
}

func (c *Client) makeFileURL(fileName string) string {
	return fmt.Sprintf("https://%v/%v/%v", c.config.Endpoint, c.config.BucketName, fileName)
}
