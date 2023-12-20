package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederLevelLimits{})
}

// SeederLevelLimits is a kind of ISeeder, You can set it as same as model.
type SeederLevelLimits migrations.Migration1689651014

// SeederName for ISeeder
func (*SeederLevelLimits) SeederName() string {
	return "LevelLimits"
}

// SeederName for ISeeder
func (s *SeederLevelLimits) TableName() string {
	return (*migrations.Migration1689651014)(nil).TableName()
}

// Default for ISeeder
func (*SeederLevelLimits) Default(db *gorm.DB) error {
	const TypeNaturalPerson = 1
	const TypeJuridicalPerson = 2
	unlimited := decimal.NewFromInt(-1)

	// Ref: https://docs.google.com/document/d/1kTjwAVY7xutcgOGjIKrcJ9GX5PYxkKBz6N_cFWHCJSo/edit#heading=h.9008won7plgt
	// TODO: Some Value is not confirmed.
	records := []SeederLevelLimits{
		// Natural Person 0 ~ 5
		{
			Type:  TypeNaturalPerson,
			Level: 0,
			// All Zero
		},
		{
			Type:                          TypeNaturalPerson,
			Level:                         1,
			MaxDepositUsdtPerTransaction:  unlimited,
			MaxDepositUsdtPerDay:          unlimited,
			MaxDepositUsdtPerMonth:        unlimited,
			MinWithdrawUsdtPerTransaction: decimal.NewFromInt(10),
		},
		{
			Type:                          TypeNaturalPerson,
			Level:                         2,
			MaxDepositTwdPerTransaction:   unlimited,
			MaxDepositTwdPerDay:           decimal.NewFromInt(150000),
			MaxDepositTwdPerMonth:         decimal.NewFromInt(500000),
			MinWithdrawTwdPerTransaction:  decimal.NewFromInt(15),
			MaxWithdrawTwdPerTransaction:  decimal.NewFromInt(100000),
			MaxWithdrawTwdPerDay:          decimal.NewFromInt(250000),
			MaxWithdrawTwdPerMonth:        decimal.NewFromInt(400000),
			MaxDepositUsdtPerTransaction:  unlimited,
			MaxDepositUsdtPerDay:          unlimited,
			MaxDepositUsdtPerMonth:        unlimited,
			MinWithdrawUsdtPerTransaction: decimal.NewFromInt(10),
			MaxWithdrawUsdtPerTransaction: decimal.NewFromInt(1000),
			MaxWithdrawUsdtPerDay:         decimal.NewFromInt(3000),
			MaxWithdrawUsdtPerMonth:       decimal.NewFromInt(5000),
		},
		{
			Type:                          TypeNaturalPerson,
			Level:                         3,
			MaxDepositTwdPerTransaction:   unlimited,
			MaxDepositTwdPerDay:           decimal.NewFromInt(350000),
			MaxDepositTwdPerMonth:         decimal.NewFromInt(800000),
			MinWithdrawTwdPerTransaction:  decimal.NewFromInt(15),
			MaxWithdrawTwdPerTransaction:  decimal.NewFromInt(300000),
			MaxWithdrawTwdPerDay:          decimal.NewFromInt(500000),
			MaxWithdrawTwdPerMonth:        decimal.NewFromInt(750000),
			MaxDepositUsdtPerTransaction:  unlimited,
			MaxDepositUsdtPerDay:          unlimited,
			MaxDepositUsdtPerMonth:        unlimited,
			MinWithdrawUsdtPerTransaction: decimal.NewFromInt(10),
			MaxWithdrawUsdtPerTransaction: decimal.NewFromInt(5000),
			MaxWithdrawUsdtPerDay:         decimal.NewFromInt(10000),
			MaxWithdrawUsdtPerMonth:       decimal.NewFromInt(25000),
		},
		{
			Type:                          TypeNaturalPerson,
			Level:                         4,
			MaxDepositTwdPerTransaction:   unlimited,
			MaxDepositTwdPerDay:           decimal.NewFromInt(1000000),
			MaxDepositTwdPerMonth:         decimal.NewFromInt(5000000),
			MinWithdrawTwdPerTransaction:  decimal.NewFromInt(15),
			MaxWithdrawTwdPerTransaction:  decimal.NewFromInt(800000),
			MaxWithdrawTwdPerDay:          decimal.NewFromInt(1500000),
			MaxWithdrawTwdPerMonth:        decimal.NewFromInt(3500000),
			MaxDepositUsdtPerTransaction:  unlimited,
			MaxDepositUsdtPerDay:          unlimited,
			MaxDepositUsdtPerMonth:        unlimited,
			MinWithdrawUsdtPerTransaction: decimal.NewFromInt(10),
			MaxWithdrawUsdtPerTransaction: decimal.NewFromInt(30000),
			MaxWithdrawUsdtPerDay:         decimal.NewFromInt(60000),
			MaxWithdrawUsdtPerMonth:       decimal.NewFromInt(150000),
		},
		{
			Type:                          TypeNaturalPerson,
			Level:                         5,
			MaxDepositTwdPerTransaction:   unlimited,
			MaxDepositTwdPerDay:           decimal.NewFromInt(3000000),
			MaxDepositTwdPerMonth:         decimal.NewFromInt(10000000),
			MinWithdrawTwdPerTransaction:  decimal.NewFromInt(15),
			MaxWithdrawTwdPerTransaction:  decimal.NewFromInt(2000000),
			MaxWithdrawTwdPerDay:          decimal.NewFromInt(5000000),
			MaxWithdrawTwdPerMonth:        decimal.NewFromInt(10000000),
			MaxDepositUsdtPerTransaction:  unlimited,
			MaxDepositUsdtPerDay:          unlimited,
			MaxDepositUsdtPerMonth:        unlimited,
			MinWithdrawUsdtPerTransaction: decimal.NewFromInt(10),
			MaxWithdrawUsdtPerTransaction: decimal.NewFromInt(100000),
			MaxWithdrawUsdtPerDay:         decimal.NewFromInt(200000),
			MaxWithdrawUsdtPerMonth:       decimal.NewFromInt(500000),
		},
		// Juridical Person 1 ~ 2
		{
			Type:                          TypeJuridicalPerson,
			Level:                         1,
			MaxDepositTwdPerTransaction:   unlimited,
			MaxDepositTwdPerDay:           decimal.NewFromInt(2000000),
			MaxDepositTwdPerMonth:         decimal.NewFromInt(10000000),
			MinWithdrawTwdPerTransaction:  decimal.NewFromInt(15),
			MaxWithdrawTwdPerTransaction:  unlimited,
			MaxWithdrawTwdPerDay:          decimal.NewFromInt(1000000),
			MaxWithdrawTwdPerMonth:        decimal.NewFromInt(5000000),
			MaxDepositUsdtPerTransaction:  unlimited,
			MaxDepositUsdtPerDay:          unlimited,
			MaxDepositUsdtPerMonth:        unlimited,
			MinWithdrawUsdtPerTransaction: decimal.NewFromInt(10),
			MaxWithdrawUsdtPerTransaction: decimal.NewFromInt(200000),
			MaxWithdrawUsdtPerDay:         decimal.NewFromInt(200000),
			MaxWithdrawUsdtPerMonth:       decimal.NewFromInt(400000),
		},
		{
			Type:                          TypeJuridicalPerson,
			Level:                         2,
			MaxDepositTwdPerTransaction:   unlimited,
			MaxDepositTwdPerDay:           decimal.NewFromInt(2000000),
			MaxDepositTwdPerMonth:         decimal.NewFromInt(10000000),
			MinWithdrawTwdPerTransaction:  decimal.NewFromInt(15),
			MaxWithdrawTwdPerTransaction:  unlimited,
			MaxWithdrawTwdPerDay:          decimal.NewFromInt(2000000),
			MaxWithdrawTwdPerMonth:        decimal.NewFromInt(10000000),
			MaxDepositUsdtPerTransaction:  unlimited,
			MaxDepositUsdtPerDay:          unlimited,
			MaxDepositUsdtPerMonth:        unlimited,
			MinWithdrawUsdtPerTransaction: decimal.NewFromInt(10),
			MaxWithdrawUsdtPerTransaction: decimal.NewFromInt(300000),
			MaxWithdrawUsdtPerDay:         decimal.NewFromInt(600000),
			MaxWithdrawUsdtPerMonth:       decimal.NewFromInt(2000000),
		},
	}

	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederLevelLimits) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
