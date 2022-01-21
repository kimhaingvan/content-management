package handler

import (
	"content-management/app/order/model"
	"content-management/app/order/service"
	"content-management/pkg/httpreq"
	"content-management/pkg/httpx"
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
	"strconv"

	"go.elastic.co/apm"

	"github.com/qor/qor"

	"github.com/qor/admin"

	"github.com/go-resty/resty/v2"

	"github.com/qor/render"
)

type OrderHandler struct {
	View         *render.Render
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) TestSample(ctx *admin.Context) {
	ctx.Writer.WriteHeader(200)
	client := resty.New()
	req1 := client.R()
	req1.Get("https://www.google.com/")
	h.View.Execute("success", map[string]interface{}{"Order": "dqwdas"}, ctx.Request, ctx.Writer)
}

func (h *OrderHandler) TestError(ctx *admin.Context) {
	test := "data test"
	err := ctx.GetDB().Save(&test).Error
	if err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, errors.New("tesst chơi"))
	}
}

func (h *OrderHandler) GetAllOrder(ctx *admin.Context) {
	var orders []model.Order
	result := ctx.GetDB().Find(&orders)
	if result.Error != nil {
		var err = result.Error.Error()
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, errors.New(err))
		return
	}
	httpx.WriteReponse(context.Background(), ctx.Writer, 200, map[string]interface{}{"Orders": orders})
}

func (h *OrderHandler) GetByIdOrder(ctx *admin.Context) {
	var order model.Order
	id, err := strconv.Atoi(ctx.Request.URL.Query().Get("id"))
	if err != nil {
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, errors.New(err.Error()))
		return
	}

	result := ctx.GetDB().Where("id = ?", id).Find(&order)
	if result.Error != nil {
		var err = result.Error.Error()
		httpx.WriteError(ctx.Request.Context(), ctx.Writer, errors.New(err))
		return
	}
	httpx.WriteReponse(context.Background(), ctx.Writer, 200, map[string]interface{}{"Order": order})
}

func (h *OrderHandler) SaveOrderHandler(order *admin.Resource) func(interface{}, *qor.Context) error {
	return func(result interface{}, ctx *qor.Context) error {
		trx := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", ctx.Request.Method, ctx.Request.RequestURI), "request")
		trxCtx := apm.ContextWithTransaction(ctx.Request.Context(), trx)
		defer trx.End()
		mediaFile, _ := result.(*model.Order)
		if httpreq.HasMediaFile(ctx.Request) {
			_, err := h.orderService.UploadFile(trxCtx, &service.UploadObjectArgs{
				Form:         ctx.Request.MultipartForm,
				MediaStorage: &mediaFile.File,
			})
			if err != nil {
				return err
			}
		}
		if err := ctx.GetDB().Save(result).Error; err != nil {
			return err
		}
		return nil
	}
}

func (h *OrderHandler) DeleteOrderHandler(order *admin.Resource) func(interface{}, *qor.Context) error {
	return func(result interface{}, ctx *qor.Context) error {
		transaction := apm.DefaultTracer.StartTransaction(fmt.Sprintf("%v %v", ctx.Request.Method, ctx.Request.RequestURI), "request")
		defer transaction.End()
		if order.HasPermission(roles.Delete, ctx) {
			if primaryQuerySQL, primaryParams := order.ToPrimaryQueryParams(ctx.ResourceID, ctx); primaryQuerySQL != "" {
				if !ctx.GetDB().First(result, append([]interface{}{primaryQuerySQL}, primaryParams...)...).RecordNotFound() {
					err := ctx.GetDB().Delete(result).Error
					return err
				}
			}
			return gorm.ErrRecordNotFound
		}
		return roles.ErrPermissionDenied
	}

}

func (h *OrderHandler) FindOneHandler(order *admin.Resource) func(interface{}, *resource.MetaValues, *qor.Context) error {
	return func(result interface{}, metaValues *resource.MetaValues, context *qor.Context) error {
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

func (h *OrderHandler) FindManyHandler(order *admin.Resource) func(interface{}, *qor.Context) error {
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
