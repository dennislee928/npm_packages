package cybavo

import (
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	"bityacht-exchange-api-server/internal/database/sql/walletsaddresses"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func WithdrawalCallbackHandler(ctx *gin.Context) {
	logger := logger.Logger.With().Str("service", "cybavo-withdrawal-callback").Logger()
	checksum := ctx.GetHeader(wallet.HeaderChecksum)
	if checksum == "" {
		logger.Warn().Msg("checksum is empty")
		ctx.Status(http.StatusBadRequest)
		return
	}

	// TODO: Validate Order ID
	// if err := wallet.Cybavo.WithdrawalCallback(checksum, ctx.Request.Body); err != nil {
	// 	logger.Warn().Err(err).Msg("withdrawal callback error")
	// 	ctx.Status(http.StatusBadRequest)
	// 	return
	// }

	ctx.String(http.StatusOK, "OK")
}

func CallbackHandler(ctx *gin.Context) {
	logger := logger.GetGinRequestLogger(ctx).With().Str("service", "cybavo-callback").Logger()
	rawReqBody := new(bytes.Buffer)
	tee := io.TeeReader(ctx.Request.Body, rawReqBody)
	checksum := ctx.GetHeader(wallet.HeaderChecksum)
	if checksum == "" {
		logger.Warn().Msg("checksum is empty")
		ctx.Status(http.StatusBadRequest)
		return
	}

	info, err := wallet.Cybavo.Callback(checksum, tee)
	if err != nil {
		logger.Warn().Err(err).Msg("callback error")
		ctx.Status(http.StatusBadRequest)
		return
	}
	logger.Debug().Any("info", info).Msg("callback")

	// return OK anyway for Cybavo health check
	defer ctx.String(http.StatusOK, "OK")

	// Check user id or fetch from wallet address
	var (
		user    users.Model
		wrapErr *errpkg.Error
	)
	if userID, err := info.UserID(); err == nil {
		// if callback info contains user id, use it
		id, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			logger.Error().Err(err).Msg("parse user id")
			return
		}
		user, wrapErr = users.GetByID(id)
	} else if info.Type == wallet.CallbackTypeDeposit {
		user, wrapErr = walletsaddresses.GetUserByAddress(info.Mainnet().BinanceNetwork(), info.ToAddress)
	}
	if wrapErr != nil {
		logger.Warn().Err(wrapErr).Msg("get user failed")
		return
	}

	switch info.Type {
	case wallet.CallbackTypeDeposit, wallet.CallbackTypeWithdraw:
		amount, err := info.AmountDecimal()
		if err != nil {
			logger.Error().Err(err).Msg("parse amount")
			return
		}

		now := time.Now().UTC()
		spotTransfer := &usersspottransfers.Model{
			TransfersID:      info.OrderID(),
			Type:             usersspottransfers.TypeCybavoAPI,
			UsersID:          user.ID,
			CurrenciesSymbol: info.CurrencyType().String(),
			Mainnet:          info.Mainnet().String(),
			FromAddress:      info.FromAddress,
			ToAddress:        info.ToAddress,
			Status:           infoStateToModelStatus(info.ProcessingState, info.State),
			Action:           infoTypeToModelAction(info.Type),
			TxID:             info.TXID,
			Serial:           sql.NullInt64{Int64: info.Serial, Valid: true},
			Amount:           amount,
			Extra:            usersspottransfers.Extra{RawMessage: rawReqBody.Bytes()},
			FinishedAt:       &now,
		}

		oldSpot, wrapErr := usersspottransfers.TakeOldAndUpsert(spotTransfer)
		if wrapErr != nil {
			logger.Error().Err(wrapErr).Msg("take old and upsert spot transfer")
			return
		}

		if !spotTransfer.Done() {
			return
		}

		if oldSpot != nil && oldSpot.Done() {
			return
		}

		if wrapErr := handleCallback(info, spotTransfer.HandlingCharge, user); wrapErr != nil {
			logger.Error().Err(wrapErr).Msg("handle failed")
			return
		}
	default:
		logger.Error().Any("info", info).Msg("bad callback type")
		return
	}
}

func handleCallback(info *wallet.CallbackInfo, handlingCharge decimal.Decimal, user users.Model) *errpkg.Error {
	switch info.ProcessingState {
	case wallet.ProcessingStateFailed:
		return handleFailedCallback(info, handlingCharge, user)

	case wallet.ProcessingStateDone:
		return handleDoneCallback(info, handlingCharge, user)

	case wallet.ProcessingStateInPool, wallet.ProcessingStateInChain:
		return nil

	default:
		return &errpkg.Error{Err: fmt.Errorf("unknown processing state: %d", info.ProcessingState)}
	}
}

func handleFailedCallback(info *wallet.CallbackInfo, handlingCharge decimal.Decimal, user users.Model) *errpkg.Error {
	val, err := info.AmountDecimal()
	if err != nil {
		return &errpkg.Error{Err: fmt.Errorf("parse amount: %w", err)}
	}

	var wrapErr *errpkg.Error
	switch info.Type {
	case wallet.CallbackTypeDeposit:
		// TODO: notify, but maybe it's not necessary due to the deposit is not done.

	case wallet.CallbackTypeWithdraw:
		if wrapErr := walletsaddresses.WithdrawFailed(user.ID, info.CurrencyType(), val.Add(handlingCharge)); wrapErr != nil {
			return wrapErr
		}

		notificationEmail := email.NewEmail(email.WithLogo())
		notificationEmail.To = []string{user.Account}
		if notificationEmail.Subject, notificationEmail.HTML, wrapErr = emailtemplates.ExecNotifyWithdrawSpotFailed(emailtemplates.NotifyWithdrawSpotFailedPayload{
			Time:             emailtemplates.FormatTime(time.Now()),
			CurrenciesSymbol: info.CurrencyType().String(),
			Mainnet:          info.Mainnet().String(),
			Amount:           val.String(),
		}); wrapErr != nil {
			logger.Logger.Err(wrapErr).Msg("failed to execute template")
			break
		}

		if wrapErr = email.SendMail(notificationEmail); wrapErr != nil {
			logger.Logger.Err(wrapErr).Msg("failed to send mail")
		}

	default:
		return &errpkg.Error{Err: fmt.Errorf("unknown callback type: %d", info.Type)}
	}

	return nil
}

func handleDoneCallback(info *wallet.CallbackInfo, handlingCharge decimal.Decimal, user users.Model) *errpkg.Error {
	val, err := info.AmountDecimal()
	if err != nil {
		return &errpkg.Error{Err: err}
	}

	var wrapErr *errpkg.Error
	switch info.Type {
	case wallet.CallbackTypeDeposit:
		if wrapErr := walletsaddresses.Deposit(info.CurrencyType().String(), info.Mainnet().BinanceNetwork(), info.ToAddress, val); wrapErr != nil {
			return wrapErr
		}

		notificationEmail := email.NewEmail(email.WithLogo())
		notificationEmail.To = []string{user.Account}
		if notificationEmail.Subject, notificationEmail.HTML, wrapErr = emailtemplates.ExecNotifyDepositSpot(emailtemplates.NotifyDepositSpotPayload{
			Time:             emailtemplates.FormatTime(time.Now()),
			CurrenciesSymbol: info.CurrencyType().String(),
			Mainnet:          info.Mainnet().String(),
			Amount:           val.String(),
		}); wrapErr != nil {
			logger.Logger.Err(wrapErr).Msg("failed to execute template")
			break
		}

		if wrapErr = email.SendMail(notificationEmail); wrapErr != nil {
			logger.Logger.Err(wrapErr).Msg("failed to send mail")
		}

	case wallet.CallbackTypeWithdraw:
		return walletsaddresses.WithdrawDone(user.ID, info.CurrencyType(), val.Add(handlingCharge))

	default:
		return &errpkg.Error{Err: fmt.Errorf("unknown callback type: %d", info.Type)}
	}

	return nil
}

func infoTypeToModelAction(t wallet.CallbackType) usersspottransfers.Action {
	switch t {
	case wallet.CallbackTypeDeposit:
		return usersspottransfers.ActionDeposit

	case wallet.CallbackTypeWithdraw:
		return usersspottransfers.ActionWithdraw

	default:
		return 0
	}
}

func infoStateToModelStatus(ps wallet.ProcessingState, ws wallet.CallbackState) usersspottransfers.Status {
	if ws == wallet.CallbackStateCancelled {
		return usersspottransfers.StatusCanceled
	}

	switch ps {
	case wallet.ProcessingStateDone:
		return usersspottransfers.StatusFinished

	case wallet.ProcessingStateInPool, wallet.ProcessingStateInChain:
		return usersspottransfers.StatusProcessing

	case wallet.ProcessingStateFailed:
		return usersspottransfers.StatusFailed

	default:
		return 0
	}
}
