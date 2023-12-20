package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederMainnets{})
}

// SeederMainnets is a kind of ISeeder, You can set it as same as model.
type SeederMainnets migrations.Migration1689670240

// SeederName for ISeeder
func (*SeederMainnets) SeederName() string {
	return "Mainnets"
}

// TableName for gorm
func (*SeederMainnets) TableName() string {
	return (*migrations.Migration1689670240)(nil).TableName()
}

// Default for ISeeder
func (*SeederMainnets) Default(db *gorm.DB) error {
	mainnets := []struct {
		CurrenciesSymbol         string
		Mainnet                  string
		Name                     string
		AddressRegex             string
		WithdrawDecimalPrecision int32
		WithdrawFee              string
		WithdrawMin              string
		WithdrawMax              string
		Status                   int32
	}{
		{
			CurrenciesSymbol:         "BTC",
			Mainnet:                  "BTC",
			Name:                     "Bitcoin",
			AddressRegex:             "^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$|^[(bc1q)|(bc1p)][0-9A-Za-z]{37,62}$",
			WithdrawDecimalPrecision: 8,
			WithdrawFee:              "0.0002",
			WithdrawMin:              "0.001",
			WithdrawMax:              "-1",
			Status:                   1,
		},
		{
			CurrenciesSymbol:         "ETH",
			Mainnet:                  "ETH",
			Name:                     "Ethereum (ERC20)",
			AddressRegex:             "^(0x)[0-9A-Fa-f]{40}$",
			WithdrawDecimalPrecision: 8,
			WithdrawFee:              "0.002",
			WithdrawMin:              "0.01",
			WithdrawMax:              "-1",
			Status:                   1,
		},
		{
			CurrenciesSymbol:         "USDC",
			Mainnet:                  "ETH",
			Name:                     "Ethereum (ERC20)",
			AddressRegex:             "^(0x)[0-9A-Fa-f]{40}$",
			WithdrawDecimalPrecision: 6,
			WithdrawFee:              "5",
			WithdrawMin:              "20",
			WithdrawMax:              "-1",
			Status:                   1,
		},
		{
			CurrenciesSymbol:         "USDT",
			Mainnet:                  "ETH",
			Name:                     "Ethereum (ERC20)",
			AddressRegex:             "^(0x)[0-9A-Fa-f]{40}$",
			WithdrawDecimalPrecision: 6,
			WithdrawFee:              "5",
			WithdrawMin:              "20",
			WithdrawMax:              "-1",
			Status:                   1,
		},
		{
			CurrenciesSymbol:         "USDT",
			Mainnet:                  "TRX",
			Name:                     "Tron (TRC20)",
			AddressRegex:             "^T[1-9A-HJ-NP-Za-km-z]{33}$",
			WithdrawDecimalPrecision: 6,
			WithdrawFee:              "2",
			WithdrawMin:              "10",
			WithdrawMax:              "-1",
			Status:                   1,
		},
	}

	records := make([]SeederMainnets, len(mainnets))
	for i, mainnet := range mainnets {
		var err error
		record := SeederMainnets{
			CurrenciesSymbol:         mainnet.CurrenciesSymbol,
			Mainnet:                  mainnet.Mainnet,
			Name:                     mainnet.Name,
			AddressRegex:             mainnet.AddressRegex,
			WithdrawDecimalPrecision: mainnet.WithdrawDecimalPrecision,
			Status:                   mainnet.Status,
		}

		if record.WithdrawFee, err = decimal.NewFromString(mainnet.WithdrawFee); err != nil {
			return err
		}
		if record.WithdrawMin, err = decimal.NewFromString(mainnet.WithdrawMin); err != nil {
			return err
		}
		if record.WithdrawMax, err = decimal.NewFromString(mainnet.WithdrawMax); err != nil {
			return err
		}

		records[i] = record
	}

	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederMainnets) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
