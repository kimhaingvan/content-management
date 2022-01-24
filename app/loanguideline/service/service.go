package service

import (
	"content-management/app/loanguideline/repository"
	"content-management/thirdparty/minio"
	"context"

	"github.com/jinzhu/gorm"
)

type Dependency struct {
	DB          *gorm.DB
	MinioClient *minio.Client
}

type LoanGuideLineService interface {
	Save(int) int
	UploadFile(ctx context.Context, args *UploadObjectArgs) (string, error)
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

func (l *loanGuideLineService) UploadFile(ctx context.Context, args *UploadObjectArgs) (string, error) {
	// Public url of file
	userMetaData := map[string]string{
		"x-amz-acl": "public-read",
	}
	fileHeaders := args.Form.File["QorResource.File"]
	res, err := l.MinioClient.UploadFile(ctx, &minio.UploadObjectArgs{
		Form:         nil,
		MediaStorage: args.MediaStorage,
		UserMetaData: userMetaData,
		FileHeaders:  fileHeaders,
	})
	if err != nil {
		return "", err
	}
	return res.URL, nil
}

func (l loanGuideLineService) Save(i int) int {
	return i
}
