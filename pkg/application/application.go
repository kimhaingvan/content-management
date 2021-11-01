package application

import (
	"content-management/registry"
	"net/http"

	"github.com/qor/middlewares"
	"github.com/qor/wildcard_router"

	"github.com/go-chi/chi"

	"github.com/qor/admin"
	"github.com/qor/assetfs"
)

// MicroAppInterface micro app interface
type MicroAppInterface interface {
	Configure()
}

type Config struct {
	Registry *registry.Registry
	Router   *chi.Mux
	Handlers []http.Handler
	AssetFS  assetfs.Interface
	Admin    *admin.Admin
}

// Application main application
type Application struct {
	*Config
}

// NewApplication new application
func NewApplication(config *Config) *Application {
	if config.Router == nil {
		config.Router = chi.NewRouter()
	}
	return &Application{
		Config: config,
	}
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.NewServeMux().ServeHTTP(w, r)
}

// NewServeMux allocates and returns a new ServeMux.
func (a *Application) NewServeMux() http.Handler {
	if len(a.Handlers) == 0 {
		return middlewares.Apply(a.Router)
	}

	wildcardRouter := wildcard_router.New()
	for _, handler := range a.Handlers {
		wildcardRouter.AddHandler(handler)
	}
	wildcardRouter.AddHandler(a.Router)

	return middlewares.Apply(wildcardRouter)
}

// Use mount router into micro app
func Use(apps ...MicroAppInterface) {
	for _, app := range apps {
		app.Configure()
	}
}
