package order

import (
	"content-management/app/order/handler"
	"content-management/app/order/model"
	"content-management/app/order/service"
	"content-management/pkg/application"
	"content-management/pkg/utils/funcmapmaker"
	"fmt"

	"github.com/qor/admin"
	"github.com/qor/render"
)

// Config order config struct
type Config struct {
	Prefix      string
	Application *application.Application
}

// App order app
type OrderServer struct {
	*Config
	orderHandler *handler.OrderHandler
	orderService service.OrderService
}

// NewOrderServer new order app
func NewOrderServer(config *Config) *OrderServer {
	orderService := config.Application.Registry.NewOrderService()
	return &OrderServer{
		Config:       config,
		orderHandler: handler.NewOrderHandler(orderService),
		orderService: orderService,
	}
}

// ConfigureApplication configure application
func (o *OrderServer) Configure() {
	// ViewPaths tính từ file main.go
	o.orderHandler.View = render.New(&render.Config{AssetFileSystem: o.Application.AssetFS.NameSpace("orders")}, "app/order/view")
	funcmapmaker.AddFuncMapMaker(o.orderHandler.View)

	qorAdmin := o.Application.Admin
	o.configure(qorAdmin)
	o.router(qorAdmin.GetRouter())
}

func (o *OrderServer) router(router *admin.Router) {
	router.Post(fmt.Sprintf("%v/TestFunc", o.Prefix), o.orderHandler.TestSample)
}

// configure configure order admin interface
func (o *OrderServer) configure(qorAdmin *admin.Admin) {
	order := qorAdmin.AddResource(
		&model.Order{},
		&admin.Config{
			Name:       "order",
			Menu:       []string{"Order Management"},
			Permission: nil,
			Themes:     nil,
			Priority:   0,
			Singleton:  false,
			Invisible:  false,
			PageCount:  0,
		},
	)

	for _, meta := range model.Metas {
		order.Meta(meta)
	}

	for _, filter := range model.Filters {
		order.Filter(filter)
	}

	for _, action := range model.Actions {
		order.Action(action)
	}

	for _, validator := range model.Validators {
		order.AddValidator(validator)
	}

	if len(model.NewAttrs) > 0 {
		order.NewAttrs(model.NewAttrs)
	}

	if len(model.IndexAttrs) > 0 {
		order.IndexAttrs(model.IndexAttrs)
	}

	if len(model.EditAttrs) > 0 {
		order.EditAttrs(model.EditAttrs)
	}

	if len(model.ShowAttrs) > 0 {
		order.ShowAttrs(model.ShowAttrs)
	}

	if len(model.SearchAttrs) > 0 {
		order.SearchAttrs(model.SearchAttrs...)
	}

	order.SaveHandler = o.orderHandler.SaveOrderHandler(order)

	order.DeleteHandler = o.orderHandler.DeleteOrderHandler(order)

	order.FindOneHandler = o.orderHandler.FindOneHandler(order)

	order.FindManyHandler = o.orderHandler.FindManyHandler(order)
}
