package bankbranchs

type Branch struct {
	BanksCode string `json:"-"`

	Code    string `json:"code"`    // 分行代號
	Chinese string `json:"chinese"` // 分行名稱（中文）
	English string `json:"english"` // 分行名稱（英文）
}
