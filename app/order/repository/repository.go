package repository

import "github.com/jinzhu/gorm"

type OrderRepository interface {
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(
	db *gorm.DB,
) OrderRepository {
	return &orderRepository{
		db: db,
	}
}