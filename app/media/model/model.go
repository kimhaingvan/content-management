package model

import (
	"github.com/qor/media/media_library"
)

type CustomMedia struct {
	media_library.MediaLibrary
	URL string `json:"url"`
}
