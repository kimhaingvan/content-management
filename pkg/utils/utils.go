package utils

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/microcosm-cc/bluemonday"
	"github.com/qor/l10n"
	"github.com/qor/session"
	"github.com/qor/session/manager"
)

// GetCurrentLocale get current locale from request
func GetCurrentLocale(req *http.Request) string {
	locale := l10n.Global
	if cookie, err := req.Cookie("locale"); err == nil {
		locale = cookie.Value
	}
	return locale
}

// URLParam get url params from request
func URLParam(name string, req *http.Request) string {
	return chi.URLParam(req, name)
}

// AddFlashMessage helper
func AddFlashMessage(w http.ResponseWriter, req *http.Request, message string, mtype string) error {
	return manager.SessionManager.Flash(w, req, session.Message{Message: template.HTML(message), Type: mtype})
}

// HTMLSanitizer HTML sanitizer
var HTMLSanitizer = bluemonday.UGCPolicy()

func FormatPrice(price interface{}) string {
	switch price.(type) {
	case float32, float64:
		return fmt.Sprintf("%0.2f", price)
	case int, uint, int32, int64, uint32, uint64:
		return fmt.Sprintf("%d.00", price)
	}
	return ""
}
