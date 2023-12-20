package userstransactions

import (
	"bityacht-exchange-api-server/internal/database/sql/receipts"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"strconv"

	"github.com/shopspring/decimal"
)

type GetListRequest struct {
	// 0 for all
	UsersID int64 `form:"usersID"`

	// * 0: All
	// * 1: Filled
	// * 2: Killed
	Status Status `form:"status" binding:"gte=0,lte=2"`

	// * 0: All
	// * 1: Buy
	// * 2: Sell
	Side Side `form:"side" binding:"gte=0,lte=2"`

	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
	EndAt   modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type GetExportRequest struct {
	UsersID    int64         `form:"usersID"`
	StatusList []Status      `form:"statusList"`
	StartAt    modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
	EndAt      modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type TransactionForManager struct {
	// 訂單編號
	TransactionsID string `json:"transactionsID"`

	// UID
	UsersID int64 `json:"usersID"`

	// 基準貨幣（交易對左）
	BaseSymbol string `json:"baseSymbol"`

	// 標價貨幣（交易對右）
	QuoteSymbol string `json:"quoteSymbol"`

	// 狀態
	// * 0: 全部
	// * 1: 已完成(Fill)
	// * 2: 交易取消(Kill)
	Status Status `json:"status"`

	// 方向
	// * 1: 買
	// * 2: 賣
	Side Side `json:"side"`

	// 數量 (基準貨幣)
	Quantity decimal.Decimal `json:"quantity"`

	// 價格
	Price decimal.Decimal `json:"price"`

	// 成交額
	Amount decimal.Decimal `json:"amount"`

	// 手續費當時匯率 (獲得的幣種)
	TwdExchangeRate decimal.Decimal `json:"twdExchangeRate"`

	// 總價值
	TwdTotalValue decimal.Decimal `json:"twdTotalValue"`

	// 手續費 (獲得的幣種)
	HandlingCharge decimal.Decimal `json:"handlingCharge"`

	// 幣安 - 訂單編號
	BinanceID int64 `json:"binanceID"`

	// 幣安 - 成交數量
	BinanceQuantity decimal.Decimal `json:"binanceQuantity"`

	// 幣安 - 成交價格
	BinancePrice decimal.Decimal `json:"binancePrice"`

	// 幣安 - 成交額
	BinanceAmount decimal.Decimal `json:"binanceAmount"`

	// 幣安 - 手續費
	BinanceHandlingCharge decimal.Decimal `json:"binanceHandlingCharge"`

	// 發票號碼
	InvoiceID string `json:"invoiceID"`

	// 發票金額 (TWD)
	InvoiceAmount decimal.Decimal `json:"invoiceAmount"`

	// 發票狀態
	// * 0: 無發票
	// * 1: 未開立
	// * 2: 開立中
	// * 3: 已開立
	// * 4: 已失敗
	InvoiceStatus receipts.Status `json:"invoiceStatus"`

	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

func GetTransactionsCSVHeaders() []string {
	return []string{"訂單編號", "狀態", "UID", "交易時間", "發票號碼", "交易對", "方向", "數量", "價格", "成交額", "總價值", "手續費", "幣安訂單編號", "幣安成交數量", "幣安成交價", "幣安成交額", "幣安手續費"}
}

func (t TransactionForManager) ToCSV() []string {

	return []string{
		t.TransactionsID,
		t.Status.Chinese(),
		strconv.FormatInt(t.UsersID, 10),
		t.CreatedAt.ToString(true),
		t.InvoiceID,
		t.BaseSymbol + "/" + t.QuoteSymbol,
		t.Side.Chinese(),
		t.Quantity.String(),
		t.Price.String(),
		t.Amount.String(),
		t.TwdTotalValue.String(),
		t.HandlingCharge.Mul(t.BinanceHandlingCharge).String(),
		strconv.FormatInt(t.BinanceID, 10),
		t.BinanceQuantity.String(),
		t.BinancePrice.String(),
		t.BinanceAmount.String(),
		t.BinanceHandlingCharge.String(),
	}
}

type GetTransactionForUserListRequest struct {
	// 狀態
	// * 0: 全部
	// * 1: 已完成(Fill)
	// * 2: 交易取消(Kill)
	Status Status `form:"status" binding:"gte=0,lte=2"`

	// 方向
	// * 0: 全部
	// * 1: 買
	// * 2: 賣
	Side Side `form:"side" binding:"gte=0,lte=2"`

	// 交易時間(開始)
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 交易時間(結束)
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type TransactionForUser struct {
	// 訂單編號
	TransactionsID string `json:"transactionsID"`

	// 基準貨幣（交易對左）
	BaseSymbol string `json:"baseSymbol"`

	// 標價貨幣（交易對右）
	QuoteSymbol string `json:"quoteSymbol"`

	// 狀態
	// * 1: 已完成(Fill)
	// * 2: 交易取消(Kill)
	Status Status `json:"status"`

	// 方向
	// * 1: 買
	// * 2: 賣
	Side Side `json:"side"`

	// 數量
	Quantity decimal.Decimal `json:"quantity"`

	// 價格
	Price decimal.Decimal `json:"price"`

	// 成交額
	Amount decimal.Decimal `json:"amount"`

	// 手續費 (台幣)
	TwdHandlingCharge decimal.Decimal `json:"handlingCharge"`

	// 台幣匯率 (獲得的幣 -> 台幣)
	TwdExchangeRate decimal.Decimal `json:"twdExchangeRate"`

	// 交易時間
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}
