package levellimits

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"fmt"
	"net/http"

	"github.com/shopspring/decimal"
)

type LevelLimit struct {
	Type  int32 `json:"-"`
	Level int32 `json:"-"`

	TWD  LimitByCurrency `json:"twd"`
	USDT LimitByCurrency `json:"usdt"`
}

type LimitByCurrency struct {
	// 儲值限制
	Deposit LimitByAction `json:"deposit"`

	// 提領限制
	Withdraw LimitByAction `json:"withdraw"`
}

type LimitByAction struct {
	// 單筆限額
	PerTransaction Limit `json:"perTransaction"`

	// 單日限額
	PerDay Limit `json:"perDay"`

	// 單月限額
	PerMonth Limit `json:"perMonth"`
}

func (lba LimitByAction) Validate(txAmount decimal.Decimal, accOfDay decimal.Decimal, accOfMonth decimal.Decimal) *errpkg.Error {
	if err := lba.PerTransaction.validate(txAmount); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAmount, Err: fmt.Errorf("bad tx amount (per transaction): %s", err.Error())}
	} else if err = lba.PerDay.validate(accOfDay.Add(txAmount)); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAmount, Err: fmt.Errorf("bad tx amount (per day): %s", err.Error())}
	} else if err = lba.PerMonth.validate(accOfMonth.Add(txAmount)); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAmount, Err: fmt.Errorf("bad tx amount (per month): %s", err.Error())}
	}

	return nil
}

type Limit struct {
	// 下限 (小於 0: 無限制, 等於 0: 0)
	Min decimal.Decimal `json:"min" swaggertype:"string"`

	// 上限 (小於 0: 無限制, 等於 0: 無權限)
	Max decimal.Decimal `json:"max" swaggertype:"string"`
}

func (l Limit) validate(val decimal.Decimal) error {
	if l.Min.GreaterThanOrEqual(decimal.Zero) && val.LessThan(l.Min) {
		return fmt.Errorf("%s < %s", val.String(), l.Min.String())
	} else if l.Max.GreaterThanOrEqual(decimal.Zero) && val.GreaterThan(l.Max) {
		return fmt.Errorf("%s > %s", val.String(), l.Min.String())
	}

	return nil
}
