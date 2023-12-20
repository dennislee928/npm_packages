package receipts

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"strconv"
)

type GetListRequest struct {
	// 發票狀態
	// 	* 0: 全部
	//	* 1: 未開立
	// 	* 2: 開立中
	// 	* 3: 已開立
	// 	* 4: 已失敗
	Status Status `form:"status" binding:"gte=0,lte=4"`

	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
	EndAt   modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type ListItem struct {
	// ID: 訂單號碼
	ID string `json:"id" binding:"required"`

	// CreatedAt: 發票建立時間
	CreatedAt modelpkg.DateTime `json:"createdAt" binding:"required" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// Status: 發票狀態
	//	* 1: 未開立
	// 	* 2: 開立中
	// 	* 3: 已開立
	// 	* 4: 已失敗
	Status Status `json:"status" binding:"required"`

	// InvoiceID: 發票號碼
	InvoiceID string `json:"invoiceID"`

	// InvoiceAmount: 發票金額
	InvoiceAmount int64 `json:"invoiceAmount" binding:"required"`

	// InvoiceIssuedAt: 發票開立時間
	InvoiceIssuedAt modelpkg.DateTime `json:"invoiceIssuedAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

type DetailItem struct {
	ListItem

	// SalesAmount: 銷售額(未稅價)
	SalesAmount int64 `json:"salesAmount" binding:"required"`

	// Tax: 營業稅額
	Tax int64 `json:"tax" binding:"required"`

	// Barcode: 載具編號
	Barcode string `json:"barcode" binding:"required"`
}

func GetDetailItemCSVHeaders() []string {
	return []string{
		"訂單編號",
		"訂單日期",
		"發票狀態",
		"發票號碼",
		"發票金額",
		"開立時間",
	}
}

func (di DetailItem) ToCSV() []string {
	return []string{
		di.ID,
		di.CreatedAt.ToString(true),
		di.Status.Chinese(),
		di.InvoiceID,
		strconv.FormatInt(di.InvoiceAmount, 10),
		di.InvoiceIssuedAt.ToString(true),
	}
}

type ExportRequest struct {
	// Status: 發票狀態
	//	* 1: 未開立
	// 	* 2: 開立中
	// 	* 3: 已開立
	// 	* 4: 已失敗
	Statuses []Status      `form:"statuses" binding:"unique,dive,gte=1,lte=4"`
	StartAt  modelpkg.Date `form:"startAt" binding:"required" swaggertype:"string" format:"date(YYYY/MM/DD)"`
	EndAt    modelpkg.Date `form:"endAt" binding:"required" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}
