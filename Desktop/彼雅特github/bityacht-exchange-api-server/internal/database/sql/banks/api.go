package banks

type Bank struct {
	Code    string `json:"code"`    // 銀行代號
	Chinese string `json:"chinese"` // 銀行名稱（中文）
	English string `json:"english"` // 銀行名稱（英文）
}
