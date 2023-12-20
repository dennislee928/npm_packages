package risks

import (
	"net/http"

	"bityacht-exchange-api-server/internal/database/sql/risks"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得所有風險設定
// @Description 取得所有風險設定
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Success 	200 {object} []risks.Risk
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs/risks [get]
func GetListHandler(ctx *gin.Context) {
	records, err := risks.Get()
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, records)
}

// @Summary 	新增一個風險設定
// @Description 新增一個風險設定
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body risks.CreateRequest true "Request Body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs/risks [post]
func CreateHandler(ctx *gin.Context) {
	var req risks.CreateRequest
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := risks.Create(req); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary 	修改一個風險設定
// @Description 修改一個風險設定
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body risks.UpdateRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs/risks/{id} [patch]
func UpdateHandler(ctx *gin.Context) {
	var req risks.UpdateRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := risks.Update(req); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	刪除一個風險設定
// @Description 刪除一個風險設定, 會同時刪除所有關聯的風險設定
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs/risks/{id} [delete]
func DeleteHandler(ctx *gin.Context) {
	var req risks.DeleteRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := risks.Delete(req.ID); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}
