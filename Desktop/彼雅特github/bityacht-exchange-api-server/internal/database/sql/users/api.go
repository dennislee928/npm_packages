package users

import (
	"bityacht-exchange-api-server/internal/database/sql/bankaccounts"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"strconv"
)

type IDURIRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0"`
}

type GetListRequest struct {
	Status  *int32        `form:"status" binding:"omitempty,gte=0,lte=5"`
	Type    *int32        `form:"type" binding:"omitempty,oneof=1 2"`
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
	EndAt   modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

// for users list
type User struct {
	// UID
	ID int64 `json:"id"`

	// 類型
	Type Type `json:"type"`

	// 狀態
	Status Status `json:"status"`

	// 會員等級
	Level int32 `json:"level"`

	// 身分證字號 or 居留證號碼 or 統一編號
	NationalID string `json:"nationalID"`

	// 護照號碼
	PassportNumber string `json:"passportNumber"`

	// 姓名 or 名稱(法人)
	FirstName string `json:"firstName"`

	// 姓氏
	LastName string `json:"lastName"`

	// E-Mail
	Account string `json:"account"`

	// 手機
	Phone string `json:"phone"`

	// 國家 or 註冊地
	CountriesCode string `json:"countriesCode"`

	// 雙重國籍
	DualNationalityCode string `json:"dualNationalityCode"`

	// 生日 or 註冊登記日
	BirthDate modelpkg.Date `json:"birthDate" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 法人性質
	JuridicalPersonNature string `json:"juridicalPersonNature"`

	// 居住地址 or 聯繫地址
	Address string `json:"address"`

	// 行業別
	IndustrialClassificationsID int64 `json:"industrialClassificationsID"`

	// 年收入
	AnnualIncome string `json:"annualIncome"`

	// 虛擬資產來源
	JuridicalPersonCryptocurrencySources string `json:"juridicalPersonCryptocurrencySources"`

	// 資金來源 or 法幣資金來源
	FundsSources string `json:"fundsSources"`

	// 投資經驗
	InvestmentExperience string `json:"investmentExperience"`

	// 使用目的
	PurposeOfUse string `json:"purposeOfUse"`

	// 被授權人姓名
	AuthorizedPersonName string `json:"authorizedPersonName"`

	// 被授權人身分證字號
	AuthorizedPersonNationalID string `json:"authorizedPersonNationalID"`

	// 被授權人聯絡電話
	AuthorizedPersonPhone string `json:"authorizedPersonPhone"`

	// 備註
	Comment string `json:"comment"`

	// 註冊時間 or 建立日期
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	Extra Extra `json:"-"`
}

func GetUserCSVHeaders() []string {
	return []string{"UID", "Status", "Level", "Email", "手機", "姓氏", "姓名", "身份", "日期"}
}

func (u User) ToCSV() []string {
	return []string{
		strconv.FormatInt(u.ID, 10),
		u.Status.Chinese(),
		strconv.Itoa(int((u.Level))),
		u.Account,
		u.Phone,
		u.LastName,
		u.FirstName,
		u.Type.Chinese(),
		u.CreatedAt.ToString(true),
	}
}

type ExportRequest struct {
	StatusList []Status      `form:"statusList"`
	StartAt    modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
	EndAt      modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type UserWithAsset struct {
	User

	MobileBarcode string               `json:"mobileBarcode"`
	Assets        []userswallets.Asset `json:"assets"`
}

type IDVStatus int32

const (
	IDVStatusNone IDVStatus = iota
	IDVStatusProcessing
	IDVStatusApproved
	IDVStatusRejected
	IDVStatusIDVRejected
)

type UserInfo struct {
	// 身份認證審查狀態
	// * 0: 未認證
	// * 1: 認證進行中
	// * 2: 已通過認證
	// * 3: 未通過認證
	// * 4: 未通過認證(可重新進行)
	IDVerificationStatus IDVStatus `json:"idVerificationStatus" binding:"required"`

	// 銀行帳戶 - 綁定狀態
	// * 0: 未綁定
	// * 1: 審查中
	// * 2: 已通過
	// * 3: 未通過
	BankAccountStatus bankaccounts.Status `json:"bankAccountStatus" binding:"required"`

	// 銀行帳戶 - 開戶銀行 (銀行代號)
	BanksCode string `json:"banksCode" binding:"required"`

	// 銀行帳戶 - 開戶銀行 (分行代號)
	BranchsCode string `json:"branchsCode" binding:"required"`

	// 銀行帳戶 - 綁定帳號
	BankAccount string `json:"bankAccount" binding:"required"`

	// 銀行帳戶 - 提交時間
	BankAccountCreatedAt modelpkg.DateTime `json:"bankAccountCreatedAt" binding:"required" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 好友邀請碼
	InviteCode string `json:"inviteCode" binding:"required"`

	// 手機載具
	MobileBarcode string `json:"mobileBarcode" binding:"required"`

	//! Deprecated (Meeting at 2023/10/2)
	// // 兩階段認證類型(Bitwise):
	// // * 0: None
	// // * 1: Email
	// // * 2: SMS
	// // * 4: Google Authenticator
	// Login2FAType TwoFAType `json:"login2FAType" binding:"required"`

	// Google Authenticator 狀態
	GoogleAuthenticator bool `json:"googleAuthenticator" binding:"required"`

	// 總受邀人數
	TotalInvited int64 `json:"totalInvited" binding:"required"`

	// 成功人數
	TotalSucceed int64 `json:"totalSucceed" binding:"required"`

	// 上次變更密碼時間
	LastChangePasswordAt modelpkg.DateTime `json:"lastChangePasswordAt" binding:"required" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// For Internal Use Only
	FinalReview       usersmodifylogs.RLStatus `json:"-"`
	IDVerificationsID int64                    `json:"-"`
}

type InviteStatus int32

const (
	InviteStatusNotFinish InviteStatus = iota + 1
	InviteStatusFinished
)

type Invitee struct {
	// 帳號
	Account string `json:"account"`

	// 狀態
	// * 1: 未完成
	// * 2: 已完成
	Status InviteStatus `json:"status"`

	// 日期
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}
