package repository

import "github.com/jinzhu/gorm"

type LoanGuideLineRepository interface {
}

type loanGuideLineRepository struct {
	db *gorm.DB
}

func NewLoanGuideLineRepository(
	db *gorm.DB,
) LoanGuideLineRepository {
	return &loanGuideLineRepository{
		db: db,
	}
}
