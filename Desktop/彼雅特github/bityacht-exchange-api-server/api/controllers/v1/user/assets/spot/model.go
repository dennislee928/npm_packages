package spot

import (
	"bityacht-exchange-api-server/internal/database/sql/userstransactions"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"
)

type CreateTransactionRequest struct {
	// 基準貨幣代號
	BaseSymbol string `json:"baseSymbol"`

	// 標價貨幣代號
	QuoteSymbol string `json:"quoteSymbol"`

	// 方向
	// * 1: 買 (花費標價 -> 獲得基準)
	// * 2: 賣 (花費基準 -> 獲得標價)
	Side userstransactions.Side `json:"side" binding:"gte=1,lte=2"`

	// 價格
	Price decimal.Decimal `json:"price"`

	// 支付數量
	PayAmount decimal.Decimal `json:"payAmount"`

	// 獲得數量
	EarnAmount decimal.Decimal `json:"earnAmount"`

	// 手續費 (單位: 獲得的幣別)
	HandlingCharge decimal.Decimal `json:"handlingCharge"`
}

func (ctr CreateTransactionRequest) Validate() *errpkg.Error {
	if ctr.Price.LessThanOrEqual(decimal.Zero) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("price should > 0")}
	}
	if ctr.PayAmount.LessThanOrEqual(decimal.Zero) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("payAmount should > 0")}
	}
	if ctr.EarnAmount.LessThanOrEqual(decimal.Zero) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("earnAmount should > 0")}
	}
	if ctr.HandlingCharge.LessThanOrEqual(decimal.Zero) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("handlingCharge should > 0")}
	}
	return nil
}
