package build

import (
	_ "content-management/core/docs"
	"content-management/pkg/application"
	"content-management/registry"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/qor/admin"
	"github.com/qor/assetfs"
	"github.com/qor/publish2"
)

func BuildApplication(r *registry.Registry) *application.Application {
	admin := admin.New(&admin.AdminConfig{
		SiteName: "MAFC CMS",
		DB:       r.DB.GormDB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
	})
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	app := application.NewApplication(&application.Config{
		Registry: r,
		Router:   router,
		AssetFS:  assetfs.AssetFS(),
		Admin:    admin,
	})
	return app
}
