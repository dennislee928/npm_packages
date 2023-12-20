package user

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/cache/memory/banners"
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	redistoken "bityacht-exchange-api-server/internal/cache/redis/token"
	redisusers "bityacht-exchange-api-server/internal/cache/redis/users"
	"bityacht-exchange-api-server/internal/cache/redis/verifications"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersloginlogs"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	passwordpkg "bityacht-exchange-api-server/internal/pkg/password"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"database/sql"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const login2FAMailExpiration = 5 * time.Minute

// @Summary 	會員登入
// @Description 會員登入
// @Tags 		User-Auth
// @Accept 		json
// @Produce		json
// @Param 		body body LoginRequest true "Request Body"
// @Success 	200 {object} LoginResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/login [post]
func LoginHandler(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	} else if err := redisusers.CheckLoginLock(ctx, req.Account); errpkg.Handler(ctx, err) {
		return
	} else if record, err := users.Login(req.Account, req.Password); errpkg.Handler(ctx, err) {
		redisusers.UpdateAttemptLogin(ctx, req.Account, false)
		return
	} else {
		redisusers.UpdateAttemptLogin(ctx, req.Account, true)
		resp := LoginResponse{ID: record.ID}

		switch record.Status {
		case usersmodifylogs.SLStatusUnverified:
			resp.AccountNotVerified = true
			ctx.JSON(http.StatusOK, resp)
		case usersmodifylogs.SLStatusEnable, usersmodifylogs.SLStatusForzen:
			twoFAType := record.Extra.GetLogin2FAType()
			resp.OnePassKey = uuid.NewString()

			if twoFAType&users.TwoFATypeEmail > 0 {
				resp.Login2FAType |= users.TwoFATypeEmail

				verificationCode := rand.NumberString(6)
				if err = verifications.IssueVerificationCode(ctx, jwt.TypeUser, record.ID, verifications.UsageLogin2FA, resp.OnePassKey+verificationCode, login2FAMailExpiration); errpkg.Handler(ctx, err) {
					return
				} else if email.IsDebug() {
					resp.VerificationCode = verificationCode
				}

				login2FAMail := email.NewEmail(email.WithLogo())
				login2FAMail.To = []string{record.Account}
				if login2FAMail.Subject, login2FAMail.HTML, err = emailtemplates.ExecVerificationLogin2FA(emailtemplates.VerificationPayload{Code: verificationCode, CodeLifeTime: strconv.FormatInt(int64(login2FAMailExpiration/time.Minute), 10)}); errpkg.Handler(ctx, err) {
					return
				}

				if err = email.SendMail(login2FAMail); errpkg.Handler(ctx, err) {
					return
				}
			}

			if resp.Login2FAType == 0 {
				if resp.TwoFactorLoginResponse, err = loginWithUsersRecord(ctx, record); errpkg.Handler(ctx, err) {
					return
				}
				resp.OnePassKey = ""
			}

			ctx.JSON(http.StatusOK, resp)
		case usersmodifylogs.SLStatusDisable:
			ctx.JSON(http.StatusUnauthorized, &errpkg.Error{Code: errpkg.CodeAccountNotAvailable})
		default:
			ctx.JSON(http.StatusInternalServerError, &errpkg.Error{Code: errpkg.CodeBadRecord, Err: errors.New("bad user status")})
		}
	}
}

func loginWithUsersRecord(ctx *gin.Context, record users.Model) (TwoFactorLoginResponse, *errpkg.Error) {
	token, userClaims, err := jwt.IssueUserToken(jwt.NewUserPayload(record, time.Now()))
	if err != nil {
		return TwoFactorLoginResponse{}, err
	}

	var redisUserReocrd redisusers.Model
	redisUserReocrd.Model.Setup(userClaims.RegisteredClaims.ID, userClaims.LoginAt)
	if err = redisusers.Login(ctx, userClaims.UserPayload.ID, redisUserReocrd); err != nil {
		return TwoFactorLoginResponse{}, err
	}

	loginLog, err := usersloginlogs.Create(record.ID, ctx.Request.UserAgent(), ctx.ClientIP(), ctx.Request.Header)
	if err != nil {
		return TwoFactorLoginResponse{}, err
	}

	go func() {
		loginMail := email.NewEmail(email.WithLogo())
		loginMail.To = []string{record.Account}
		if loginMail.Subject, loginMail.HTML, err = emailtemplates.ExecNotifyLogin(emailtemplates.NotifyLoginPayload{Account: record.Account, Device: loginLog.Device, IP: loginLog.IP, Browser: loginLog.Browser, Location: loginLog.Location, Time: emailtemplates.FormatTime(loginLog.CreatedAt)}); errpkg.Handler(ctx, err) {
			return
		}

		if err := email.SendMail(loginMail); errpkg.Handler(ctx, err) {
			logger := logger.GetGinRequestLogger(ctx)
			logger.Err(err.Err).Msg("Send Login Mail Failed")
		}

		// Set RegisterIP if empty
		if record.Extra.RegisterIP == "" {
			if err := users.SetRegisterIP(record.ID, ctx.ClientIP()); err != nil {
				logger := logger.GetGinRequestLogger(ctx)
				logger.Err(err.Err).Msg("SetRegisterIP Failed")
			}
		}
	}()

	return TwoFactorLoginResponse{
		AccessToken:           token,
		RefreshToken:          redisUserReocrd.RefreshToken,
		RefreshTokenExpiredAt: &redisUserReocrd.RefreshTokenExpiredAt,
	}, nil
}

// @Summary 	會員兩階段登入
// @Description 會員兩階段登入
// @Tags 		User-Auth
// @Accept 		json
// @Produce		json
// @Param 		body body TwoFactorLoginRequest true "Request Body"
// @Success 	200 {object} TwoFactorLoginResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/2fa-login [post]
func TwoFactorLoginHandler(ctx *gin.Context) {
	var req TwoFactorLoginRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	} else if err := verifications.Verify(ctx, jwt.TypeUser, req.ID, verifications.UsageLogin2FA, req.OnePassKey+req.VerificationCode); errpkg.Handler(ctx, err) {
		return
	} else if record, err := users.GetByID(req.ID); errpkg.Handler(ctx, err) {
		return
	} else if resp, err := loginWithUsersRecord(ctx, record); errpkg.Handler(ctx, err) {
		return
	} else {
		ctx.JSON(http.StatusOK, resp)
	}
}

const forgotPasswordExpiration = 3 * time.Minute

// @Summary 	忘記密碼
// @Description 忘記密碼
// @Tags 		User-Auth
// @Accept 		json
// @Param 		body body ForgotPasswordRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/forgot-password [post]
func ForgotPasswordHandler(ctx *gin.Context) {
	var req ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	record, err := users.GetByAccount(req.Account)
	if err != nil {
		switch err.Code {
		case errpkg.CodeRecordNotFound:
			ctx.Status(http.StatusOK)
		default:
			errpkg.Handler(ctx, err)
		}

		return
	}

	verificationCode := rand.NumberString(6)
	if err := verifications.IssueVerificationCode(ctx, jwt.TypeUser, record.ID, verifications.UsageForgotPassword, verificationCode, forgotPasswordExpiration); errpkg.Handler(ctx, err) {
		return
	}

	resetPasswordMail := email.NewEmail(email.WithLogo())
	resetPasswordMail.To = []string{req.Account}
	if resetPasswordMail.Subject, resetPasswordMail.HTML, err = emailtemplates.ExecVerificationResetPassword(emailtemplates.VerificationPayload{Code: verificationCode, CodeLifeTime: strconv.FormatInt(int64(forgotPasswordExpiration/time.Minute), 10)}); errpkg.Handler(ctx, err) {
		return
	}

	if err = email.SendMail(resetPasswordMail); errpkg.Handler(ctx, err) {
		return
	}

	if email.IsDebug() {
		ctx.JSON(http.StatusOK, gin.H{"verificationCode": verificationCode})
	} else {
		ctx.Status(http.StatusOK)
	}
}

// @Summary 	驗證重設密碼
// @Description 驗證重設密碼
// @Tags 		User-Auth
// @Accept 		json
// @Produce		json
// @Param 		body body VerifyResetPasswordRequest true "Request Body"
// @Success 	200 {object} VerifyResetPasswordResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/verify-reset-password [post]
func VerifyResetPasswordHandler(ctx *gin.Context) {
	var req VerifyResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	record, err := users.GetByAccount(req.Account)
	if err != nil {
		switch err.Code {
		case errpkg.CodeRecordNotFound:
			ctx.JSON(http.StatusBadRequest, &errpkg.Error{Code: errpkg.CodeBadAuthorizationToken})
		default:
			errpkg.Handler(ctx, err)
		}

		return
	}

	if err := verifications.Verify(ctx, jwt.TypeUser, record.ID, verifications.UsageForgotPassword, req.VerificationCode); errpkg.Handler(ctx, err) {
		return
	}

	token, claims, err := jwt.IssuePreverifyToken(jwt.PreverifyPayload{AccountType: jwt.TypeUser, ID: record.ID, Usage: jwt.PreverifyUsageForgotPassword})
	if errpkg.Handler(ctx, err) {
		return
	} else if err := redistoken.SetForPreverify(ctx, claims); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, VerifyResetPasswordResponse{Token: token})
}

// @Summary 	重設密碼
// @Description 重設密碼
// @Tags 		User-Auth
// @Accept 		json
// @Param 		body body ResetPasswordRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/reset-password [post]
func ResetPasswordHandler(ctx *gin.Context) {
	var req ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	} else if err := passwordpkg.StrengthValidate(req.Password); errpkg.Handler(ctx, err) {
		return
	}

	claims, err := jwt.ValidatePreverify(jwt.TypeUser, 0, jwt.PreverifyUsageForgotPassword, req.Token)
	if errpkg.Handler(ctx, err) {
		return
	} else if err = redistoken.CheckForPreverify(ctx, claims); errpkg.Handler(ctx, err) {
		return
	} else if err = users.UpdatePassword(false, claims.PreverifyPayload.ID, "", req.Password); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	刷新 Token
// @Description 刷新 Token
// @Tags 		User-Auth
// @Accept 		json
// @Produce		json
// @Param 		body body RefreshTokenRequest true "Request Body"
// @Success 	200 {object} RefreshTokenResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/token [post]
func RefreshTokenHandler(ctx *gin.Context) {
	var req RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	record, err := users.GetByID(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	switch record.Status {
	case usersmodifylogs.SLStatusEnable, usersmodifylogs.SLStatusForzen:
		var redisUserReocrd redisusers.Model
		redisUserReocrd.Model.Setup(uuid.NewString(), time.Time{})
		if err = redisusers.Refresh(ctx, record.ID, req.RefreshToken, &redisUserReocrd); errpkg.Handler(ctx, err) {
			return
		}

		token, _, err := jwt.IssueUserToken(jwt.NewUserPayload(record, redisUserReocrd.LoginAt), jwt.WithJWTID(redisUserReocrd.Token))
		if errpkg.Handler(ctx, err) {
			return
		}

		ctx.JSON(http.StatusOK, TwoFactorLoginResponse{
			AccessToken:           token,
			RefreshToken:          redisUserReocrd.RefreshToken,
			RefreshTokenExpiredAt: &redisUserReocrd.RefreshTokenExpiredAt,
		})
	default:
		ctx.JSON(http.StatusUnauthorized, &errpkg.Error{Code: errpkg.CodeAccountNotAvailable})
	}
}

// @Summary 	取得現貨趨勢
// @Description 取得現貨趨勢
// @Tags 		User-Index
// @Produce		json
// @Success 	200 {object} []spottrend.SpotTrend
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/spot-trend [get]
func GetSpotTrendHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, spottrend.GetSpotTrendsForIndexResp())
}

// @Summary 	會員登出
// @Description 會員登出
// @Tags 		User-Auth
// @Security	BearerAuth
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/logout [post]
func LogoutHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err = redisusers.Logout(ctx, claims.UserPayload.ID, claims.RegisteredClaims.ID); err != nil {
		logger := logger.GetGinRequestLogger(ctx)
		logger.Err(err).Msg("redisusers logout error")
	}

	ctx.Status(http.StatusOK)
}

const verifyMailExpiration = 5 * time.Minute

// @Summary 	會員註冊
// @Description 會員註冊
// @Tags 		User-Index
// @Accept 		json
// @Produce		json
// @Param 		body body RegisterRequest true "Request Body"
// @Success 	201 {object} RegisterResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/register [post]
func RegisterHandler(ctx *gin.Context) {
	var req RegisterRequest
	var inviterID int64

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	if req.InviteCode != "" {
		var err *errpkg.Error

		inviterID, err = users.ParseInviteCode(req.InviteCode)
		if errpkg.Handler(ctx, err) {
			return
		}
	}

	record, err := users.CreateNaturalPerson(req.Account, req.Password, inviterID, ctx.ClientIP())
	if errpkg.Handler(ctx, err) {
		return
	}
	resp := RegisterResponse{
		ID: record.ID,
	}

	verificationCode := rand.NumberString(6)
	if err := verifications.IssueVerificationCode(ctx, jwt.TypeUser, record.ID, verifications.UsageVerifyEmail, verificationCode, verifyMailExpiration); errpkg.Handler(ctx, err) {
		return
	} else if email.IsDebug() {
		resp.VerificationCode = verificationCode
	}

	registerMail := email.NewEmail(email.WithLogo())
	registerMail.To = []string{record.Account}
	if registerMail.Subject, registerMail.HTML, err = emailtemplates.ExecVerificationRegister(emailtemplates.VerificationPayload{Code: verificationCode, CodeLifeTime: strconv.FormatInt(int64(verifyMailExpiration/time.Minute), 10)}); errpkg.Handler(ctx, err) {
		return
	} else if err = email.SendMail(registerMail); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// @Summary 	重發 Email 驗證碼
// @Description 重發 Email 驗證碼
// @Tags 		User-Auth
// @Accept 		json
// @Param 		body body ResendEmailVerificationCodeRequest true "Request Body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/resend-verify [post]
func ResendEmailVerificationCodeHandler(ctx *gin.Context) {
	var req ResendEmailVerificationCodeRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	record, err := users.GetByID(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	verificationCode := rand.NumberString(6)
	if err := verifications.IssueVerificationCode(ctx, jwt.TypeUser, record.ID, verifications.UsageVerifyEmail, verificationCode, verifyMailExpiration); errpkg.Handler(ctx, err) {
		return
	}

	verificationMail := email.NewEmail(email.WithLogo())
	verificationMail.To = []string{record.Account}
	if verificationMail.Subject, verificationMail.HTML, err = emailtemplates.ExecVerificationRegister(emailtemplates.VerificationPayload{Code: verificationCode, CodeLifeTime: strconv.FormatInt(int64(verifyMailExpiration/time.Minute), 10)}); errpkg.Handler(ctx, err) {
		return
	} else if err = email.SendMail(verificationMail); errpkg.Handler(ctx, err) {
		return
	}

	if email.IsDebug() {
		ctx.JSON(http.StatusCreated, gin.H{"verificationCode": verificationCode})
	} else {
		ctx.Status(http.StatusCreated)
	}
}

// @Summary 	Email 驗證
// @Description Email 驗證
// @Tags 		User-Auth
// @Accept 		json
// @Param 		body body VerifyEmailRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/verify [post]
func VerifyEmailHandler(ctx *gin.Context) {
	var req VerifyEmailRequest

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	} else if err := verifications.Verify(ctx, jwt.TypeUser, req.ID, verifications.UsageVerifyEmail, req.VerificationCode); errpkg.Handler(ctx, err) {
		return
	}

	userRecord, err := users.EmailVerified(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)

	notifyPayload := emailtemplates.NotifyRegisteredPayload{
		Time: emailtemplates.FormatTime(time.Now()),
		URL:  configs.Config.Email.BitYachtFrontendURL,
	}
	if kycURL, err := url.JoinPath(configs.Config.Email.BitYachtFrontendURL, "Members", strconv.FormatInt(userRecord.ID, 10), "verify"); err != nil {
		errLogger := logger.GetGinRequestLogger(ctx)
		errLogger.Err(err).Msg("get kyc url failed")
	} else {
		notifyPayload.URL = kycURL
	}

	registerSuccessfulMail := email.NewEmail(email.WithLogo())
	registerSuccessfulMail.To = []string{userRecord.Account}
	if registerSuccessfulMail.Subject, registerSuccessfulMail.HTML, err = emailtemplates.ExecNotifyRegistered(notifyPayload); errpkg.Handler(ctx, err) {
		return
	} else if err = email.SendMail(registerSuccessfulMail); errpkg.Handler(ctx, err) {
		return
	}
}

// @Summary 	取得 Banner 列表
// @Description 取得 Banner 列表
// @Tags 		User-Index
// @Produce		json
// @Success 	200 {object} []banners.Banner
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/banners [get]
func GetBannersHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, banners.GetBannersResp())
}

// @Summary 	取得會員資訊
// @Description 取得會員資訊
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Produce		json
// @Success 	200 {object} users.UserInfo
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/info [get]
func GetUserInfoHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	userInfo, err := users.GetUserInfoByID(claims.UserPayload.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	if userInfo.IDVerificationStatus, err = idverifications.GetIDVStatus(userInfo.FinalReview, sql.NullInt64{Int64: userInfo.IDVerificationsID, Valid: true}); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}

// @Summary 	取得會員登入紀錄
// @Description 取得會員登入紀錄
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Produce		json
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]usersloginlogs.Log}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/login-logs [get]
func GetUserLoginLogsHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	if resp.Data, err = usersloginlogs.GetLogsByUser(claims.UserPayload.ID, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
