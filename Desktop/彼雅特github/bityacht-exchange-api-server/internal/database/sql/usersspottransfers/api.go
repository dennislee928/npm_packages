package usersspottransfers

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"errors"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type GetListRequest struct {
	// UID
	UsersID int64 `form:"usersID"`

	// * 0: 全部
	// * 1: 處理中
	// * 2: 已完成
	// * 3: 已失敗
	// * 4: 已取消
	// * 5: 審核中
	Status Status `form:"status" binding:"gte=0,lte=5"`

	// 幣種
	Coin string `form:"coin"`

	// 交易時間（開始）
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 交易時間（結束）
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type GetExportRequest struct {
	// UID
	UsersID int64 `form:"usersID"`

	// 狀態
	StatusList []Status `form:"statusList"`

	// 交易時間（開始）
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 交易時間（結束）
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type Transfer struct {
	// 交易編號
	TransfersID string `json:"id"`

	// UID
	UsersID int64 `json:"usersID"`

	// 幣種
	CurrenciesSymbol string `json:"currenciesSymbol"`

	// 主網
	Mainnet string `json:"mainnet"`

	// 狀態
	// * 1: 處理中
	// * 2: 已完成
	// * 3: 已失敗
	// * 4: 已取消
	// * 5: 審核中
	Status Status `json:"status"`

	// 類型
	// * 1: 發送
	// * 2: 接收
	Action Action `json:"action"`

	// TXID
	TxID string `json:"txID"`

	// 數量
	Amount decimal.Decimal `json:"amount"`

	// 手續費
	HandlingCharge decimal.Decimal `json:"handlingCharge"`

	// 完成時間
	FinishedAt modelpkg.DateTime `json:"finishedAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 交易時間
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

func GetTransferCSVHeaders() []string {
	return []string{"交易編號", "狀態", "類型", "幣種", "主網", "數量", "手續費", "TXID", "UID", "交易時間", "完成時間"}
}

func (t Transfer) ToCSV() []string {
	return []string{
		t.TransfersID,
		t.Status.Chinese(),
		t.Action.Chinese(),
		t.CurrenciesSymbol,
		t.Mainnet,
		t.Amount.String(),
		t.HandlingCharge.String(),
		t.TxID,
		strconv.FormatInt(t.UsersID, 10),
		t.CreatedAt.ToString(true),
		t.FinishedAt.ToString(true),
	}
}

type AegisTransfer struct {
	// From DB
	TransfersID string
	ToAddress   string
	Amount      decimal.Decimal

	// From User
	TxID            string
	Status          Status
	FinishedAt      modelpkg.DateTime
	FinishedAtValid bool
}

func GetAegisTransferCSVHeaders() []string {
	return []string{"訂單編號", "Address", "Amount", "TXID", "狀態", "完成時間"}
}

func (at AegisTransfer) ToCSV() []string {
	return []string{
		at.TransfersID,
		at.ToAddress,
		at.Amount.String(),
		"",
		"",
		"",
	}
}

func (at *AegisTransfer) FromCSV(rawRecord []string) *errpkg.Error {
	if len(rawRecord) < 6 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadCSVContent, Err: errors.New("bad record length")}
	}

	at.TransfersID = rawRecord[0]
	if at.TransfersID == "" {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadCSVContent, Err: errors.New("bad transfers id")}
	}

	at.ToAddress = rawRecord[1]
	at.Amount, _ = decimal.NewFromString(rawRecord[2])
	at.TxID = rawRecord[3]
	at.Status.FromChinese(rawRecord[4])

	if rawRecord[5] != "" {
		if err := at.FinishedAt.Parse(modelpkg.JSONDateTimeFormat, rawRecord[5]); err != nil {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadCSVContent, Err: err.Err}
		}
		at.FinishedAtValid = true
	}

	return nil
}

type GetSpotTransferForUserRequest struct {
	// 幣種
	CurrenciesSymbol string `form:"currenciesSymbol"`

	// 狀態
	// * 0: 全部
	// * 1: 處理中
	// * 2: 已完成
	// * 3: 已失敗
	// * 4: 已取消
	// * 5: 審核中
	Status Status `form:"status" binding:"gte=0,lte=5"`

	// 類型
	// * 0: 全部
	// * 1: 發送
	// * 2: 接收
	Action Action `form:"action" binding:"gte=0,lte=2"`

	// 交易時間（開始）
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 交易時間（結束）
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type SpotTransferForUser struct {
	// 訂單編號
	TransfersID string `json:"transfersID"`

	// 幣種
	CurrenciesSymbol string `json:"currenciesSymbol"`

	// 狀態
	// * 1: 處理中
	// * 2: 已完成
	// * 3: 已失敗
	// * 4: 已取消
	// * 5: 審核中
	Status Status `json:"status"`

	// 類型
	// * 1: 接收
	// * 2: 發送
	Action Action `json:"action"`

	// TXID
	TxID string `json:"txID"`

	// 數量
	Amount decimal.Decimal `json:"amount"`

	// 手續費
	HandlingCharge decimal.Decimal `json:"handlingCharge"`

	// 交易時間
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type AccWithdrawValuation struct {
	// 今日累計發送
	AccWithdrawInDay decimal.Decimal `json:"accWithdrawInDay"`

	// 本月累計發送
	AccWithdrawInMonth decimal.Decimal `json:"accWithdrawInMonth"`
}
