package build

import (
	"content-management/core/config"
	"content-management/pkg/application"
	"content-management/pkg/database"

	"github.com/gorilla/mux"
	"github.com/qor/assetfs"

	"github.com/qor/admin"
	"github.com/qor/publish2"
)

func BuildApplication(cfg config.Config) *application.Application {
	database := database.New(cfg.Databases.PostgresConfig)
	admin := admin.New(&admin.AdminConfig{
		SiteName: "MAFC CMS",
		DB:       database.GormDB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
	})

	router := mux.NewRouter()

	app := application.New(&application.Config{
		Router:   router,
		Admin:    admin,
		DB:       database.GormDB,
		AssetFS:  assetfs.AssetFS(),
		Handlers: nil,
	})
	return app
}
