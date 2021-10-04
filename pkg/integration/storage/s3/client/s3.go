package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	AwsS3Region        string `yaml:"aws_s3_region"    valid:"required"`
	AwsS3Bucket        string `yaml:"aws_s3_bucket"    valid:"required"`
	AwsAccessKey       string `yaml:"aws_access_key"    valid:"required"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"    valid:"required"`
	AwsSessionToken    string `yaml:"aws_session_token"`
}

type Client struct {
	cfg             Config
	awsS3Client     *s3.Client
	awsS3Uploader   *manager.Uploader
	awsS3Downloader *manager.Downloader
}

func New(cfg Config) *Client {
	awsS3Client := s3.New(s3.Options{
		Region:      cfg.AwsS3Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AwsAccessKey, cfg.AwsSecretAccessKey, cfg.AwsSessionToken)),
	})

	uploader := manager.NewUploader(awsS3Client)
	downloader := manager.NewDownloader(awsS3Client)
	return &Client{
		cfg:             cfg,
		awsS3Client:     awsS3Client,
		awsS3Uploader:   uploader,
		awsS3Downloader: downloader,
	}
}

func (c *Client) UploadFile(ctx context.Context, req *UploadFileRequest) (*UploadFileResponse, error) {
	res, err := c.awsS3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(c.cfg.AwsS3Bucket),
		Key:    aws.String(req.FileName),
		Body:   req.File,
	})
	if err != nil {
		return nil, err
	}

	return &UploadFileResponse{
		Location:  res.Location,
		VersionID: res.VersionID,
		UploadID:  res.UploadID,
	}, nil
}

func (c *Client) DownloadFile(ctx context.Context, req *UploadFileRequest) {

}
