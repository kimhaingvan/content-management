package httpx

import (
	"content-management/pkg/errorx"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func ParseRequest(r *http.Request, p interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&p)
	if err != io.EOF {
		return err
	}
	return nil
}

func WriteError(ctx context.Context, w http.ResponseWriter, err error) {
	errIn := errorx.ToErrorInterface(err)
	statusCode := errIn.GetCode()
	jsonErr := errorx.ToErrorJSON(errIn)
	errBody, err := json.Marshal(&jsonErr)
	if err != nil {
		errBody = []byte("{\"type\": \"internal\", \"msg\": \"There was an error but it could not be serialized into JSON\"}") // fallback
	}
	w.Header().Set("Content-Type", "application/json") // Error responses are always JSON
	w.Header().Set("Content-Length", strconv.Itoa(len(errBody)))
	w.WriteHeader(statusCode) // set HTTP status code and send response

	_, writeErr := w.Write(errBody)
	if writeErr != nil {
		_ = writeErr
	}
}

func WriteReponse(ctx context.Context, w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Can not marshal response`))
		return
	}
	w.WriteHeader(status)
	w.Write(response)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
