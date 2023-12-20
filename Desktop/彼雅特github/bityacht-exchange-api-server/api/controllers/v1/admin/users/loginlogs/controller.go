package loginlogs

import (
	"fmt"
	"net/http"
	"time"

	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersloginlogs"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得該用戶登入紀錄
// @Description 取得該用戶登入紀錄
// @Tags 		Admin-login-logs
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		id path int true "user id"
// @Success 	200 {object} modelpkg.GetResponse{data=[]usersloginlogs.Log}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/login-logs [get]
func GetListHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
		err  *errpkg.Error
	)
	if resp.Data, err = usersloginlogs.GetLogsByUser(req.ID, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出該用戶登入紀錄
// @Description 匯出該用戶登入紀錄
// @Tags 		Admin-login-logs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Param 		query query ExportRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/login-logs/export [get]
func GetExportHandler(ctx *gin.Context) {
	var req ExportRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := usersloginlogs.GetExport(req.ID, req.StartAt, req.EndAt)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("login_%s.csv", time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, usersloginlogs.GetLogCSVHeaders(), data)
}
