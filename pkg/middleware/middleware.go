package middleware

import (
	"content-management/pkg/httpx"
	intLog "content-management/pkg/log"
	"fmt"
	"net/http"
	"time"

	"github.com/openzipkin/zipkin-go"
)

func APILoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := httpx.NewLoggingResponseWriter(w)
		span := zipkin.SpanFromContext(r.Context())
		defer func() {
			intLog.Info(fmt.Sprintf("API infomation: %v [%v]", r.RequestURI, r.Method), span, map[string]interface{}{
				"method": r.Method,
				"path":   r.RequestURI,
				"status": lrw.StatusCode,
			})
		}()
		next.ServeHTTP(lrw, r)
	})
}

func ZipkinMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		span := zipkin.SpanFromContext(r.Context())
		lrw := httpx.NewLoggingResponseWriter(w)
		defer func() {
			span.Finish()
			span.Tag("Method", r.Method)
			span.Tag("URL", r.URL.String())
			span.Tag("Status Code", fmt.Sprintf("%v", lrw.StatusCode))
			span.Tag("Time (Nanoseconds)", fmt.Sprintf("%v", time.Since(start).Nanoseconds()))
		}()
		ctx := zipkin.NewContext(r.Context(), span)
		r = r.WithContext(ctx)
		next.ServeHTTP(lrw, r)
	})
}
