package minio

import (
	"mime/multipart"

	"github.com/qor/media/media_library"
)

type UploadObjectResponse struct {
	URL string
}

type UploadObjectArgs struct {
	Form         *multipart.Form
	MediaStorage *media_library.MediaLibraryStorage
	UserMetaData map[string]string
	FileHeaders  []*multipart.FileHeader
}
