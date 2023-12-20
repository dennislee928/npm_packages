package invite

import (
	"bityacht-exchange-api-server/internal/database/sql/userscommissions"

	"github.com/shopspring/decimal"
)

type GetInfoRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0"`
}

type GetInfoResponse struct {
	InviteCode   string          `json:"inviteCode"`   // 邀請碼
	TotalInvited int64           `json:"totalInvited"` // 邀請人數
	TotalSucceed int64           `json:"totalSucceed"` // 達成人數
	TotalReward  decimal.Decimal `json:"totalReward"`  // 累積獎勵
	NotWithdrew  decimal.Decimal `json:"notWithdrew"`  // 尚未領取
}

type GetRequest = GetInfoRequest

type GetRewardsRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 類型
	// * 0: 全部
	// * 1: 返佣
	// * 2: 提領
	Action userscommissions.Action `form:"action" binding:"gte=0,lte=3"`
}

type GetExportRewardsHandler = GetRewardsRequest
