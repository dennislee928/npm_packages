package assets

import (
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	"bityacht-exchange-api-server/internal/database/sql/currencies"
	"bityacht-exchange-api-server/internal/database/sql/usersvaluationhistories"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// @Summary 	取得我的資產
// @Description 取得我的資產
// @Tags 		User-Assets
// @Security	BearerAuth
// @Produce		json
// @Param 		disableHistories query bool false "Default: false"
// @Success 	200 {object} GetResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/assets [get]
func GetHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	disableHistories, err := modelpkg.GetBoolFromQuery(ctx, "disableHistories", false)
	if errpkg.Handler(ctx, err) {
		return
	}

	var resp GetResponse

	currencyInfoList := spottrend.GetCurrencyInfoList()
	assetForUserMap, err := userswallets.GetAssetForUserMapByUser(claims.UserPayload.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	if !disableHistories {
		aYearAgo := time.Now().AddDate(0, 0, -365)
		historiesStartAt := time.Date(aYearAgo.Year(), aYearAgo.Month(), aYearAgo.Day(), 0, 0, 0, 0, time.UTC)
		valuationHistories, err := usersvaluationhistories.GetHistoriesByUser(claims.UserPayload.ID, historiesStartAt, 364)
		if errpkg.Handler(ctx, err) {
			return
		}

		resp.Histories = make([]decimal.Decimal, 365)
		for _, history := range valuationHistories {
			index := int(history.Date.Sub(historiesStartAt) / time.Hour / 24)
			if index >= len(resp.Histories) {
				reqLogger := logger.GetGinRequestLogger(ctx)
				reqLogger.Warn().Any("valuation history", history).Int64("users id", claims.UserPayload.ID).Msg("bad date")
				continue
			}

			resp.Histories[index] = history.Valuation
		}
	}

	resp.AssetInfos = make([]AssetInfo, 0, len(currencyInfoList))
	for _, currencyInfo := range currencyInfoList {
		assetForUser := assetForUserMap[currencyInfo.Symbol]
		assetForUser.CurrenciesSymbol = currencyInfo.Symbol

		switch currencyInfo.Type {
		case currencies.TypeFiat:
			resp.FiatBalance = resp.FiatBalance.Add(assetForUser.FreeAmount.Mul(currencyInfo.ToTWDRateFromMax))
		case currencies.TypeCrypto:
			valuation := assetForUser.FreeAmount.Mul(currencyInfo.ToTWDRateFromMax)
			resp.CryptocurrencyValuation = resp.CryptocurrencyValuation.Add(valuation)

			resp.AssetInfos = append(resp.AssetInfos, AssetInfo{
				AssetForUser:   assetForUser,
				CurrenciesName: currencyInfo.Name,
				Valuation:      valuation,
			})
		}
	}

	if !disableHistories {
		resp.Histories[len(resp.Histories)-1] = resp.FiatBalance.Add(resp.CryptocurrencyValuation)
	}

	ctx.JSON(http.StatusOK, resp)
}
