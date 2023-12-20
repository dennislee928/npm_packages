package mainnets

import "github.com/shopspring/decimal"

type Mainnet struct {
	// 幣種
	CurrenciesSymbol string `json:"currenciesSymbol"`

	// 鏈
	Mainnet string `json:"mainnet"`

	// 鏈名
	Name string `json:"name"`

	// 手續費
	WithdrawFee decimal.Decimal `json:"withdrawFee"`

	// 最小提幣數量
	WithdrawMin decimal.Decimal `json:"withdrawMin"`

	// 最大提幣數量
	// WithdrawMax      decimal.Decimal `json:"withdrawMax"`
}

type UpdateRequest struct {
	Currency string `uri:"Currency" swaggerignore:"true" binding:"gt=0"`
	Mainnet  string `uri:"Mainnet" swaggerignore:"true" binding:"gt=0"`

	// 手續費 (單位 = Currency)
	WithdrawFee decimal.Decimal `json:"withdrawFee"`

	// 最小提幣數量 (單位 = Currency)
	WithdrawMin decimal.Decimal `json:"withdrawMin"`
}
