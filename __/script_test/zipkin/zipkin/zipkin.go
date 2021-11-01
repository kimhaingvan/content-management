package zipkin

import (
	"content-management/core/config"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/openzipkin/zipkin-go/reporter"

	"github.com/openzipkin/zipkin-go"
	reporterhttp "github.com/openzipkin/zipkin-go/reporter/http"
)

func NewReporter(zipkinURL string) reporter.Reporter {
	// Create reporter to send data to zipkin
	endpointURL := fmt.Sprintf("%v/api/v2/spans", zipkinURL)
	reporter := reporterhttp.NewReporter(endpointURL)
	return reporter
}

func NewTracer(reporter reporter.Reporter) (*zipkin.Tracer, error) {
	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint(config.GetAppConfig().ApplicationName, "localhost:"+"8080")
	if err != nil {
		return nil, errors.New("Can not register zipkin tracer")
	}

	// Sampler tells you which traces are going to be sampled or not. In this case we will record 100% (1.00) of traces.
	sampler, err := zipkin.NewCountingSampler(1)
	if err != nil {
		return nil, err
	}

	// initialize our tracer
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	if err != nil {
		return nil, err
	}
	return tracer, nil
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
