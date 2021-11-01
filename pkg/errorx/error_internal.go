package errorx

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

// err ...
type internalError struct {
	Code     int
	Err      error
	Message  string
	Original string
	OrigFile string
	OrigLine int
	Meta     map[string]string
}

func DefaultErrorMessage(code int) string {
	switch code {
	case http.StatusOK:
		return ""
	case http.StatusNotFound:
		return "Không tìm thấy."
	case http.StatusBadRequest:
		return "Có lỗi xảy ra."
	case http.StatusInternalServerError:
		return "Lỗi không xác định."
	case http.StatusUnauthorized:
		return "Vui lòng đăng nhập (hoặc đăng ký nếu chưa có tài khoản)."
	case http.StatusForbidden:
		return "Không tìm thấy hoặc cần quyền truy cập."
	}
	return "Lỗi không xác định."
}

func newError(code int, message string, err error) *internalError {
	if message == "" {
		message = DefaultErrorMessage(code)
	}
	if err != nil {
		// Overwrite *apiError
		if xerr, ok := err.(*internalError); ok {
			// Keep original message
			meta := map[string]string{}
			if xerr.Original == "" {
				xerr.Original = xerr.Message
			}
			if xerr.Original != "" {
				meta["orig"] = xerr.Original
			}
			if xerr.Err != nil {
				meta["cause"] = xerr.Err.Error()
			}
			xerr.Code = code
			xerr.Message = message
			xerr.Meta = meta
			return xerr
		}
	}

	// Always include the original location
	_, file, line, _ := runtime.Caller(2)
	xerr := &internalError{
		Err:      err,
		Code:     code,
		Message:  message,
		Original: "",
		OrigFile: file,
		OrigLine: line,
		Meta:     map[string]string{},
	}
	fmt.Println("ERROR LOCATION:", xerr.Location())
	return xerr
}

// Error return ErrorInterface
func Error(code int, message string, err error) ErrorInterface {
	return newError(code, message, err)
}

// Error return ErrorInterface with args...
func Errorf(code int, err error, message string, args ...interface{}) ErrorInterface {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return newError(code, message, err)
}

func (t *internalError) GetCode() int {
	return t.Code
}

func (t *internalError) Msg() string {
	return t.Message
}

func (t *internalError) GetMeta(key string) string {
	meta := t.Meta
	if meta != nil {
		return meta[key]
	}
	return ""
}

func (t *internalError) WithMeta(key string, val string) ErrorInterface {
	t.Meta[key] = val
	return t
}

func (t *internalError) MetaMap() map[string]string {
	return t.Meta
}

func (t *internalError) Error() string {
	var b strings.Builder
	b.WriteString(t.Message)
	if t.Err != nil {
		b.WriteString(" cause=")
		b.WriteString(t.Err.Error())
	}
	if t.Original != "" {
		b.WriteString(" origi=")
		b.WriteString(t.Original)
	}
	for k, v := range t.Meta {
		b.WriteByte(' ')
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(v)
	}
	return b.String()
}

func (t *internalError) Location() string {
	return fmt.Sprintf("%v:%v", t.OrigFile, t.OrigLine)
}

func ToErrorInterface(err error) ErrorInterface {
	if err == nil {
		return nil
	}
	if xerr, ok := err.(ErrorInterface); ok {
		return xerr
	}
	xerr, ok := err.(*internalError)
	if !ok {
		xerr = newError(http.StatusInternalServerError, "", err)
	}

	meta := map[string]string{}
	for k, v := range xerr.Meta {
		meta[k] = v
	}
	if xerr.Err != nil {
		xerr.Meta["cause"] = xerr.Err.Error()
	}
	if xerr.Original != "" {
		xerr.Meta["orig"] = xerr.Original
	}

	return xerr
}
