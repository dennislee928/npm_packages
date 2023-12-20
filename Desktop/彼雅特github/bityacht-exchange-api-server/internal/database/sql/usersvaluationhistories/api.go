package usersvaluationhistories

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/shopspring/decimal"
)

type History struct {
	Date      modelpkg.Date
	Valuation decimal.Decimal
}
