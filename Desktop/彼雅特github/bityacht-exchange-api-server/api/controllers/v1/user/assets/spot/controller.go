package spot

import (
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/database/sql/transactionpairs"
	"bityacht-exchange-api-server/internal/database/sql/userscommissions"
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	"bityacht-exchange-api-server/internal/database/sql/userstransactions"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/exchange"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	dbsql "database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/adshao/go-binance/v2"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// @Summary 	取得錢包歷史紀錄(加密貨幣)
// @Description 取得錢包歷史紀錄(加密貨幣)
// @Tags 		User-Assets
// @Security	BearerAuth
// @Produce		json
// @Param 		query query usersspottransfers.GetSpotTransferForUserRequest false "query"
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]usersspottransfers.SpotTransferForUser}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/assets/spot/histories [get]
func GetHistoriesHandler(ctx *gin.Context) {
	var req usersspottransfers.GetSpotTransferForUserRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	if resp.Data, err = usersspottransfers.GetSpotTransfersForUser(claims.UserPayload.ID, req, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	取得現貨交易紀錄
// @Description 取得現貨交易紀錄
// @Tags 		User-Assets
// @Security	BearerAuth
// @Produce		json
// @Param 		query query userstransactions.GetTransactionForUserListRequest false "query"
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]userstransactions.TransactionForUser}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/assets/spot/transactions [get]
func GetTransactionsHandler(ctx *gin.Context) {
	var req userstransactions.GetTransactionForUserListRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	if resp.Data, err = userstransactions.GetTransactionForUserListByUser(claims.UserPayload.ID, req, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	取得現貨相關選項
// @Description 取得現貨相關選項
// @Tags 		User-Assets
// @Security	BearerAuth
// @Produce		json
// @Success 	200 {object} sqlcache.SpotOptionsResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/assets/spot/options [get]
func GetOptionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sqlcache.GetSpotOptionsResponse())
}

// @Summary 	取得現貨買賣價格資訊
// @Description 取得現貨買賣價格資訊
// @Tags 		User-Assets
// @Security	BearerAuth
// @Produce		json
// @Success 	200 {object} []spottrend.SpotTrend
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/assets/spot/trend [get]
func GetTrendHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, spottrend.GetSpotTrendsResp())
}

// @Summary 	現貨買賣
// @Description 現貨買賣
// @Tags 		User-Assets
// @Security	BearerAuth
// @Accept 		json
// @Produce		json
// @Param 		body body CreateTransactionRequest true "Request Body"
// @Success 	200 {object} userstransactions.TransactionForUser
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/assets/spot/transactions [post]
func CreateTransactionHandler(ctx *gin.Context) {
	var req CreateTransactionRequest

	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := req.Validate(); errpkg.Handler(ctx, err) {
		return
	}

	spotTrend, err := spottrend.GetSpotTrend(req.BaseSymbol, req.QuoteSymbol)
	if errpkg.Handler(ctx, err) {
		return
	} else if spotTrend.Status != transactionpairs.StatusEnable {
		errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadTransactionPairStatus, errors.New("bad status"))
		return
	}

	newTxRecord := userstransactions.Model{
		UsersID:     claims.UserPayload.ID,
		BaseSymbol:  req.BaseSymbol,
		QuoteSymbol: req.QuoteSymbol,
		Side:        req.Side,
		Price:       req.Price,
	}

	// Check Req Calculation
	var (
		earnPrecision                             int32
		earnCurrencyInfo                          spottrend.CurrencyInfo
		exchangeAmount, notional, binanceQuantity decimal.Decimal
		paySymbol                                 string
		binanceSide                               binance.SideType
	)

	switch req.Side {
	case userstransactions.SideBuy: // Pay: quote, Earn: base
		paySymbol = req.QuoteSymbol
		earnPrecision = spotTrend.BasePrecision
		exchangeAmount = req.PayAmount.Div(req.Price) // = Pay * (1 / Price)
		earnCurrencyInfo, err = spottrend.GetCurrencyInfo(req.BaseSymbol)
		binanceSide = binance.SideTypeBuy

		newTxRecord.Quantity = req.EarnAmount
		newTxRecord.Amount = req.PayAmount
		notional = req.PayAmount
		binanceQuantity = req.EarnAmount.Add(req.HandlingCharge)
	case userstransactions.SideSell: // Pay: base, Earn: quote
		paySymbol = req.BaseSymbol
		earnPrecision = spotTrend.QuotePrecision
		exchangeAmount = req.PayAmount.Mul(req.Price) // Pay * Price
		earnCurrencyInfo, err = spottrend.GetCurrencyInfo(req.QuoteSymbol)
		binanceSide = binance.SideTypeSell

		newTxRecord.Quantity = req.PayAmount
		newTxRecord.Amount = req.EarnAmount
		notional = exchangeAmount
		binanceQuantity = req.PayAmount
	default:
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("bad validator of req.Side")})
		return
	}

	if errpkg.Handler(ctx, err) { // Handle the err of spottrend.GetCurrencyInfo
		return
	} else if (spotTrend.MinNotional.GreaterThan(decimal.Zero) && notional.LessThan(spotTrend.MinNotional)) ||
		(spotTrend.MaxNotional.GreaterThan(decimal.Zero) && notional.GreaterThan(spotTrend.MaxNotional)) {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad notional")})
		return
	}

	// Set Quantity by StepSize
	if spotTrend.StepSize.GreaterThan(decimal.Zero) {
		binanceQuantity = binanceQuantity.Div(spotTrend.StepSize).RoundCeil(0).Mul(spotTrend.StepSize)
	}
	// Check Quantity by MinQuantity & MaxQuantity
	if spotTrend.MinQuantity.GreaterThan(decimal.Zero) && binanceQuantity.LessThan(spotTrend.MinQuantity) {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("quantity not enough")})
		return
	} else if spotTrend.MaxQuantity.GreaterThan(decimal.Zero) && binanceQuantity.GreaterThan(spotTrend.MaxQuantity) {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("quantity too large")})
		return
	}

	// exchangeAmount * (HandlingChargeRate / 100)
	newTxRecord.HandlingCharge = exchangeAmount.Mul(spotTrend.HandlingChargeRate).Div(decimal.NewFromInt(100)).RoundUp(earnPrecision)
	earnAmount := exchangeAmount.Sub(newTxRecord.HandlingCharge).RoundDown(earnPrecision) // earn = exchangeAmount - HandlingCharge

	if !req.HandlingCharge.Equal(newTxRecord.HandlingCharge) {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadCalcOfTransaction, Err: errors.New("handling charge")})
		return
	} else if !req.EarnAmount.Equal(earnAmount) {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadCalcOfTransaction, Err: errors.New("earn amount")})
		return
	}

	newTxRecord.TwdExchangeRate = earnCurrencyInfo.OriToTWDPrice
	newTxRecord.TwdTotalValue = req.EarnAmount.Mul(earnCurrencyInfo.OriToTWDPrice).Round(0)

	// Start Transaction
	walletNeedRollback := true
	reqLogger := logger.GetGinRequestLogger(ctx)
	if err := userswallets.LockAmount(claims.ID(), paySymbol, req.PayAmount); errpkg.Handler(ctx, err) {
		return
	}
	defer func() {
		if walletNeedRollback {
			if err := userswallets.UnlockAmount(claims.ID(), paySymbol, req.PayAmount); err != nil {
				reqLogger.Err(err.Err).Msg("failed to unlock amount")
			}
		}
	}()

	orderResp, err := exchange.Binance.CreateOrder(ctx, newTxRecord.BaseSymbol+newTxRecord.QuoteSymbol, binanceSide, binance.OrderTypeLimit, binanceQuantity, newTxRecord.Price)
	if errpkg.Handler(ctx, err) {
		return
	} else if byteOrderResp, err := json.Marshal(orderResp); err != nil {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err})
		return
	} else {
		newTxRecord.Extra.BinanceResp = string(byteOrderResp)
	}

	var commissionRecord *userscommissions.Model
	switch orderResp.Status {
	case binance.OrderStatusTypeFilled:
		var binanceHandlingCharge decimal.Decimal
		for _, fill := range orderResp.Fills {
			binanceHandlingCharge = binanceHandlingCharge.Add(fill.Commission)
		}

		newTxRecord.Status = userstransactions.StatusFilled
		newTxRecord.BinanceID = orderResp.OrderID
		newTxRecord.BinancePrice = orderResp.Price
		newTxRecord.BinanceQuantity = orderResp.ExecutedQuantity
		newTxRecord.BinanceAmount = orderResp.CummulativeQuoteQuantity
		newTxRecord.BinanceHandlingCharge = binanceHandlingCharge

		if claims.UserPayload.InviterID != 0 {
			// (0.2 * HandlingCharge) -> USDT
			commissionAmount := newTxRecord.HandlingCharge.Mul(earnCurrencyInfo.OriToUSDTPrice).Mul(decimal.New(2, -1)).RoundDown(2)

			if commissionAmount.GreaterThan(decimal.Zero) {
				commissionRecord = &userscommissions.Model{
					UsersID:     claims.InviterID,
					FromUsersID: dbsql.NullInt64{Int64: claims.ID(), Valid: true},
					Action:      userscommissions.ActionDeposit,
					Amount:      commissionAmount,
				}
			}
		}
	case binance.OrderStatusTypeCanceled, binance.OrderStatusTypePendingCancel, binance.OrderStatusTypeRejected, binance.OrderStatusTypeExpired:
		newTxRecord.HandlingCharge = decimal.Zero
		newTxRecord.Status = userstransactions.StatusKilled
	default:
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: errors.New("bad order status"), Data: orderResp.Status})
		return
	}

	receiptRecord, err := userstransactions.Create(&newTxRecord, commissionRecord, paySymbol, req.PayAmount, earnCurrencyInfo.Symbol, req.EarnAmount)
	if errpkg.Handler(ctx, err) {
		return
	}
	walletNeedRollback = newTxRecord.Status != userstransactions.StatusFilled

	resp := userstransactions.TransactionForUser{
		TransactionsID:  newTxRecord.TransactionsID,
		BaseSymbol:      newTxRecord.BaseSymbol,
		QuoteSymbol:     newTxRecord.QuoteSymbol,
		Status:          newTxRecord.Status,
		Side:            newTxRecord.Side,
		Quantity:        newTxRecord.Quantity,
		Price:           newTxRecord.Price,
		Amount:          newTxRecord.Amount,
		TwdExchangeRate: newTxRecord.TwdExchangeRate,
		CreatedAt:       modelpkg.DateTime{Time: newTxRecord.CreatedAt},
	}
	if receiptRecord != nil {
		resp.TwdHandlingCharge = decimal.NewFromInt(receiptRecord.InvoiceAmount)
	}

	ctx.JSON(http.StatusOK, resp)

	// Send Transaction Result Mail
	txMail := email.NewEmail(email.WithLogo())
	txMail.To = []string{claims.UserPayload.Account}
	mailPayload := emailtemplates.TransactionPayload{
		Time:           emailtemplates.FormatTime(newTxRecord.CreatedAt),
		TransactionsID: newTxRecord.TransactionsID,
		BaseSymbol:     newTxRecord.BaseSymbol,
		QuoteSymbol:    newTxRecord.QuoteSymbol,
		Side:           newTxRecord.Side.Chinese(),
		Quantity:       newTxRecord.Quantity.String(),
	}

	switch newTxRecord.Status {
	case userstransactions.StatusFilled:
		txMail.Subject, txMail.HTML, err = emailtemplates.ExecTransactionSuccessful(mailPayload)
	case userstransactions.StatusKilled:
		txMail.Subject, txMail.HTML, err = emailtemplates.ExecTransactionFailed(mailPayload)
	default:
		reqLogger.Warn().Int32("Status", int32(newTxRecord.Status)).Msg("bad status of transaction record")
		return
	}

	if err != nil {
		reqLogger.Err(err).Msg("execute transaction template error")
	} else if err = email.SendMail(txMail); err != nil {
		reqLogger.Err(err).Msg("send mail error")
	}
}
