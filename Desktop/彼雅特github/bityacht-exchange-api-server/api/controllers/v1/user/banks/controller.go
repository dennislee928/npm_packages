package banks

import (
	"net/http"

	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得銀行&分行選項
// @Description 取得銀行&分行選項
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Success 	200 {object} sqlcache.BankOptionsResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/banks/options [get]
func GetOptionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sqlcache.GetBankOptionsResponse())
}

// @Summary 	更新銀行帳戶
// @Description 更新銀行帳戶
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Produce		json
// @Param 		body body UpsertAccountRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/banks/account [put]
func UpsertAccountHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	var req UpsertAccountRequest
	if err := ctx.BindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	bankAccountsRecord, err := req.ToModel()
	if errpkg.Handler(ctx, err) {
		return
	}

	if err := users.UpsertBankAccount(claims.ID(), bankAccountsRecord); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	刪除銀行帳戶
// @Description 刪除銀行帳戶
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Produce		json
// @Success 	204
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/banks/account [delete]
func DeleteAccountHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	if err := users.DeleteBankAccount(claims.ID()); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusNoContent)
}
