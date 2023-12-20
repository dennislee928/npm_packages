package sqlcache

import (
	"bityacht-exchange-api-server/internal/database/sql/bankbranchs"
	"bityacht-exchange-api-server/internal/database/sql/banks"
	"bityacht-exchange-api-server/internal/database/sql/countries"
	"bityacht-exchange-api-server/internal/database/sql/industrialclassifications"
)

type IDVOptionsResponse struct {
	Countries []countries.Country            `json:"countries"` // 國籍 & 雙重國籍 列表
	ICs       []industrialclassifications.IC `json:"ics"`       // 行業別列表
}

type SpotOptionsResponse struct {
	Mainnets map[string]map[string]string `json:"mainnets"` // map[幣別]map[主網] => 主網名稱
}

type BankOptionsResponse struct {
	Banks []BankInfo `json:"banks"`
}

type BankInfo struct {
	banks.Bank

	Branchs []bankbranchs.Branch `json:"branchs"`
}
