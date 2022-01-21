package loanguideline

import (
	"content-management/app/loanguideline/handler"
	"content-management/model"
	"content-management/pkg/application"
	"content-management/pkg/utils/funcmapmaker"
	"fmt"
	"github.com/qor/admin"
	"github.com/qor/render"
)

// Config Loan Guide Line config struct
type Config struct {
	Prefix      string
	Application *application.Application
}

// App Loan guideline app
type LoanGuideLineServer struct {
	*Config
	handler *handler.LoanGuidelineHandler
}

// NewLoanGuideLineServer new loan guideline app
func NewLoanGuideLineServer(config *Config) *LoanGuideLineServer {
	loanGuideLineService := config.Application.Registry.NewLoanGuideLineService()
	return &LoanGuideLineServer{
		Config:  config,
		handler: handler.NewLoanGuideLineHandler(loanGuideLineService),
	}
}

// ConfigureApplication configure application
func (o *LoanGuideLineServer) Configure() {
	// ViewPaths tính từ file main.go
	o.handler.View = render.New(&render.Config{AssetFileSystem: o.Application.AssetFS.NameSpace("orders")}, "app/loanguideline/view")
	funcmapmaker.AddFuncMapMaker(o.handler.View)

	qorAdmin := o.Application.Admin
	o.configure(qorAdmin)
	o.router(qorAdmin.GetRouter())
}

func (o *LoanGuideLineServer) router(router *admin.Router) {
	//router.Post(fmt.Sprintf("%v/TestFunc", o.Prefix), o.handler.TestSample)
	//router.Post(fmt.Sprintf("%v/TestError", o.Prefix), o.handler.TestError)
	router.Get(fmt.Sprintf("%v", o.Prefix), o.handler.GetAllLoanGuidelineHandler)
	router.Get(fmt.Sprintf("%v", o.Prefix), o.handler.GetByIdLoanGuidelineHandler)

}

// configure configure loan guideline admin interface
func (o *LoanGuideLineServer) configure(qorAdmin *admin.Admin) {
	qorAdmin.DB.AutoMigrate(&model.LoanGuideline{})
	loanGuideline := qorAdmin.AddResource(
		&model.LoanGuideline{},
		&admin.Config{
			Name:       "Loan guideline",
			Menu:       []string{"My Finance"},
			Permission: nil,
			Themes:     nil,
			Priority:   0,
			Singleton:  false,
			Invisible:  false,
			PageCount:  0,
		},
	)

	loanGuideline.UseTheme("fancy")


	for _, meta := range Metas {
		loanGuideline.Meta(meta)
	}

	//for _, filter := range model.Filters {
	//	order.Filter(filter)
	//}
	//
	//for _, action := range model.Actions {
	//	order.Action(action)
	//}
	//
	//for _, validator := range model.Validators {
	//	order.AddValidator(validator)
	//}
	//
	//if len(model.NewAttrs) > 0 {
	//	order.NewAttrs(model.NewAttrs)
	//}
	//
	//if len(model.IndexAttrs) > 0 {
	//	order.IndexAttrs(model.IndexAttrs)
	//}
	//
	//if len(model.EditAttrs) > 0 {
	//	order.EditAttrs(model.EditAttrs)
	//}
	//
	//if len(model.ShowAttrs) > 0 {
	//	order.ShowAttrs(model.ShowAttrs)
	//}
	//
	//if len(model.SearchAttrs) > 0 {
	//	order.SearchAttrs(model.SearchAttrs...)
	//}

	loanGuideline.SaveHandler = o.handler.SaveLoanGuidelineHandler(loanGuideline)
	//
	//order.DeleteHandler = o.orderHandler.DeleteOrderHandler(order)
	//
	//order.FindOneHandler = o.orderHandler.FindOneHandler(order)
	//
	//order.FindManyHandler = o.orderHandler.FindManyHandler(order)
}
