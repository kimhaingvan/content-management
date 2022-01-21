package handler

import (
	"content-management/app/loanguideline/service"
	"content-management/model"
	"content-management/pkg/httpx"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/k0kubun/pp"

	"github.com/qor/admin"
	"github.com/qor/qor"
	"go.elastic.co/apm"

	"github.com/qor/render"
)

type LoanGuidelineHandler struct {
	View                 *render.Render
	loanGuideLineService service.LoanGuideLineService
}

func NewLoanGuideLineHandler(loanGuideLineService service.LoanGuideLineService) *LoanGuidelineHandler {
	return &LoanGuidelineHandler{
		loanGuideLineService: loanGuideLineService,
	}
}

func (h *LoanGuidelineHandler) SaveLoanGuidelineHandler(loanGuideline *admin.Resource) func(interface{}, *qor.Context) error {
	return func(result interface{}, ctx *qor.Context) error {
		trx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", ctx.Request.Method, ctx.Request.RequestURI), "request")
		//trxCtx := apm.ContextWithTransaction(ctx.Request.Context(), trx)
		defer trx.End()
		mediaFile, _ := result.(*model.LoanGuideline)
		pp.Println(mediaFile)
		//if httpreq.HasMediaFile(ctx.Request) {
		//	_, err := h.orderService.UploadFile(trxCtx, &service.UploadObjectArgs{
		//		Form:         ctx.Request.MultipartForm,
		//		MediaStorage: &mediaFile.File,
		//	})
		//	if err != nil {
		//		return err
		//	}
		//}
		if err := ctx.GetDB().Save(result).Error; err != nil {
			return err
		}
		return nil
	}
}

func (h *LoanGuidelineHandler) GetAllLoanGuidelineHandler(ctx *admin.Context) {
	var data []model.LoanGuideline
	result := ctx.GetDB().Find(&data)
	if result.Error != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, errors.New("no data"))
	}
	httpx.WriteReponse(context.Background(), ctx.Writer, 200, map[string]interface{}{"Orders": data})
}

func (h *LoanGuidelineHandler) GetByIdLoanGuidelineHandler(ctx *admin.Context) {
	var data []model.LoanGuideline
	id, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
	if err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, errors.New("no data"))
		return
	}
	result := ctx.GetDB().Where("id = ?", id).Find(&data)
	if result.Error != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, errors.New("no data"))
	}
	httpx.WriteReponse(context.Background(), ctx.Writer, 200, map[string]interface{}{"Orders": data})
}
