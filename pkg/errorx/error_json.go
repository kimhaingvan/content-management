package errorx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ErrorJSON sẽ là cấu trúc Error API trả về cho client
type ErrorJSON struct {
	Code string            `json:"code"`
	Msg  string            `json:"msg"`
	Meta map[string]string `json:"meta,omitempty"`
}

func ToErrorJSON(errInterface ErrorInterface) *ErrorJSON {
	return &ErrorJSON{
		Code: fmt.Sprint(errInterface.GetCode()),
		Msg:  errInterface.Msg(),
		Meta: errInterface.MetaMap(),
	}
}

func (e *ErrorJSON) Error() (s string) {
	if len(e.Meta) == 0 {
		return e.Msg
	}
	b := strings.Builder{}
	b.WriteString(e.Msg)
	b.WriteString(" (")
	for _, v := range e.Meta {
		b.WriteString(v)
		break
	}
	b.WriteString(")")
	return b.String()
}

func WriteError(ctx context.Context, resp http.ResponseWriter, err error) {
	errIn := ToErrorInterface(err)
	statusCode := errIn.GetCode()
	jsonErr := ToErrorJSON(errIn)
	respBody, err := json.Marshal(&jsonErr)
	if err != nil {
		respBody = []byte("{\"type\": \"internal\", \"msg\": \"There was an error but it could not be serialized into JSON\"}") // fallback
	}
	resp.Header().Set("Content-Type", "application/json") // Error responses are always JSON
	resp.Header().Set("Content-Length", strconv.Itoa(len(respBody)))
	resp.WriteHeader(statusCode) // set HTTP status code and send response

	_, writeErr := resp.Write(respBody)
	if writeErr != nil {
		_ = writeErr
	}
}
