package service

import (
	"content-management/app/loanguideline/repository"
	"content-management/thirdparty/minio"

	"github.com/jinzhu/gorm"
)

type Dependency struct {
	DB          *gorm.DB
	MinioClient *minio.Client
}

type LoanGuideLineService interface {
}

type loanGuideLineService struct {
	*Dependency
	loanGuideLineRepo repository.LoanGuideLineRepository
}

func NewLoanGuideLineService(d *Dependency) LoanGuideLineService {
	return &loanGuideLineService{
		Dependency:        d,
		loanGuideLineRepo: repository.NewLoanGuideLineRepository(d.DB),
	}
}
