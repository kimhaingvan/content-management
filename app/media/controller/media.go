package controller

import (
	"time"

	"github.com/qor/admin"

	"github.com/k0kubun/pp"

	"github.com/qor/render"
)

type Controller struct {
	View *render.Render
}

func (c *Controller) GetMediaFile(ctx *admin.Context) {
	pp.Println("GetMedia")
	ctx.Writer.WriteHeader(200)
	time.Sleep(1 * time.Second)
}
