package model

import "github.com/qor/admin"

var OrderActions = []*admin.Action{
	{
		Name:        "Import Product",
		URLOpenType: "slideout",
		URL: func(record interface{}, context *admin.Context) string {
			return "/admin/workers/new?job=Import Products"
		},
		Modes: []string{"collection"},
	},
}
