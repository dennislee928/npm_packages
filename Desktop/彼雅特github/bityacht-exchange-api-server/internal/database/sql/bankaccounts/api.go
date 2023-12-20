package bankaccounts

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

type BankAccount struct {
	ID int64 `json:"id"`

	// 銀行代碼
	BanksCode string `json:"banksCode"`

	// 分行代碼
	BranchsCode string `json:"branchsCode"`

	// 帳戶名稱
	Name string `json:"name"`

	// 帳號
	Account string `json:"account"`

	// 封面照片
	CoverImage string `json:"coverImage"`

	// 綁定狀態
	// * 0: 未綁定
	// * 1: 審查中
	// * 2: 已通過
	// * 3: 未通過
	Status Status `json:"status"`

	// 綁定時間
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 審核時間
	AuditTime modelpkg.DateTime `json:"auditTime" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}
