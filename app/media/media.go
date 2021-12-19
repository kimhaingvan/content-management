package media

import (
	"content-management/app/media/controller"
	"content-management/app/media/model"
	"content-management/pkg/application"
	"content-management/pkg/httpreq"
	minio2 "content-management/pkg/integration/minio"
	"content-management/pkg/utils/funcmapmaker"
	"errors"
	"fmt"
	"net/http"

	"go.elastic.co/apm"

	"github.com/qor/media/media_library"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/render"
	"github.com/qor/roles"
)

// app
type MediaMicroApp struct {
	config      *Config
	controller  *controller.Controller
	minioClient *minio2.Client
}

// Config order config struct
type Config struct {
	Prefix      string
	MinioClient *minio2.Client
}

// New new order app
func New(config *Config) *MediaMicroApp {
	if config.Prefix == "" {
		config.Prefix = "/"
	}

	return &MediaMicroApp{
		config:      config,
		controller:  &controller.Controller{},
		minioClient: config.MinioClient,
	}
}

// ConfigureApplication configure application
func (app *MediaMicroApp) ConfigureApplication(application *application.Application) {
	// ViewPaths tính từ file main.go
	app.controller.View = render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("medias")}, "app/media/view")
	funcmapmaker.AddFuncMapMaker(app.controller.View)

	qorAdmin := application.Admin
	app.ConfigureAdmin(qorAdmin)

	qorAdmin.GetRouter().Post("/GetMedia", app.controller.GetMediaFile)
	application.Router.PathPrefix(app.config.Prefix).Handler(
		http.StripPrefix(
			"/",
			application.Handler,
		),
	)

}

// ConfigureAdmin configure admin interface
func (o *MediaMicroApp) ConfigureAdmin(Admin *admin.Admin) {
	Admin.DB.AutoMigrate(&model.CustomMedia{})
	media := Admin.AddResource(
		&model.CustomMedia{},
		&admin.Config{
			Menu: []string{"Media Management"},
		},
	)

	for _, meta := range model.MediaMetas {
		media.Meta(meta)
	}
	media.UseTheme("media_library")

	media.SaveHandler = func(result interface{}, cont *qor.Context) error {
		trx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", cont.Request.Method, cont.Request.RequestURI), "request")
		trxCtx := apm.ContextWithTransaction(cont.Request.Context(), trx)
		if (cont.GetDB().NewScope(result).PrimaryKeyZero() &&
			media.HasPermission(roles.Create, cont)) || // has create permission
			media.HasPermission(roles.Update, cont) { // has update permission

			// Save file in CMC Cloud
			mediaFile, _ := result.(*model.CustomMedia)
			if httpreq.HasMediaFile(cont.Request) {
				res, err := o.minioClient.UploadFileToMinioFromRequest(trxCtx, cont.Request, mediaFile.File)
				if err != nil {
					return err
				}
				mediaFile.File.Url = res.URL
				mediaFile.URL = res.URL
			}
			// Save file infomations in local database
			if err := cont.GetDB().Save(mediaFile).Error; err != nil {
				return err
			}
			return nil
		}
		return roles.ErrPermissionDenied
	}

	media.DeleteHandler = func(result interface{}, context *qor.Context) error {
		trx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", context.Request.Method, context.Request.RequestURI), "request")
		trxCtx := apm.ContextWithTransaction(context.Request.Context(), trx)
		if media.HasPermission(roles.Delete, context) {
			if primaryQuerySQL, primaryParams := media.ToPrimaryQueryParams(context.ResourceID, context); primaryQuerySQL != "" {
				if !context.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).RecordNotFound() {
					// Delete file in local database
					err := context.GetDB().Delete(result).Error
					if err != nil {
						return err
					}

					// Initialize minio client object and delete file.
					db := context.GetDB().First(result).Value
					mediaFile := db.(*media_library.MediaLibrary)
					fileName := mediaFile.File.FileName
					err = o.minioClient.RemoveObject(trxCtx, fileName)
					return err
				}
			}
			return gorm.ErrRecordNotFound
		}
		return roles.ErrPermissionDenied
	}

	media.FindOneHandler = func(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		if media.HasPermission(roles.Read, context) {
			var (
				primaryQuerySQL string
				primaryParams   []interface{}
			)

			if metaValues == nil {
				primaryQuerySQL, primaryParams = media.ToPrimaryQueryParams(context.ResourceID, context)
			} else {
				primaryQuerySQL, primaryParams = media.ToPrimaryQueryParamsFromMetaValue(metaValues, context)
			}

			if primaryQuerySQL != "" {
				if metaValues != nil {
					if destroy := metaValues.Get("_destroy"); destroy != nil {
						if fmt.Sprint(destroy.Value) != "0" && media.HasPermission(roles.Delete, context) {
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

	media.FindManyHandler = func(result interface{}, context *qor.Context) error {
		if media.HasPermission(roles.Read, context) {
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
