package service

import (
	"content-management/app/order/repository"
	"content-management/thirdparty/minio"
	"context"

	"github.com/jinzhu/gorm"
)

type OrderService interface {
	Save(int) int
	UploadFile(ctx context.Context, args *UploadObjectArgs) (string, error)
}

type orderService struct {
	*Dependency
	orderRepo repository.OrderRepository
}

type Dependency struct {
	DB          *gorm.DB
	MinioClient *minio.Client
}

func NewOrderService(d *Dependency) OrderService {
	return &orderService{
		Dependency: d,
		orderRepo:  repository.NewOrderRepository(d.DB),
	}
}

func (o *orderService) UploadFile(ctx context.Context, args *UploadObjectArgs) (string, error) {
	// Public url of file
	userMetaData := map[string]string{
		"x-amz-acl": "public-read",
	}
	fileHeaders := args.Form.File["QorResource.File"]
	res, err := o.MinioClient.UploadFile(ctx, &minio.UploadObjectArgs{
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

func (o orderService) Save(i int) int {
	return i
}
