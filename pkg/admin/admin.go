package admin

import (
	"cms-project/pkg/database"
	"github.com/k0kubun/pp"
	"github.com/qor/admin"
	"github.com/qor/publish2"
)

func New(db *database.Database) *admin.Admin {
	admin := admin.New(&admin.AdminConfig{
		SiteName: "MAFC CMS",
		//Auth:     auth.AdminAuth{},
		DB: db.DB.Set(publish2.VisibleMode, publish2.ModeOff).Set(publish2.ScheduleMode, publish2.ModeOff),
	})
	pp.Println("NEW: ", admin)
	return admin
}
