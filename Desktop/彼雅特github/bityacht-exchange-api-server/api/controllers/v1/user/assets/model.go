package assets

import (
	"bityacht-exchange-api-server/internal/database/sql/userswallets"

	"github.com/shopspring/decimal"
)

type GetResponse struct {
	FiatBalance             decimal.Decimal   `json:"fiatBalance"`             // 法幣餘額
	CryptocurrencyValuation decimal.Decimal   `json:"cryptocurrencyValuation"` // 數位資產估值
	AssetInfos              []AssetInfo       `json:"assetInfos"`              // 資產清單
	Histories               []decimal.Decimal `json:"histories,omitempty"`     // 總資產估值
}

type AssetInfo struct {
	userswallets.AssetForUser

	CurrenciesName string          `json:"currenciesName"` // 幣種名稱
	Valuation      decimal.Decimal `json:"valuation"`      // 當前估值(TWD)
}
