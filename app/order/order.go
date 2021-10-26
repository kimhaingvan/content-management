package order

import (
	"content-management/app/order/controler"
	"content-management/app/order/model"
	"content-management/pkg/application"
	"content-management/pkg/integration/storage/s3/client"
	"content-management/pkg/integration/storage/s3/driver"
	"content-management/pkg/utils/funcmapmaker"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/qor/qor/resource"

	"github.com/jinzhu/gorm"

	"github.com/qor/roles"

	"github.com/qor/qor"

	"github.com/qor/admin"
	"github.com/qor/render"
)

// New new order app
func New(config *Config) *OrderMicroApp {
	if config.Prefix == "" {
		config.Prefix = "/"
	}
	s3Driver := driver.New(client.Config{
		AwsS3Region:        config.AwsS3Region,
		AwsS3Bucket:        config.AwsS3Bucket,
		AwsAccessKey:       config.AwsAccessKey,
		AwsSecretAccessKey: config.AwsSecretAccessKey,
		AwsSessionToken:    config.AwsSessionToken,
	})

	return &OrderMicroApp{
		config: config,
		controller: &controler.Controller{
			S3Driver: s3Driver,
		},
	}
}

// App order app
type OrderMicroApp struct {
	config     *Config
	controller *controler.Controller
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
	app.controller.View = render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("orders")}, "app/order/views")
	funcmapmaker.AddFuncMapMaker(app.controller.View)
	admin := application.Admin
	app.ConfigureAdmin(application.Admin)
	application.Router.HandleFunc("/TestFunc", app.controller.TestFunc).Methods(http.MethodPost)

	application.Router.PathPrefix(app.config.Prefix).Handler(
		http.StripPrefix(
			strings.TrimSuffix(app.config.Prefix, "/"),
			admin.NewServeMux("/"),
		),
	)
}

// ConfigureAdmin configure admin interface
func (*OrderMicroApp) ConfigureAdmin(Admin *admin.Admin) {
	// Add Order
	Admin.DB.AutoMigrate(&model.Order{})
	order := Admin.AddResource(
		&model.Order{},
		&admin.Config{
			Menu: []string{"Order Management"},
		},
	)
	order.Meta(&admin.Meta{Name: "ShippedAt", Type: "date"})
	order.Meta(&admin.Meta{Name: "Description", Config: &admin.RichEditorConfig{Plugins: []admin.RedactorPlugin{
		{Name: "medialibrary", Source: "/admin/assets/javascripts/qor_redactor_medialibrary.js"},
		{Name: "table", Source: "/vendors/redactor_table.js"},
	},
		Settings: map[string]interface{}{
			"medialibraryUrl": "/admin/product_images",
		},
	}})
	order.SaveHandler = func(result interface{}, cont *qor.Context) error {
		if (cont.GetDB().NewScope(result).PrimaryKeyZero() &&
			order.HasPermission(roles.Create, cont)) || // has create permission
			order.HasPermission(roles.Update, cont) { // has update permission
			if err := cont.GetDB().Save(result).Error; err != nil {
				return err
			}
			return nil
		}
		return roles.ErrPermissionDenied
	}

	order.DeleteHandler = func(result interface{}, context *qor.Context) error {
		if order.HasPermission(roles.Delete, context) {
			if primaryQuerySQL, primaryParams := order.ToPrimaryQueryParams(context.ResourceID, context); primaryQuerySQL != "" {
				if !context.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).RecordNotFound() {
					err := context.GetDB().Delete(result).Error
					if err != nil {
					}
					return err
				}
			}
			return gorm.ErrRecordNotFound
		}
		return roles.ErrPermissionDenied
	}

	order.FindOneHandler = func(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		if order.HasPermission(roles.Read, context) {
			var (
				primaryQuerySQL string
				primaryParams   []interface{}
			)

			if metaValues == nil {
				primaryQuerySQL, primaryParams = order.ToPrimaryQueryParams(context.ResourceID, context)
			} else {
				primaryQuerySQL, primaryParams = order.ToPrimaryQueryParamsFromMetaValue(metaValues, context)
			}

			if primaryQuerySQL != "" {
				if metaValues != nil {
					if destroy := metaValues.Get("_destroy"); destroy != nil {
						if fmt.Sprint(destroy.Value) != "0" && order.HasPermission(roles.Delete, context) {
							context.GetDB().Delete(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...)
							return resource.ErrProcessorSkipLeft
						}
					}
				}
				err := context.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).Error
				if err != nil {
					return err
				}
				return nil
			}
			return errors.New("failed to find")
		}
		return roles.ErrPermissionDenied
	}

	order.FindManyHandler = func(result interface{}, context *qor.Context) error {
		if order.HasPermission(roles.Read, context) {
			db := context.GetDB()
			if _, ok := db.Get("qor:getting_total_count"); ok {
				if err := context.GetDB().Count(result).Error; err != nil {
					return err
				}
				return nil
			}
			if err := context.GetDB().Set("gorm:order_by_primary_key", "DESC").Find(result).Error; err != nil {
				return err
			}
			return nil
		}
		return roles.ErrPermissionDenied
	}
}
