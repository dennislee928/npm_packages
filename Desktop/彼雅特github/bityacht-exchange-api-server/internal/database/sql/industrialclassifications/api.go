package industrialclassifications

// IC is abbreviation of Industrial Classification
type IC struct {
	ID      int64  `json:"id"`
	Code    string `json:"code"`    // 中華民國 行業統計分類 代號
	Chinese string `json:"chinese"` // 中文
	English string `json:"english"` // 英文
}
