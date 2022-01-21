package registry

import "content-management/app/loanguideline/service"

// Customer ...
func (r *Registry) NewLoanGuideLineService() service.LoanGuideLineService {
	return service.NewLoanGuideLineService(&service.Dependency{
		DB:          r.DB.GormDB,
		MinioClient: r.MinioClient,
	})
}
