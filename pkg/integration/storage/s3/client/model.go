package client

import "mime/multipart"

type UploadFileRequest struct {
	File     multipart.File
	FileName string
}

type UploadFileResponse struct {
	Location  string
	VersionID *string
	UploadID  string
}
