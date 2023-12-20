package transactionpairs

import "github.com/shopspring/decimal"

type TransactionPair struct {
	BaseCurrenciesSymbol            string          `json:"baseCurrenciesSymbol"`
	QuoteCurrenciesSymbol           string          `json:"quoteCurrenciesSymbol"`
	SpreadsOfBuy                    decimal.Decimal `json:"spreadsOfBuy" swaggertype:"string" format:"Decimal(27, 9)"`
	SpreadsOfSell                   decimal.Decimal `json:"spreadsOfsell" swaggertype:"string" format:"Decimal(27, 9)"`
	HandlingChargeRate              decimal.Decimal `json:"handlingChargeRate" swaggertype:"string" format:"Decimal(27, 9)"`
	BaseCurrenciesDecimalPrecision  int32           `json:"baseCurrenciesDecimalPrecision"`
	QuoteCurrenciesDecimalPrecision int32           `json:"quoteCurrenciesDecimalPrecision"`
	Status                          Status          `json:"status"`
}

type UpdateRequest struct {
	BaseCurrenciesSymbol  string          `json:"baseCurrenciesSymbol"`
	QuoteCurrenciesSymbol string          `json:"quoteCurrenciesSymbol"`
	SpreadsOfBuy          decimal.Decimal `json:"spreadsOfBuy" swaggertype:"string" format:"Decimal(27, 9)"`
	SpreadsOfSell         decimal.Decimal `json:"spreadsOfsell" swaggertype:"string" format:"Decimal(27, 9)"`
	HandlingChargeRate    decimal.Decimal `json:"handlingChargeRate" swaggertype:"string" format:"Decimal(27, 9)"`
	Status                Status          `json:"status"`
}
