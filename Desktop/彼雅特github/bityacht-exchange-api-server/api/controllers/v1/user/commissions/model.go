package commissions

import (
	"bityacht-exchange-api-server/internal/database/sql/userscommissions"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/shopspring/decimal"
)

type GetCommissionsResponse struct {
	modelpkg.Paginator

	Data        []userscommissions.Commission `json:"data"`
	TotalReward decimal.Decimal               `json:"totalReward"` // 累計獎勵
	NotWithdrew decimal.Decimal               `json:"notWithdrew"` // 尚未提領
}
