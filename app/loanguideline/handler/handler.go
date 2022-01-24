package handler

import (
	"bytes"
	"content-management/app/loanguideline"
	"content-management/app/loanguideline/service"
	"content-management/model"
	"content-management/pkg/errorx"
	"content-management/pkg/httpx"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/k0kubun/pp"

	"github.com/qor/qor/resource"
	"github.com/qor/roles"

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

func (h *LoanGuidelineHandler) SaveLoanGuidelineHandler(loanGuideline *admin.Resource) func(interface{}, *qor.Context) error {
	return func(result interface{}, ctx *qor.Context) error {
		trx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", ctx.Request.Method, ctx.Request.RequestURI), "request")
		trxCtx := apm.ContextWithTransaction(ctx.Request.Context(), trx)
		defer trx.End()
		l, ok := result.(*model.LoanGuideline)
		if !ok {
			return errorx.ErrInvalidParameter(errors.New("Cannot convert parameters"))
		}
		return h.loanGuideLineService.CreateLoanGuideline(trxCtx, l)
	}
}

func (h *LoanGuidelineHandler) FindOneHandler(order *admin.Resource) func(interface{}, *resource.MetaValues, *qor.Context) error {
	return func(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
		pp.Println(result)
		if order.HasPermission(roles.Read, context) {
			var (
				primaryQuerySQL string
				primaryParams   []interface{}
			)
			if metaValues == nil {
				primaryQuerySQL, primaryParams = order.ToPrimaryQueryParams(context.ResourceID, context)
			} else {
				primaryQuerySQL, primaryParams = order.ToPrimaryQueryParamsFromMetaValue(metaValues, context)
			}

			if primaryQuerySQL != "" {
				if metaValues != nil {
					if destroy := metaValues.Get("_destroy"); destroy != nil {
						if fmt.Sprint(destroy.Value) != "0" && order.HasPermission(roles.Delete, context) {
							context.GetDB().Delete(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...)
							return resource.ErrProcessorSkipLeft
						}
					}
				}
				err := context.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).Error
				return err
			}
			return errors.New("failed to find")
		}
		return roles.ErrPermissionDenied
	}
}

func (h *LoanGuidelineHandler) FindManyHandler(order *admin.Resource) func(interface{}, *qor.Context) error {
	return func(result interface{}, ctx *qor.Context) error {
		transaction := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", ctx.Request.Method, ctx.Request.RequestURI), "request")
		defer transaction.End()
		if order.HasPermission(roles.Read, ctx) {
			db := ctx.GetDB()
			if _, ok := db.Get("qorm:getting_total_count"); ok {
				if err := ctx.GetDB().Count(result).Error; err != nil {
					return err
				}
				return nil
			}
			if err := ctx.GetDB().Set("gorm:order_by_primary_key", "DESC").Find(result).Error; err != nil {
				return err
			}
			return nil
		}
		return roles.ErrPermissionDenied
	}

}
