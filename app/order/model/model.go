package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	gorm.Model
	UserID         *uint
	PaymentAmount  float32
	Description    string `sql:"size:2000"`
	TrackingNumber *string
	ShippedAt      *time.Time
	ReturnedAt     *time.Time
	CancelledAt    *time.Time
}
