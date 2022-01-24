package build

import (
	"content-management/app/loanguideline/handler"
	"content-management/lib/qor/metax"
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
	router.Get(fmt.Sprintf("%v/get_list", o.Prefix), o.handler.GetListHandler)
	router.Get(fmt.Sprintf("%v/get", o.Prefix), o.handler.GetHandler)

}

// configure configure loan guideline admin interface
func (o *LoanGuideLineServer) configure(qorAdmin *admin.Admin) {
	qorAdmin.DB.AutoMigrate(&model.Media{}, &model.LoanGuideline{})
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
			PageCount:  20,
		},
	)

	for _, meta := range guidelineMetas {
		loanGuideline.Meta(meta)
	}

	mediasRes := loanGuideline.Meta(&admin.Meta{Name: "Medias"}).Resource
	for _, meta := range MediaGuidelineMetas {
		mediasRes.Meta(meta)
	}

	//for _, filter := range model.Filters {
	//	order.Filter(filter)
	//}

	//for _, action := range model.Actions {
	//	order.Action(action)
	//}

	//for _, validator := range model.Validators {
	//	order.AddValidator(validator)
	//}

	//if len(model.NewAttrs) > 0 {
	//	order.NewAttrs(model.NewAttrs)
	//}

	//if len(model.IndexAttrs) > 0 {
	//	order.IndexAttrs(model.IndexAttrs)
	//}

	//if len(model.EditAttrs) > 0 {
	//	order.EditAttrs(model.EditAttrs)
	//}

	//if len(model.ShowAttrs) > 0 {
	//	order.ShowAttrs(model.ShowAttrs)
	//}

	//if len(model.SearchAttrs) > 0 {
	//	order.SearchAttrs(model.SearchAttrs...)
	//}

	loanGuideline.SaveHandler = o.handler.SaveLoanGuidelineHandler(loanGuideline)
	//order.DeleteHandler = o.orderHandler.DeleteOrderHandler(order)
	loanGuideline.FindOneHandler = o.handler.FindOneHandler(loanGuideline)
	loanGuideline.FindManyHandler = o.handler.FindManyHandler(loanGuideline)
}

var (
	guidelineMetas = []*admin.Meta{
		{
			Name:   "Type",
			Label:  "Loại",
			Type:   metax.SelectOne.String(),
			Config: &admin.SelectOneConfig{Collection: model.LoanGuidelineTypes},
		},
		{
			Name:  "HTMLCode",
			Label: "Html code",
			Type:  metax.String.String(),
		},
	}

	MediaGuidelineMetas = []*admin.Meta{
		{
			Name: "URL",
			Type: metax.Readonly.String(),
		},
	}
)
