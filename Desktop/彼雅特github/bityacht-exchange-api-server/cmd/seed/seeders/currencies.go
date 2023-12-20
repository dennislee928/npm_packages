package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"

	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederCurrencies{})
}

// SeederCurrencies is a kind of ISeeder, You can set it as same as model.
type SeederCurrencies migrations.Migration1689663148

// SeederName for ISeeder
func (*SeederCurrencies) SeederName() string {
	return "Currencies"
}

// TableName for gorm
func (*SeederCurrencies) TableName() string {
	return (*migrations.Migration1689663148)(nil).TableName()
}

// Default for ISeeder
func (*SeederCurrencies) Default(db *gorm.DB) error {
	// Ref: https://www.binance.com/en/trade-rule
	records := []SeederCurrencies{
		{Symbol: "TWD", Name: "TWD", Type: 1, DecimalPrecision: 0},
		{Symbol: "USDT", Name: "Tether", Type: 2, DecimalPrecision: 8},
		{Symbol: "USDC", Name: "USD Coin", Type: 2, DecimalPrecision: 8},
		{Symbol: "BTC", Name: "Bitcoin", Type: 2, DecimalPrecision: 8},
		{Symbol: "ETH", Name: "Ethereum", Type: 2, DecimalPrecision: 6},
	}

	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederCurrencies) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
