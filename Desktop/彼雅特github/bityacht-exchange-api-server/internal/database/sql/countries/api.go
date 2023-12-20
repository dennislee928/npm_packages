package countries

type Country struct {
	Code    string `json:"code"`    // ISO-3166 Alpha-3
	Chinese string `json:"chinese"` // 中文
	English string `json:"english"` // 英文
	Locale  string `json:"locale"`  // 語言(For KryptoGO)
}
