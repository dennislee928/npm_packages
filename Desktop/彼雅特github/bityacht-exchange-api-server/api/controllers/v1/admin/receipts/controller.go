package receipts

import (
	"bityacht-exchange-api-server/internal/database/sql/receipts"
	"bityacht-exchange-api-server/internal/service"
	"fmt"
	"net/http"
	"time"

	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/sourcegraph/conc/iter"
)

// @Summary 	取得所有發票列表
// @Description 取得所有發票列表
// @Tags 		Admin-Receipt
// @Produce		json
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		searcher query modelpkg.Searcher false "searcher"
// @Param		query query receipts.GetListRequest false "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]receipts.ListItem}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/receipts [get]
func GetListHandler(ctx *gin.Context) {
	var req receipts.GetListRequest

	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	searcher := modelpkg.GetSearcherFromQuery(ctx)

	var err *errpkg.Error
	if resp.Data, err = receipts.GetList(&resp.Paginator, req, searcher); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	發票開立
// @Description 發票開立
// @Tags 		Admin-Receipt
// @Accept 		json
// @Produce		json
// @Param 		body body IssueRequest true "Request Body"
// @Success 	202
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router		/admin/receipts/issue [post]
func IssueHandler(ctx *gin.Context) {
	var req IssueRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	logger := logger.GetGinRequestLogger(ctx)
	go func() {
		iter.ForEach(req.IDs, func(id *string) {
			if wrapErr := service.IssueReceipt(*id); wrapErr != nil && wrapErr.Err != nil {
				logger.Err(wrapErr.Err).Str("id", *id).Msg("issue receipt failed")
				return
			}
			logger.Info().Str("id", *id).Msg("issue receipt success")
		})
	}()

	ctx.Status(http.StatusAccepted)
}

// @Summary 	匯出發票列表 (CSV)
// @Description 匯出發票列表 (CSV)
// @Tags 		Admin-Receipt
// @Accept 		json
// @Param 		query query ExportRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router		/admin/receipts/export [get]
func ExportCSVHandler(ctx *gin.Context) {
	var req receipts.ExportRequest

	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := receipts.GetExport(req)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("invoiceList_%s.csv", time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, receipts.GetDetailItemCSVHeaders(), data)
}

// @Summary 	取得發票詳細資料
// @Description 取得發票詳細資料
// @Tags 		Admin-Receipt
// @Produce		json
// @Param 		ID path string true "Receipt ID"
// @Success 	200 {object} receipts.DetailItem
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router		/admin/receipts/:ID [get]
func GetDetailHandler(ctx *gin.Context) {
	var req ViewRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	record, wrapErr := receipts.GetDetailByID(req.ID)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	ctx.JSON(http.StatusOK, record)
}
