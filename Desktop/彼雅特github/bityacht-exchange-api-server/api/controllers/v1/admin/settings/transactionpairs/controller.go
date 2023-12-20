package transactionpairs

import (
	"net/http"

	"bityacht-exchange-api-server/internal/database/sql/transactionpairs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得交易對列表
// @Description 取得交易對列表
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]transactionpairs.TransactionPair}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/transactionpairs [get]
func GetHandler(ctx *gin.Context) {
	var (
		err  *errpkg.Error
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	)

	if resp.Data, err = transactionpairs.GetList(&resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	編輯交易對
// @Description 編輯交易對
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param 		body body transactionpairs.UpdateRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/transactionpairs [patch]
func UpdateHandler(ctx *gin.Context) {
	var req transactionpairs.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := transactionpairs.Update(ctx, req); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}
