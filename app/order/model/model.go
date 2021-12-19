package model

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/qor/media/media_library"
)

type Order struct {
	gorm.Model
	UserID         uint
	PaymentAmount  float32
	Description    string `sql:"size:2000"`
	TrackingNumber *string
	File           media_library.MediaLibraryStorage `sql:"size:4294967295;" media_library:"url:/system/{{class}}/{{primary_key}}/{{column}}.{{extension}}"`
	ShippedAt      *time.Time
	ReturnedAt     *time.Time
	CancelledAt    *time.Time
	DeliveryMethod DeliveryMethod
	OrderItems     []OrderItem

	//publish2.Version
	//publish2.Schedule
	//publish2.Visible
	//l10n.Locale
	//sorting.SortingDESC
}

type DeliveryMethod struct {
	gorm.Model
	OrderID uint
	Name    string
	Price   float32
}

type PaymentMethod = string

const (
	COD        PaymentMethod = "COD"
	AmazonPay  PaymentMethod = "AmazonPay"
	CreditCard PaymentMethod = "CreditCard"
)

type OrderItem struct {
	gorm.Model
	OrderID         uint
	SizeVariationID uint `cartitem:"SizeVariationID"`
	Quantity        uint `cartitem:"Quantity"`
	Price           float32
	DiscountRate    uint
}

var (
	// DraftState draft state
	DraftState = "draft"
)
