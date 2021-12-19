package model

import (
	"content-management/lib/qor/metax"

	"github.com/qor/admin"
)

var OrderFilters = []*admin.Filter{
	{

		Name:       "UserID",
		Label:      "Media Type",
		Operations: []string{"contains"},
		Type:       metax.SelectOne.String(),
		Config:     &admin.SelectOneConfig{Collection: []string{"1231", "2"}},
	},
	{
		Name: "Description",
		Type: "string",
	},
}
