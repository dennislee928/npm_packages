package admin

import (
	redismanagers "bityacht-exchange-api-server/internal/cache/redis/managers"
	"bityacht-exchange-api-server/internal/database/sql/managers"
	"bityacht-exchange-api-server/internal/pkg/email"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	passwordpkg "bityacht-exchange-api-server/internal/pkg/password"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 	Admin 登入
// @Description Admin 登入
// @Tags 		Admin
// @Accept 		json
// @Produce 	json
// @Param 		body body LoginRequest true "Request Body"
// @Success 	200 {object} LoginResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/login [post]
func LoginHandler(ctx *gin.Context) {
	var req LoginRequest
	var resp LoginResponse
	var managerClaims jwt.ManagerClaims

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	record, err := managers.Login(req.Account, req.Password)
	if errpkg.Handler(ctx, err) {
		return
	} else if resp.AccessToken, managerClaims, err = jwt.IssueManagerToken(jwt.ManagerPayload{ManagersRolesID: record.ManagersRolesID, ID: record.ID, Name: record.Name}); errpkg.Handler(ctx, err) {
		return
	}

	var redisManagerReocrd redismanagers.Model
	redisManagerReocrd.Setup(managerClaims.RegisteredClaims.ID, time.Now())
	if err := redismanagers.Login(ctx, record.ID, redisManagerReocrd); errpkg.Handler(ctx, err) {
		return
	}

	resp.NeedChangePassword = record.Extra.NeedChangePassword
	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	Admin 登出
// @Description Admin 登出
// @Tags 		Admin
// @Security	BearerAuth
// @Produce 	json
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/logout [post]
func LogoutHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := redismanagers.Logout(ctx, claims.ManagerPayload.ID, claims.RegisteredClaims.ID); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	忘記密碼
// @Description 忘記密碼
// @Tags 		Admin
// @Accept 		json
// @Produce 	json
// @Param 		body body ForgotPasswordRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/forgot-password [post]
func ForgotPasswordHandler(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}
	record, err := managers.ResetPasswordByEamil(req.Account)
	if errpkg.Handler(ctx, err) {
		return
	}
	if record.ID > 0 {
		passwordMail := email.NewEmail()
		passwordMail.To = []string{req.Account}
		passwordMail.Subject = "BitYacht 兑幣所忘記密碼函"
		passwordMail.Text = []byte(fmt.Sprintf("%s 您好：%s 為您的重設密碼，請於首次登入後變更密碼。", record.Name, record.Password))

		if err = email.SendMail(passwordMail); errpkg.Handler(ctx, err) {
			return
		}
	}
	ctx.Status(http.StatusOK)
}

// @Summary 	Jwt token檢查
// @Description Jwt token檢查
// @Tags 		Admin
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/auth [get]
func AuthHandler(ctx *gin.Context) {
	var resp AuthResponse
	resp.IsValid = false

	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusOK, resp)
		return
	}
	splitedAuthHeader := strings.Split(authHeader, " ")
	if len(splitedAuthHeader) != 2 {
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if splitedAuthHeader[0] != "Bearer" {
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if claims, err := jwt.ValidateManager(splitedAuthHeader[1]); err != nil {
		ctx.JSON(http.StatusOK, resp)
		return
	} else if err := redismanagers.Validate(ctx, claims.ManagerPayload.ID, claims.RegisteredClaims.ID); err != nil {
		ctx.JSON(http.StatusOK, resp)
		return
	} else {
		resp.IsValid = true
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	更新密碼
// @Description 更新密碼
// @Tags 		Admin
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body UpdatePasswordRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/password [PATCH]
func UpdatePasswordHandler(ctx *gin.Context) {
	var req UpdatePasswordRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	} else if err := passwordpkg.StrengthValidate(req.Password); errpkg.Handler(ctx, err) {
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}
	if err = managers.UpdatePassword(claims.ManagerPayload.ID, req.Password); errpkg.Handler(ctx, err) {
		return
	}
	ctx.Status(http.StatusOK)
}
