package idv

import (
	"bityacht-exchange-api-server/configs"
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	redistoken "bityacht-exchange-api-server/internal/cache/redis/token"
	"bityacht-exchange-api-server/internal/cache/redis/verifications"
	"bityacht-exchange-api-server/internal/database/sql/duediligences"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/kyc"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bityacht-exchange-api-server/internal/pkg/sms"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func countriesCodeIsTW(countriesCode string) bool {
	return countriesCode == "TWN"
}

// @Summary 	取得身份認證選項
// @Description 取得身份認證選項
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Produce		json
// @Success 	200 {object} sqlcache.IDVOptionsResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/idv/options [get]
func GetOptionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sqlcache.GetUserIDVOptionsResponse())
}

// @Summary 	確認手機是否已被使用
// @Description 確認手機是否已被使用
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept		json
// @Param 		body body CheckPhoneRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/idv/check-phone [post]
func CheckPhoneHandler(ctx *gin.Context) {
	var req CheckPhoneRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := users.CheckPhoneExist(claims.ID(), req.Phone.String()); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	身份認證 - 發送手機驗證碼 & 檢查身份認證表單
// @Description 身份認證 - 發送手機驗證碼 & 檢查身份認證表單
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept		json
// @Param 		body body IssuePhoneVerificationCodeRequest true "Request Body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/idv/issue-phone-verify [post]
func IssuePhoneVerificationCodeHandler(ctx *gin.Context) {
	var req IssuePhoneVerificationCodeRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if _, err := req.CreateIDVerificationRequest.Validate(); errpkg.Handler(ctx, err) {
		return
	} else if err = users.CheckNationalIDExist(claims.ID(), req.CreateIDVerificationRequest.NationalID); errpkg.Handler(ctx, err) {
		return
	} else if err := users.CheckPhoneExist(claims.ID(), req.Phone.String()); errpkg.Handler(ctx, err) {
		return
	}

	userRecord, err := users.GetByID(claims.ID())
	if errpkg.Handler(ctx, err) {
		return
	}

	idvStatus, err := idverifications.GetIDVStatus(userRecord.FinalReview, userRecord.IDVerificationsID)
	if errpkg.Handler(ctx, err) {
		return
	}

	switch idvStatus {
	case users.IDVStatusNone: // Do Nothing
	case users.IDVStatusProcessing, users.IDVStatusIDVRejected:
		errpkg.HandlerWithCode(ctx, http.StatusForbidden, errpkg.CodeBadAction, errors.New("idv is on going"))
		return
	case users.IDVStatusApproved, users.IDVStatusRejected:
		errpkg.HandlerWithCode(ctx, http.StatusForbidden, errpkg.CodeBadAction, errors.New("idv is done"))
		return
	default:
		errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeBadRecord, errors.New("bad idv status"))
		return
	}
	verificationCode := rand.NumberString(6)

	if err = verifications.IssueVerificationCode(ctx, jwt.TypeUser, claims.ID(), verifications.UsageVerifyPhoneForIDV, req.Phone.GetLocalString()+verificationCode, 5*time.Minute); errpkg.Handler(ctx, err) {
		return
	} else if err := sms.Send(ctx, sms.Message{
		Phone:        req.Phone.GetLocalString(),
		Message:      fmt.Sprintf("【BitYacht 兌幣所】 身份驗證，您的手機驗證碼：%s 請回到網站將您的驗證碼填入，以繼續進行身份驗證作業。", verificationCode),
		ReceiverName: fmt.Sprintf("[%d] %s %s", claims.ID(), claims.FirstName, claims.LastName),
	}); errpkg.Handler(ctx, err) {
		return
	}

	if !configs.Config.SMS.Enable {
		ctx.JSON(http.StatusCreated, gin.H{"verificationCode": verificationCode})
	} else {
		ctx.Status(http.StatusCreated)
	}
}

// @Summary 	身份認證 - 驗證手機
// @Description 身份認證 - 驗證手機
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept		json
// @Param 		body body VerifyPhoneRequest true "Request Body"
// @Success 	200 {object} VerifyPhoneResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/idv/verify-phone [post]
func VerifyPhoneHandler(ctx *gin.Context) {
	var req VerifyPhoneRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := verifications.Verify(ctx, jwt.TypeUser, claims.ID(), verifications.UsageVerifyPhoneForIDV, req.Phone.GetLocalString()+req.VerificationCode); errpkg.Handler(ctx, err) {
		return
	}

	token, preverifyClaims, err := jwt.IssuePreverifyToken(jwt.PreverifyPayload{
		AccountType: jwt.TypeUser,
		ID:          claims.ID(),
		Usage:       jwt.PreverifyUsageVerifyPhoneForIDV,
		Phone:       req.Phone.String(),
	})
	if errpkg.Handler(ctx, err) {
		return
	} else if err := redistoken.SetForPreverify(ctx, preverifyClaims); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, VerifyPhoneResponse{PhoneToken: token})
}

// @Summary 	申請身份認證
// @Description 申請身份認證
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		body body CreateIDVerificationRequest true "Request Body"
// @Success 	201 {object} CreateIDVerificationResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/idv [post]
func CreateIDVerificationHandler(ctx *gin.Context) {
	var req CreateIDVerificationRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	userCountry, err := req.Validate()
	if errpkg.Handler(ctx, err) {
		return
	}

	preverifyClaims, err := jwt.ValidatePreverify(jwt.TypeUser, claims.ID(), jwt.PreverifyUsageVerifyPhoneForIDV, req.PhoneToken)
	if errpkg.Handler(ctx, err) {
		return
	} else if err = redistoken.CheckForPreverify(ctx, preverifyClaims); errpkg.Handler(ctx, err) {
		return
	}

	isForeigner := !countriesCodeIsTW(req.CountriesCode)
	idvRecord := req.ToModel(claims.ID(), isForeigner)
	ddRecord := duediligences.Model{UsersID: claims.ID(), Type: duediligences.TypeCreateByIDV}
	userRecordUpdateMap := req.ToUpdateMap(isForeigner, preverifyClaims.Phone)

	if err := duediligences.CreateIDVAndDD(req.NationalID, preverifyClaims.Phone, &idvRecord, &ddRecord, userRecordUpdateMap); errpkg.Handler(ctx, err) {
		return
	} else if isForeigner {
		var err *errpkg.Error
		idvMail := email.NewEmail(email.WithLogo())
		idvMail.To = []string{claims.UserPayload.Account}
		if idvMail.Subject, idvMail.HTML, err = emailtemplates.ExecIDVOnGoing(emailtemplates.IDVOnGoingPayload{Time: emailtemplates.FormatTime(idvRecord.CreatedAt)}); errpkg.Handler(ctx, err) {
			return
		} else if err = email.SendMail(idvMail); errpkg.Handler(ctx, err) {
			return
		}

		ctx.JSON(http.StatusCreated, CreateIDVerificationResponse{})
		return
	}

	// Taiwanese
	kycResp, err := kyc.KryptoGO.InitIDV(ctx, userCountry, claims.ID(), idvRecord.ID, ddRecord.ID, req.LastName+req.FirstName, req.BirthDate, req.NationalID, nil)
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, CreateIDVerificationResponse{IDVerificationURL: kycResp.URL})

	if err := idverifications.UpdateTaskID(idvRecord.ID, strconv.FormatInt(kycResp.IDVTaskID, 10)); err != nil {
		errLogger := logger.GetGinRequestLogger(ctx)
		errLogger.Err(err).Msg("UpdateTaskID error")
	}
}

func checkIDVCanUpdateImage(usersID int64, isForeigner bool) (users.Model, *errpkg.Error) {
	userRecord, err := users.GetByID(usersID)
	if err != nil {
		return users.Model{}, err
	}

	if countriesCode := userRecord.GetCountriesCode(); countriesCode == "" || countriesCodeIsTW(countriesCode) == isForeigner {
		return users.Model{}, &errpkg.Error{HttpStatus: http.StatusForbidden, Code: errpkg.CodeBadAction, Err: errors.New("bad user country")}
	}

	idvStatus, err := idverifications.GetIDVStatus(userRecord.FinalReview, userRecord.IDVerificationsID)
	if err != nil {
		return users.Model{}, err
	}

	if userRecord.GetNationalID() == "" {
		return users.Model{}, &errpkg.Error{HttpStatus: http.StatusForbidden, Code: errpkg.CodeBadAction, Err: errors.New("bad user national id")}
	}

	switch idvStatus {
	case users.IDVStatusIDVRejected: // Do nothing
	case users.IDVStatusNone, users.IDVStatusProcessing, users.IDVStatusRejected, users.IDVStatusApproved:
		return users.Model{}, &errpkg.Error{HttpStatus: http.StatusForbidden, Code: errpkg.CodeBadAction, Err: errors.New("bad idv status")}
	default:
		return users.Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadRecord, Err: errors.New("bad idv status")}
	}

	return userRecord, nil
}

// @Summary 	取得 KryptoGO IDV URL (僅限本國人)
// @Description 取得 KryptoGO IDV URL (僅限本國人)
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Success 	201 {object} CreateIDVerificationResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/idv/krypto-go-url [get]
func GetKryptoGoURLHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	userRecord, err := checkIDVCanUpdateImage(claims.ID(), false)
	if errpkg.Handler(ctx, err) {
		return
	}

	userCountry, err := sqlcache.GetCountry(userRecord.CountriesCode.String)
	if err != nil {
		errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeMemoryCacheError, err.Err)
		return
	}

	userNationalID := userRecord.NationalID.String
	userFullName := userRecord.LastName + userRecord.FirstName
	if userFullName == "" {
		errpkg.HandlerWithCode(ctx, http.StatusForbidden, errpkg.CodeBadAction, errors.New("bad user full name"))
		return
	}

	idvRecord := idverifications.Model{UsersID: claims.ID(), Type: idverifications.TypeKryptoGO}
	ddRecord := duediligences.Model{UsersID: claims.ID(), Type: duediligences.TypeCreateByIDV}
	userRecordUpdateMap := make(map[string]any)

	if err := duediligences.CreateIDVAndDD(userNationalID, "", &idvRecord, &ddRecord, userRecordUpdateMap); errpkg.Handler(ctx, err) {
		return
	}

	kycResp, err := kyc.KryptoGO.InitIDV(ctx, userCountry, claims.ID(), idvRecord.ID, ddRecord.ID, userFullName, modelpkg.NewDate(userRecord.BirthDate), userNationalID, nil)
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusCreated, CreateIDVerificationResponse{IDVerificationURL: kycResp.URL})

	if err := idverifications.UpdateTaskID(idvRecord.ID, strconv.FormatInt(kycResp.IDVTaskID, 10)); err != nil {
		errLogger := logger.GetGinRequestLogger(ctx)
		errLogger.Err(err).Msg("UpdateTaskID error")
	}
}

// @Summary 	更新 IDV 相片 (僅限外國人)
// @Description 更新 IDV 相片 (僅限外國人)
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		body body UpdateIDVImageRequest true "Request Body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/idv/image [patch]
func UpdateIDVImageHandler(ctx *gin.Context) {
	var req UpdateIDVImageRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := req.Validate(); errpkg.Handler(ctx, err) {
		return
	}

	userRecord, err := checkIDVCanUpdateImage(claims.ID(), true)
	if errpkg.Handler(ctx, err) {
		return
	}

	idvRecord := req.ToModel(claims.ID())
	ddRecord := duediligences.Model{UsersID: claims.ID(), Type: duediligences.TypeCreateByIDV}
	userRecordUpdateMap := make(map[string]any)

	if err := duediligences.CreateIDVAndDD(userRecord.GetNationalID(), "", &idvRecord, &ddRecord, userRecordUpdateMap); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusCreated)
}
