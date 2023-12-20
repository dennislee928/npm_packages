package userswithdrawalwhitelist

type Record struct {
	ID      int64  `json:"id"`
	Mainnet string `json:"mainnet"` // 主網
	Address string `json:"address"` // 地址
	Extra   Extra  `json:"extra"`
}

func GetCSVHeaders() []string {
	return []string{"主網", "地址", "備註"}
}

func (r Record) ToCSV() []string {
	return []string{
		r.Mainnet,
		r.Address,
		r.Extra.Memo,
	}
}
