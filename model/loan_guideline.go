package model

import "time"

var (
	LoanGuidelineTypes = []string{
		"Video", "Image",
	}
)

type LoanGuideline struct {
	ID       int `gorm:"primary_key"`
	Type     string
	HTMLCode string
	CSSCode  string
	Medias   Medias

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}
