package order

import (
	"content-management/app/order/controller"
	"content-management/app/order/model"
	"content-management/pkg/application"
	"content-management/pkg/integration/storage/s3/client"
	"content-management/pkg/integration/storage/s3/driver"
	"content-management/pkg/utils/funcmapmaker"

	"github.com/qor/admin"
	"github.com/qor/render"
)

// New new order app
func New(config *Config) *OrderMicroApp {
	if config.Prefix == "" {
		config.Prefix = "/admin"
	}
	s3Driver := driver.New(client.Config{
		AwsS3Region:        config.AwsS3Region,
		AwsS3Bucket:        config.AwsS3Bucket,
		AwsAccessKey:       config.AwsAccessKey,
		AwsSecretAccessKey: config.AwsSecretAccessKey,
		AwsSessionToken:    config.AwsSessionToken,
	})

	return &OrderMicroApp{
		Config: config,
		Controller: &controller.Controller{
			S3Driver: s3Driver,
		},
	}
}

// App order app
type OrderMicroApp struct {
	Config     *Config
	Controller *controller.Controller
}

// Config order config struct
type Config struct {
	Prefix             string
	AwsS3Region        string
	AwsS3Bucket        string
	AwsAccessKey       string
	AwsSecretAccessKey string
	AwsSessionToken    string
}

// ConfigureApplication configure application
func (app *OrderMicroApp) ConfigureApplication(application *application.Application) {
	// ViewPaths tính từ file main.go
	app.Controller.View = render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("orders")}, "app/order/views")
	funcmapmaker.AddFuncMapMaker(app.Controller.View)
	admin := application.Admin
	app.ConfigureAdmin(application.Admin)

	application.Router.Post(app.Config.Prefix+"/ExtraFunc", app.Controller.ExtraFunc)
	application.Router.Mount(app.Config.Prefix, admin.NewServeMux(app.Config.Prefix))
}

// ConfigureAdmin configure admin interface
func (*OrderMicroApp) ConfigureAdmin(Admin *admin.Admin) {
	// Add Order
	Admin.DB.AutoMigrate(&model.Order{})
	order := Admin.AddResource(&model.Order{}, &admin.Config{Menu: []string{"Order Management"}})
	order.Meta(&admin.Meta{Name: "ShippedAt", Type: "date"})
	order.Meta(&admin.Meta{Name: "Description", Config: &admin.RichEditorConfig{Plugins: []admin.RedactorPlugin{
		{Name: "medialibrary", Source: "/admin/assets/javascripts/qor_redactor_medialibrary.js"},
		{Name: "table", Source: "/vendors/redactor_table.js"},
	},
		Settings: map[string]interface{}{
			"medialibraryUrl": "/admin/product_images",
		},
	}})
}
