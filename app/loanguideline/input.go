package loanguideline

import (
	"content-management/pkg/errorx"

	"github.com/go-playground/validator/v10"
)

type GetListRequest struct {
	Limit  *int `json:"limit" validate:"required"`
	Offset *int `json:"offset" validate:"required"`
}

func (i *GetListRequest) Validate() error {
	if err := validator.New().Struct(i); err != nil {
		return errorx.ErrInvalidParameter(err)
	}
	return nil
}

type GetRequest struct {
	ID *int `json:"id" validate:"required"`
}

func (i *GetRequest) Validate() error {
	if err := validator.New().Struct(i); err != nil {
		return errorx.ErrInvalidParameter(err)
	}
	return nil
}
