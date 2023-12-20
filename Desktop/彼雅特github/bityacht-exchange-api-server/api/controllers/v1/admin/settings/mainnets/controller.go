package mainnets

import (
	"bityacht-exchange-api-server/internal/database/sql/mainnets"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得主網列表
// @Description 取得主網列表
// @Tags 		Admin-settings
// @Accept 		json
// @Produce 	json
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]mainnets.Mainnet}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/mainnets [get]
func GetHandler(ctx *gin.Context) {
	var (
		err  *errpkg.Error
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	)

	if resp.Data, err = mainnets.GetMainnetList(&resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	編輯主網
// @Description 編輯主網
// @Tags 		Admin-settings
// @Accept 		json
// @Produce 	json
// @Param 		currency path int true "幣種"
// @Param 		mainnet path int true "主網"
// @Param 		body body UpdateRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/mainnets/{currency}/{mainnet} [patch]
func UpdateHandler(ctx *gin.Context) {
	var req mainnets.UpdateRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	if err := mainnets.UpdateMainnet(req); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}
