package build

import (
	"content-management/pkg/application"
	"content-management/pkg/database"

	"github.com/gorilla/mux"

	"github.com/qor/admin"
	"github.com/qor/publish2"
)

type Output struct {
	Db *database.Database
}

func buildApplication(db *database.Database) *application.Application {
	admin := admin.New(&admin.AdminConfig{
		SiteName: "MAFC CMS",
		//Auth:     auth.AdminAuth{},
		DB: db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
	})
	app := application.New(&application.AppConfig{
		Router: mux.NewRouter(),
		Admin:  admin,
		DB:     db.DB,
	})
	return app
}
