package media

import (
	"content-management/app/media/controller"
	"content-management/pkg/application"
	"content-management/pkg/utils/funcmapmaker"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/qor/media/media_library"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/render"
	"github.com/qor/roles"
)

// App order app
type OrderMicroApp struct {
	config      *Config
	controller  *controller.Controller
	minioClient *minio.Client
}

// Config order config struct
type Config struct {
	Prefix          string
	Endpoint        string
	SecretAccessKey string
	AccessKey       string
	BucketName      string
}

// New new order app
func New(config *Config) *OrderMicroApp {
	if config.Prefix == "" {
		config.Prefix = "/"
	}

	return &OrderMicroApp{
		config:     config,
		controller: &controller.Controller{},
	}
}

// ConfigureApplication configure application
func (app *OrderMicroApp) ConfigureApplication(application *application.Application) {
	// ViewPaths tính từ file main.go
	app.controller.View = render.New(&render.Config{AssetFileSystem: application.AssetFS.NameSpace("media")}, "app/media/views")
	funcmapmaker.AddFuncMapMaker(app.controller.View)
	admin := application.Admin
	app.ConfigureAdmin(application.Admin)

	application.Router.PathPrefix(app.config.Prefix).Handler(
		http.StripPrefix(
			strings.TrimSuffix(app.config.Prefix, "/"),
			admin.NewServeMux("/"),
		),
	)
}

// ConfigureAdmin configure admin interface
func (o *OrderMicroApp) ConfigureAdmin(Admin *admin.Admin) {
	// Add Order
	Admin.DB.AutoMigrate(&media_library.MediaLibrary{})
	media := Admin.AddResource(
		&media_library.MediaLibrary{},
		&admin.Config{
			Menu: []string{"Media Management"},
		},
	)
	media.SaveHandler = func(result interface{}, cont *qor.Context) error {
		if (cont.GetDB().NewScope(result).PrimaryKeyZero() &&
			media.HasPermission(roles.Create, cont)) || // has create permission
			media.HasPermission(roles.Update, cont) { // has update permission
			var err error
			var contentType, fileName string
			var fileSize int64
			var file multipart.File
			var mediaFile *media_library.MediaLibrary

			// Save file in CMC Cloud
			if cont.Request != nil && cont.Request.MultipartForm != nil && cont.Request.MultipartForm.File["QorResource.File"] != nil {
				fileHeaders := cont.Request.MultipartForm.File["QorResource.File"]
				if len(fileHeaders) > 0 && fileHeaders[0] != nil {
					fileHeader := fileHeaders[0]
					fileName = fileHeader.Filename
					fileSize = fileHeader.Size
					contentTypes := fileHeader.Header["Content-Type"]
					if len(contentTypes) > 0 && contentTypes[0] != "" {
						contentType = contentTypes[0]
					}
				}

				var ok bool
				mediaFile, ok = result.(*media_library.MediaLibrary)
				if !ok {
					return errors.New("Can not convert file")
				}

				mediaFile.File.SelectedType = contentType
				mediaFile.SelectedType = contentType

				// Public url of file
				userMetaData := map[string]string{
					"x-amz-acl": "public-read",
				}

				fileURL := fmt.Sprintf("https://%v/%v/%v", o.config.Endpoint, o.config.BucketName, fileName)
				mediaFile.File.Url = fileURL
				mediaFile.File.Description = fileURL

				file, err = mediaFile.File.FileHeader.Open()
				if err != nil {
					return err
				}
				defer file.Close()
				// Initialize minio client object.
				minioClient, err := minio.New(o.config.Endpoint, &minio.Options{
					Creds:  credentials.NewStaticV4(o.config.AccessKey, o.config.SecretAccessKey, ""),
					Secure: true,
				})
				if err != nil {
					return err
				}

				var f io.Reader
				f = file

				// Upload file to CMC Cloud (Minio)
				_, err = minioClient.PutObject(cont.Request.Context(), o.config.BucketName, fileName, f, fileSize, minio.PutObjectOptions{
					ContentType:  contentType,
					UserMetadata: userMetaData,
				})
				if err != nil {
					return err
				}
			}

			// Save file infomations in local database
			if err = cont.GetDB().Save(mediaFile).Error; err != nil {
				return err
			}
			return nil
		}
		return roles.ErrPermissionDenied
	}

	media.DeleteHandler = func(result interface{}, context *qor.Context) error {
		if media.HasPermission(roles.Delete, context) {
			if primaryQuerySQL, primaryParams := media.ToPrimaryQueryParams(context.ResourceID, context); primaryQuerySQL != "" {
				if !context.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).RecordNotFound() {
					// Delete file in local database
					err := context.GetDB().Delete(result).Error
					if err != nil {
						return err
					}

					// Initialize minio client object and delete file.
					minioClient, err := minio.New(o.config.Endpoint, &minio.Options{
						Creds:  credentials.NewStaticV4(o.config.AccessKey, o.config.SecretAccessKey, ""),
						Secure: true,
					})
					if err != nil {
						return err
					}
					db := context.GetDB().First(result).Value
					mediaFile := db.(*media_library.MediaLibrary)
					fileName := mediaFile.File.FileName
					err = minioClient.RemoveObject(context.Request.Context(), o.config.BucketName, fileName, minio.RemoveObjectOptions{})
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
