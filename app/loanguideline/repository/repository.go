package repository

import (
	"content-management/app/loanguideline"
	"content-management/model"
	"content-management/pkg/errorx"
	"context"
	"errors"

	"github.com/jinzhu/gorm"
)

type LoanGuideLineRepository interface {
	CreateLoanGuideline(context.Context, *model.LoanGuideline) error
	GetListLoanGuideline(query *loanguideline.GetListQuery) ([]*model.LoanGuideline, error)
	GetFirstBy(*model.LoanGuideline) (*model.LoanGuideline, error)
}

type loanGuideLineRepository struct {
	db *gorm.DB
}

func (r *loanGuideLineRepository) GetFirstBy(loanGuideline *model.LoanGuideline) (*model.LoanGuideline, error) {
	ct := &model.LoanGuideline{}
	if err := r.db.
		Preload("Medias").
		Where(loanGuideline).
		First(ct).
		Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errorx.ErrDatabase(errors.New("Loan guideline can not found"))
		}
		return nil, errorx.ErrDatabase(err)
	}

	return ct, nil
}

func (r *loanGuideLineRepository) GetListLoanGuideline(query *loanguideline.GetListQuery) ([]*model.LoanGuideline, error) {
	e := []*model.LoanGuideline{}
	if err := r.db.
		Preload("Medias").
		Order("created_at ASC").
		Limit(query.Limit).
		Offset(query.Offset).
		Find(&e).
		Error; err != nil {
		return nil, errorx.ErrDatabase(err)
	}
	return e, nil
}

func NewLoanGuideLineRepository(db *gorm.DB) LoanGuideLineRepository {
	return &loanGuideLineRepository{
		db: db,
	}
}

func (r *loanGuideLineRepository) CreateLoanGuideline(ctx context.Context, loanGuideline *model.LoanGuideline) error {
	if err := r.db.Save(loanGuideline).Error; err != nil {
		return errorx.ErrDatabase(err)
	}
	return nil
}
