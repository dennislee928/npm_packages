package userskycs

import (
	"net/http"

	"bityacht-exchange-api-server/internal/database/sql/risks"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersrisks"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得該使用者所有風險索引
// @Description 取得該使用者所有風險索引
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "users id"
// @Success 	200 {object} []int64
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/risks [get]
func GetHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}
	if resp, err := usersrisks.GetRisksIDsByUserID(req.ID); errpkg.Handler(ctx, err) {
		return
	} else {
		ctx.JSON(http.StatusOK, resp)
	}
}

// @Summary 	更新該使用者所有風險索引
// @Description 更新該使用者所有風險索引, 會先刪除所有風險索引, 再新增
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Param 		body body risks.UpdateRisksRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/risks [post]
func UpdateHandler(ctx *gin.Context) {
	var req risks.UpdateRisksRequest
	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}
	if err := risks.UpdateRisks(req); errpkg.Handler(ctx, err) {
		return
	}
	ctx.Status(http.StatusOK)
}
