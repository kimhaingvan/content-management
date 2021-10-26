package controler

import (
	"content-management/pkg/integration/storage/s3/driver"
	"net/http"

	"github.com/k0kubun/pp"
	"github.com/qor/render"
)

type Controller struct {
	View     *render.Render
	S3Driver *driver.S3Driver
}

func (c *Controller) TestFunc(w http.ResponseWriter, req *http.Request) {
	pp.Println("TestFunc")
	//parentSpan := opentracing.SpanFromContext(req.Context())
	//childSpan := zipkin.Tracer.StartSpan(
	//	"child",
	//	opentracing.ChildOf(parentSpan.Context()),
	//)
	//defer childSpan.Finish()
	//ctx := opentracing.ContextWithSpan(req.Context(), childSpan)

	//req.ParseMultipartForm(10 << 20)
	////Get a file from the form input name "file"
	//file, header, err := req.FormFile("file")
	//if err != nil {
	//	http.Error(w, "Something went wrong retrieving the file from the form", http.StatusInternalServerError)
	//	return
	//}
	//defer file.Close()
	//uploadFileArgs := &driver.UploadFileArgs{
	//	File:     file,
	//	FileName: header.Filename,
	//}
	//
	//_, err = c.S3Driver.UploadFile(ctx, uploadFileArgs)
	//if err != nil {
	//	panic(err)
	//}
	c.View.Execute("success", map[string]interface{}{"Order": "dqwdas"}, req, w)
}
