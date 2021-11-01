package service

import (
	"mime/multipart"

	"github.com/qor/media/media_library"
)

type UploadObjectArgs struct {
	Form         *multipart.Form
	MediaStorage *media_library.MediaLibraryStorage
}
