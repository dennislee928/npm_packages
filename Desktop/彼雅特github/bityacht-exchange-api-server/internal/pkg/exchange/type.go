package exchange

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"net/http"

	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
)

type SymbolStatus int32

const (
	SymbolStatusUnknown SymbolStatus = iota
	SymbolStatusPreTrading
	SymbolStatusTrading
	SymbolStatusPostTrading
	SymbolStatusEndOfDay
	SymbolStatusHalt
	SymbolStatusAuctionMatch
	SymbolStatusBreak
)

type ExchangeInfo struct {
	Symbol                     string       `json:"symbol"`
	Status                     SymbolStatus `json:"status"`
	BaseAsset                  string       `json:"baseAsset"`
	BaseAssetPrecision         int          `json:"baseAssetPrecision"`
	QuoteAsset                 string       `json:"quoteAsset"`
	QuotePrecision             int          `json:"quotePrecision"`
	QuoteAssetPrecision        int          `json:"quoteAssetPrecision"`
	BaseCommissionPrecision    int32        `json:"baseCommissionPrecision"`
	QuoteCommissionPrecision   int32        `json:"quoteCommissionPrecision"`
	OrderTypes                 []string     `json:"orderTypes"`
	IcebergAllowed             bool         `json:"icebergAllowed"`
	OcoAllowed                 bool         `json:"ocoAllowed"`
	QuoteOrderQtyMarketAllowed bool         `json:"quoteOrderQtyMarketAllowed"`
	IsSpotTradingAllowed       bool         `json:"isSpotTradingAllowed"`
	IsMarginTradingAllowed     bool         `json:"isMarginTradingAllowed"`
	Permissions                []string     `json:"permissions"`

	MinNotional decimal.Decimal `json:"minNotional"`
	MaxNotional decimal.Decimal `json:"maxNotional"`
	MinPrice    decimal.Decimal `json:"minPrice"`
	MaxPrice    decimal.Decimal `json:"maxPrice"`
	TickSize    decimal.Decimal `json:"tickSize"`
	MinQuantity decimal.Decimal `json:"minQty"`
	MaxQuantity decimal.Decimal `json:"maxQty"`
	StepSize    decimal.Decimal `json:"stepSize"`
}

const (
	BinanceSymbolStatusPreTrading   = "PRE_TRADING"
	BinanceSymbolStatusTrading      = "TRADING"
	BinanceSymbolStatusPostTrading  = "POST_TRADING"
	BinanceSymbolStatusEndOfDay     = "END_OF_DAY"
	BinanceSymbolStatusHalt         = "HALT"
	BinanceSymbolStatusAuctionMatch = "AUCTION_MATCH"
	BinanceSymbolStatusBreak        = "BREAK"
)

func parseDecimalFromBinance(val interface{}) (decimal.Decimal, *errpkg.Error) {
	if val, ok := val.(string); !ok {
		return decimal.Zero, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI}
	} else if val == "" {
		return decimal.Zero, nil
	} else if decimalVal, err := decimal.NewFromString(val); err != nil {
		return decimal.Zero, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: err}
	} else {
		return decimalVal, nil
	}
}

func (ei *ExchangeInfo) fromBinance(src binance.Symbol) *errpkg.Error {
	ei.Symbol = src.Symbol
	ei.BaseAsset = src.BaseAsset
	ei.BaseAssetPrecision = src.BaseAssetPrecision
	ei.QuoteAsset = src.QuoteAsset
	ei.QuotePrecision = src.QuotePrecision
	ei.QuoteAssetPrecision = src.QuoteAssetPrecision
	ei.BaseCommissionPrecision = src.BaseCommissionPrecision
	ei.QuoteCommissionPrecision = src.QuoteCommissionPrecision
	ei.OrderTypes = src.OrderTypes
	ei.IcebergAllowed = src.IcebergAllowed
	ei.OcoAllowed = src.OcoAllowed
	ei.QuoteOrderQtyMarketAllowed = src.QuoteOrderQtyMarketAllowed
	ei.IsSpotTradingAllowed = src.IsSpotTradingAllowed
	ei.IsMarginTradingAllowed = src.IsMarginTradingAllowed
	ei.Permissions = src.Permissions

	switch src.Status {
	case BinanceSymbolStatusPreTrading:
		ei.Status = SymbolStatusPreTrading
	case BinanceSymbolStatusTrading:
		ei.Status = SymbolStatusTrading
	case BinanceSymbolStatusPostTrading:
		ei.Status = SymbolStatusPostTrading
	case BinanceSymbolStatusEndOfDay:
		ei.Status = SymbolStatusEndOfDay
	case BinanceSymbolStatusHalt:
		ei.Status = SymbolStatusHalt
	case BinanceSymbolStatusAuctionMatch:
		ei.Status = SymbolStatusAuctionMatch
	case BinanceSymbolStatusBreak:
		ei.Status = SymbolStatusBreak
	default:
		ei.Status = SymbolStatusUnknown
	}

	var err *errpkg.Error
	for _, filter := range src.Filters {
		filterType, ok := filter["filterType"].(string)
		if !ok {
			continue
		}

		switch filterType {
		case "PRICE_FILTER":
			if ei.MinPrice, err = parseDecimalFromBinance(filter["minPrice"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad min price")
				}
				return err
			} else if ei.MaxPrice, err = parseDecimalFromBinance(filter["maxPrice"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad max price")
				}
				return err
			} else if ei.TickSize, err = parseDecimalFromBinance(filter["tickSize"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad tick size")
				}
				return err
			}
		case "NOTIONAL":
			if ei.MinNotional, err = parseDecimalFromBinance(filter["minNotional"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad min notional")
				}
				return err
			} else if ei.MaxNotional, err = parseDecimalFromBinance(filter["maxNotional"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad max notional")
				}
				return err
			}
		case "LOT_SIZE":
			if ei.MinQuantity, err = parseDecimalFromBinance(filter["minQty"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad min qty")
				}
			} else if ei.MaxQuantity, err = parseDecimalFromBinance(filter["maxQty"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad max qty")
				}
			} else if ei.StepSize, err = parseDecimalFromBinance(filter["stepSize"]); err != nil {
				if err.Err == nil {
					err.Err = errors.New("bad step size")
				}
			}
		}
	}

	return nil
}

type BookTicker struct {
	Symbol      string          `json:"symbol"`
	BidPrice    decimal.Decimal `json:"bidPrice"`
	BidQuantity decimal.Decimal `json:"bidQty"`
	AskPrice    decimal.Decimal `json:"askPrice"`
	AskQuantity decimal.Decimal `json:"askQty"`
}

func (bt *BookTicker) fromBinance(src *binance.BookTicker) *errpkg.Error {
	if src == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("nil book ticker")}
	}

	bt.Symbol = src.Symbol

	var err *errpkg.Error
	if bt.BidPrice, err = parseDecimalFromBinance(src.BidPrice); err != nil {
		return err
	} else if bt.BidQuantity, err = parseDecimalFromBinance(src.BidQuantity); err != nil {
		return err
	} else if bt.AskPrice, err = parseDecimalFromBinance(src.AskPrice); err != nil {
		return err
	} else if bt.AskQuantity, err = parseDecimalFromBinance(src.AskQuantity); err != nil {
		return err
	}

	return nil
}

func (bt *BookTicker) fromMax(symbol string, src *MaxTicker) *errpkg.Error {
	if src == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("nil book ticker")}
	}

	bt.Symbol = symbol
	bt.BidPrice = src.Buy
	bt.BidQuantity = src.Vol
	bt.AskPrice = src.Sell
	bt.AskQuantity = src.Vol

	return nil
}

type CreateOrderResponse struct {
	Symbol                   string          `json:"symbol"`
	OrderID                  int64           `json:"orderId"`
	ClientOrderID            string          `json:"clientOrderId"`
	TransactTime             int64           `json:"transactTime"`
	Price                    decimal.Decimal `json:"price"`
	OrigQuantity             decimal.Decimal `json:"origQty"`
	ExecutedQuantity         decimal.Decimal `json:"executedQty"`
	CummulativeQuoteQuantity decimal.Decimal `json:"cummulativeQuoteQty"`
	IsIsolated               bool            `json:"isIsolated"` // for isolated margin

	Status      binance.OrderStatusType `json:"status"`
	TimeInForce binance.TimeInForceType `json:"timeInForce"`
	Type        binance.OrderType       `json:"type"`
	Side        binance.SideType        `json:"side"`

	// for order response is set to FULL
	Fills                 []Fill          `json:"fills"`
	MarginBuyBorrowAmount decimal.Decimal `json:"marginBuyBorrowAmount"` // for margin
	MarginBuyBorrowAsset  decimal.Decimal `json:"marginBuyBorrowAsset"`
}

type Fill struct {
	TradeID         int             `json:"tradeId"`
	Price           decimal.Decimal `json:"price"`
	Quantity        decimal.Decimal `json:"qty"`
	Commission      decimal.Decimal `json:"commission"`
	CommissionAsset string          `json:"commissionAsset"`
}

func (cor *CreateOrderResponse) fromBinance(src *binance.CreateOrderResponse) *errpkg.Error {
	if src == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("nil create order response")}
	}

	cor.Symbol = src.Symbol
	cor.OrderID = src.OrderID
	cor.ClientOrderID = src.ClientOrderID
	cor.TransactTime = src.TransactTime
	cor.IsIsolated = src.IsIsolated
	cor.Status = src.Status
	cor.TimeInForce = src.TimeInForce
	cor.Type = src.Type
	cor.Side = src.Side

	var err *errpkg.Error
	cor.Fills = make([]Fill, len(src.Fills))
	for index, fill := range src.Fills {
		target := &cor.Fills[index]
		target.TradeID = fill.TradeID
		target.CommissionAsset = fill.CommissionAsset

		if target.Price, err = parseDecimalFromBinance(fill.Price); err != nil {
			return err
		} else if target.Quantity, err = parseDecimalFromBinance(fill.Quantity); err != nil {
			return err
		} else if target.Commission, err = parseDecimalFromBinance(fill.Commission); err != nil {
			return err
		}
	}

	if cor.Price, err = parseDecimalFromBinance(src.Price); err != nil {
		return err
	} else if cor.OrigQuantity, err = parseDecimalFromBinance(src.OrigQuantity); err != nil {
		return err
	} else if cor.ExecutedQuantity, err = parseDecimalFromBinance(src.ExecutedQuantity); err != nil {
		return err
	} else if cor.CummulativeQuoteQuantity, err = parseDecimalFromBinance(src.CummulativeQuoteQuantity); err != nil {
		return err
	} else if cor.MarginBuyBorrowAmount, err = parseDecimalFromBinance(src.MarginBuyBorrowAmount); err != nil {
		return err
	} else if cor.MarginBuyBorrowAsset, err = parseDecimalFromBinance(src.MarginBuyBorrowAsset); err != nil {
		return err
	}

	return nil
}
