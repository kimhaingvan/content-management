package admin

import (
	"content-management/pkg/database"

	"github.com/qor/admin"
	"github.com/qor/publish2"
)

func New(db *database.Database) *admin.Admin {
	admin := admin.New(&admin.AdminConfig{
		SiteName: "MAFC CMS",
		//Auth:     auth.AdminAuth{},
		DB: db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
	})
	return admin
}
