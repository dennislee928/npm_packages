package userscommissions

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/shopspring/decimal"
)

type Commission struct {
	// 帳號
	Account string `json:"account"`

	// 類型
	// * 1: 返佣
	// * 2: 提領
	Action Action `json:"action"`

	// 數量 (USDT)
	Amount decimal.Decimal `json:"amount"`

	// 日期
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

func GetCommissionCSVHeaders() []string {
	return []string{"日期", "類型", "數量", "帳號"}
}

func (c Commission) ToCSV() []string {
	return []string{
		c.CreatedAt.ToString(true),
		c.Action.Chinese(),
		c.Amount.String(),
		c.Account,
	}
}
