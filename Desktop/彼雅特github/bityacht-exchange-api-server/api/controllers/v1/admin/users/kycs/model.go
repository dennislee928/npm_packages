package userskycs

import (
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"mime/multipart"
)

type UpdateKryptoGOTaskIDRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	TaskID string `json:"taskID" binding:"required,gt=0"` // KryptoGO 單號
}

type UpdateKryptoReviewRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 狀態
	// * 2: 已通過
	// * 3: 已拒絕
	Result usersmodifylogs.RLStatus `json:"result" binding:"required,oneof=2 3"`

	// 備註
	Comment string `json:"comment"`
}

type UpdateNameCheckRequest struct {
	ID   int64                 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`
	File *multipart.FileHeader `form:"file" swaggerignore:"true"`

	// 狀態
	// * 2: 已通過
	// * 3: 已拒絕
	Result usersmodifylogs.RLStatus `form:"result" binding:"required,oneof=2 3"`
}

type UpdateComplianceReviewRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 狀態
	// * 2: 已通過
	// * 3: 已拒絕
	Result usersmodifylogs.RLStatus `json:"result" binding:"required,oneof=2 3"`

	// 備註
	Comment string `json:"comment"`
}

type UpdateFinalReviewRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 狀態
	// * 2: 已通過
	// * 3: 已拒絕
	// * 4: 身分驗證審核未通過
	Result usersmodifylogs.RLStatus `json:"result" binding:"required,oneof=2 3 4"`

	// 備註
	Comment string `json:"comment"`
}

type UpdateIDVAuditStatusRequest struct {
	ID int64 `uri:"ID" binding:"gt=0" swaggerignore:"true"`

	// 認證結果
	// * 2: 通過
	// * 3: 不通過
	AuditStatus idverifications.AuditStatus `json:"auditStatus" binding:"gte=2,lte=3"`

	// 備註
	Comment string `json:"comment"`
}
