package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/internal/database/sql/banks"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederBanks{})
}

// SeederBanks is a kind of ISeeder, You can set it as same as model.
type SeederBanks banks.Model

// SeederName for ISeeder
func (*SeederBanks) SeederName() string {
	return "Banks"
}

// TableName for gorm
func (*SeederBanks) TableName() string {
	return (*migrations.Migration1697001488)(nil).TableName()
}

// Default for ISeeder
func (*SeederBanks) Default(db *gorm.DB) error {
	records := []SeederBanks{
		{Code: "004", Chinese: "臺灣銀行", English: ""},
		{Code: "005", Chinese: "臺灣土地銀行", English: ""},
		{Code: "006", Chinese: "合作金庫商業銀行", English: ""},
		{Code: "007", Chinese: "第一商業銀行", English: ""},
		{Code: "008", Chinese: "華南商業銀行", English: ""},
		{Code: "009", Chinese: "彰化商業銀行", English: ""},
		{Code: "011", Chinese: "上海商業儲蓄銀行", English: ""},
		{Code: "012", Chinese: "台北富邦商業銀行", English: ""},
		{Code: "013", Chinese: "國泰世華商業銀行", English: ""},
		{Code: "016", Chinese: "高雄銀行", English: ""},
		{Code: "017", Chinese: "兆豐國際商業銀行", English: ""},
		{Code: "018", Chinese: "全國農業金庫", English: ""},
		{Code: "048", Chinese: "王道商業銀行", English: ""},
		{Code: "050", Chinese: "臺灣中小企業銀行", English: ""},
		{Code: "052", Chinese: "渣打國際商業銀行", English: ""},
		{Code: "053", Chinese: "台中商業銀行", English: ""},
		{Code: "054", Chinese: "京城銀行", English: ""},
		{Code: "081", Chinese: "滙豐（台灣）商業銀行", English: ""},
		{Code: "101", Chinese: "瑞興商業銀行", English: ""},
		{Code: "102", Chinese: "華泰商業銀行", English: ""},
		{Code: "103", Chinese: "臺灣新光商業銀行", English: ""},
		{Code: "108", Chinese: "陽信商業銀行", English: ""},
		{Code: "118", Chinese: "板信商業銀行", English: ""},
		{Code: "147", Chinese: "三信商業銀行", English: ""},
		{Code: "700", Chinese: "中華郵政股份有限公司", English: ""},
		{Code: "803", Chinese: "聯邦商業銀行", English: ""},
		{Code: "805", Chinese: "遠東國際商業銀行", English: ""},
		{Code: "806", Chinese: "元大商業銀行", English: ""},
		{Code: "807", Chinese: "永豐商業銀行", English: ""},
		{Code: "808", Chinese: "玉山商業銀行", English: ""},
		{Code: "809", Chinese: "凱基商業銀行", English: ""},
		{Code: "810", Chinese: "星展（台灣）商業銀行", English: ""},
		{Code: "812", Chinese: "台新國際商業銀行", English: ""},
		{Code: "816", Chinese: "安泰商業銀行", English: ""},
		{Code: "822", Chinese: "中國信託商業銀行", English: ""},
		// {Code: "823", Chinese: "將來商業銀行", English: ""},
	}
	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederBanks) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
