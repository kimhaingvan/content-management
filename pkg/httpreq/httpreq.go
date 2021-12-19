package httpreq

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.elastic.co/apm"

	"github.com/go-resty/resty/v2"
)

type RestyConfig struct {
	Client *http.Client
}

type Resty struct {
	resty.Client
}

func NewResty(cfg RestyConfig) *Resty {
	httpClient := &http.Client{} // make a new client
	if cfg.Client != nil {
		*httpClient = *cfg.Client // copy the provided client
	}
	client := &Resty{}
	if cfg.Client == nil {
		client.Client = *resty.New()
	} else {
		client.Client = *resty.NewWithClient(cfg.Client)
	}
	return client
}

func IsNullJsonRaw(data json.RawMessage) bool {
	return len(data) == 0 ||
		len(data) == 4 && string(data) == "null"
}

type SendRequestArgs struct {
	URL                 string
	Req                 interface{}
	Resp                interface{}
	Headers             map[string]string
	QueryParams         map[string]string
	Method              string
	HandleResponseFunc  func(context.Context, *resty.Response, interface{}) error
	ExternalServiceName string
}

func SendRequest(ctx context.Context, args SendRequestArgs) error {
	span, ctx := apm.StartSpan(ctx, fmt.Sprintf("%v %v", args.Method, args.URL), args.ExternalServiceName)
	defer span.End()
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	rcfg := RestyConfig{Client: client}
	rClient := NewResty(rcfg)

	var (
		res *resty.Response
		req *resty.Request
		err error
	)

	req = rClient.R().
		SetHeaders(args.Headers).
		SetQueryParams(args.QueryParams).
		SetBody(args.Req)

	switch args.Method {
	case http.MethodPost:
		res, err = req.Post(args.URL)
	case http.MethodGet:
		res, err = req.Get(args.URL)
	case http.MethodPut:
		res, err = req.Put(args.URL)
	case http.MethodDelete:
		res, err = req.Delete(args.URL)
	default:
		panic(fmt.Sprintf("unsupported method %v", args.Method))
	}
	if err != nil {
		return err
	}
	if args.HandleResponseFunc != nil {
		handleFunc := args.HandleResponseFunc
		return handleFunc(ctx, res, args.Resp)
	}
	return err
}

func HasMediaFile(request *http.Request) bool {
	return request != nil && request.MultipartForm != nil && request.MultipartForm.File["QorResource.File"] != nil
}
