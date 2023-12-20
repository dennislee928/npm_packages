package receipt

import "fmt"

type ezErrCode int

const (
	ezErrCodeSuccess         ezErrCode = 0
	ezErrCodeInvalidToken    ezErrCode = -3
	ezErrCodeInvalidAppKey   ezErrCode = -5
	ezErrCodeInvalidAppCode  ezErrCode = -10
	ezErrCodeNotAllowedToekn ezErrCode = -20
)

func (code ezErrCode) Err() error {
	switch code {
	case ezErrCodeSuccess:
		return nil
	case ezErrCodeInvalidToken:
		return fmt.Errorf("(%d) invalid token", code)
	case ezErrCodeInvalidAppKey:
		return fmt.Errorf("(%d) invalid app key", code)
	case ezErrCodeInvalidAppCode:
		return fmt.Errorf("(%d) invalid app code", code)
	case ezErrCodeNotAllowedToekn:
		return fmt.Errorf("(%d) not allowed token", code)
	default:
		return fmt.Errorf("unknown error code: %d", code)
	}
}

type ListResp[T any] struct {
	List    []T    `json:"list"`
	Entries uint64 `json:"entries"`
}

type InvNumberRespItem struct {
	// InID: 字軌識別碼
	InID uint64 `json:"inID"`

	// Period: 發票字軌期別，格式為西元年 YYYYMM。例如 2019 年 3~4 月，期別為 201903
	Period string `json:"period"`

	// Lead: 字軌名稱
	Lead string `json:"lead"`

	// StartNo: 發票號碼組的起始序號
	StartNo uint64 `json:"startNo"`

	// EndNo: 發票號碼組的最末一個序號
	EndNo uint64 `json:"endNo"`

	// InvType: 該發票號碼組所適用的稅額計算方式
	//  * 7: 一般稅額計算之電子發票
	//  * 8: 特種稅額計算之電子發票
	InvType int `json:"invType"`

	// BizType: 字軌類別
	//	* 0: 不限用途
	//  * 1: 可用來開立 B2C 的字軌
	//  * 2: 可用來開立 B2B 的字軌
	BizType int `json:"bizType"`

	// Memo: 備註
	Memo string `json:"memo"`

	// IsClosed: 是否已關閉
	//	* 0: 未關閉
	//	* 1: 已關閉
	IsClosed int `json:"isClosed"`

	// ClosedCode: 字軌關閉代碼
	//	* 0: 使用中或可再次被開啟
	//	* 1: 未來期別，不可開啟
	//	* 2: 過期字軌，不可開啟
	//	* 3: 字軌用盡，不可開啟
	ClosedCode int `json:"closedCode"`

	// Platform: 字軌所屬，平台若回傳 null 或無此屬性，表示該分段字軌尚未設定所屬平台
	//	* 1: 易發票
	//	* 100: 其他平台
	Platform *int `json:"platform"`

	// SgoID: 商標識別碼，代表此字軌開出的發票皆使用這個商標
	SgoID *uint64 `json:"sgoID"`
}

type CreateOrderRespItem struct {
	// OrderID: 銷售訂單識別碼
	OrderID uint64 `json:"orderID"`

	// OrderNo: 訂單編號
	OrderNo string `json:"orderNo"`

	// BoNo: 買方訂單號碼
	BoNo string `json:"boNo"`

	// Title: 訂單標題
	Title string `json:"title"`

	// AccCode: 訂單存取認證碼。當購買者為訪客時，會回傳此屬性，用來讓訪客認證具有存取該訂單的資格。
	// 當購買者不是訪客時，就不會有此屬性。
	AccCode *string `json:"accCode"`

	// VoidOrderID: 作廢訂單的識別碼
	VoidOrderID *uint64 `json:"voidOrderID"`

	// Channel: 銷售通路
	Channel *struct {
		// ChannelID: 通路識別碼
		ChannelID uint64 `json:"channelID"`

		// ChannelCode: 通路代碼
		ChannelCode string `json:"channelCode"`

		// ChannelName: 通路名稱
		ChannelName string `json:"channelName"`

		// IsPlatform: 此通路是否為電商平台
		IsPlatform bool `json:"isPlatform"`
	} `json:"channel"`

	// ShipInfo: 取貨方式
	ShipInfo *struct {
		// DevNo: 送貨單號
		DevNo string `json:"devNo"`

		// Pickup: 取貨方式
		//	* 1: 到店取貨（自取）
		//	* 2: 寄送
		//	* 10: 超商取貨
		Pickup int `json:"pickup"`
	} `json:"shipInfo"`

	// SalesAmount: 此訂單的銷售總金額（稅前）
	SalesAmount float64 `json:"salesAmount"`

	// TaxAmount: 此訂單的總稅額
	TaxAmount float64 `json:"taxAmount"`

	// ItemCount: 訂單中的商品數
	ItemCount uint64 `json:"itemCount"`

	// XportFee: 運費
	XportFee float64 `json:"xportFee"`

	// Currency: 幣別代碼
	Currency string `json:"currency"`

	// TRCode: 稅率代碼
	// 若訂單為應稅或混合稅時，此屬性之意義如下：
	//	* 0: 一般稅額(5%)
	//	* 1: 特種飲食業-有陪侍，特種稅額(25%)
	//	* 2: 特種飲食業-有娛樂節目，特種稅額(15%)
	//	* 3: 金融業-其他專屬本業，特種稅額(2%)
	//	* 4: 金融業-保險業再保費，特種稅額(1%)
	//	* 5: 金融業-其他非專屬本業，特種稅額(5%)
	//	* 6: 金融業-本業，特種稅額(5%)
	TRCode int `json:"trCode"`

	// TaxType: 課稅別
	//	* 1: 應稅
	//	* 2: 零稅率
	//	* 3: 免稅
	//	* 9: 混合稅
	TaxType int `json:"taxType"`

	// Status: 訂單狀態
	//	* -1: 作廢
	//	* 4: 草稿（購物車）
	//	* 6: 已成立的訂單
	//	* 51: 付款中
	//	* 52: 完成付款
	//	* 90: 訂單付款失敗
	Status int `json:"status"`

	// Credit4: 信用卡號後四碼
	// 僅適用於信用卡交易，於上傳 B2C 發票時可參考
	Credit4 string `json:"credit4"`

	// Remarks: 訂單備註
	Remarks string `json:"remarks"`

	// CreateTime: 訂單成立時間
	// 格式: YYYY-MM-DD HH:mm:ss
	CreateTime string `json:"createTime"`

	// Buyer: 買受者資訊
	// 格式為 JSON 物件
	Buyer *struct {
		// UserID: 買受者的使用者識別碼
		UserID uint64 `json:"userID"`

		// AccName: 買受者的帳號名稱
		AccName string `json:"accName"`

		// Nid: 買受者的統一編號
		//建立訂單時必須提供 buyerID 參數才會有此回傳屬性
		//也就是 B2B 訂單才會回傳買受者的統一編號
		Nid *string `json:"nid"`

		// Email: 買受者之電子郵件
		Email string `json:"email"`

		// Name: 買受者姓名
		Name string `json:"name"`

		// Phone: 買受者電話
		Phone string `json:"phone"`

		// Addr: 買受者地址
		Addr string `json:"addr"`

		// IsGuest: 建立 B2C 訂單可以不需要輸入買方帳號
		// 此時回傳的買方資訊，其中的 userID, accName 及 nid 等屬性不具參考性
		// 因為僅是訪客帳號
		// 若訂單買方為訪客帳號，此屬性（_isGuest_）將回傳 true，以協助開發者識別
		IsGuest bool `json:"isGuest"`
	} `json:"buyer"`

	// ProdList: 商品清單
	ProdList []struct {
		// SoiID: 購買品項識別碼
		SoiID uint64 `json:"soiID"`

		// PdID: 商品識別碼
		PdID uint64 `json:"pdID"`

		// ProdNo: 商品編號
		ProdNo string `json:"prodNo"`

		// Title: 商品名稱
		Title string `json:"title"`

		// Sales: 單項商品銷售金額（稅前）
		Sales float64 `json:"sales"`

		// SaleTax: 單項商品稅額
		SaleTax float64 `json:"saleTax"`

		// Qty: 購買數量
		Qty uint64 `json:"qty"`

		// Unit: 商品單位名稱
		Unit string `json:"unit"`

		// TaxType: 課稅別
		//	* 1: 應稅
		//	* 2: 零稅
		//	* 3: 免稅
		TaxType int `json:"taxType"`

		// Remarks: 額外說明
		Remarks string `json:"remarks"`
	} `json:"prodList"`
}

type SetCarrierRespItem struct {
	// CarrierType: 標示未來開立 B2C 發票時，發票的載具類型
	CarrierType int `json:"carrierType"`

	// CarrierInfo: 提供載具的額外資訊
	// 當 `carrierType` = 1 時，此處表示買受者的會員帳號
	// 當 `carrierType` = 2 時，此處表示買受者的手機條碼
	// 當 `carrierType` = 3 時，此處表示買受人的自然人憑證編號
	// 當 `carrierType` = 10 時，此處表示買受人的統一編號
	// 當 `carrierType` = 20（境外電商） 時，此處表示買受人的電子郵件帳號
	CarrierInfo string `json:"carrierInfo"`

	// Charity: 發票所要捐贈的慈善機構的捐贈碼
	Charity string `json:"charity"`

	// SendTo: 紙本證明連的寄送資訊
	// 若載具非紙本證明，不需填寫此參數
	SendTo *struct {
		// Name: 郵寄紙本的收件人姓名
		Name string `json:"name"`

		// Phone: 郵寄紙本的收件人電話
		Phone string `json:"phone"`

		// Addr: 郵寄紙本的收件人地址
		Addr string `json:"addr"`
	} `json:"sendTo"`
}

type CreateB2CInvoiceRespItem struct {
	// InvID: 發票識別碼
	InvID uint64 `json:"invID"`

	// InvNo: 發票號碼
	InvNo string `json:"invNo"`

	// OrderID: 訂單識別碼
	OrderID uint64 `json:"orderID"`

	// OrderNo: 訂單編號
	OrderNo string `json:"orderNo"`

	// SalesAmount: 銷售總金額（稅前）
	SalesAmount float64 `json:"salesAmount"`

	// TaxAmount: 總稅額
	TaxAmount float64 `json:"taxAmount"`

	// TaxType: 課稅別
	//	* 1: 應稅
	//	* 2: 零稅率
	//	* 3: 免稅
	//	* 9: 混合稅
	TaxType int `json:"taxType"`

	// TaxRate: 稅率
	TaxRate float64 `json:"taxRate"`

	// IsCustom: 通關方式
	//	* 0: 非經海關出口
	//	* 1: 經海關出口
	IsCustom *int `json:"isCustom"`

	// CreateUser: 發票開立者資訊
	CreateUser *struct {
		// UserID: 使用者識別碼
		UserID uint64 `json:"userID"`

		// IconID: 使用者大頭貼識別碼
		IconID uint64 `json:"iconID"`

		// DspName: 使用者暱稱
		DspName string `json:"dspName"`

		// Email: 使用者電子郵件
		Email string `json:"email"`
	} `json:"createUser"`

	// InvoiceTime: 發票開立時間
	// 格式: YYYY-MM-DD HH:mm:ss
	InvoiceTime string `json:"invoiceTime"`

	// CarrierType: 發票載具類型
	//	* 1: 會員載具
	//	* 2: 手機條碼載具
	//	* 3: 自然人憑證
	//	* 5: 捐贈(捐贈碼)
	//	* 10: 電子發票證明聯
	//	* 20: 跨境電商電子郵件載具
	CarrierType int `json:"carrierType"`

	// CarrierInfo: 載具的額外資訊
	// 當 `carrierType` = 1 時，此屬性代表買受者的會員帳號
	// 當 `carrierType` = 2 時，此屬性代表買受者的手機條碼
	// 當 `carrierType` = 3 時，此屬性代表買受人的自然人憑證編號
	// 當 `carrierType` = 10 時，此屬性代表買受人的統一編號（可能為空字串）
	// 當 `carrierType` = 20（境外電商） 時，此屬性代表買受人的電子郵件帳號
	CarrierInfo string `json:"carrierInfo"`

	// Charity: 若發票要做愛心捐贈時，此屬性顯示所要捐贈的慈善機構代碼
	Charity string `json:"charity"`

	// Remarks: 備註
	Remarks string `json:"remarks"`

	// Buyer: 買受者資訊
	// 格式為 JSON 物件
	Buyer *struct {
		// UserID: 買受者的使用者識別碼
		UserID uint64 `json:"userID"`

		// Name: 買受者姓名
		Name string `json:"name"`

		// Addr: 買受者地址
		Addr string `json:"addr"`

		// Phone: 買受者電話
		Phone string `json:"phone"`

		// IsGuest: 建立 B2C 訂單可以不需要輸入買方帳號
		// 此時回傳的買方資訊，其中的 userID, accName 及 nid 等屬性不具參考性
		// 因為僅是訪客帳號
		// 若訂單買方為訪客帳號，此屬性（_isGuest_）將回傳 true，以協助開發者識別
		IsGuest bool `json:"isGuest"`
	} `json:"buyer"`
}

type ConfirmOrderRespItem struct {
	// OrderID: 訂單識別碼
	//
	// NOTE: 易平台在此欄位似乎有問題，回傳的值為字串，但其他地方都是數字
	// 此處先忽略此值，實際上在此專案中也不會用到 `ConfirmOrderRespItem` 的
	// OrderID
	// OrderID uint64/string `json:"orderID"`

	// OrderNo: 訂單編號
	OrderNo string `json:"orderNo"`

	// BoNo: 買方訂單號碼
	BoNo string `json:"boNo"`

	// Title: 訂單標題
	Title string `json:"title"`
}

type OrderListItem struct {
	OrderID uint64 `json:"orderID"`
}

type CheckMobileCodeResp struct {
	// 手機條碼是否存在。0: 否，1: 是
	IsExist int `json:"isExist"`
}
