package handler

import (
	"bytes"
	"content-management/app/loanguideline"
	"content-management/app/loanguideline/service"
	"content-management/pkg/httpx"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/qor/admin"
)

type LoanGuidelineHandler struct {
	loanGuideLineService service.LoanGuideLineService
}

func NewLoanGuideLineHandler(loanGuideLineService service.LoanGuideLineService) *LoanGuidelineHandler {
	return &LoanGuidelineHandler{
		loanGuideLineService: loanGuideLineService,
	}
}

// GetListHandler godoc
// @Summary GetListHandler loan_guideline
// @Description get list handler loan guideline
// @host localhost:8081
// @Accept  json
// @Produce  json
// @param limit offset body loanguideline.GetListRequest true "Request".
// @Success 200 {array} map[string]interface{}
// @Router /admin/loan_guideline/get_list [get]
func (h *LoanGuidelineHandler) GetListHandler(ctx *admin.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	if err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}
	req := new(loanguideline.GetListRequest)

	if err = json.Unmarshal(b, req); err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}

	if err = req.Validate(); err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}

	res, err := h.loanGuideLineService.GetListLoanGuideline(ctx.Request.Context(), req)
	if err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}
	httpx.WriteReponse(context.Background(), ctx.Writer, http.StatusOK, res)
}

// GetHandler godoc
// @Summary GetHandler loan guideline
// @Description Get handler loan guideline
// @host localhost:8081
// @Accept  json
// @Produce  json
// @param id body loanguideline.GetRequest true "Request".
//@Success 200 {array} map[string]interface{}
// @Router /admin/loan_guideline/get [get]
func (h *LoanGuidelineHandler) GetHandler(ctx *admin.Context) {
	b, err := ioutil.ReadAll(ctx.Request.Body)
	ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	if err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}
	req := new(loanguideline.GetRequest)

	if err = json.Unmarshal(b, req); err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}

	if err = req.Validate(); err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}

	res, err := h.loanGuideLineService.GetLoanGuideline(ctx.Request.Context(), req)
	if err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, err)
		return
	}
	httpx.WriteReponse(context.Background(), ctx.Writer, http.StatusOK, res)
}
