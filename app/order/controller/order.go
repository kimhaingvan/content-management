package controller

import (
	"content-management/pkg/integration/storage/s3/driver"
	"context"
	"net/http"

	"github.com/qor/render"
)

type Controller struct {
	View     *render.Render
	S3Driver *driver.S3Driver
}

// Cart shopping cart
func (c *Controller) ExtraFunc(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(10 << 20)

	// Get a file from the form input name "file"
	file, header, err := req.FormFile("file")
	if err != nil {
		http.Error(w, "Something went wrong retrieving the file from the form", http.StatusInternalServerError)
		return
	}
	defer file.Close()
	uploadFileArgs := &driver.UploadFileArgs{
		File:     file,
		FileName: header.Filename,
	}

	_, err = c.S3Driver.UploadFile(context.Background(), uploadFileArgs)
	if err != nil {
		panic(err)
	}
	c.View.Execute("success", map[string]interface{}{"Order": "dqwdas"}, req, w)
}
