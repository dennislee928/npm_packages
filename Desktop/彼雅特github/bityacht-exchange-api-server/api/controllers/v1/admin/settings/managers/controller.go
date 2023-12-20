package managers

import (
	"errors"
	"fmt"
	"net/http"

	redismanagers "bityacht-exchange-api-server/internal/cache/redis/managers"
	"bityacht-exchange-api-server/internal/database/sql/managers"
	"bityacht-exchange-api-server/internal/pkg/email"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

// @Summary 	建立管理員帳號
// @Description 建立管理員帳號
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body CreateRequest true "Request Body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/managers [post]
func CreateHandler(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	record, err := managers.Create(ctx, req.Account, req.Name, req.ManagersRolesID)
	if errpkg.Handler(ctx, err) {
		return
	}

	passwordMail := email.NewEmail()
	passwordMail.To = []string{record.Account}
	passwordMail.Subject = "BitYacht 兑幣所管理者密碼函"
	passwordMail.Text = []byte(fmt.Sprintf("管理者 %s 您好：%s 為您的預設密碼，請於首次登入後變更密碼。", record.Name, record.Password))

	if err = email.SendMail(passwordMail); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary 	取得管理員列表
// @Description 取得管理員列表
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]managers.Manager}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/managers [get]
func GetHandler(ctx *gin.Context) {
	var (
		err  *errpkg.Error
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	)

	if resp.Data, err = managers.GetManagerList(&resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	編輯管理員
// @Description 編輯管理員
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param 		id path int true "managers id"
// @Param 		body body managers.UpdateRequest true "Request Body, only admin can modify other manager's account."
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/managers/{id} [patch]
func UpdateHandler(ctx *gin.Context) {
	var req managers.UpdateRequest

	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if req.Account == "" && req.Name == "" && req.Password == "" && req.ManagersRolesID == 0 {
		ctx.JSON(http.StatusBadRequest, errpkg.Error{Code: errpkg.CodeBadBody, Err: errors.New("body cannot be empty")})
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	record, err := managers.Update(ctx, req)
	if errpkg.Handler(ctx, err) {
		return
	}

	if record.Password != "" && !(req.ID == claims.ManagerPayload.ID) { // Send New Password to Manager
		passwordMail := email.NewEmail()
		passwordMail.To = []string{record.Account}
		passwordMail.Subject = "BitYacht 兑幣所管理者密碼異動函"
		passwordMail.Text = []byte(fmt.Sprintf("管理者 %s 您好：超級管理者已將您的密碼設定為%s。", record.Name, record.Password))

		if err = email.SendMail(passwordMail); errpkg.Handler(ctx, err) {
			return
		}
	}
	if record.Status == managers.StatusDisable {
		if err := redismanagers.ForceLogout(ctx, req.ID); err != nil {
			errLogger := logger.GetGinRequestLogger(ctx)
			errLogger.Err(err).Msg("force logout failed")
		}
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	刪除管理員
// @Description 刪除管理員
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param 		id path int true "managers id"
// @Success 	204
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/managers/{id} [delete]
func DeleteHandler(ctx *gin.Context) {
	var req DeleteRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if req.ID == claims.ManagerPayload.ID { // Cannot Delete Self
		ctx.JSON(http.StatusBadRequest, errpkg.Error{Code: errpkg.CodeBadParam, Err: errors.New("cannot delete self")})
		return
	} else if err = managers.Delete(req.ID); errpkg.Handler(ctx, err) {
		return
	} else if err = redismanagers.ForceLogout(ctx, req.ID); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusNoContent)
}
