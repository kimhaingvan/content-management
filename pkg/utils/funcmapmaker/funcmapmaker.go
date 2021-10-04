package funcmapmaker

import (
	"cms-project/pkg/utils"
	"html/template"
	"net/http"

	"github.com/qor/render"
	"github.com/qor/session"
	"github.com/qor/session/manager"
)

// AddFuncMapMaker add FuncMapMaker to view
func AddFuncMapMaker(view *render.Render) *render.Render {
	oldFuncMapMaker := view.FuncMapMaker
	view.FuncMapMaker = func(render *render.Render, req *http.Request, w http.ResponseWriter) template.FuncMap {
		funcMap := template.FuncMap{}
		if oldFuncMapMaker != nil {
			funcMap = oldFuncMapMaker(render, req, w)
		}

		funcMap["raw"] = func(str string) template.HTML {
			return template.HTML(utils.HTMLSanitizer.Sanitize(str))
		}

		funcMap["flashes"] = func() []session.Message {
			return manager.SessionManager.Flashes(w, req)
		}

		funcMap["current_locale"] = func() string {
			return utils.GetCurrentLocale(req)
		}

		funcMap["format_price"] = func(price interface{}) string {
			return utils.FormatPrice(price)
		}

		return funcMap
	}

	return view
}
