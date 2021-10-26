package client

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string `yaml:"endpoint" `
	SecretAccessKey string `yaml:"secret_access_key"    valid:"required"`
	AccessKey       string `yaml:"access_key"    valid:"required"`
	BucketName      string `yaml:"bucket_name"    valid:"required"`
}
type Client struct {
	cfg         Config
	minioClient *minio.Client
}

func New(cfg Config) *Client {
	useSSL := true
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &Client{
		cfg:         cfg,
		minioClient: minioClient,
	}
}

func (c *Client) MakeBucket(ctx context.Context, req *MakeBucketRequest) error {
	return c.minioClient.MakeBucket(ctx, req.Name, minio.MakeBucketOptions{Region: req.Region, ObjectLocking: req.ObjectLocking})
}

func (c *Client) RemoveBucket(ctx context.Context, req *RemoveBucketRequest) error {
	return c.minioClient.RemoveBucket(ctx, req.Name)
}

func (c *Client) PutObject(ctx context.Context, req *PutObjectRequest) (minio.UploadInfo, error) {
	return c.minioClient.PutObject(context.Background(), req.BucketName, req.ObjectName, req.File, 1024, minio.PutObjectOptions{ContentType: "application/octet-stream"})
}
