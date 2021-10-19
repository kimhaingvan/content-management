package main

import (
	intzipkin "content-management/__/script_test/zipkin/zipkin"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/openzipkin/zipkin-go"

	"github.com/gorilla/mux"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
)

var zipkinClient *zipkinhttp.Client

func main() {
	reporter := intzipkin.NewReporter(os.Getenv("ZIPKIN_URL"))
	defer reporter.Close()

	tracer, err := intzipkin.NewTracer(reporter)
	// create global zipkin http server middleware
	serverMiddleware := zipkinhttp.NewServerMiddleware(
		tracer, zipkinhttp.TagResponseSize(true),
	)

	// create global zipkin traced http client
	zipkinClient, err = zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/abc", HomeHandlerFactory(zipkinClient))
	r.HandleFunc("/foo", FooHandler)
	r.Use(serverMiddleware) // name for request span

	log.Fatal(http.ListenAndServe(":8080", r))
}

func FooHandler(w http.ResponseWriter, r *http.Request) {
	span := zipkin.SpanFromContext(r.Context())
	ctx := zipkin.NewContext(context.Background(), span)
	zipkinCtx := intzipkin.SetSpanNameByRoute(ctx, r)
	n, err := http.NewRequestWithContext(zipkinCtx, "POST", "http://example.com", nil)
	res, err := zipkinClient.DoWithAppSpan(n, "HomeHandlerFactory")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res.StatusCode > 399 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	res.Body.Close()
}

func HomeHandlerFactory(client *zipkinhttp.Client) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		span := zipkin.SpanFromContext(r.Context())
		ctx := zipkin.NewContext(context.Background(), span)
		ctx = SetSpanNameByRoute(ctx, r)
		n, err := http.NewRequestWithContext(ctx, "POST", "http://example.com", nil)
		res, err := client.DoWithAppSpan(n, "HomeHandlerFactory")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if res.StatusCode > 399 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		res.Body.Close()
	}
}

func SetSpanNameByRoute(ctx context.Context, r *http.Request) context.Context {
	if span := zipkin.SpanFromContext(ctx); span != nil {
		if route := mux.CurrentRoute(r); route != nil {
			if routePath, err := route.GetPathTemplate(); err == nil {
				zipkin.TagHTTPRoute.Set(span, routePath)
				span.SetName(r.Method + " " + routePath)
			}
		}
	}
	return ctx
}
