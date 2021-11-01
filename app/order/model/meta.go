package model

import (
	"content-management/lib/qor/metax"
	"fmt"

	"github.com/qor/admin"
)

var OrderMetas = []*admin.Meta{
	{
		Label: "Shipped at",
		Name:  "ShippedAt",
		Type:  "date",
	},
	{

		Name: "DeliveryMethod",
		Type: "select_one",
		Config: &admin.SelectOneConfig{
			AllowBlank: true,
			Collection: func(_ interface{}, context *admin.Context) (options [][]string) {
				var methods []DeliveryMethod
				context.GetDB().Find(&methods)

				for _, m := range methods {
					idStr := fmt.Sprintf("%d", m.ID)
					var option = []string{idStr, fmt.Sprintf("%s (%0.2f) руб", m.Name, m.Price)}
					options = append(options, option)
				}

				return options
			},
		},
	},
	{
		Name: "Description",
		Type: metax.Readonly.String(),
		//Config: &admin.RichEditorConfig{
		//	Plugins: []admin.RedactorPlugin{
		//		{Name: "medialibrary", Source: "/admin/assets/javascripts/qor_redactor_medialibrary.js"},
		//		{Name: "table", Source: "/vendors/redactor_table.js"},
		//	},
		//	Settings: map[string]interface{}{
		//		"medialibraryUrl": "/admin/product_images",
		//	},
		//},
	},
}
