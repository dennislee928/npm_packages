package usersmodifylogs

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
)

// RLStatus is Status of ReviewLog
type RLStatus int32

// Status List of ReviewLog
const (
	RLStatusUnknown RLStatus = iota
	RLStatusPending
	RLStatusApproved
	RLStatusRejected
	RLStatusIDVRejected
	RLStatusToBeReview
	RLIDVStatusPending
	RLIDVStatusAccepted
	RLIDVStatusRejected
)

func (s RLStatus) Chinese() string {
	switch s {
	case RLStatusUnknown:
		return "未知"
	case RLStatusPending, RLIDVStatusPending:
		return "待審核"
	case RLStatusApproved:
		return "已通過"
	case RLStatusRejected:
		return "已拒絕"
	case RLStatusIDVRejected:
		return "身分驗證審核未通過"
	case RLStatusToBeReview:
		return "待複核"
	case RLIDVStatusAccepted:
		return "通過"
	case RLIDVStatusRejected:
		return "未通過"
	}

	return badStatusChinese
}

// RLType is SubType of ReviewLog
type RLType int32

// SubType List of ReviewLog
const (
	RLTypeKryptoReview RLType = iota + 1
	RLTypeNameCheckUploadPDF
	RLTypeInternalReview
	RLTypeComplianceReview
	RLTypeFinalReview
	RLTypeUploadResultImage
	RLTypeIDVStatus
)

func (t RLType) Chinese() string {
	switch t {
	case RLTypeKryptoReview:
		return "KryptoGO"
	case RLTypeNameCheckUploadPDF:
		return "姓名檢核排除評估"
	case RLTypeInternalReview:
		return "內部風險審核"
	case RLTypeComplianceReview:
		return "法遵審查"
	case RLTypeFinalReview:
		return "最終審查"
	case RLTypeUploadResultImage:
		return "上傳認證結果"
	case RLTypeIDVStatus:
		return "認證結果"
	}

	return badStatusChinese
}

type ReviewLog struct {
	// 日期
	CreatedAt modelpkg.DateTime `json:"createdAt" swaggertype:"string" format:"dateTime(YYYY/MM/DD HH:mm:SS)"`

	// 類別
	// * 1: KryptoGO 複核
	// * 2: 姓名檢核排除評估
	// * 3: 內部風險審核
	// * 4: 法遵審查
	// * 5: 最終審查
	// * 6: 上傳認證結果(外籍人士)
	// * 7: 認證結果(外籍人士)
	SubType RLType `json:"type"`

	// 狀態
	// * 0: 未知
	// * 1: 待審核
	// * 2: 已通過
	// * 3: 已拒絕
	// * 4: 身分驗證審核未通過
	// * 5: 待複核
	// * 6: 待審核
	// * 7: 通過
	// * 8: 未通過
	Status RLStatus `json:"status"`

	// 備註
	Comment string `json:"comment"`

	// 管理者
	ManagersName string `json:"managersName"`
}

func GetReviewLogCSVHeaders() []string {
	return []string{"日期", "類別", "狀態", "備註", "管理者"}
}

func (rl ReviewLog) ToCSV() []string {
	var status string
	if rl.SubType != RLTypeUploadResultImage {
		status = rl.Status.Chinese()
	}

	return []string{
		rl.CreatedAt.ToString(true),
		rl.SubType.Chinese(),
		status,
		rl.Comment,
		rl.ManagersName,
	}
}
