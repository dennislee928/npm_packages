package spottrend

import (
	"bityacht-exchange-api-server/internal/database/sql/transactionpairs"
	"bityacht-exchange-api-server/internal/pkg/exchange"

	"github.com/shopspring/decimal"
)

type SpotTrendForIndex struct {
	Symbol      string            `json:"symbol"`      // 幣別
	Name        string            `json:"name"`        // 幣別全名
	BuyPrice    decimal.Decimal   `json:"buyPrice"`    // 買入價格 (台幣)
	SellPrice   decimal.Decimal   `json:"sellPrice"`   // 賣出價格 (台幣)
	UpsAndDowns decimal.Decimal   `json:"upsAndDowns"` // 24H 漲跌 （%）
	Trend       []decimal.Decimal `json:"trend"`       // 7日價格走勢 (不一定為台幣)
}

type SpotTrend struct {
	// 基準貨幣代號
	BaseSymbol string `json:"baseSymbol"`

	// 基準貨幣全名
	BaseName string `json:"baseName"`

	// 基準貨幣精準度 (小數後幾位)
	BasePrecision int32 `json:"basePrecision"`

	// 標價貨幣代號
	QuoteSymbol string `json:"quoteSymbol"`

	// 標價貨幣全名
	QuoteName string `json:"quoteName"`

	// 標價貨幣精準度 (小數後幾位)
	QuotePrecision int32 `json:"quotePrecision"`

	// 買入價格 (單位=標價貨幣)
	BuyPrice decimal.Decimal `json:"buyPrice"`

	// 賣出價格 (單位=標價貨幣)
	SellPrice decimal.Decimal `json:"sellPrice"`

	// 手續費 (%)
	HandlingChargeRate decimal.Decimal `json:"handlingChargeRate"`

	// 24H 漲跌 （%）
	UpsAndDowns decimal.Decimal `json:"upsAndDowns"`

	// 7日價格走勢
	Trend []decimal.Decimal `json:"trend"`

	// 最小單位 (Price * Quantity)
	MinNotional decimal.Decimal `json:"minNotional"`

	// 最大單位 (Price * Quantity)
	MaxNotional decimal.Decimal `json:"maxNotional"`

	// 狀態
	// * 0: 暫時關閉
	// * 1: 正常
	Status transactionpairs.Status `json:"status"`

	//! For Internal Use Only
	OriBuyPrice   decimal.Decimal `json:"-"`
	OriSellPrice  decimal.Decimal `json:"-"`
	SpreadsOfBuy  decimal.Decimal `json:"-"`
	SpreadsOfSell decimal.Decimal `json:"-"`
	MinPrice      decimal.Decimal `json:"-"`
	MaxPrice      decimal.Decimal `json:"-"`
	TickSize      decimal.Decimal `json:"-"`
	MinQuantity   decimal.Decimal `json:"-"`
	MaxQuantity   decimal.Decimal `json:"-"`
	StepSize      decimal.Decimal `json:"-"`
}

func (st *SpotTrend) setExchangeInfo(exchangeInfo *exchange.ExchangeInfo, oldSt *SpotTrend) {
	if exchangeInfo != nil {
		st.MinNotional = exchangeInfo.MinNotional
		st.MaxNotional = exchangeInfo.MaxNotional
		st.MinPrice = exchangeInfo.MinPrice
		st.MaxPrice = exchangeInfo.MaxPrice
		st.TickSize = exchangeInfo.TickSize
		st.MinQuantity = exchangeInfo.MinQuantity
		st.MaxQuantity = exchangeInfo.MaxQuantity
		st.StepSize = exchangeInfo.StepSize
	} else if oldSt != nil {
		st.MinNotional = oldSt.MinNotional
		st.MaxNotional = oldSt.MaxNotional
		st.MinPrice = oldSt.MinPrice
		st.MaxPrice = oldSt.MaxPrice
		st.TickSize = oldSt.TickSize
		st.MinQuantity = oldSt.MinQuantity
		st.MaxQuantity = oldSt.MaxQuantity
		st.StepSize = oldSt.StepSize
	}
}
