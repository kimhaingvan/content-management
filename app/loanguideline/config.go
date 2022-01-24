package loanguideline

import (
	"content-management/app/loanguideline/model"
	"content-management/lib/qor/metax"

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
