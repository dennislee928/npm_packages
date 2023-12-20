package users

import (
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

type CreateRequest struct {
	// Email
	Account string `json:"account" binding:"required,email"`

	// 名稱
	Name string `json:"name" binding:"required"`

	// 統一編號
	NationalID string `json:"nationalID"`

	// 註冊地
	Country string `json:"country"`

	// 註冊登記日
	BirthDate modelpkg.Date `json:"birthDate" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 法人性質
	JuridicalPersonNature string `json:"juridicalPersonNature"`

	// 營業地址
	Address string `json:"address"`

	// 聯絡電話
	Phone string `json:"phone"`

	// 行業別
	IndustrialClassificationsID int64 `json:"industrialClassificationsID"`

	// 法幣資金來源
	JuridicalPersonFundsSources string `json:"juridicalPersonFundsSources"`

	// 虛擬資產來源
	JuridicalPersonCryptocurrencySources string `json:"juridicalPersonCryptocurrencySources"`

	// 被授權人姓名
	AuthorizedPersonName string `json:"authorizedPersonName"`

	// 被授權人身分證字號
	AuthorizedPersonNationalID string `json:"authorizedPersonNationalID"`

	// 被授權人聯絡電話
	AuthorizedPersonPhone string `json:"authorizedPersonPhone"`

	// 其他資訊（備註）
	Comment string `json:"comment"`
}

type UpdateLevelRequest struct {
	ID    int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
	Level int32 `json:"level" binding:"required,gte=1,lte=5"` // 等級 (法人: 1~2, 自然人: 2~5)
}

type UpdateStatusRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 狀態
	// * 0: 未啟用
	// * 1: 已啟用
	// * 2: 已停用
	// * 3: 凍結中
	Status usersmodifylogs.SLStatus `json:"status" binding:"gte=0,lte=3"`

	// 備註
	Comment string `json:"comment"`
}

type DeleteWithdrawalWhitelistRequest struct {
	ID      int64 `uri:"WhitelistID" binding:"required,gt=0" swaggerignore:"true"`
	UsersID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
}
