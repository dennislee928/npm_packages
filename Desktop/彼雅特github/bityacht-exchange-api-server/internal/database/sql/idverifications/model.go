package idverifications

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/datauri"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	dbsql "database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of id_verifications table
const TableName = "id_verifications"

type Model struct {
	ID             int64
	UsersID        int64           `gorm:"not null;index"`
	Type           Type            `gorm:"not null"`
	TaskID         string          `gorm:"not null;default:''"`
	IDImage        []byte          `gorm:"type:mediumblob"`
	IDBackImage    []byte          `gorm:"type:mediumblob"`
	PassportImage  []byte          `gorm:"type:mediumblob"`
	FaceImage      []byte          `gorm:"type:mediumblob"`
	IDAndFaceImage []byte          `gorm:"type:mediumblob"`
	ResultImage    []byte          `gorm:"type:mediumblob"`
	State          State           `gorm:"not null;default:0"`
	AuditStatus    AuditStatus     `gorm:"not null;default:0"`
	Detail         json.RawMessage `gorm:"type:json;not null;default:'{}'"` // Save the original response, don't Use Detail as Type
	CreatedAt      time.Time       `gorm:"not null;default:UTC_TIMESTAMP()"`
	AuditTime      dbsql.NullTime
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func (m Model) UpdateImagesByURL(ctx *gin.Context, req UpdateImagesByURLRequest, updateMap map[string]any) {
	updateInfos := []struct {
		UpdateMapKey string
		ImageURL     string
		OldImageSize int
	}{
		{"id_image", req.IDImage, len(m.IDImage)},
		{"id_back_image", req.IDBackImage, len(m.IDBackImage)},
		{"passport_image", req.PassportImage, len(m.PassportImage)},
		{"face_image", req.FaceImage, len(m.FaceImage)},
		{"id_and_face_image", req.IDAndFaceImage, len(m.IDAndFaceImage)},
	}

	errLogger := logger.GetGinRequestLogger(ctx)
	for _, info := range updateInfos {
		if info.ImageURL == "" || info.OldImageSize != 0 {
			continue
		} else if data := datauri.DownloadImage(errLogger, info.ImageURL); len(data) != 0 {
			updateMap[info.UpdateMapKey] = data
		}
	}
}

func GetByID(id int64) (Model, *errpkg.Error) {
	var record Model
	if err := sql.DB().Where("`id` = ?", id).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Model{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func GetLatestByUsersID(tx *gorm.DB, usersID int64) (Model, *errpkg.Error) {
	var record Model

	if tx == nil {
		tx = sql.DB()
	}

	if err := tx.Select("`t`.*").
		Table(fmt.Sprintf("`%s` AS `u`", users.TableName)).
		Where("`u`.`id` = ?", usersID).
		Joins(fmt.Sprintf("INNER JOIN `%s` AS `t` ON `u`.`id_verifications_id` = `t`.`id`", TableName)).
		Take(&record).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

//! Deprecated (Meeting at 2023/11/01)
// func UpdateResultImageByUsersID(managersID int64, usersID int64, resultImage []byte) *errpkg.Error {
// 	reviewLog := usersmodifylogs.Model{
// 		ManagersID: dbsql.NullInt64{Int64: managersID, Valid: true},
// 		UsersID:    usersID,
// 		Type:       usersmodifylogs.TypeReviewLog,
// 		SubType:    int32(usersmodifylogs.RLTypeUploadResultImage),
// 	}

// 	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
// 		record, err := GetLatestByUsersID(tx.Scopes(modelpkg.WithLock(true)), usersID)
// 		if err != nil {
// 			return err
// 		}

// 		if err := tx.Model(&record).Update("result_image", resultImage).Error; err != nil {
// 			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
// 		} else if err := tx.Create(&reviewLog).Error; err != nil {
// 			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
// 		}

// 		return nil
// 	}); err != nil {
// 		return err.(*errpkg.Error)
// 	}

// 	return nil
// }

func UpdateAuditStatusByUser(managersID int64, usersID int64, auditStatus AuditStatus, comment string) *errpkg.Error {
	reviewLog := usersmodifylogs.Model{
		ManagersID: dbsql.NullInt64{Int64: managersID, Valid: true},
		UsersID:    usersID,
		Type:       usersmodifylogs.TypeReviewLog,
		SubType:    int32(usersmodifylogs.RLTypeIDVStatus),
		Status:     int32(auditStatus.ToRLStatus()),
		Comment:    comment,
	}

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		record, err := GetLatestByUsersID(tx.Scopes(modelpkg.WithLock(true)), usersID)
		if err != nil {
			return err
		}

		if record.Type != TypeManual {
			return &errpkg.Error{HttpStatus: http.StatusForbidden, Code: errpkg.CodeBadAction, Err: errors.New("bad type")}
		}

		if record.AuditStatus == auditStatus {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeRecordNoChange}
		}

		if err := tx.Model(&record).Update("audit_status", auditStatus).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if err = tx.Create(&reviewLog).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func CreateAndGetUserWithTx(tx *gorm.DB, record *Model, nationalID string, userRecordUpdateMap map[string]any) (users.Model, *errpkg.Error) {
	var userRecord users.Model

	if tx == nil {
		return userRecord, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("tx is nil")}
	} else if err := tx.Where("`id` = ? AND `deleted_at` IS NULL", record.UsersID).Clauses(clause.Locking{Strength: "UPDATE"}).Take(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userRecord, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRecordNotFound}
		}
		return userRecord, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	userRecordUpdateMap["national_id"] = nationalID
	if record.Type == TypeManual {
		if userRecord.NameCheck == usersmodifylogs.RLStatusUnknown {
			userRecordUpdateMap["name_check"] = usersmodifylogs.RLStatusPending
		}
		if userRecord.ComplianceReview == usersmodifylogs.RLStatusUnknown {
			userRecordUpdateMap["compliance_review"] = usersmodifylogs.RLStatusPending
		}
		if userRecord.FinalReview == usersmodifylogs.RLStatusUnknown {
			userRecordUpdateMap["final_review"] = usersmodifylogs.RLStatusPending
		}
	}

	if err := tx.Create(record).Error; err != nil {
		return userRecord, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	userRecordUpdateMap["id_verifications_id"] = record.ID

	return userRecord, nil
}

func UpdateTaskID(id int64, taskID string) *errpkg.Error {
	if err := sql.DB().Table(TableName).Where("`id` = ?", id).Update("task_id", taskID).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

func UpdateFromCallback(ctx *gin.Context, idvType Type, newRecord Model, updateImageReq UpdateImagesByURLRequest) (bool, *errpkg.Error) {
	var record Model
	var needSendNotify bool

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", newRecord.ID).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if record.UsersID != newRecord.UsersID {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("users id not match")}
		} else if record.Type != idvType {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("record type not match")}
		}

		updateMap := make(map[string]any, 9)
		if record.TaskID == "" {
			updateMap["task_id"] = newRecord.TaskID
		} else if record.TaskID != newRecord.TaskID {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("task id not match")}
		}
		record.UpdateImagesByURL(ctx, updateImageReq, updateMap)
		needSendNotify = record.AuditStatus == AuditStatusUnknown &&
			((record.State == StateUnknown && string(record.Detail) == "{}") || // First Callback
				newRecord.AuditStatus == AuditStatusRejected) // Unknown -> Rejected

		if newRecord.State != StateUnknown && record.State != newRecord.State {
			updateMap["state"] = newRecord.State
		}
		if record.AuditStatus != newRecord.AuditStatus {
			updateMap["audit_status"] = newRecord.AuditStatus
		}
		if newRecord.AuditTime.Valid && (!record.AuditTime.Valid || !newRecord.AuditTime.Time.Equal(record.AuditTime.Time)) {
			updateMap["audit_time"] = newRecord.AuditTime
		}
		updateMap["detail"] = newRecord.Detail

		if newRecord.State == StateReject || newRecord.AuditStatus == AuditStatusRejected {
			if err := tx.Table(users.TableName).Where("`id` = ?", record.UsersID).Update("final_review", int32(usersmodifylogs.RLStatusIDVRejected)).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
		}

		if err := tx.Model(&record).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return false, err.(*errpkg.Error)
	}

	return needSendNotify, nil
}

func GetIDVStatus(finalReview usersmodifylogs.RLStatus, ddsID dbsql.NullInt64) (users.IDVStatus, *errpkg.Error) {
	switch finalReview {
	case usersmodifylogs.RLStatusUnknown:
		return users.IDVStatusNone, nil
	case usersmodifylogs.RLStatusApproved, usersmodifylogs.RLStatusToBeReview:
		return users.IDVStatusApproved, nil
	case usersmodifylogs.RLStatusRejected:
		return users.IDVStatusRejected, nil
	case usersmodifylogs.RLStatusIDVRejected:
		return users.IDVStatusIDVRejected, nil
	case usersmodifylogs.RLStatusPending:
		if !ddsID.Valid || ddsID.Int64 == 0 {
			return users.IDVStatusNone, nil
		}

		record, err := GetByID(ddsID.Int64)
		if err != nil {
			return users.IDVStatusNone, err
		} else if record.State == StateReject {
			return users.IDVStatusIDVRejected, nil
		}

		switch record.AuditStatus {
		case AuditStatusUnknown:
			if time.Since(record.CreatedAt) >= 2*time.Hour {
				return users.IDVStatusNone, nil
			}

			return users.IDVStatusProcessing, nil
		case AuditStatusPending, AuditStatusAccepted:
			return users.IDVStatusProcessing, nil
		case AuditStatusRejected:
			return users.IDVStatusIDVRejected, nil
		}

		return users.IDVStatusNone, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadRecord, Err: errors.New("bad idv audit status")}
	}

	return users.IDVStatusNone, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadRecord, Err: errors.New("bad user final review")}
}
