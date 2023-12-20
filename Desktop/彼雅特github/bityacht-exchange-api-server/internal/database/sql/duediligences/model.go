package duediligences

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/kyc"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	dbsql "database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of due_diligences table
const TableName = "due_diligences"

// Model of due_diligences table
type Model struct {
	ID              int64
	UsersID         int64           `gorm:"not null;default:0"`
	Type            Type            `gorm:"not null"`
	TaskID          string          `gorm:"not null;default:''"`
	SanctionMatched Bool            `gorm:"not null;default:0"`
	PotentialRisk   int64           `gorm:"not null;default:0"`
	AuditAccepted   Bool            `gorm:"not null;default:0"`
	Comment         string          `gorm:"not null;default:''"`
	Detail          json.RawMessage `gorm:"type:json;not null;default:'{}'"` // Save the original response, don't Use Detail as Type
	CreatedAt       time.Time       `gorm:"not null;default:UTC_TIMESTAMP()"`
	AuditTime       dbsql.NullTime
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func NewFromKryptoGOTaskSummary(resp kyc.TaskSummaryResponse) (Model, *errpkg.Error) {
	var (
		record = Model{
			TaskID:        strconv.FormatInt(resp.TaskID, 10),
			PotentialRisk: int64(resp.Report.PotentialRisk),
			Comment:       resp.Report.Comment,
		}
		outputErr *errpkg.Error
	)

	if resp.Progress == 100 {
		record.SanctionMatched.Set(resp.Report.SanctionMatched)
	}

	if t := kyc.ParseKryptoGOTimestamp(resp.AuditTime); !t.IsZero() {
		record.AuditTime = dbsql.NullTime{Time: t, Valid: true}
		record.AuditAccepted.Set(resp.Report.Accepted)
	}

	return record, outputErr
}

func CreateForUser(record *Model) *errpkg.Error {
	var userRecord users.Model

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", record.UsersID).Scopes(modelpkg.WithNotDeleted(), modelpkg.WithLock(true)).Take(&userRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: err}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if err := tx.Create(record).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if err = tx.Model(&userRecord).Update("due_diligences_id", record.ID).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func UpdateTaskID(id int64, taskID string) *errpkg.Error {
	if err := sql.DB().Table(TableName).Where("`id` = ?", id).Update("task_id", taskID).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

func UpdateFromCallback(newRecord Model) *errpkg.Error {
	var record Model

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", newRecord.ID).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if record.UsersID != newRecord.UsersID {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("users id not match")}
		}

		updateMap := make(map[string]any, 7)
		if record.TaskID == "" {
			updateMap["task_id"] = newRecord.TaskID
		} else if record.TaskID != newRecord.TaskID {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("task id not match")}
		}

		if record.SanctionMatched != newRecord.SanctionMatched {
			updateMap["sanction_matched"] = newRecord.SanctionMatched
		}
		if record.PotentialRisk != newRecord.PotentialRisk {
			updateMap["potential_risk"] = newRecord.PotentialRisk
		}
		if record.AuditAccepted != newRecord.AuditAccepted {
			updateMap["audit_accepted"] = newRecord.AuditAccepted
		}
		if record.Comment != newRecord.Comment {
			updateMap["comment"] = newRecord.Comment
		}
		if newRecord.AuditTime.Valid && !newRecord.AuditTime.Time.IsZero() && (!record.AuditTime.Valid || record.AuditTime.Time.Equal(newRecord.AuditTime.Time)) {
			updateMap["audit_time"] = newRecord.AuditTime
		}
		updateMap["detail"] = newRecord.Detail

		if err := tx.Model(&record).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func GetLatestByUsersID(usersID int64) (Model, *errpkg.Error) {
	var record Model
	if err := sql.DB().Table(fmt.Sprintf("`%s` AS `u`", users.TableName)).
		Where("`u`.`deleted_at` IS NULL AND `u`.`id` = ?", usersID).
		Select("`t`.*").
		Joins(fmt.Sprintf("INNER JOIN `%s` AS `t` ON `u`.`due_diligences_id` = `t`.`id`", TableName)).
		Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}

		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func UpdateAuditInfo(managersID *int64, record *Model, accepted bool, comment string) *errpkg.Error {
	var auditAccepted Bool

	logRecord := usersmodifylogs.Model{
		UsersID: record.UsersID,
		Type:    usersmodifylogs.TypeReviewLog,
		SubType: int32(usersmodifylogs.RLTypeKryptoReview),
		Status:  int32(usersmodifylogs.RLStatusRejected),
		Comment: comment,
	}

	if managersID != nil {
		logRecord.ManagersID = dbsql.NullInt64{Int64: *managersID, Valid: true}
	}

	auditAccepted.Set(accepted)
	if accepted {
		logRecord.Status = int32(usersmodifylogs.RLStatusApproved)
	}

	updateMap := map[string]any{
		"audit_accepted": auditAccepted,
		"comment":        comment,
		"audit_time":     time.Now(),
	}

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(record).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if err := tx.Create(&logRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func GetWithDDList(paginator *modelpkg.Paginator, req GetWithDDListRequest, searcher modelpkg.Searcher) ([]Review, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]Review, 0)
	query := sql.DB().Table(fmt.Sprintf("`%s` as t", TableName)).
		Scopes(modelpkg.WithStartDateAndEnd("`t`.`created_at`", true, req.StartAt, req.EndAt)).
		Joins(fmt.Sprintf("INNER JOIN `%s` AS `u` ON `t`.`id` = `u`.`due_diligences_id`", users.TableName))

	if req.FinalReview > 0 {
		query = query.Where("`u`.`final_review` = ?", req.FinalReview)
	}
	if req.ComplianceReview > 0 {
		query = query.Where("`u`.`compliance_review` = ?", req.ComplianceReview)
	}
	query = searcher.AddToQuery(query, []string{"`t`.`users_id`", "`u`.`first_name`", "`u`.`last_name`"}).Session(&gorm.Session{})

	if err := query.Select(
		"`t`.`id` AS `due_diligences_id`",
		"`t`.`users_id` AS `users_id`",
		"`t`.`created_at` AS `created_at`",
		"`u`.`final_review` AS `final_review`",
		"`u`.`last_name` AS `last_name`",
		"`u`.`first_name` AS `first_name`",
		"`u`.`countries_code` AS `countries_code`",
		"`u`.`internal_risks_total` AS `internal_risks_total`",
		"`u`.`compliance_review` AS `compliance_review`",
		"`u`.`final_review_time` AS `final_review_time`").
		Scopes(modelpkg.WithPaginator(paginator)).
		Order("`t`.`created_at` DESC, `t`.`id` DESC").
		Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetWithDDExport(req ExportWithDDRequest) ([]csv.Record, *errpkg.Error) {
	var (
		records []Review
		query   = sql.DB().Table(fmt.Sprintf("`%s` as t", TableName)).
			Scopes(modelpkg.WithStartDateAndEnd("`t`.`created_at`", true, req.StartAt, req.EndAt)).
			Joins(fmt.Sprintf("INNER JOIN `%s` AS `u` ON `t`.`id` = `u`.`due_diligences_id`", users.TableName))
	)

	if len(req.StatusList) > 0 {
		query = query.Where("`u`.`final_review` IN ?", req.StatusList)
	}

	if err := query.Select(
		"`t`.`id` AS `due_diligences_id`",
		"`t`.`users_id` AS `users_id`",
		"`t`.`created_at` AS `created_at`",
		"`u`.`final_review` AS `final_review`",
		"`u`.`last_name` AS `last_name`",
		"`u`.`first_name` AS `first_name`",
		"`u`.`countries_code` AS `countries_code`",
		"`u`.`internal_risks_total` AS `internal_risks_total`",
		"`u`.`compliance_review` AS `compliance_review`",
		"`u`.`final_review_time` AS `final_review_time`").
		Order("`t`.`created_at` DESC, `t`.`id` DESC").
		Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}
	return output, nil
}

func GetAnnualWithDDList(paginator *modelpkg.Paginator, req GetAnnualWithDDListRequest, searcher modelpkg.Searcher) ([]Review, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]Review, 0)
	query := sql.DB().Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Joins(fmt.Sprintf("INNER JOIN `%s` AS `u` ON `t`.`id` = `u`.`due_diligences_id`", users.TableName)).
		Scopes(modelpkg.WithStartDateAndEnd("`t`.`created_at`", true, req.StartAt, req.EndAt)).
		Where("`u`.`final_review` = ?", int32(usersmodifylogs.RLStatusToBeReview))

	if req.ComplianceReview > 0 {
		query = query.Where("`u`.`compliance_review` = ?", req.ComplianceReview)
	}
	query = searcher.AddToQuery(query, []string{"`t`.`users_id`", "`u`.`first_name`", "`u`.`last_name`"}).Session(&gorm.Session{})

	if err := query.Select(
		"`t`.`id` AS `due_diligences_id`",
		"`t`.`users_id` AS `users_id`",
		"`t`.`created_at` AS `created_at`",
		"`u`.`final_review` AS `final_review`",
		"`u`.`last_name` AS `last_name`",
		"`u`.`first_name` AS `first_name`",
		"`u`.`countries_code` AS `countries_code`",
		"`u`.`internal_risks_total` AS `internal_risks_total`",
		"`u`.`compliance_review` AS `compliance_review`",
		"`u`.`final_review_time` AS `final_review_time`",
	).Scopes(modelpkg.WithPaginator(paginator)).
		Order("`t`.`created_at` DESC, `t`.`id` DESC").
		Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err = query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetAnnualWithDDExport(req ExportAnnualWithDDRequest) ([]csv.Record, *errpkg.Error) {
	records := make([]Review, 0)

	if err := sql.DB().Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Scopes(modelpkg.WithStartDateAndEnd("`t`.`created_at`", true, req.StartAt, req.EndAt)).
		Joins(fmt.Sprintf("INNER JOIN `%s` AS `u` ON `t`.`id` = `u`.`due_diligences_id`", users.TableName)).
		Where("`u`.`final_review` = ?", int32(usersmodifylogs.RLStatusToBeReview)).
		Select(
			"`t`.`id` AS `due_diligences_id`",
			"`t`.`users_id` AS `users_id`",
			"`t`.`created_at` AS `created_at`",
			"`u`.`final_review` AS `final_review`",
			"`u`.`last_name` AS `last_name`",
			"`u`.`first_name` AS `first_name`",
			"`u`.`countries_code` AS `countries_code`",
			"`u`.`internal_risks_total` AS `internal_risks_total`",
			"`u`.`compliance_review` AS `compliance_review`",
			"`u`.`final_review_time` AS `final_review_time`").
		Order("`t`.`created_at` DESC, `t`.`id` DESC").
		Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}

	return output, nil
}

func GetWithDD(usersID int64) (UserWithDD, *errpkg.Error) {
	var record UserWithDD
	if err := sql.DB().Table(fmt.Sprintf("`%s` AS `u`", users.TableName)).Select(
		"`u`.`national_id` AS `national_id`",
		"`u`.`passport_number` AS `passport_number`",
		"`u`.`last_name` AS `last_name`",
		"`u`.`first_name` AS `first_name`",
		"`u`.`name_check` AS `name_check`",
		"`u`.`name_check_pdf_name` AS `name_check_pdf_name`",
		"`u`.`internal_risks_total` AS `internal_risks_total`",
		"`u`.`compliance_review` AS `compliance_review`",
		"`u`.`compliance_review_comment` AS `compliance_review_comment`",
		"`u`.`final_review` AS `final_review`",
		"`u`.`final_review_notice` AS `final_review_notice`",
		"`u`.`final_review_comment` AS `final_review_comment`",
		"`u`.`phone` AS `phone`",
		"`u`.`birth_date` AS `birth_date`",
		"`u`.`countries_code` AS `countries_code`",
		"`u`.`dual_nationality_code` AS `dual_nationality_code`",
		"`u`.`address` AS `address`",
		"`u`.`industrial_classifications_id` AS `industrial_classifications_id`",
		"`u`.`annual_income` AS `annual_income`",
		"`u`.`funds_sources` AS `funds_sources`",
		"`u`.`purpose_of_use` AS `purpose_of_use`",
		"`t`.`task_id` AS `task_id`",
		"`t`.`potential_risk` AS `potential_risk`",
		"`t`.`sanction_matched` AS `sanction_matched`",
		"`t`.`audit_accepted` AS `audit_accepted`",
		"`i`.`type` AS `idv_type`",
		"`i`.`task_id` AS `idv_task_id`",
		"`i`.`id_image` AS `id_image`",
		"`i`.`id_back_image` AS `id_back_image`",
		"`i`.`passport_image` AS `passport_image`",
		"`i`.`face_image` AS `face_image`",
		"`i`.`id_and_face_image` AS `id_and_face_image`",
		// "`i`.`result_image` AS `result_image`", //! Deprecated (Meeting at 2023/11/01)
		"`i`.`state` AS `idv_state`",
		"`i`.`audit_status` AS `idv_audit_status`",
		"`i`.`created_at` AS `created_at`").
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `t` ON `u`.`due_diligences_id` = `t`.`id`", TableName)).
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `i` ON `u`.`id_verifications_id` = `i`.`id`", idverifications.TableName)).
		Where("`u`.`id` = ?", usersID).
		Take(&record).Error; err != nil {
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func ResentKryptoGo(ctx *gin.Context, usersID int64) (kyc.CreateTasksParams, *Model, *errpkg.Error) {
	var (
		createTasksParams kyc.CreateTasksParams
		userRecord        users.Model
		ddRecord          = &Model{
			UsersID: usersID,
			Type:    TypeManualResend,
		}
	)

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("`id` = ? AND `deleted_at` IS NULL", usersID).Take(&userRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if !userRecord.IDVerificationsID.Valid {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("user id verifications id is null")}
		} else if userRecord.CountriesCode.String != "TWN" {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("user country is not Taiwan")}
		}

		if err := tx.Table(TableName).Create(ddRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if err := tx.Model(&userRecord).Update("due_diligences_id", ddRecord.ID).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		createTasksParams = kyc.CreateTasksParams{
			Name:              userRecord.LastName + userRecord.FirstName,
			Birthday:          userRecord.BirthDate,
			Citizenship:       userRecord.CountriesCode.String,
			CallBackURL:       kyc.GetKryptoGODDCallbackURL(usersID, ddRecord.ID),
			CustomerReference: fmt.Sprintf("dds_%020d", ddRecord.ID),
		}
		return nil

	}); err != nil {
		return createTasksParams, nil, err.(*errpkg.Error)
	}

	return createTasksParams, ddRecord, nil
}

func CreateIDVAndDD(nationalID string, checkPhoneNum string, idvRecord *idverifications.Model, ddRecord *Model, userRecordUpdateMap map[string]any) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		userRecord, err := idverifications.CreateAndGetUserWithTx(tx, idvRecord, nationalID, userRecordUpdateMap)
		if err != nil {
			return err
		}

		if err := tx.Create(ddRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		userRecordUpdateMap["due_diligences_id"] = ddRecord.ID

		if err := tx.Model(userRecord).Updates(userRecordUpdateMap).Error; err != nil {
			var sqlErr *mysql.MySQLError

			if errors.As(err, &sqlErr) && sqlErr.Number == 1062 {
				return &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodeNationalIDDuplicated}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if checkPhoneNum != "" {
			var countByPhone int64
			if err := tx.Table(users.TableName).Where("`phone` = ?", checkPhoneNum).Count(&countByPhone).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			} else if countByPhone != 1 {
				return &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodePhoneNumberDuplicated}
			}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}
