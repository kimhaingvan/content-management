package registry

import "content-management/app/order/service"

// Customer ...
func (r *Registry) NewOrderService() service.OrderService {
	return service.NewOrderService(&service.Dependency{
		DB:          r.DB.GormDB,
		MinioClient: r.MinioClient,
	})
}
