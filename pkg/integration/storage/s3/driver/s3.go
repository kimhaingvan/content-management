package driver

import (
	"content-management/pkg/integration/storage/s3/client"
	"context"
	"mime/multipart"
)

type S3Driver struct {
	s3Client *client.Client
}

func New(cfg client.Config) *S3Driver {
	c := S3Driver{
		s3Client: client.New(cfg),
	}
	return &c
}

func (d *S3Driver) UploadFile(ctx context.Context, args *UploadFileArgs) (*UploadFileResponse, error) {
	req := &client.UploadFileRequest{
		File:     args.File,
		FileName: args.FileName,
	}
	res, err := d.s3Client.UploadFile(ctx, req)
	if err != nil {
		return nil, err
	}
	return &UploadFileResponse{
		Location:  res.Location,
		VersionID: res.VersionID,
		UploadID:  res.UploadID,
	}, nil
}

type UploadFileArgs struct {
	File     multipart.File
	FileName string
}

type UploadFileResponse struct {
	Location  string
	VersionID *string
	UploadID  string
}
