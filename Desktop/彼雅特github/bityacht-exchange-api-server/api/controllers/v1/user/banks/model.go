package banks

import (
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/database/sql/bankaccounts"
	"bityacht-exchange-api-server/internal/pkg/datauri"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

type UpsertAccountRequest struct {
	Name        string `json:"name" binding:"gt=0"`               // 帳戶名稱
	BanksCode   string `json:"banksCode" binding:"gt=0,number"`   // 銀行代號
	BranchsCode string `json:"branchsCode" binding:"gt=0,number"` // 分行代號
	Account     string `json:"account" binding:"gt=0,number"`     // 帳號
	CoverImage  string `json:"coverImage"`                        // 封面照片
}

func (uar UpsertAccountRequest) ToModel() (*bankaccounts.Model, *errpkg.Error) {
	record := new(bankaccounts.Model)

	if _, err := sqlcache.GetBankBranch(uar.BanksCode, uar.BranchsCode); err != nil {
		return nil, err
	}

	if len(uar.CoverImage) != 0 {
		if err := datauri.ValidateImage(uar.CoverImage); err != nil {
			return nil, err
		}
		record.CoverImage = []byte(uar.CoverImage)
	}

	record.BanksCode = uar.BanksCode
	record.BranchsCode = uar.BranchsCode
	record.Name = uar.Name
	record.Account = uar.Account

	return record, nil
}
