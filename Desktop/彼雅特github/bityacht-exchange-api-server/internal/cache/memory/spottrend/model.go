package spottrend

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/database/sql/currencies"
	"bityacht-exchange-api-server/internal/database/sql/transactionpairs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/exchange"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

var cache cacheStruct

type cacheStruct struct {
	BinanceQuoteToMaxQuotePrice decimal.Decimal
	SpotTrendsForIndex          []SpotTrendForIndex
	SpotTrendsForIndexResp      json.RawMessage

	CurrencyInfoMap  map[string]CurrencyInfo
	CurrencyInfoList []CurrencyInfo
	SpotTrendMap     map[string]SpotTrend
	SpotTrendsResp   json.RawMessage

	mux sync.RWMutex
}

type CurrencyInfo struct {
	currencies.Model

	ToTWDRateFromMax decimal.Decimal
	//! For Internal Use Only
	OriToTWDPrice  decimal.Decimal `json:"-"`
	OriToUSDTPrice decimal.Decimal `json:"-"`
}

var (
	symbolsFromBinance = []string{"BTC", "ETH"}
	symbolsFromMax     = []string{"USDT", "USDC"}
)

const (
	binanceQuote = "USDT"
	maxQuote     = "TWD"
)

func getTxPairStr(baseSymbol string, quoteSymbol string) string {
	return baseSymbol + "/" + quoteSymbol
}

func updateCurrencyInfo(ctx context.Context, serviceLogger zerolog.Logger, txPairMap map[string]transactionpairs.Model) (map[string]*exchange.BookTicker, *errpkg.Error) {
	allCurrencies, err := currencies.GetAll()
	if err != nil {
		return nil, err
	}

	// Update CurrencyInfo
	var (
		newBinanceQuoteToMaxQuotePrice decimal.Decimal
		oldCurrencyInfoMap             = GetCurrencyInfoMap()
		priceMap                       = make(map[string]*exchange.BookTicker, len(allCurrencies)*2)
		newCurrencyInfoMap             = make(map[string]CurrencyInfo, len(allCurrencies))
		newCurrencyInfoList            = make([]CurrencyInfo, 0, len(allCurrencies))
	)

	for _, currency := range allCurrencies {
		newCurrencyInfo := CurrencyInfo{Model: currency}

		switch currency.Type {
		case currencies.TypeFiat:
			if currency.Symbol == maxQuote {
				newCurrencyInfo.ToTWDRateFromMax = decimal.NewFromInt(1)
			} else {
				serviceLogger.Warn().Any("currency", currency).Msg("not twd fiat")
				continue
			}
		case currencies.TypeCrypto:
			// Get TWD price
			if price, err := getPriceFromExchange(ctx, serviceLogger.Error(), exchangeSrcMax, currency.Symbol, maxQuote, priceMap); err != nil {
				// Try To Get old Info
				if info, ok := oldCurrencyInfoMap[currency.Symbol]; ok {
					newCurrencyInfo.ToTWDRateFromMax = info.ToTWDRateFromMax
					newCurrencyInfo.OriToTWDPrice = info.OriToTWDPrice
				}
			} else {
				_, newCurrencyInfo.ToTWDRateFromMax, _ = getBuyAndSellPrice(serviceLogger, currency.Symbol, maxQuote, txPairMap, nil, nil, nil, price, decimal.NewFromInt(1))
				newCurrencyInfo.OriToTWDPrice = price.AskPrice
			}

			// Get USDT price
			if currency.Symbol == binanceQuote {
				newCurrencyInfo.OriToUSDTPrice = decimal.NewFromInt(1)
				newBinanceQuoteToMaxQuotePrice = newCurrencyInfo.ToTWDRateFromMax
			} else if price, err := getPriceFromExchange(ctx, serviceLogger.Error(), exchangeSrcBinance, currency.Symbol, binanceQuote, priceMap); err != nil {
				// Try To Get old Info
				if info, ok := oldCurrencyInfoMap[currency.Symbol]; ok {
					newCurrencyInfo.OriToUSDTPrice = info.OriToUSDTPrice
				}
			} else {
				newCurrencyInfo.OriToUSDTPrice = price.AskPrice
			}
		}

		newCurrencyInfoList = append(newCurrencyInfoList, newCurrencyInfo)
		newCurrencyInfoMap[currency.Symbol] = newCurrencyInfo
	}

	if newBinanceQuoteToMaxQuotePrice.LessThanOrEqual(decimal.Zero) {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallMaxAPI, Err: errors.New("binance quote to max quote price <= 0")}
	}

	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.BinanceQuoteToMaxQuotePrice = newBinanceQuoteToMaxQuotePrice
	cache.CurrencyInfoList = newCurrencyInfoList
	cache.CurrencyInfoMap = newCurrencyInfoMap

	return priceMap, nil
}

func updateSpotTrendsForIndex(ctx context.Context, serviceLogger zerolog.Logger, txPairMap map[string]transactionpairs.Model, currencyInfoMap map[string]CurrencyInfo, priceMap map[string]*exchange.BookTicker) *errpkg.Error {
	var (
		err                       *errpkg.Error
		newSpotTrendsForIndexResp json.RawMessage
		newSpotTrendsForIndex     = make([]SpotTrendForIndex, len(symbolsFromBinance)+len(symbolsFromMax))
	)

	cache.mux.RLock()
	oldSpotTrendsForIndex := make([]SpotTrendForIndex, len(cache.SpotTrendsForIndex))
	copy(oldSpotTrendsForIndex, cache.SpotTrendsForIndex)
	binanceQuoteToMaxQuotePrice := cache.BinanceQuoteToMaxQuotePrice
	cache.mux.RUnlock()

	// Update SpotTrendsForIndex
	for i, symbol := range symbolsFromMax {
		i = i + len(symbolsFromBinance)
		var oldSpotTrend *SpotTrendForIndex
		if len(oldSpotTrendsForIndex) != 0 {
			oldSpotTrend = &oldSpotTrendsForIndex[i]
		}

		newSpotTrend := &newSpotTrendsForIndex[i]
		newSpotTrend.Symbol = symbol
		newSpotTrend.Name = getCurrencyNameFromMap(serviceLogger, symbol, oldSpotTrend, currencyInfoMap)

		price, _ := getPriceFromExchange(ctx, serviceLogger.Error(), exchangeSrcMax, symbol, maxQuote, priceMap)
		newSpotTrend.BuyPrice, newSpotTrend.SellPrice, err = getBuyAndSellPrice(serviceLogger, symbol, maxQuote, txPairMap, nil, nil, nil, price, decimal.NewFromInt(1))
		if err != nil {
			return err
		}

		histories := getHistoryFromExchange(ctx, serviceLogger.Error(), exchangeSrcMax, symbol, maxQuote)
		if newSpotTrend.UpsAndDowns, newSpotTrend.Trend, err = getUpDownsAndTrend(nil, oldSpotTrend, histories); err != nil {
			return err
		}
	}

	for i, symbol := range symbolsFromBinance {
		var oldSpotTrend *SpotTrendForIndex
		if len(oldSpotTrendsForIndex) != 0 {
			oldSpotTrend = &oldSpotTrendsForIndex[i]
		}

		newSpotTrend := &newSpotTrendsForIndex[i]
		newSpotTrend.Symbol = symbol
		newSpotTrend.Name = getCurrencyNameFromMap(serviceLogger, symbol, oldSpotTrend, currencyInfoMap)

		price, _ := getPriceFromExchange(ctx, serviceLogger.Error(), exchangeSrcBinance, symbol, binanceQuote, priceMap)
		if newSpotTrend.BuyPrice, newSpotTrend.SellPrice, err = getBuyAndSellPrice(serviceLogger, symbol, binanceQuote, txPairMap, nil, oldSpotTrend, nil, price, binanceQuoteToMaxQuotePrice); err != nil {
			return err
		}

		histories := getHistoryFromExchange(ctx, serviceLogger.Error(), exchangeSrcBinance, symbol, binanceQuote)
		if newSpotTrend.UpsAndDowns, newSpotTrend.Trend, err = getUpDownsAndTrend(nil, oldSpotTrend, histories); err != nil {
			return err
		}
	}

	if rawBytes, err := json.Marshal(newSpotTrendsForIndex); err != nil {
		serviceLogger.Err(err).Msg("marshal newSpotTrendsForIndex error")
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err}
	} else {
		newSpotTrendsForIndexResp = rawBytes
	}

	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.SpotTrendsForIndex = newSpotTrendsForIndex
	cache.SpotTrendsForIndexResp = newSpotTrendsForIndexResp

	return nil
}

func updateSpotTrends(ctx context.Context, serviceLogger zerolog.Logger, allTxPairs []transactionpairs.Model, txPairMap map[string]transactionpairs.Model, currencyInfoMap map[string]CurrencyInfo, priceMap map[string]*exchange.BookTicker) *errpkg.Error {
	var newSpotTrendsResp json.RawMessage

	cache.mux.RLock()
	oldSpotTrendMap := make(map[string]SpotTrend, len(cache.SpotTrendMap))
	for k, v := range cache.SpotTrendMap {
		oldSpotTrendMap[k] = v
	}
	cache.mux.RUnlock()

	// Update SpotTrends
	newSpotTrends := make([]SpotTrend, 0, len(allTxPairs))
	newSpotTrendMap := make(map[string]SpotTrend, len(allTxPairs))
	for _, transactionPairs := range allTxPairs {
		baseCurrencyInfo, ok := currencyInfoMap[transactionPairs.BaseCurrenciesSymbol]
		if !ok || baseCurrencyInfo.Type == currencies.TypeFiat {
			continue
		}

		quoteCurrencyInfo, ok := currencyInfoMap[transactionPairs.QuoteCurrenciesSymbol]
		if !ok || quoteCurrencyInfo.Type == currencies.TypeFiat {
			continue
		}

		newSpotTrend := SpotTrend{
			BaseSymbol:         transactionPairs.BaseCurrenciesSymbol,
			BaseName:           baseCurrencyInfo.Name,
			BasePrecision:      transactionPairs.BaseCurrenciesDecimalPrecision,
			QuoteSymbol:        transactionPairs.QuoteCurrenciesSymbol,
			QuoteName:          quoteCurrencyInfo.Name,
			QuotePrecision:     transactionPairs.QuoteCurrenciesDecimalPrecision,
			HandlingChargeRate: transactionPairs.HandlingChargeRate,
			Status:             transactionPairs.Status,

			//! For Internal Use Only
			SpreadsOfBuy:  transactionPairs.SpreadsOfBuy,
			SpreadsOfSell: transactionPairs.SpreadsOfSell,
		}

		var oldSpotTrend *SpotTrend
		if spotTrend, ok := oldSpotTrendMap[getTxPairStr(newSpotTrend.BaseSymbol, newSpotTrend.QuoteSymbol)]; ok {
			oldSpotTrend = &spotTrend
			newSpotTrend.Status = oldSpotTrend.Status
		}

		var needWarnLog bool
		exchangeInfo, err := exchange.Binance.GetExchangeInfo(ctx, newSpotTrend.BaseSymbol+newSpotTrend.QuoteSymbol)
		newSpotTrend.setExchangeInfo(exchangeInfo, oldSpotTrend)
		if err == nil {
			if exchangeInfo.Status != exchange.SymbolStatusTrading {
				newSpotTrend.Status = transactionpairs.StatusDisable
			}
		} else {
			if needWarnLog = needLog(err); needWarnLog {
				serviceLogger.Warn().AnErr("err", err).Str("base", newSpotTrend.BaseSymbol).Str("quote", newSpotTrend.QuoteSymbol).Msg("binance get exchange info error")
			}
			if oldSpotTrend == nil {
				newSpotTrend.Status = transactionpairs.StatusDisable
			}
		}

		price, err := getPriceFromExchange(ctx, serviceLogger.Warn(), exchangeSrcBinance, newSpotTrend.BaseSymbol, newSpotTrend.QuoteSymbol, priceMap)
		if err == nil {
			newSpotTrend.OriBuyPrice = price.BidPrice
			newSpotTrend.OriSellPrice = price.AskPrice
		} else if oldSpotTrend != nil {
			newSpotTrend.OriBuyPrice = oldSpotTrend.BuyPrice
			newSpotTrend.OriSellPrice = oldSpotTrend.SellPrice
		}

		if newSpotTrend.BuyPrice, newSpotTrend.SellPrice, err = getBuyAndSellPrice(serviceLogger, newSpotTrend.BaseSymbol, newSpotTrend.QuoteSymbol, txPairMap, oldSpotTrend, nil, exchangeInfo, price, decimal.NewFromInt(1)); err != nil {
			if needWarnLog {
				serviceLogger.Warn().AnErr("err", err).Str("base", newSpotTrend.BaseSymbol).Str("quote", newSpotTrend.QuoteSymbol).Msg("binance get buy and sell price error")
			}
			newSpotTrend.Status = transactionpairs.StatusDisable
		}

		histories := getHistoryFromExchange(ctx, serviceLogger.Warn(), exchangeSrcBinance, newSpotTrend.BaseSymbol, newSpotTrend.QuoteSymbol)
		if newSpotTrend.UpsAndDowns, newSpotTrend.Trend, err = getUpDownsAndTrend(oldSpotTrend, nil, histories); err != nil {
			if needWarnLog {
				serviceLogger.Warn().AnErr("err", err).Str("base", newSpotTrend.BaseSymbol).Str("quote", newSpotTrend.QuoteSymbol).Msg("get ups downs and trend error")
			}
			newSpotTrend.Status = transactionpairs.StatusDisable
		}

		newSpotTrends = append(newSpotTrends, newSpotTrend)
		newSpotTrendMap[getTxPairStr(newSpotTrend.BaseSymbol, newSpotTrend.QuoteSymbol)] = newSpotTrend
	}

	if rawBytes, err := json.Marshal(newSpotTrends); err != nil {
		serviceLogger.Err(err).Msg("marshal newSpotTrends error")
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err}
	} else {
		newSpotTrendsResp = rawBytes
	}

	cache.mux.Lock()
	defer cache.mux.Unlock()

	cache.SpotTrendMap = newSpotTrendMap
	cache.SpotTrendsResp = newSpotTrendsResp

	return nil
}

func Update(ctx context.Context) *errpkg.Error {
	var serviceLogger = logger.Logger.With().Str("service", "Update of memory spot trend").Logger()

	allTxPairs, err := transactionpairs.GetAll()
	if err != nil {
		return err
	}

	txPairMap := make(map[string]transactionpairs.Model, len(allTxPairs))
	for _, txPair := range allTxPairs {
		txPairMap[getTxPairStr(txPair.BaseCurrenciesSymbol, txPair.QuoteCurrenciesSymbol)] = txPair
	}

	priceMap, err := updateCurrencyInfo(ctx, serviceLogger, txPairMap)
	if err != nil {
		return err
	}

	newCurrencyInfoMap := GetCurrencyInfoMap()
	if err := updateSpotTrendsForIndex(ctx, serviceLogger, txPairMap, newCurrencyInfoMap, priceMap); err != nil {
		return err
	}

	if err := updateSpotTrends(ctx, serviceLogger, allTxPairs, txPairMap, newCurrencyInfoMap, priceMap); err != nil {
		return err
	}

	return nil
}

func GetSpotTrendsResp() json.RawMessage {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.SpotTrendsResp
}

func GetSpotTrendsForIndexResp() json.RawMessage {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.SpotTrendsForIndexResp
}

func GetCurrencyInfoList() []CurrencyInfo {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.CurrencyInfoList
}

func GetCurrencyInfo(symbol string) (CurrencyInfo, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	if info, ok := cache.CurrencyInfoMap[symbol]; ok {
		return info, nil
	}

	return CurrencyInfo{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCurrencyNotFound}
}

func GetCurrencyInfoMap() map[string]CurrencyInfo {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	return cache.CurrencyInfoMap
}

func GetSpotTrend(base string, quote string) (SpotTrend, *errpkg.Error) {
	cache.mux.RLock()
	defer cache.mux.RUnlock()

	if spotTrend, ok := cache.SpotTrendMap[getTxPairStr(base, quote)]; ok {
		return spotTrend, nil
	}

	return SpotTrend{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeTransactionPairNotFound}
}

func getCurrencyNameFromMap(serviceLogger zerolog.Logger, symbol string, oldSpotTrendForIndex *SpotTrendForIndex, currencyInfoMap map[string]CurrencyInfo) string {
	if info, ok := currencyInfoMap[symbol]; !ok {
		serviceLogger.Err(errors.New("record not found in currencies table")).Str("symbol", symbol).Msg("update spot trend get currency error")

		if oldSpotTrendForIndex != nil {
			return oldSpotTrendForIndex.Name
		}
	} else {
		return info.Name
	}

	return symbol
}

func getBuyAndSellPrice(
	serviceLogger zerolog.Logger,
	base string,
	quote string,
	txPairMap map[string]transactionpairs.Model,
	oldSpotTrend *SpotTrend,
	oldSpotTrendForIndex *SpotTrendForIndex,
	exchangeInfo *exchange.ExchangeInfo,
	price *exchange.BookTicker,
	multiplier decimal.Decimal,
) (buyPrice decimal.Decimal, sellPrice decimal.Decimal, err *errpkg.Error) {
	if price == nil {
		if oldSpotTrend != nil {
			return oldSpotTrend.BuyPrice, oldSpotTrend.SellPrice, nil
		}
		if oldSpotTrendForIndex != nil {
			return oldSpotTrendForIndex.BuyPrice, oldSpotTrendForIndex.SellPrice, nil
		}

		return decimal.Zero, decimal.Zero, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeUpdateCache, Err: errors.New("getBuyAndSellPrice")}
	} else {
		var (
			spreadsOfBuy  = decimal.New(101, -2) // Default: (100 + 1)%
			spreadsOfSell = decimal.New(99, -2)  // Default: (100 - 1)%
		)

		if txPair, ok := txPairMap[getTxPairStr(base, quote)]; ok {
			spreadsOfBuy = decimal.NewFromInt(1).Add(txPair.SpreadsOfBuy.Div(decimal.NewFromInt(100)))
			spreadsOfSell = decimal.NewFromInt(1).Sub(txPair.SpreadsOfSell.Div(decimal.NewFromInt(100)))
		} else {
			serviceLogger.Warn().Str("base", base).Str("quote", quote).Msg("update spot trend tx pair not found")
		}

		// symbol -> binanceQuote -> maxQuote
		buyPrice = price.BidPrice.Mul(multiplier).Mul(spreadsOfBuy)
		sellPrice = price.AskPrice.Mul(multiplier).Mul(spreadsOfSell)

		if exchangeInfo != nil {
			if exchangeInfo.TickSize.GreaterThan(decimal.Zero) {
				buyPrice = buyPrice.Div(exchangeInfo.TickSize).RoundCeil(0).Mul(exchangeInfo.TickSize)
				sellPrice = sellPrice.Div(exchangeInfo.TickSize).RoundCeil(0).Mul(exchangeInfo.TickSize)
			}

			if exchangeInfo.MinPrice.GreaterThan(decimal.Zero) {
				if buyPrice.LessThan(exchangeInfo.MinPrice) {
					buyPrice = exchangeInfo.MinPrice
				}
				if sellPrice.LessThan(exchangeInfo.MinPrice) {
					sellPrice = exchangeInfo.MinPrice
				}
			}

			if exchangeInfo.MaxPrice.GreaterThan(decimal.Zero) {
				if buyPrice.GreaterThan(exchangeInfo.MaxPrice) {
					buyPrice = exchangeInfo.MaxPrice
				}
				if sellPrice.GreaterThan(exchangeInfo.MaxPrice) {
					sellPrice = exchangeInfo.MaxPrice
				}
			}
		}

		return
	}
}

func getUpDownsAndTrend(oldSpotTrend *SpotTrend, oldSpotTrendForIndex *SpotTrendForIndex, histories []decimal.Decimal) (decimal.Decimal, []decimal.Decimal, *errpkg.Error) {
	if historiesLength := len(histories); historiesLength == 0 {
		if oldSpotTrend != nil {
			return oldSpotTrend.UpsAndDowns, oldSpotTrend.Trend, nil
		}
		if oldSpotTrendForIndex != nil {
			return oldSpotTrendForIndex.UpsAndDowns, oldSpotTrendForIndex.Trend, nil
		}

		return decimal.Zero, make([]decimal.Decimal, 0), &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeUpdateCache, Err: errors.New("getUpDownsAndTrend")}
	} else {
		index := historiesLength - 25 // -1 -24
		if index < 0 {
			index = 0
		}

		return (histories[historiesLength-1].Sub(histories[index])).Mul(decimal.NewFromInt(100)).Div(histories[index]), histories, nil
	}
}

type exchangeSrc int8

const (
	exchangeSrcBinance exchangeSrc = iota + 1
	exchangeSrcMax
)

func needLog(err *errpkg.Error) bool {
	if err == nil {
		return false
	}

	return !configs.Config.Exchange.Binance.UseTestnet || !exchange.Binance.IsInvalidSymbol(err.Err)
}

func getPriceFromExchange(ctx context.Context, logEvent *zerolog.Event, src exchangeSrc, base string, quote string, priceMap map[string]*exchange.BookTicker) (*exchange.BookTicker, *errpkg.Error) {
	key := fmt.Sprintf("%d_%s/%s", src, base, quote)

	price, ok := priceMap[key]
	if ok {
		return price, nil
	}

	var (
		err *errpkg.Error
		msg string
	)
	switch src {
	case exchangeSrcBinance:
		price, err = exchange.Binance.GetPrice(ctx, base+quote)
		msg = "binance get price error"
	case exchangeSrcMax:
		price, err = exchange.Max.GetPrice(ctx, base+quote)
		msg = "max get price error"
	}

	if err != nil {
		if needLog(err) {
			logEvent.Err(err).Str("base", base).Str("quote", quote).Msg(msg)
		}
		return nil, err
	}

	priceMap[key] = price
	return price, nil
}

func getHistoryFromExchange(ctx context.Context, logEvent *zerolog.Event, src exchangeSrc, base string, quote string) []decimal.Decimal {
	var (
		histories []decimal.Decimal
		err       *errpkg.Error
		msg       string
	)
	switch src {
	case exchangeSrcBinance:
		histories, err = exchange.Binance.GetHistoryPrice(ctx, base+quote)
		msg = "binance get history price error"
	case exchangeSrcMax:
		histories, err = exchange.Max.GetHistoryPrice(ctx, base+quote)
		msg = "max get history price error"
	}

	if err != nil {
		if needLog(err) {
			logEvent.Err(err).Str("base", base).Str("quote", maxQuote).Msg(msg)
		}
		return histories
	}

	return histories
}
