package settings

import (
	"bityacht-exchange-api-server/internal/cache/redis/verifications"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/googleauthenticator"
	"bityacht-exchange-api-server/internal/pkg/invoice"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bityacht-exchange-api-server/internal/pkg/receipt"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 	修改密碼
// @Description 修改密碼
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body UpdatePasswordRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/settings/password [patch]
func UpdatePasswordHandler(ctx *gin.Context) {
	var req UpdatePasswordRequest
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)

	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := users.UpdatePassword(true, claims.UserPayload.ID, req.OldPassword, req.NewPassword); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	修改手機載具
// @Description 修改手機載具
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body UpdateMobileBarcodeRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/settings/mobile-barcode [patch]
func UpdateMobileBarcodeHandler(ctx *gin.Context) {
	var req UpdateMobileBarcodeRequest
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)

	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	if req.MobileBarcode != "" {
		if err := invoice.ValidateMobileBarcode(req.MobileBarcode); errpkg.Handler(ctx, err) {
			return
		} else if err = receipt.EZ.CheckMobileCode(ctx, req.MobileBarcode); errpkg.Handler(ctx, err) {
			return
		}
	}

	if err := users.UpdateExtra(claims.UserPayload.ID, &req.MobileBarcode, nil, nil, nil); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

//! Deprecated (Meeting at 2023/10/2)
// // @Summary 	修改兩階段認證類型
// // @Description 修改兩階段認證類型
// // @Tags 		User-MemberCenter
// // @Security	BearerAuth
// // @Accept 		json
// // @Param 		body body UpdateLogin2FATypeRequest true "Request Body"
// // @Success 	200
// // @Failure 	400 {object} errpkg.JsonError
// // @Failure 	500 {object} errpkg.JsonError
// // @Router 		/user/settings/2fa [patch]
// func UpdateLogin2FATypeHandler(ctx *gin.Context) {
// 	var req UpdateLogin2FATypeRequest

// 	if claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx); errpkg.Handler(ctx, err) {
// 		return
// 	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
// 		return
// 	} else if err := users.UpdateExtra(claims.UserPayload.ID, nil, &req.Login2FAType); errpkg.Handler(ctx, err) {
// 		return
// 	}

// 	ctx.Status(http.StatusOK)
// }

const (
	verifyMailExpiration          = 5 * time.Minute
	verificationCodeEnablePrefix  = "enable"
	verificationCodeDisablePrefix = "disable"
)

// @Summary 	發送驗證信 - 更新兩階段認證
// @Description 發送驗證信 - 更新兩階段認證
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept 		json
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/settings/issue-withdraw-2fa-verify [post]
func IssueWithdraw2FAVerifyHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	emailPayload := emailtemplates.UpdateWithdraw2FAPayload{
		VerificationPayload: emailtemplates.VerificationPayload{
			Code:         rand.NumberString(6),
			CodeLifeTime: strconv.FormatInt(int64(verifyMailExpiration/time.Minute), 10),
		},
		Action: "啟用",
	}
	redisCode := verificationCodeEnablePrefix + emailPayload.Code

	userRecord, err := users.GetByID(claims.ID())
	if errpkg.Handler(ctx, err) {
		return
	}
	if userRecord.Extra.IsEnableWithdrawGA2FA() {
		emailPayload.Action = "停用"
		redisCode = verificationCodeDisablePrefix + emailPayload.Code
	}

	if err := verifications.IssueVerificationCode(ctx, jwt.TypeUser, claims.ID(), verifications.UsageUpdateWithdraw2FA, redisCode, verifyMailExpiration); errpkg.Handler(ctx, err) {
		return
	}

	verificationMail := email.NewEmail(email.WithLogo())
	verificationMail.To = []string{claims.Account}
	if verificationMail.Subject, verificationMail.HTML, err = emailtemplates.ExecVerificationUpdateWithdraw2FA(emailPayload); errpkg.Handler(ctx, err) {
		return
	} else if err = email.SendMail(verificationMail); errpkg.Handler(ctx, err) {
		return
	}

	if email.IsDebug() {
		ctx.JSON(http.StatusCreated, gin.H{"verificationCode": emailPayload.Code, "action": emailPayload.Action})
	} else {
		ctx.Status(http.StatusCreated)
	}
}

// @Summary 	取得兩階段認證資訊
// @Description 取得兩階段認證資訊
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept 		json
// @Param 		query query GetWithdraw2FAInfoRequest true "query"
// @Success 	200 {object} GetWithdraw2FAInfoResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/settings/withdraw-2fa-info [get]
func GetWithdraw2FAInfoHandler(ctx *gin.Context) {
	var req GetWithdraw2FAInfoRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp      GetWithdraw2FAInfoResponse
		redisCode = verificationCodeEnablePrefix + req.VerificationCode
	)

	userRecord, err := users.GetByID(claims.ID())
	if errpkg.Handler(ctx, err) {
		return
	}
	if userRecord.Extra.IsEnableWithdrawGA2FA() { // Enable -> Disable
		redisCode = verificationCodeDisablePrefix + req.VerificationCode
	}

	if err := verifications.Verify(ctx, jwt.TypeUser, claims.ID(), verifications.UsageUpdateWithdraw2FA, redisCode); errpkg.Handler(ctx, err) {
		return
	}

	if !userRecord.Extra.IsEnableWithdrawGA2FA() { // Disable -> Enable
		if resp.Secret, resp.QRCode, err = googleauthenticator.GenerateSecret(claims.Account); errpkg.Handler(ctx, err) {
			return
		}
	}

	if resp.Token, _, err = jwt.IssuePreverifyToken(jwt.PreverifyPayload{
		AccountType: jwt.TypeUser,
		ID:          claims.ID(),
		Usage:       jwt.PreverifyUsageUpdateWithdraw2FA,
		GASecret:    resp.Secret,
	}); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	更新兩階段認證
// @Description 更新兩階段認證
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body UpdateWithdraw2FAInfoRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/settings/withdraw-2fa [patch]
func UpdateWithdraw2FAHandler(ctx *gin.Context) {
	var req UpdateWithdraw2FAInfoRequest
	reqTime := time.Now().Unix()

	userClaims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	preverifyClaims, err := jwt.ValidatePreverify(jwt.TypeUser, userClaims.ID(), jwt.PreverifyUsageUpdateWithdraw2FA, req.Token)
	if errpkg.Handler(ctx, err) {
		return
	}

	secret := preverifyClaims.GASecret
	updateWithdraw2FAType := users.TwoFATypeGoogleAuthenticator
	if secret == "" { // Enable -> Disable
		userRecord, err := users.GetByID(userClaims.ID())
		if errpkg.Handler(ctx, err) {
			return
		}

		if userRecord.Extra.Withdraw2FAType&users.TwoFATypeGoogleAuthenticator == 0 {
			errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadAction, errors.New("ga is disable"))
			return
		} else if userRecord.Extra.GoogleAuthenticatorSecret == "" {
			errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeBadRecord, errors.New("ga secret is empty"))
			return
		}
		secret = userRecord.Extra.GoogleAuthenticatorSecret
		updateWithdraw2FAType = -updateWithdraw2FAType
	}

	if err := googleauthenticator.VerifyTOTP(secret, reqTime, req.VerificationCode); errpkg.Handler(ctx, err) {
		return
	}

	if err := users.UpdateExtra(userClaims.ID(), nil, nil, &updateWithdraw2FAType, &preverifyClaims.GASecret); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}
