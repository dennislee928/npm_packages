package suspicioustransactions

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	suspicioustransactionspkg "bityacht-exchange-api-server/internal/pkg/suspicioustransactions"
	"errors"
	"net/http"
	"strconv"
)

type GetSuspiciousTXListRequest struct {
	// 狀態
	// * 0: 全部
	// * 1: 待審核
	// * 2: 通過
	// * 3: 駁回
	DedicatedReviewResult DedicatedReviewResult `form:"dedicatedReviewResult" binding:"gte=0,lte=3"`

	// 可疑樣態
	// * 0: 全部
	// * 1: 多筆提領態樣
	// * 2: 多筆同額態樣
	// * 3: 同一出金地址態樣
	// * 4: 迅速轉出態樣
	// * 5: 迅速買賣態樣
	// * 6: 小額接收大額轉出態樣
	Type suspicioustransactionspkg.Type `form:"type" binding:"gte=0,lte=6"`

	// 掃描時間（開始）
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 掃描時間（結束）
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type ExportSuspiciousTXCSVRequest struct {
	// 狀態 (Empty = All)
	// * 1: 待審核
	// * 2: 通過
	// * 3: 駁回
	DedicatedReviewResults []DedicatedReviewResult `form:"dedicatedReviewResults"`

	// 可疑樣態 (Empty = All)
	// * 1: 多筆提領態樣
	// * 2: 多筆同額態樣
	// * 3: 同一出金地址態樣
	// * 4: 迅速轉出態樣
	// * 5: 迅速買賣態樣
	// * 6: 小額接收大額轉出態樣
	Types []suspicioustransactionspkg.Type `form:"types"`

	// 掃描時間（開始）
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 掃描時間（結束）
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type SuspiciousTX struct {
	// 案件單號
	ID int64 `json:"id"`

	// UID
	UsersID int64 `json:"usersID"`

	// 狀態 (= 專責審查結果)
	// * 1: 待審核
	// * 2: 通過
	// * 3: 駁回
	DedicatedReviewResult DedicatedReviewResult `json:"dedicatedReviewResult"`

	// 訂單編號
	OrderID string `json:"orderID"`

	// 姓
	LastName string `json:"lastName"`

	// 名
	FirstName string `json:"firstName"`

	// E-Mail
	Email string `json:"email"`

	// 可疑樣態
	// * 1: 多筆提領態樣
	// * 2: 多筆同額態樣
	// * 3: 同一出金地址態樣
	// * 4: 迅速轉出態樣
	// * 5: 迅速買賣態樣
	// * 6: 小額接收大額轉出態樣
	Type suspicioustransactionspkg.Type `json:"type"`

	// 掃描時間
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

func GetSuspiciousTXCSVHeaders() []string {
	return []string{"案件單號", "UID", "狀態", "訂單編號", "姓", "名", "E-Mail", "可疑樣態", "掃描時間"}
}

func (s SuspiciousTX) ToCSV() []string {
	return []string{
		strconv.FormatInt(s.ID, 10),
		strconv.FormatInt(s.UsersID, 10),
		s.DedicatedReviewResult.Chinese(),
		s.OrderID,
		s.LastName,
		s.FirstName,
		s.Email,
		s.Type.Chinese(),
		s.CreatedAt.ToString(true),
	}
}

type SuspiciousTXDetail struct {
	SuspiciousTX

	// 手機
	Phone string `json:"phone"`

	// 身份證字號
	NationalID string `json:"nationalID"`

	// 國籍
	CountriesCode string `json:"countriesCode"`

	// 雙重國籍
	DualNationalityCode string `json:"dualNationalityCode"`

	// 行業別
	IndustrialClassificationsID int64 `json:"industrialClassificationsID"`

	// 年收入
	AnnualIncome string `json:"annualIncome"`

	// 資金來源
	FundsSources string `json:"fundsSources"`

	// 使用目的
	PurposeOfUse string `json:"purposeOfUse"`

	// 註冊時間
	RegisterAt modelpkg.DateTime `json:"registerAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 註冊IP
	RegisterIP string `json:"registerIP"`

	// 交易資訊
	Informations suspicioustransactionspkg.Information `json:"informations"`

	// 資訊審核(檔案)
	InformationReviewFiles Files `json:"informationReviewFiles"`

	// 資訊審核(備註)
	InformationReviewComment string `json:"informationReviewComment"`

	// 風控審查(結果)
	// * 1: 待審核
	// * 2: 非可疑交易
	// * 3: 有可疑交易風險
	RiskReviewResult RiskReviewResult `json:"riskReviewResult"`

	// 風控審查(檔案)
	RiskReviewFiles Files `json:"riskReviewFiles"`

	// 專責審查(備註)
	DedicatedReviewComment string `json:"dedicatedReviewComment"`

	// 呈報調查局日期
	ReportMJIBAt modelpkg.DateTime `json:"reportMJIBAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

type UpdateType int32

const (
	UpdateTypeInformationReviewComment UpdateType = iota + 1
	UpdateTypeRiskReviewResult
	UpdateTypeDedicatedReview
	UpdateTypeReportMJIBAt
)

type UpdateRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 更新類型
	// * 1: 資訊審核 (所需欄位：備註)
	// * 2: 風控審查結果 (所需欄位：風控審查結果)
	// * 3: 專責審查 (所需欄位：專責審查結果、備註)
	// * 4: 呈報調查局日期 (所需欄位：呈報調查局日期)
	Type UpdateType `json:"type" binding:"required,gte=1,lte=4"`

	// 備註 (資訊審核 or 專責審查)
	Comment string `json:"comment"`

	// 風控審查結果
	// * 1: 待審核
	// * 2: 非可疑交易
	// * 3: 有可疑交易風險
	RiskReviewResult RiskReviewResult `json:"riskReviewResult"`

	// 專責審查結果
	// * 1: 待審核
	// * 2: 通過
	// * 3: 駁回
	DedicatedReviewResult DedicatedReviewResult `json:"dedicatedReviewResult"`

	// 呈報調查局日期
	ReportMJIBAt modelpkg.DateTime `json:"reportMJIBAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`
}

func (ur UpdateRequest) Validate() *errpkg.Error {
	switch ur.Type {
	case UpdateTypeInformationReviewComment: // Check Nothing
	case UpdateTypeRiskReviewResult:
		if err := ur.RiskReviewResult.Validate(); err != nil {
			return err
		}
	case UpdateTypeDedicatedReview:
		if err := ur.DedicatedReviewResult.Validate(); err != nil {
			return err
		}
	case UpdateTypeReportMJIBAt:
		if ur.ReportMJIBAt.IsZero() {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad report mjib at")}
		}
	default:
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad type")}
	}

	return nil
}

type UpdateFilesType int32

const (
	UpdateFilesTypeInformationReview UpdateFilesType = iota + 1
	UpdateFilesTypeRiskReview
)
