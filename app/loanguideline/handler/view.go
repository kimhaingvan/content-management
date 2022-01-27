package handler

import (
	"content-management/app/loanguideline/service"
	"content-management/model"
	"content-management/pkg/errorx"
	"errors"
	"fmt"
	"github.com/k0kubun/pp"
	"github.com/qor/admin"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/render"
	"github.com/qor/roles"
	"go.elastic.co/apm"
)

type LoanGuidelineViewHandler struct {
	View                 *render.Render
	loanGuideLineService service.LoanGuideLineService
}

func NewLoanGuideLineViewHandler(loanGuideLineService service.LoanGuideLineService) *LoanGuidelineViewHandler {
	return &LoanGuidelineViewHandler{
		loanGuideLineService: loanGuideLineService,
	}
}

func (h *LoanGuidelineViewHandler) SaveLoanGuidelineHandler(loanGuideline *admin.Resource) func(interface{}, *qor.Context) error {
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

func (h *LoanGuidelineViewHandler) FindOneHandler(order *admin.Resource) func(interface{}, *resource.MetaValues, *qor.Context) error {
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

func (h *LoanGuidelineViewHandler) FindManyHandler(order *admin.Resource) func(interface{}, *qor.Context) error {
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
