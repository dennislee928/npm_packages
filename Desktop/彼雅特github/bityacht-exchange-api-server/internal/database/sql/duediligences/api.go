package duediligences

import (
	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"fmt"
	"strconv"
)

type GetWithDDListRequest struct {
	// 最終審核
	// * 0: 全部
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	// * 4: 身分驗證審核未通過
	// * 5: 待複核
	FinalReview usersmodifylogs.RLStatus `form:"finalReview" binding:"gte=0,lte=5"`

	// 法遵審查
	// * 0: 全部
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	ComplianceReview usersmodifylogs.RLStatus `form:"complianceReview" binding:"gte=0,lte=3"`

	// 送審日期(開始)
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 送審日期(結束)
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type Review struct {
	// 審查單號
	DueDiligencesID int64 `json:"dueDiligencesID" binding:"required"`

	// 使用者 UID
	UsersID int64 `json:"usersID" binding:"required"`

	// 審核結果
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	// * 4: 身分驗證審核未通過
	// * 5: 待複核
	FinalReview usersmodifylogs.RLStatus `json:"finalReview"`

	// 姓氏
	LastName string `json:"lastName"`

	// 姓名
	FirstName string `json:"firstName"`

	// 國籍
	CountriesCode string `json:"countriesCode"`

	// 內部風險審核
	InternalRisksTotal int64 `json:"internalRisksTotal"`

	// 法遵審查
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	ComplianceReview usersmodifylogs.RLStatus `json:"complianceReview"`

	// 審查日期
	FinalReviewTime modelpkg.DateTime `json:"finalReviewTime" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 送審日期
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

func GetReviewCSVHeaders() []string {
	return []string{"審查單號", "UID", "審核結果", "姓氏", "姓名", "國籍", "風險評估", "法遵審查", "送審日期", "審查日期"}
}

func (r Review) ToCSV() []string {
	countryString := r.CountriesCode
	if country, err := sqlcache.GetCountry(r.CountriesCode); err == nil {
		countryString = fmt.Sprintf("%s / %s", country.Chinese, country.English)
	}

	return []string{
		strconv.FormatInt(r.DueDiligencesID, 10),
		strconv.FormatInt(r.UsersID, 10),
		r.FinalReview.Chinese(),
		r.LastName,
		r.FirstName,
		countryString,
		strconv.FormatInt(r.InternalRisksTotal, 10),
		r.ComplianceReview.Chinese(),
		r.CreatedAt.ToString(true),
		r.FinalReviewTime.ToString(true),
	}
}

type GetAnnualWithDDListRequest struct {
	// * 0: 全部
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	ComplianceReview usersmodifylogs.RLStatus `form:"complianceReview" binding:"gte=0,lte=3"`

	// 送審日期(開始)
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 送審日期(結束)
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type ExportWithDDRequest struct {
	// 審核結果
	StatusList []usersmodifylogs.RLStatus `form:"statusList"`

	// 送審日期(開始)
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 送審日期(結束)
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type ExportAnnualWithDDRequest struct {
	// 送審日期(開始)
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 送審日期(結束)
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type UserWithDD struct {
	// 身分證字號 or 居留證號碼
	NationalID string `json:"nationalID"`

	// 護照號碼
	PassportNumber string `json:"passportNumber"`

	// 姓氏
	LastName string `json:"lastName"`

	// 姓名
	FirstName string `json:"firstName"`

	// 身分證 or 居留證 (正面)
	IDImage string `json:"idImage"`

	// 身分證 or 居留證 (背面)
	IDBackImage string `json:"idBackImage"`

	// 護照
	PassportImage string `json:"passportImage"`

	// 臉部照片
	FaceImage string `json:"faceImage"`

	// 身分證 or 居留證 + 臉部照片
	IDAndFaceImage string `json:"idAndFaceImage"`

	//! Deprecated (Meeting at 2023/11/01)
	// // 認證結果
	// // ResultImage string `json:"resultImage"`

	// IDV 類型
	// * 1: 人工審查（外國人）
	// * 2: KryptoGO（本國人）
	IDVType idverifications.Type `json:"idvType" gorm:"column:idv_type"`

	// IDV 單號
	IDVTaskID string `json:"idvTaskID" gorm:"column:idv_task_id"`

	// IDV 自動驗證狀態
	// * 0: 未知
	// * 1: 通過
	// * 2: 需人工驗證
	// * 3: 不通過
	// * 4: 驗證中
	// * 5: 驗證初始化
	IDVState idverifications.State `json:"idvState" gorm:"column:idv_state"`

	// 本國人: IDV 驗證結果； 外國人: 認證結果
	// * 0: 未知
	// * 1: 未驗證
	// * 2: 接受
	// * 3: 拒絕
	IDVAuditStatus idverifications.AuditStatus `json:"idvAuditStatus" gorm:"column:idv_audit_status"`

	// KryptoGO 單號
	TaskID string `json:"kryptoID"`

	// KryptoGO 風險評分
	PotentialRisk int64 `json:"kryptoPotentialRisk"`

	// KryptoGO 管制名單
	// * 0: 未知
	// * 1: 未命中
	// * 2: 命中
	SanctionMatched Bool `json:"kryptoSanctionMatched"`

	// KryptoGO 複核
	// * 0: 未複核
	// * 1: 拒絕
	// * 2: 通過
	AuditAccepted Bool `json:"kryptoAuditAccepted"`

	// 姓名檢核排除評估
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	NameCheck usersmodifylogs.RLStatus `json:"nameCheck"`

	// 姓名檢核排除評估 PDF 名稱
	NameCheckPdfName string `json:"nameCheckPdfName"`

	// 內部風險審核
	InternalRisksTotal int64 `json:"internalRisksTotal"`

	// 法遵審查
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	ComplianceReview usersmodifylogs.RLStatus `json:"complianceReview"`

	// 法遵審查備註
	ComplianceReviewComment string `json:"complianceReviewComment"`

	// 最終審查
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	// * 4: 身分驗證審核未通過
	// * 5: 待複核
	FinalReview usersmodifylogs.RLStatus `json:"finalReview"`

	// 最終審查 通知訊息
	FinalReviewNotice string `json:"finalReviewNotice"`

	// 最終審查 緣由備註
	FinalReviewComment string `json:"finalReviewComment"`

	// 手機
	Phone string `json:"phone"`

	// 生日
	BirthDate modelpkg.Date `json:"birthDate" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 國家
	CountriesCode string `json:"countriesCode"`

	// 雙重國籍
	DualNationalityCode string `json:"dualNationalityCode"`

	// 地址
	Address string `json:"address"`

	// 行業別 ID
	IndustrialClassificationsID int64 `json:"industrialClassificationsID"`

	// 年收入
	AnnualIncome string `json:"annualIncome"`

	// 資金來源
	FundsSources string `json:"fundsSources"`

	// 使用目的
	PurposeOfUse string `json:"purposeOfUse"`

	// 申請時間
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}
