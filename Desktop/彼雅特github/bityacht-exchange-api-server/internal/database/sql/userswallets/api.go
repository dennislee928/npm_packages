package userswallets

import (
	"github.com/shopspring/decimal"
)

type Asset struct {
	CurrenciesSymbol string          `json:"currenciesSymbol"`
	FreeAmount       decimal.Decimal `json:"freeAmount"`
}

type AssetForUser struct {
	CurrenciesSymbol string          `json:"currenciesSymbol"` // 幣種
	FreeAmount       decimal.Decimal `json:"freeAmount"`       // 持有數量
	LockedAmount     decimal.Decimal `json:"lockedAmount"`     // 處理中數額
}
