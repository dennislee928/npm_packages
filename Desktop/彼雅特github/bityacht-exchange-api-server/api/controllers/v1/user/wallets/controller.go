package wallets

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/cache/redis/verifications"
	"bityacht-exchange-api-server/internal/database/sql/suspicioustransactions"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	"bityacht-exchange-api-server/internal/database/sql/userswithdrawalwhitelist"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bityacht-exchange-api-server/internal/pkg/mmdb"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bityacht-exchange-api-server/internal/pkg/sms"
	suspicioustransactionspkg "bityacht-exchange-api-server/internal/pkg/suspicioustransactions"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

const withdraw2FAMailExpiration = 5 * time.Minute

// @Summary 	用戶提幣
// @Description 用戶提幣
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		body body WithdrawRequest true "body"
// @Success 	200 {object} WithdrawResponse  "2FA email sent"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/withdraw [post]
func WithdrawHandler(ctx *gin.Context) {
	var req WithdrawRequest

	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	if usersRecord, wrapErr := users.GetByID(claims.ID()); errpkg.Handler(ctx, wrapErr) {
		return
	} else if wrapErr = usersRecord.AllowWithdraw(); errpkg.Handler(ctx, wrapErr) {
		return
	}

	whitelist, err := userswithdrawalwhitelist.GetByIDAndUser(req.WhitelistID, claims.ID())
	if errpkg.Handler(ctx, err) {
		return
	}

	amount := req.Amount
	if amount.IsNegative() {
		errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadAmount, fmt.Errorf("amount must be positive: %s", amount))
		return
	}

	param := withdrawParam{
		claims:   claims,
		currency: req.CurrencyType,
		mainnet:  req.Mainnet,
		address:  whitelist.Address,
		amount:   amount,
	}

	// 2FA check
	param.userRecord, wrapErr = users.GetByID(claims.ID())
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if issueWithdraw2FA(ctx, param) {
		return
	}

	withdraw(ctx, param)
}

// @Summary 	用戶 2FA 提幣
// @Description 用戶 2FA 提幣
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		body body Withdraw2FARequest true "body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/2fa-withdraw [post]
func TwoFactorWithdrawHandler(ctx *gin.Context) {
	var req Withdraw2FARequest
	reqTime := time.Now()

	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	data, wrapErr := verifications.GetWithdrawVerification(ctx, reqTime.Unix(), claims.ID(), req.OnePassKey, req.EmailVerificationCode, req.SMSVerificationCode, req.GAVerificationCode)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	amount, err := decimal.NewFromString(data.Amount)
	if err != nil {
		errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeBadAmount, fmt.Errorf("invalid amount: %s", data.Amount))
		return
	}

	withdraw(ctx, withdrawParam{
		claims:   claims,
		currency: data.Currency,
		mainnet:  data.Mainnet,
		address:  data.Address,
		amount:   amount,
	})
}

type withdrawParam struct {
	claims     *jwt.UserClaims
	currency   wallet.CurrencyType
	mainnet    wallet.Mainnet
	address    string
	amount     decimal.Decimal
	userRecord users.Model
}

// issueWithdraw2FA will send the response to client if 2FA is enable or error occurred, and return true.
func issueWithdraw2FA(ctx *gin.Context, param withdrawParam) bool {
	var (
		resp        WithdrawResponse
		wrapErr     *errpkg.Error
		verifyParam = verifications.Withdraw{
			ID:       param.claims.ID(),
			Currency: param.currency,
			Mainnet:  param.mainnet,
			Address:  param.address,
			Amount:   param.amount.String(),
		}
	)

	twoFAType := param.userRecord.Extra.GetWithdraw2FAType()
	if twoFAType&users.TwoFATypeEmail > 0 {
		resp.TwoFAType |= users.TwoFATypeEmail
		verifyParam.EmailCode = rand.NumberString(6)

		withdraw2FAMail := email.NewEmail(email.WithLogo())
		withdraw2FAMail.To = []string{param.claims.Account}
		if withdraw2FAMail.Subject, withdraw2FAMail.HTML, wrapErr = emailtemplates.ExecVerificationWithdraw2FA(emailtemplates.VerificationPayload{Code: verifyParam.EmailCode, CodeLifeTime: strconv.FormatInt(int64(withdraw2FAMailExpiration/time.Minute), 10)}); errpkg.Handler(ctx, wrapErr) {
			return true
		}

		if wrapErr = email.SendMail(withdraw2FAMail); errpkg.Handler(ctx, wrapErr) {
			return true
		}
	}
	if twoFAType&users.TwoFATypeSMS > 0 {
		resp.TwoFAType |= users.TwoFATypeSMS
		verifyParam.SMSCode = rand.NumberString(6)
		userPhone := modelpkg.TWCellPhone(param.userRecord.Phone)

		if err := sms.Send(ctx, sms.Message{
			Phone:        userPhone.GetLocalString(),
			Message:      fmt.Sprintf("【BitYacht 兌幣所】 提幣驗證，您的手機驗證碼：%s 請回到網站將您的驗證碼填入，以繼續進行提幣作業。", verifyParam.SMSCode),
			ReceiverName: fmt.Sprintf("[%d] %s %s", param.userRecord.ID, param.userRecord.FirstName, param.userRecord.LastName),
		}); errpkg.Handler(ctx, err) {
			return true
		}
	}
	if param.userRecord.Extra.IsEnableWithdrawGA2FA() {
		resp.TwoFAType |= users.TwoFATypeGoogleAuthenticator
		verifyParam.GASecret = param.userRecord.Extra.GoogleAuthenticatorSecret
	}

	if resp.TwoFAType == 0 {
		return false
	}

	resp.OnePassKey, wrapErr = verifications.IssueWithdrawVerification(ctx.Request.Context(), verifyParam, withdraw2FAMailExpiration)
	if errpkg.Handler(ctx, wrapErr) {
		return true
	}

	if email.IsDebug() {
		resp.EmailVerificationCode = verifyParam.EmailCode
	}
	if !configs.Config.SMS.Enable {
		resp.SMSVerificationCode = verifyParam.SMSCode
	}

	ctx.JSON(http.StatusOK, resp)
	return true
}

func withdraw(ctx *gin.Context, param withdrawParam) {
	userID := param.claims.ID()

	currencyInfo, wrapErr := spottrend.GetCurrencyInfo(param.currency.String())
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	mainnetRecord, wrapErr := sqlcache.GetMainnetRecord(param.currency.String(), param.mainnet.BinanceNetwork())
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	ip := ctx.ClientIP()
	location, err := mmdb.LookupCity(ip)
	if errpkg.Handler(ctx, err) {
		return
	}

	if wrapErr := userswallets.LockAmount(userID, param.currency.String(), param.amount); errpkg.Handler(ctx, wrapErr) {
		return
	}

	defer func() {
		if wrapErr != nil {
			if wrapErr := userswallets.UnlockAmount(userID, param.currency.String(), param.amount); wrapErr != nil {
				logger.Logger.Err(wrapErr).Msg("failed to unlock amount")
			}
		}
	}()

	amountDeductFee := param.amount.Sub(mainnetRecord.WithdrawFee)
	spotTransfer := &usersspottransfers.Model{
		Type:             usersspottransfers.TypeCybavoAPI,
		UsersID:          userID,
		CurrenciesSymbol: param.currency.String(),
		Mainnet:          param.mainnet.String(),
		ToAddress:        param.address,
		Status:           usersspottransfers.StatusProcessing,
		Action:           usersspottransfers.ActionWithdraw,
		Amount:           amountDeductFee,
		Valuation:        param.amount.Mul(currencyInfo.OriToUSDTPrice),
		HandlingCharge:   mainnetRecord.WithdrawFee,
		Extra: usersspottransfers.Extra{
			ToUSDTPrice: currencyInfo.OriToUSDTPrice,
			ToTWDPrice:  currencyInfo.OriToTWDPrice,
			IP:          ip,
			Location:    location,
		},
	}

	if spotTransfer.Valuation.GreaterThanOrEqual(decimal.New(1000, 0)) { // Aegis Manual
		spotTransfer.Type = usersspottransfers.TypeAegisManual
		spotTransfer.Status = usersspottransfers.StatusReviewing

		if wrapErr := usersspottransfers.Create(spotTransfer); errpkg.Handler(ctx, wrapErr) {
			return
		}
	} else { // Cybavo API
		if wrapErr = usersspottransfers.Create(spotTransfer); errpkg.Handler(ctx, wrapErr) {
			return
		}

		if wrapErr = wallet.Cybavo.WithdrawAssets(
			ctx.Request.Context(),
			param.currency,
			param.mainnet,
			wallet.WithdrawPayload{
				Requests: []wallet.WithdrawReq{
					{
						OrderID: spotTransfer.TransfersID,
						Address: param.address,
						Amount:  amountDeductFee.String(),
						UserID:  strconv.FormatInt(userID, 10),
					},
				},
			}); errpkg.Handler(ctx, wrapErr) {
			spotTransfer.Status = usersspottransfers.StatusFailed
			if _, updateErr := usersspottransfers.TakeOldAndUpsert(spotTransfer); updateErr != nil {
				logger.Logger.Err(updateErr).Msg("failed to update status to failed")
			}
			return
		}
	}

	ctx.Status(http.StatusCreated)

	withdraw2FAMail := email.NewEmail(email.WithLogo())
	withdraw2FAMail.To = []string{param.claims.Account}
	if withdraw2FAMail.Subject, withdraw2FAMail.HTML, wrapErr = emailtemplates.ExecNotifyWithdrawSpot(emailtemplates.NotifyWithdrawSpotPayload{
		Time:             emailtemplates.FormatTime(time.Now()),
		CurrenciesSymbol: param.currency.String(),
		Mainnet:          param.mainnet.String(),
		Address:          param.address,
		Amount:           amountDeductFee.String(),
	}); wrapErr != nil {
		logger.Logger.Err(wrapErr).Msg("failed to execute template")
		return
	}

	if wrapErr = email.SendMail(withdraw2FAMail); wrapErr != nil {
		logger.Logger.Err(wrapErr).Msg("failed to send mail")
	}

	detectResults, wrapErr := suspicioustransactionspkg.Detect(suspicioustransactionspkg.ActionWithdrawCryptocurrency, spotTransfer)
	if wrapErr != nil {
		logger.Logger.Err(wrapErr).Int64("users id", userID).Str("transfers id", spotTransfer.TransfersID).Msg("failed to detect suspicious transactions")
		return
	}

	if wrapErr := suspicioustransactions.CreateFromResults(userID, spotTransfer.TransfersID, detectResults); wrapErr != nil {
		logger.Logger.Err(wrapErr).Int64("users id", userID).Str("transfers id", spotTransfer.TransfersID).Any("results", detectResults).Msg("failed to create suspicious records")
		return
	}
}

// @Summary 	取得提幣資訊
// @Description 取得提幣資訊
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		query query GetDepositAddressRequest true "query"
// @Success 	200 {object} GetCoinInfoResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/withdraw-info [get]
func GetWithdrawInfoHandler(ctx *gin.Context) {
	var req GetDepositAddressRequest

	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	mainnetInfo, err := sqlcache.GetMainnetRecord(req.CurrencyType.String(), req.Mainnet.BinanceNetwork())
	if errpkg.Handler(ctx, err) {
		return
	}

	levelLimit, wrapErr := sqlcache.GetLevelLimit(claims.UserPayload.Type, claims.UserPayload.Level)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	accWithdrawValuation, wrapErr := usersspottransfers.GetAccWithdrawValuationByUser(nil, claims.ID())
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	ctx.JSON(http.StatusOK, GetCoinInfoResponse{
		Name:                 mainnetInfo.Name,
		Network:              mainnetInfo.Mainnet,
		WithdrawFee:          mainnetInfo.WithdrawFee.String(),
		WithdrawMax:          mainnetInfo.WithdrawMax.String(),
		WithdrawMin:          mainnetInfo.WithdrawMin.String(),
		AccWithdrawValuation: accWithdrawValuation,
		LevelLimit:           levelLimit.USDT.Withdraw,
	})
}

// @Summary 	取得提幣白名單
// @Description 取得提幣白名單
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Produce		json
// @Param 		query query GetWithdrawalWhitelistRequest true "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]userswithdrawalwhitelist.Record}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/withdrawal-whitelist [get]
func GetWithdrawalWhitelistHandler(ctx *gin.Context) {
	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	var req GetWithdrawalWhitelistRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	var resp modelpkg.GetResponse
	if resp.Data, wrapErr = userswithdrawalwhitelist.GetRecordsByUserAndMainnet(claims.ID(), req.Mainnet); errpkg.Handler(ctx, wrapErr) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	提幣白名單 - 新增
// @Description 提幣白名單 - 新增
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Accept		json
// @Param 		body body CreateWithdrawalWhitelistRequest true "body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/withdrawal-whitelist [post]
func CreateWithdrawalWhitelistHandler(ctx *gin.Context) {
	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	var req CreateWithdrawalWhitelistRequest
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	recordCount, wrapErr := userswithdrawalwhitelist.CountByUserAndMainnet(claims.ID(), req.Mainnet)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if recordCount >= 10 {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusForbidden, Code: errpkg.CodeOverWithdrawalWhitelistLimit})
		return
	}

	record := req.ToModel(claims.ID())
	if wrapErr := wallet.Cybavo.AddWithdrawalWhitelistEntry(ctx, claims.ID(), req.Mainnet, req.Address); errpkg.Handler(ctx, wrapErr) {
		return
	}

	if wrapErr := userswithdrawalwhitelist.Create(&record); errpkg.Handler(ctx, wrapErr) {
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary 	提幣白名單 - 刪除
// @Description 提幣白名單 - 刪除
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		id path int true "ID in withdrawal whitelist"
// @Success 	204
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/withdrawal-whitelist/{id} [delete]
func DeleteWithdrawalWhitelistHandler(ctx *gin.Context) {
	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	var req DeleteWithdrawalWhitelistRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	record, wrapErr := userswithdrawalwhitelist.Delete(req.ID, claims.ID())
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	ctx.Status(http.StatusNoContent)

	if record != nil {
		errLogger := logger.GetGinRequestLogger(ctx)

		m, ok := wallet.ParseMainnet(record.Mainnet)
		if !ok {
			errLogger.Error().Any("record", *record).Msg("parse mainnet failed")
			return
		}

		if err := wallet.Cybavo.RemoveWithdrawalWhitelistEntry(context.Background(), claims.ID(), m, record.Address); err != nil {
			errLogger.Err(err.Err).Any("record", *record).Msg("remove withdrawal whitelist from cybavo failed")
		}
	}
}
