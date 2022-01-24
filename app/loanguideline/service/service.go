package service

import (
	"content-management/app/loanguideline"
	"content-management/app/loanguideline/repository"
	"content-management/model"
	"content-management/pkg/errorx"
	"content-management/thirdparty/minio"
	"context"
	"errors"
	"mime/multipart"

	"github.com/qor/media/media_library"

	"github.com/jinzhu/gorm"
)

type LoanGuideLineService interface {
	UploadFile(ctx context.Context, args *UploadObjectArgs) (*UploadObjectResponse, error)
	CreateLoanGuideline(ctx context.Context, args *model.LoanGuideline) error
	GetListLoanGuideline(ctx context.Context, args *loanguideline.GetListRequest) (*loanguideline.GetListResponse, error)
	GetLoanGuideline(ctx context.Context, args *loanguideline.GetRequest) (*loanguideline.LoanGuideline, error)
}

type Dependency struct {
	DB          *gorm.DB
	MinioClient *minio.Client
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

func (s *loanGuideLineService) UploadFile(ctx context.Context, args *UploadObjectArgs) (*UploadObjectResponse, error) {
	fileHeader, ok := args.Media.FileHeader.(*multipart.FileHeader)
	if !ok {
		return nil, nil
	}
	userMetaData := map[string]string{
		"x-amz-acl": "public-read",
	}
	res, err := s.MinioClient.UploadFile(ctx, &minio.UploadObjectArgs{
		MediaStorage: &media_library.MediaLibraryStorage{
			OSS:   args.Media,
			Sizes: args.Media.GetSizes(),
		},
		UserMetaData: userMetaData,
		FileHeaders: []*multipart.FileHeader{
			fileHeader,
		},
	})
	if err != nil {
		return nil, err
	}
	return &UploadObjectResponse{
		URL: res.URL,
	}, nil
}

func (s *loanGuideLineService) CreateLoanGuideline(ctx context.Context, args *model.LoanGuideline) error {
	for _, m := range args.Medias {
		if m.Thumbnail.FileHeader != nil {
			res, err := s.UploadFile(ctx, &UploadObjectArgs{
				Media: m.Thumbnail,
			})
			if err != nil {
				return errorx.ErrInternal(errors.New("Cannot upload media"))
			}
			m.URL = res.URL
			m.Thumbnail.Url = res.URL
		}
	}

	return s.loanGuideLineRepo.CreateLoanGuideline(ctx, args)
}

func (s *loanGuideLineService) GetListLoanGuideline(ctx context.Context, args *loanguideline.GetListRequest) (*loanguideline.GetListResponse, error) {
	l, err := s.loanGuideLineRepo.GetListLoanGuideline(&loanguideline.GetListQuery{
		Limit:  *args.Limit,
		Offset: *args.Offset,
	})
	if err != nil {
		return nil, err
	}
	return &loanguideline.GetListResponse{
		LoanGuidelines: Convert_model_LoanGuidelines_to_service_LoanGuidelines(l),
	}, nil
}

func (s *loanGuideLineService) GetLoanGuideline(ctx context.Context, args *loanguideline.GetRequest) (*loanguideline.LoanGuideline, error) {
	res, err := s.loanGuideLineRepo.GetFirstBy(&model.LoanGuideline{
		ID: *args.ID,
	})
	if err != nil {
		return nil, err
	}
	return Convert_model_LoanGuideline_to_service_LoanGuideline(res), nil
}
