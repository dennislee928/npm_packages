package seeders

import (
	"bityacht-exchange-api-server/cmd/migrate/migrations"
	"bityacht-exchange-api-server/cmd/seed"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func init() {
	seed.Register(&SeederTransactionPairs{})
}

// SeederTransactionPairs is a kind of ISeeder, You can set it as same as model.
type SeederTransactionPairs migrations.Migration1689670823

// SeederName for ISeeder
func (*SeederTransactionPairs) SeederName() string {
	return "TransactionPairs"
}

// TableName for gorm
func (*SeederTransactionPairs) TableName() string {
	return (*migrations.Migration1689670823)(nil).TableName()
}

// Default for ISeeder
func (*SeederTransactionPairs) Default(db *gorm.DB) error {
	transactionPairs := []struct {
		BaseCurrenciesSymbol            string
		QuoteCurrenciesSymbol           string
		SpreadsOfBuy                    string
		SpreadsOfSell                   string
		HandlingChargeRate              string
		BaseCurrenciesDecimalPrecision  int32
		QuoteCurrenciesDecimalPrecision int32
		Status                          int32
	}{
		{
			BaseCurrenciesSymbol:            "BTC",
			QuoteCurrenciesSymbol:           "USDT",
			SpreadsOfBuy:                    "0.5",
			SpreadsOfSell:                   "0.5",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          1,
		},
		{
			BaseCurrenciesSymbol:            "ETH",
			QuoteCurrenciesSymbol:           "USDT",
			SpreadsOfBuy:                    "0.5",
			SpreadsOfSell:                   "0.5",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          1,
		},
		{
			BaseCurrenciesSymbol:            "USDC",
			QuoteCurrenciesSymbol:           "USDT",
			SpreadsOfBuy:                    "0.5",
			SpreadsOfSell:                   "0.5",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          1,
		},
		{
			BaseCurrenciesSymbol:            "BTC",
			QuoteCurrenciesSymbol:           "USDC",
			SpreadsOfBuy:                    "0.5",
			SpreadsOfSell:                   "0.5",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          1,
		},
		{
			BaseCurrenciesSymbol:            "ETH",
			QuoteCurrenciesSymbol:           "USDC",
			SpreadsOfBuy:                    "0.5",
			SpreadsOfSell:                   "0.5",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          1,
		},
		{
			BaseCurrenciesSymbol:            "BTC",
			QuoteCurrenciesSymbol:           "TWD",
			SpreadsOfBuy:                    "1.0",
			SpreadsOfSell:                   "1.0",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          0,
		},
		{
			BaseCurrenciesSymbol:            "ETH",
			QuoteCurrenciesSymbol:           "TWD",
			SpreadsOfBuy:                    "1.0",
			SpreadsOfSell:                   "1.0",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          0,
		},
		{
			BaseCurrenciesSymbol:            "USDT",
			QuoteCurrenciesSymbol:           "TWD",
			SpreadsOfBuy:                    "1.0",
			SpreadsOfSell:                   "1.0",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          0,
		},
		{
			BaseCurrenciesSymbol:            "USDC",
			QuoteCurrenciesSymbol:           "TWD",
			SpreadsOfBuy:                    "1.0",
			SpreadsOfSell:                   "1.0",
			HandlingChargeRate:              "0.1",
			BaseCurrenciesDecimalPrecision:  8,
			QuoteCurrenciesDecimalPrecision: 8,
			Status:                          0,
		},
	}

	records := make([]SeederTransactionPairs, len(transactionPairs))
	for i, transactionPair := range transactionPairs {
		var err error
		record := SeederTransactionPairs{
			BaseCurrenciesSymbol:            transactionPair.BaseCurrenciesSymbol,
			QuoteCurrenciesSymbol:           transactionPair.QuoteCurrenciesSymbol,
			BaseCurrenciesDecimalPrecision:  transactionPair.BaseCurrenciesDecimalPrecision,
			QuoteCurrenciesDecimalPrecision: transactionPair.QuoteCurrenciesDecimalPrecision,
			Status:                          transactionPair.Status,
		}

		if record.SpreadsOfBuy, err = decimal.NewFromString(transactionPair.SpreadsOfBuy); err != nil {
			return err
		}
		if record.SpreadsOfSell, err = decimal.NewFromString(transactionPair.SpreadsOfSell); err != nil {
			return err
		}
		if record.HandlingChargeRate, err = decimal.NewFromString(transactionPair.HandlingChargeRate); err != nil {
			return err
		}

		records[i] = record
	}

	return db.Create(&records).Error
}

// Fake for ISeeder
func (*SeederTransactionPairs) Fake(db *gorm.DB) error {
	return seed.ErrNotImplement
}
