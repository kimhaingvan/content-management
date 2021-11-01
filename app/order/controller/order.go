package controller

import (
	"time"

	"github.com/qor/admin"

	"github.com/go-resty/resty/v2"

	"github.com/qor/render"
)

type Controller struct {
	View *render.Render
}

func (c *Controller) TestFunc(ctx *admin.Context) {
	ctx.Writer.WriteHeader(200)
	time.Sleep(1 * time.Second)
	//span := opentracing.SpanFromContext(req.Context())
	//childSpan := zipkin.Tracer.StartSpan(
	//	"child",
	//	opentracing.ChildOf(span.Context()),
	//)
	//defer childSpan.Finish()
	client := resty.New()
	req1 := client.R()
	_, err := req1.Get("https://www.google.com/")
	if err != nil {
	}
	c.View.Execute("success", map[string]interface{}{"Order": "dqwdas"}, ctx.Request, ctx.Writer)
}
