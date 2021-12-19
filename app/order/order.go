package order

import (
	"content-management/app/order/controller"
	"content-management/app/order/model"
	"content-management/pkg/application"
	"content-management/pkg/httpreq"
	minio2 "content-management/pkg/integration/minio"
	"content-management/pkg/utils/funcmapmaker"
	"errors"
	"fmt"
	"net/http"

	"go.elastic.co/apm"

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

	return &OrderMicroApp{
		config:      config,
		controller:  &controller.Controller{},
		minioClient: config.MinioClient,
	}
}

// App order app
type OrderMicroApp struct {
	config      *Config
	controller  *controller.Controller
	minioClient *minio2.Client
}

// Config order config struct
type Config struct {
	Prefix      string
	MinioClient *minio2.Client
}

// ConfigureApplication configure application
func (o *OrderMicroApp) ConfigureApplication(application *application.Application) {
	// ViewPaths tính từ file main.go
	o.controller.View = render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("orders")}, "app/order/view")
	funcmapmaker.AddFuncMapMaker(o.controller.View)

	qorAdmin := application.Admin
	o.ConfigureAdmin(qorAdmin)

	qorAdmin.GetRouter().Post("/TestFunc", o.controller.TestFunc)
	application.Router.PathPrefix(o.config.Prefix).Handler(
		http.StripPrefix(
			"/",
			application.Handler,
		),
	)

}

// ConfigureAdmin configure order admin interface
func (o *OrderMicroApp) ConfigureAdmin(qorAdmin *admin.Admin) {
	// Add Order
	qorAdmin.DB.AutoMigrate(&model.Order{}, &model.OrderItem{}, &model.DeliveryMethod{})
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

	//order.UseTheme("sorting_mode")
	for _, meta := range model.OrderMetas {
		order.Meta(meta)
	}

	for _, filter := range model.OrderFilters {
		order.Filter(filter)
	}

	for _, action := range model.OrderActions {
		order.Action(action)
	}

	for _, validator := range model.OrderValidators {
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

	order.SaveHandler = func(result interface{}, cont *qor.Context) error {
		trx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", cont.Request.Method, cont.Request.RequestURI), "request")
		trxCtx := apm.ContextWithTransaction(cont.Request.Context(), trx)
		defer trx.End()
		if (cont.GetDB().NewScope(result).PrimaryKeyZero() &&
			order.HasPermission(roles.Create, cont)) || // has create permission
			order.HasPermission(roles.Update, cont) { // has update permission

			mediaFile, _ := result.(*model.Order)
			if httpreq.HasMediaFile(cont.Request) {
				_, err := o.minioClient.UploadFileToMinioFromRequest(trxCtx, cont.Request, mediaFile.File)
				if err != nil {
					return err
				}
			}

			if err := cont.GetDB().Save(result).Error; err != nil {
				return err
			}
			return nil
		}
		return roles.ErrPermissionDenied
	}

	order.DeleteHandler = func(result interface{}, context *qor.Context) error {
		transaction := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", context.Request.Method, context.Request.RequestURI), "request")
		defer transaction.End()
		if order.HasPermission(roles.Delete, context) {
			if primaryQuerySQL, primaryParams := order.ToPrimaryQueryParams(context.ResourceID, context); primaryQuerySQL != "" {
				if !context.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).RecordNotFound() {
					err := context.GetDB().Delete(result).Error
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
				return err
			}
			return errors.New("failed to find")
		}
		return roles.ErrPermissionDenied
	}

	order.FindManyHandler = func(result interface{}, context *qor.Context) error {
		if order.HasPermission(roles.Read, context) {
			db := context.GetDB()
			if _, ok := db.Get("qorm:getting_total_count"); ok {
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
