package service

import (
	"content-management/app/loanguideline"
	"content-management/model"

	"github.com/qor/media/oss"
)

type UploadObjectArgs struct {
	Media oss.OSS
}

type UploadObjectResponse struct {
	URL string
}

type CreateLoanGuidelineArgs struct {
	Type     string
	HTMLCode *string
	CSSCode  *string
	VideoURL *string
	Medias   []*model.Media
}

func Convert_model_LoanGuideline_to_service_LoanGuideline(args *model.LoanGuideline) *loanguideline.LoanGuideline {
	if args == nil {
		return nil
	}
	return &loanguideline.LoanGuideline{
		ID:        args.ID,
		HTMLCode:  args.HTMLCode,
		CSSCode:   args.CSSCode,
		Medias:    Convert_model_Medias_to_service_Medias(args.Medias),
		CreatedAt: args.CreatedAt,
		UpdatedAt: args.UpdatedAt,
	}
}

func Convert_model_LoanGuidelines_to_service_LoanGuidelines(args []*model.LoanGuideline) []*loanguideline.LoanGuideline {
	l := make([]*loanguideline.LoanGuideline, 0)
	for _, v := range args {
		l = append(l, Convert_model_LoanGuideline_to_service_LoanGuideline(v))
	}
	return l
}

func Convert_model_Media_to_service_Media(args *model.Media) *loanguideline.Media {
	if args == nil {
		return nil
	}
	return &loanguideline.Media{
		ID:              args.ID,
		URL:             args.URL,
		Description:     args.Description,
		LoanGuidelineID: args.LoanGuidelineID,
		CreatedAt:       args.CreatedAt,
		UpdatedAt:       args.UpdatedAt,
	}
}

func Convert_model_Medias_to_service_Medias(args []*model.Media) []*loanguideline.Media {
	l := make([]*loanguideline.Media, 0)
	for _, v := range args {
		l = append(l, Convert_model_Media_to_service_Media(v))
	}
	return l
}
