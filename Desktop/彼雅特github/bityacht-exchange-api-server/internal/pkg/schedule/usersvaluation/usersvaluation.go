package usersvaluation

import (
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/schedulelogs"
	"bityacht-exchange-api-server/internal/database/sql/usersvaluationhistories"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

const scheduleKeyPrefix = "UsersValuation-"

func CalcUsersValuation(job gocron.Job) {
	jobLogger := logger.Logger.With().Str("service", "users valuation schedule").Logger()

	now := time.Now()
	lastDate := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC)
	scheduleKey := scheduleKeyPrefix + lastDate.Format(time.DateOnly)

	scheduleLogRecord, err := schedulelogs.Create(scheduleKey, strings.Join(job.Tags(), ", "))
	if err != nil {
		jobLogger.Err(err.Err).Msg("create schedule logs failed")
		return
	}

	scheduleLogRecord.Result.Error = "Unknown"
	scheduleLogRecord.Status = schedulelogs.StatusAbort
	defer func() {
		if err := schedulelogs.Update(&scheduleLogRecord); err != nil {
			jobLogger.Err(err.Err).Msg("update schedule logs failed")
		}
	}()

	wallets, err := userswallets.GetAllFreeAmountGreaterThanZero()
	if err != nil {
		scheduleLogRecord.Result.Error = err.Error()
		jobLogger.Err(err.Err).Msg("[userswallets] get all free amount greater than zero failed")
		return
	}

	currencyInfoMap := spottrend.GetCurrencyInfoMap()
	scheduleLogRecord.Result.CurrencyInfo = currencyInfoMap
	valuationHistories := make([]*usersvaluationhistories.Model, 0, len(wallets))
	valuationHistoryMap := make(map[int64]*usersvaluationhistories.Model)

	for _, wallet := range wallets {
		currencyInfo, ok := currencyInfoMap[wallet.CurrenciesSymbol]
		if !ok {
			jobLogger.Warn().Any("wallet", wallet).Msg("currency info not found")
			continue
		}

		valuationHistory := valuationHistoryMap[wallet.UsersID]
		if valuationHistory == nil {
			valuationHistory = &usersvaluationhistories.Model{
				UsersID: wallet.UsersID,
				Date:    lastDate,
			}
			valuationHistoryMap[wallet.UsersID] = valuationHistory
			valuationHistories = append(valuationHistories, valuationHistory)
		}
		valuationHistory.Valuation = valuationHistory.Valuation.Add(wallet.FreeAmount.Mul(currencyInfo.ToTWDRateFromMax))
	}

	if err := sql.DB().Create(&valuationHistories).Error; err != nil {
		scheduleLogRecord.Result.Error = err.Error()
		jobLogger.Err(err).Msg("[db] create users valuation histories failed")
		return
	}

	scheduleLogRecord.Result.Error = ""
	scheduleLogRecord.Status = schedulelogs.StatusFinished
}
