package model

import (
	"content-management/lib/qor/metax"

	"github.com/qor/admin"
)

var MediaMetas = []*admin.Meta{
	{
		Name:  "URL",
		Label: "File URL",
		Type:  metax.Readonly.String(),
	},
}
