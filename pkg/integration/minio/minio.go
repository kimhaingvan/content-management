package minio

import (
	"content-management/pkg/errorx"
	"content-management/pkg/log"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/qor/media/media_library"
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

func (c *Client) UploadFileToMinioFromRequest(ctx context.Context, request *http.Request, mediaFile media_library.MediaLibraryStorage) (*UploadObjectResponse, error) {
	trx := apm.TransactionFromContext(ctx)
	minioServiceSpan := trx.StartSpan(fmt.Sprintf("Upload file to Minio server"), "External service.Minio", nil)
	defer minioServiceSpan.End()
	var (
		err                   error
		contentType, fileName string
		fileSize              int64
		file                  multipart.File
	)
	fileHeaders := request.MultipartForm.File["QorResource.File"]
	if len(fileHeaders) > 0 && fileHeaders[0] != nil {
		fileHeader := fileHeaders[0]
		fileName = fileHeader.Filename
		fileSize = fileHeader.Size
		contentTypes := fileHeader.Header["Content-Type"]
		if len(contentTypes) > 0 && contentTypes[0] != "" {
			contentType = contentTypes[0]
		}
	}

	mediaFile.SelectedType = contentType

	// Public url of file
	userMetaData := map[string]string{
		"x-amz-acl": "public-read",
	}

	fileURL := c.makeFileURL(fileName)
	mediaFile.Url = fileURL

	file, err = mediaFile.FileHeader.Open()
	if err != nil {
		return nil, errorx.Errorf(http.StatusInternalServerError, err, "Can not open file")
	}
	defer file.Close()

	// Upload file to CMC Cloud (Minio)
	_, err = c.minioClient.PutObject(ctx, c.config.BucketName, fileName, file.(io.Reader), fileSize, minio.PutObjectOptions{
		ContentType:  contentType,
		UserMetadata: userMetaData,
	})
	if err != nil {
		return nil, errorx.Errorf(http.StatusInternalServerError, err, "Can not upload file to minio server")
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
		return errorx.Errorf(http.StatusInternalServerError, err, "Can not remove object in minio server")
	}
	return nil
}

func (c *Client) makeFileURL(fileName string) string {
	return fmt.Sprintf("https://%v/%v/%v", c.config.Endpoint, c.config.BucketName, fileName)
}
