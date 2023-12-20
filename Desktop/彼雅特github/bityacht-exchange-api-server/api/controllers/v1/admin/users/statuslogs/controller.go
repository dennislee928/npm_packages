package statuslogs

import (
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得該用戶資料狀態異動紀錄
// @Description 取得該用戶資料狀態異動紀錄
// @Tags 		Admin-status-logs
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		id path int true "user id"
// @Success 	200 {object} modelpkg.GetResponse{data=[]usersmodifylogs.StatusLog}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/statuslogs [get]
func GetHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
		err  *errpkg.Error
	)
	if resp.Data, err = usersmodifylogs.GetStatusLogList(&resp.Paginator, req.ID); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出用戶資料狀態異動紀錄
// @Description 匯出用戶資料狀態異動紀錄
// @Tags 		Admin-status-logs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Accept 		json
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/statuslogs/export [get]
func ExportHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	data, err := usersmodifylogs.GetStatusLogExport(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("status_logs_%d.csv", req.ID)
	csv.ExportCSVFile(ctx, filename, usersmodifylogs.GetStatusLogCSVHeaders(), data)
}
