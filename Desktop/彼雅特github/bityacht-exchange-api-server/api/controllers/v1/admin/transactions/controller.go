package transactions

import (
	"fmt"
	"net/http"
	"time"

	"bityacht-exchange-api-server/internal/database/sql/userstransactions"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得交易訂單列表
// @Description 取得交易訂單列表
// @Tags 		Admin-transactions
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		searcher query modelpkg.Searcher false "searcher"
// @Param		query query userstransactions.GetListRequest false "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]userstransactions.TransactionForManager}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/transactions [get]
func GetListHandler(ctx *gin.Context) {
	var req userstransactions.GetListRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	searcher := modelpkg.GetSearcherFromQuery(ctx)

	var err *errpkg.Error
	if resp.Data, err = userstransactions.GetList(req, &resp.Paginator, searcher); errpkg.Handler(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出交易訂單列表
// @Description 匯出交易訂單列表
// @Tags 		Admin-transactions
// @Security	BearerAuth
// @Accept 		json
// @Param 		query query userstransactions.GetExportRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/transactions/export [get]
func ExportHandler(ctx *gin.Context) {
	var req userstransactions.GetExportRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := userstransactions.GetExport(req)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("transactions_%s.csv", time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, userstransactions.GetTransactionsCSVHeaders(), data)
}
