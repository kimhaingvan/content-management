package loanguideline

import (
	"content-management/lib/qor/metax"
	"content-management/model"

	"github.com/qor/admin"
)

var Metas = []*admin.Meta{
	{
		Name:   "Type",
		Label:  "Loại",
		Type:   metax.SelectOne.String(),
		Config: &admin.SelectOneConfig{Collection: model.LoanGuidelineTypes},
	},
	{
		Name:  "HTMLCode",
		Label: "Html code",
		Type:  metax.String.String(),
	},
}
