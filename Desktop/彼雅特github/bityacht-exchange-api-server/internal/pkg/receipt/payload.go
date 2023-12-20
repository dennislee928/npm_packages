package receipt

import "errors"

const (
	testBaseURL = "https://tryapi.ezreceipt.cc"
	prodBaseURL = "https://api.ezreceipt.cc"
)

const (
	headerAppCode = "x-deva-appcode"
	headerAppKey  = "x-deva-appkey"
	headerToken   = "x-deva-token" // #nosec G101: the api of ezreceipt requires token header
)

type ProdItem struct {
	// ProdNo: 商品編號，長度最多 32 碼。
	// 若未輸入商品編號，系統將視為未知商品。未知商品仍然可以列入訂單中
	// 如果有商品編號，但資料庫中找不到對應的商品且用戶有啟動自動建立商品的功能（預設為關閉）
	// 系統會根據輸入的資料自動建立該項商品。否則應為已知商品
	// 輸入的已知商品如果未包含完整之商品資訊如商品名稱、售價等資訊，系統會自動帶入
	// 過去所記錄的商品資訊
	ProdNo string `json:"prodNo,omitempty"`

	// Title: 商品名稱，長度不得超過 128 個字元
	// 若此參數未給但本平台曾記錄該商品資訊，系統將依過去記錄的資訊自動填入商品名稱
	// 否則系統將回報錯誤
	Title string `json:"title"`

	// Sales: 商品單價，可記錄到小數點以下四位
	// 該價格為含稅與否，依 `incTax` 參數決定
	Sales float64 `json:"sales,omitempty"`

	// SalesAmount: 商品價格小計（商品單價乘以數量），可記錄到小數點以下四位
	// 該價格為含稅與否，依 `incTax` 參數決定
	// 若有設定此欄位，則單價參數 (`sales`) 將自動忽略
	SalesAmount float64 `json:"salesAmount,omitempty"`

	// IncTax: 用來標示單價 (`sales`) 與小計 (`salesAmount`) 的價格為含稅價或未稅價
	// 預設為未稅。若商品課稅別為零稅率或免稅時，此參數將被強制設為 false，即未稅價
	//	* true: 含稅價格
	// 	* false: 未稅價格
	IncTax bool `json:"incTax"`

	// Tax: 單項商品稅額。可記錄到小數點以下四位
	// 該稅額為單價稅額或小計稅額，依下列規則判定：
	// 如果商品價格給的是小計 (`salesAmount`)，此欄位即為小計稅額，否則是單價稅額
	// 如果此欄位未填寫，將由本平台計算稅額
	// 如果此欄位有值，則商品稅額以此欄位為準，本平台將不做任何稅額調整或檢查
	Tax float64 `json:"tax,omitempty"`

	// Qty: 購買數量，必要
	Qty uint64 `json:"qty"`

	// Unit: 商品單位名稱。字串長度請勿超過 6 個字元
	Unit string `json:"unit,omitempty"`

	// McType: 商品類型
	// 系統判定原則為先查詢店家的商品資訊（以商品編號查詢），如果有找到就採用該商品的類型；
	// 如果找不到，則判定為一般商品 (mcType = 0)
	// 當商品類型為「折扣商品」時，價格才能小於零
	//	* 0: 一般商品
	//	* 3: 運費
	//	* 100: 折扣商品
	//	* -1: 類型不明，交由系統判定。
	McType int `json:"mcType,omitempty"`

	// TaxType: 課稅別
	// 若為混合稅率則必須指定課稅別，否則預設為應稅項目
	//	* 1: 應稅
	//	* 2: 零稅率
	//	* 3: 免稅
	TaxType int `json:"taxType,omitempty"`

	// Remarks: 備註說明。輸入的字串長度不能超過 256 個字元
	Remarks string `json:"remarks,omitempty"`
}

// CreateOrderPayload represents the payload of CreateOrder.
// uri: sales/order/create
type CreateOrderPayload struct {
	// OrderNo: 訂單編號。長度請勿超過 32 個字元。若未輸入則由系統自行依規則編列流水號。
	// NOTE: 若無特殊需求請保持此欄為空，使用系統自行編列的流水號
	OrderNo string `json:"orderNo,omitempty"`

	// VoidOrderID: 退換貨的原訂單識別碼。進行退換貨時，此參數為必要的輸入
	VoidOrderID uint64 `json:"voidOrderID,omitempty"`

	// Title: 訂單標題。長度請勿超過 64 個字元
	Title string `json:"title"`

	// BuyerID: 買受人的帳號識別碼。若未來要開立 B2B 發票時，此參數為必要
	BuyerID uint64 `json:"buyerID,omitempty"`

	// ChannelID: 通路識別碼。可以從 API: sales/channel/list 取得銷售通路清單
	// 若此參數未設定，將設定為預設通路
	ChannelID uint64 `json:"channelID,omitempty"`

	// BoNo: 買方訂單號碼。長度請勿超過 32 個字元
	BoNo string `json:"boNo,omitempty"`

	// DevNo: 送貨單號。長度請勿超過 16 個字元
	DevNo string `json:"devNo,omitempty"`

	// Pickup: 取貨方式
	//	* 1: 到店取貨（自取）
	//	* 2: 寄送
	//	* 10: 超商取貨。
	Pickup int `json:"pickup,omitempty"`

	// Currency: 訂單交易幣別的幣別代碼。幣別代碼請參考技術文件中的幣別代碼表
	// 境外電商應填寫此欄位。預設為新台幣(NTD)
	Currency string `json:"currency,omitempty"`

	// TRCode: 稅率代碼。若本次銷售為應稅或混合稅時，請填入此參數
	//	* 0: 一般稅額(5%)
	//	* 1: 特種飲食業-有陪侍，特種稅額(25%)
	//	* 2: 特種飲食業-有娛樂節目，特種稅額(15%)
	//	* 3: 金融業-其他專屬本業，特種稅額(2%)
	//	* 4: 金融業-保險業再保費，特種稅額(1%)
	//	* 5: 金融業-其他非專屬本業，特種稅額(5%)
	//	* 6: 金融業-本業，特種稅額(5%)
	TRCode int `json:"trCode,omitempty"`

	// IsCustom: 通關方式。若未來開立發票時為零稅率發票（ taxType = 2）
	// 或是混合稅發票中有零稅率項目時，此參數必要
	//	* 0: 非經海關出口
	//	* 1: 經海關出口
	IsCustom int `json:"isCustom,omitempty"`

	// Remarks: 訂單備註。長度請勿超過 40 個字元
	Remarks string `json:"remarks,omitempty"`

	// CreateTime: 訂單建立時間。預設為呼叫此 API 的時間
	// 格式: YYYY-MM-DD HH:mm:ss
	CreateTime string `json:"createTime,omitempty"`

	// ProdList 訂單中的商品清單。若 confirm 參數設為 true，則此「商品清單」成為必要參數。此參數為陣列，陣列包含商品內容
	ProdList []ProdItem `json:"prodList"`

	// Credit4: 信用卡號後四碼。僅適用於信用卡交易，於上傳 B2C 發票時可參考
	Credit4 string `json:"credit4,omitempty"`

	// Confirm: 若設為 true，表示同時確認本訂單已成立
	// 否則本訂單仍視為草稿，後續需以 /sales/order/confirm 確認訂單
	Confirm bool `json:"confirm"`
}

// SetCarrierPayload represents the payload of SetCarrier.
// uri: sales/order/setCarrier
type SetCarrierPayload struct {
	// AccCode: 16 碼的訂單存取認證碼。當訂單買方為訪客時，需要此認證碼才能存取該訂單的資料。
	AccCode string `json:"accCode,omitempty"`

	// CarrierType: 標示未來開立 B2C 發票時，發票的載具類型。
	// 請注意，財政部規定 B2C 發票必須確認載具資訊，因此若將載具設定為 -1.未知載具
	// 後續將無法完成 B2C 發票開立。
	// 允許的數值有：
	//	* -1: 未知載具
	//	* 1: 會員載具
	//	* 2: 手機條碼
	//	* 3: 自然人憑證
	//	* 5: 捐贈碼捐贈
	//	* 10: 紙本證明聯
	//	* 20: 境外電商專用載具
	CarrierType int `json:"carrierType"`

	// CarrierInfo: 提供載具的額外資訊
	// 當 carrierType = 1 時，請於此處輸入買受者的會員帳號
	// 當 carrierType = 2 時，請於此處輸入買受者的手機條碼
	// 當 carrierType = 3 時，此處需輸入買受人的自然人憑證編號
	// 當 carrierType = 10 時，此處可輸入買受人的統一編號（無則免填）
	// 當 carrierType = 20（境外電商） 時，此處需輸入買受人的電子郵件帳號
	CarrierInfo string `json:"carrierInfo"`

	// BuyerNID: 買受人統一編號
	// 這個參數只有在 carrierType = 1,2,3 時會被參考
	// 當 carrierType = 10 時，若需輸入統一編號，請參考 carrierInfo 的說明
	BuyerNID string `json:"buyerNID,omitempty"`

	// BuyerName: 買受人名稱
	// 列印紙本發票時可能會使用到
	BuyerName string `json:"buyerName,omitempty"`

	// Charity: 發票所要捐贈的慈善機構的捐贈碼
	// 當 carrierType = 5 時，此輸入參數為必要
	// 發票捐贈後將無法再列印為紙本
	Charity string `json:"charity,omitempty"`

	// ReqProof: 使用手機條碼載具且加註統一編號時，是否列印發票紙本
	// 設定為 true 代表要將發票紙本印出
	// 這個參數只有在 carrierType = 2 且 buyerNID 有輸入時會被參考
	ReqProof bool `json:"reqProof,omitempty"`

	// IsNonprofit: 客戶是否為非營利事業單位
	// 非營利事業單位通常代表是政府機關或非營利事業機構
	// 預設為 false
	IsNonprofit bool `json:"isNonprofit"`

	// SendTo: 紙本證明連的寄送資訊
	// 若載具非紙本證明，不需填寫此參數
	SendTo *struct {
		// Name: 郵寄紙本的收件人姓名，長度不可超過 64 個字元
		Name string `json:"name"`

		// Phone: 郵寄紙本的收件人電話，長度不可超過 16 個字元
		Phone string `json:"phone"`

		// Addr: 郵寄紙本的收件人地址，長度不可超過 128 個字元
		Addr string `json:"addr"`
	} `json:"sendTo,omitempty"`
}

type CreateB2CInvoicePayload struct {
	// InvNo: 發票號碼，必須為十碼（字軌加號碼
	// 若要由系統自動選擇號碼，則此參數勿給。若是發票註銷後重開，則此參數將被忽略
	InvNo string `json:"invNo,omitempty"`

	// AutoInvNo: 如果未預設字軌、或是預設字軌號碼已用盡時，是不是由系統自動選擇合適的字軌
	// 如果已有指定發票號碼 (invNo 參數），此參數將被忽略不考慮。
	AutoInvNo bool `json:"autoInvNo,omitempty"`

	// Title: 買受人名稱。當買受人在系統上無會員帳號，但仍有需要標示買受人名稱時，可記錄於此欄位。
	// 字串長度不可超過 60 個字元。
	Title string `json:"title,omitempty"`

	// IsCustom: 通關方式
	// 若為零稅率發票（ taxType = 2）或是混合稅發票中有零稅率項目時，此參數必要
	//	* 0: 非經海關出口
	//	* 1: 經海關出口
	IsCustom int `json:"isCustom,omitempty"`

	// InvoiceTime: 發票開立時間，預設為呼叫此 API 的時間
	// 若是發票註銷重開將沿用原發票開立時間，此參數將被忽略
	// 參數格式為：YYYY-MM-DD HH:mm:ss
	InvoiceTime string `json:"invoiceTime,omitempty"`

	// Remarks: 註解。可用來附加任何額外之資訊。長度請勿超過 200 個字元
	Remarks string `json:"remarks,omitempty"`

	// AccCode: 16 碼的訂單存取認證碼
	// 當使用者為一般消費者（非店家人員）或訪客帳號時，需要此認證碼才能開立發票
	// 此使用情境通常是消費者更改載具後，將發票注銷重開的情境
	AccCode string `json:"accCode,omitempty"`
}

type ConfirmOrderPayload struct {
	// AccCode: 16 碼的訂單存取認證碼。當訂單買方為訪客時，需要此認證碼才能存取該訂單的資料
	AccCode string `json:"accCode,omitempty"`

	// Credit4: 信用卡號後四碼。僅適用於信用卡交易，於上傳 B2C 發票時可參考
	Credit4 string `json:"credit4,omitempty"`

	// CreateTime: 訂單建立時間。若未設定，則沿用原草稿建立的時間。格式: YYYY-MM-DD HH:mm:ss
	CreateTime string `json:"createTime,omitempty"`

	// PaidTime: 訂單付款時間
	// 若此時間有設定，表示該筆訂單已付清，也因此 confirm 參數將被強制為 true
	// 若尚未付清，請勿填寫此參數
	// 此處設定的時間不得早於訂單成立時間 ( createTime )
	// 格式: YYYY-MM-DD HH:mm:ss
	PaidTime string `json:"paidTime,omitempty"`

	// Discounts: 折扣內容。此參數為一個陣列，陣列值是描述折扣內容的物件
	Discounts []struct {
		// Title: 折扣的說明文字，可用來說明提供折扣的原因或內容
		// 程度請勿超過 128 個字元。有提供折扣時，此為必要輸入
		Title string `json:"title"`

		// Amount: 折扣金額。請以正數表示。例如折扣金額為 35 元時，請輸入 35，而不是 -35
		// 有提供折扣時，此為必要輸入
		Amount float64 `json:"amount"`
	} `json:"discounts,omitempty"`
}

type FastCreateB2CPayload struct {
	ID         string
	Title      string
	ProdItem   ProdItem
	CarrierNum string
}

func (p FastCreateB2CPayload) Validate() error {
	if p.ID == "" {
		return errors.New("id is empty")
	}

	if p.Title == "" {
		return errors.New("title is empty")
	}

	if p.ProdItem.Title == "" {
		return errors.New("title is empty")
	}

	if p.ProdItem.Sales == 0 {
		return errors.New("sales is zero")
	}

	if p.ProdItem.Qty == 0 {
		return errors.New("qty is zero")
	}

	return nil
}

// CheckMobileCodePayload represents the payload of CheckMobileCode.
// uri: openTax/carrier/checkMobileCode
type CheckMobileCodePayload struct {
	// 手機條碼。格式應為「/」加上7碼英數字，共8碼
	MobileCode string `json:"mobileCode"`
}
