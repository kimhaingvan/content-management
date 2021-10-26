package application

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jinzhu/gorm"
	"github.com/qor/admin"
	"github.com/qor/assetfs"
	"github.com/qor/middlewares"
	"github.com/qor/publish2"
	"github.com/qor/wildcard_router"
)

// MicroAppInterface micro app interface
type MicroAppInterface interface {
	ConfigureApplication(*Application)
}

// Application main application
type Application struct {
	*AppConfig
}

type AppConfig struct {
	Router   *mux.Router
	Handlers []http.Handler
	AssetFS  assetfs.Interface
	Admin    *admin.Admin
	DB       *gorm.DB
}

// New new application
func New(config *AppConfig) *Application {
	if config == nil {
		config = &AppConfig{}
	}

	if config.Router == nil {
		config.Router = mux.NewRouter()
	}

	if config.AssetFS == nil {
		config.AssetFS = assetfs.AssetFS()
	}

	if config.Admin == nil {
		admin := admin.New(&admin.AdminConfig{
			SiteName: "MAFC CMS",
			//Auth:     auth.AdminAuth{},
			DB: config.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
		})
		config.Admin = admin
	}
	return &Application{
		AppConfig: config,
	}
}

// Use mount router into micro app
func (application *Application) Use(apps ...MicroAppInterface) {
	for _, app := range apps {
		app.ConfigureApplication(application)
	}
}

// NewServeMux allocates and returns a new ServeMux.
func (application *Application) NewServeMux() http.Handler {
	if len(application.Handlers) == 0 {
		return middlewares.Apply(application.Router)
	}

	wildcardRouter := wildcard_router.New()
	for _, handler := range application.Handlers {
		wildcardRouter.AddHandler(handler)
	}
	wildcardRouter.AddHandler(application.Router)

	return middlewares.Apply(wildcardRouter)
}
