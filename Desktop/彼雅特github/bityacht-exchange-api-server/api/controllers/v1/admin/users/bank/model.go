package bank

import (
	"bityacht-exchange-api-server/internal/database/sql/bankaccounts"
	"bityacht-exchange-api-server/internal/database/sql/bankbranchs"
	"bityacht-exchange-api-server/internal/database/sql/banks"
)

type GetRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0"`
}

type GetResponse struct {
	bankaccounts.BankAccount

	// 銀行資訊
	BankInfo banks.Bank `json:"bankInfo"`

	// 分行資訊
	BranchInfo bankbranchs.Branch `json:"branchInfo"`
}

type PatchRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 帳戶 ID
	BankAccountsID int64 `json:"id" binding:"required,gt=0"`

	// 帳戶狀態
	// * 1: 審核中
	// * 2: 已綁定
	// * 3: 未通過
	Status bankaccounts.Status `json:"status" binding:"required,gte=1,lte=3"`

	// 備註
	Comment string `json:"comment"`
}
