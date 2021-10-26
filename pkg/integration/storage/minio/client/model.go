package client

import (
	"mime/multipart"
)

type MakeBucketRequest struct {
	Name          string
	Region        string
	ObjectLocking bool
}

type RemoveBucketRequest struct {
	Name string
}

type PutObjectRequest struct {
	BucketName string
	ObjectName string
	File       multipart.File
}
