package model

import (
	"github.com/jinzhu/gorm"
	"github.com/qor/media/media_library"
)

var (
	LoanGuidelineTypes = []string{
		"html", "css", "video", "image",
	}
)

type LoanGuideline struct {
	gorm.Model
	Type            string
	HTMLCode        *string
	CSSCode         *string
	VideoURL        *string
	GuidelineImage  media_library.MediaLibraryStorage    `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
	GuidelineImages []*media_library.MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
	//PaymentAmount  float32
	//Description    string `sql:"size:2000"`
	//TrackingNumber *string
	//File           media_library.MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
	//ShippedAt      *time.Time
	//ReturnedAt     *time.Time
	//CancelledAt    *time.Time
	//DeliveryMethod DeliveryMethod
	//OrderItems     []OrderItem

	//publish2.Version
	//publish2.Schedule
	//publish2.Visible
	//l10n.Locale
	//sorting.SortingDESC
}
