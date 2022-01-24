package model

import (
	"time"

	"github.com/qor/media/oss"
)

type Media struct {
	ID              int
	URL             string
	Thumbnail       oss.OSS
	Description     string
	LoanGuidelineID *int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

type Medias []*Media
