package bank

import (
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得銀行帳戶資訊
// @Description 取得銀行帳戶資訊
// @Tags 		Admin-users
// @Security	BearerAuth
// @Produce		json
// @Param 		id path int true "users id"
// @Success 	200 {object} GetResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/bank [get]
func GetHandler(ctx *gin.Context) {
	var req GetRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp GetResponse
		err  *errpkg.Error
	)

	if resp.BankAccount, err = users.GetBankAccountByID(req.ID); errpkg.Handler(ctx, err) {
		return
	}

	if resp.BankAccount.BanksCode != "" {
		resp.BankInfo, resp.BranchInfo, _ = sqlcache.GetBankAndBranch(resp.BankAccount.BanksCode, resp.BankAccount.BranchsCode)
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	更新銀行帳戶資訊
// @Description 更新銀行帳戶資訊
// @Tags 		Admin-users
// @Security	BearerAuth
// @Accept		json
// @Param 		id path int true "users id"
// @Param 		body body PatchRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/bank [patch]
func PatchHandler(ctx *gin.Context) {
	var req PatchRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	if err := users.UpdateBankAccountStatusByID(req.ID, claims.ID(), req.BankAccountsID, req.Status, req.Comment); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	取得用戶銀行帳戶審核狀態異動紀錄
// @Description 取得用戶銀行帳戶審核狀態異動紀錄
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		id path int true "user id"
// @Success 	200 {object} modelpkg.GetResponse{data=[]usersmodifylogs.BankAccountLog}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/bank/logs [get]
func GetLogHandler(ctx *gin.Context) {
	var req GetRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
		err  *errpkg.Error
	)

	if resp.Data, err = usersmodifylogs.GetBankAccountLogList(&resp.Paginator, req.ID); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出用戶銀行帳戶審核狀態異動紀錄
// @Description 匯出用戶銀行帳戶審核狀態異動紀錄
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Accept 		json
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/bank/logs/export [get]
func ExportLogHandler(ctx *gin.Context) {
	var req GetRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	data, err := usersmodifylogs.GetBankAccountLogExport(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("bank_account_logs_%d.csv", req.ID)
	csv.ExportCSVFile(ctx, filename, usersmodifylogs.GetBankAccountLogHeaders(), data)
}
