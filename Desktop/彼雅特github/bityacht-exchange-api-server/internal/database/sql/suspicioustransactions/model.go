package suspicioustransactions

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	suspicioustransactionspkg "bityacht-exchange-api-server/internal/pkg/suspicioustransactions"
	dbsql "database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// TableName of suspicious_transactions table
const TableName = "suspicious_transactions"

// Model of suspicious_transactions table
type Model struct {
	ID                       int64
	UsersID                  int64                                 `gorm:"not null"`
	Type                     suspicioustransactionspkg.Type        `gorm:"not null"`
	OrderID                  string                                `gorm:"not null;default:''"`
	Informations             suspicioustransactionspkg.Information `gorm:"type:json;not null;default:'{}'"`
	InformationReviewFiles   Files                                 `gorm:"type:json;not null;default:'null'"`
	InformationReviewComment string                                `gorm:"not null;default:''"`
	RiskReviewResult         RiskReviewResult                      `gorm:"not null;default:0"`
	RiskReviewFiles          Files                                 `gorm:"type:json;not null;default:'null'"`
	DedicatedReviewResult    DedicatedReviewResult                 `gorm:"not null;default:0"`
	DedicatedReviewComment   string                                `gorm:"not null;default:''"`
	CreatedAt                time.Time                             `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`
	ReportMJIBAt             dbsql.NullTime
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func CreateFromResults(usersID int64, orderID string, results []suspicioustransactionspkg.Result) *errpkg.Error {
	if len(results) == 0 {
		return nil
	}

	records := make([]Model, len(results))
	for i, result := range results {
		records[i] = Model{
			UsersID:               usersID,
			Type:                  result.Type,
			OrderID:               orderID,
			Informations:          result.Information,
			RiskReviewResult:      RiskReviewResultPending,
			DedicatedReviewResult: DedicatedReviewResultPending,
		}
	}

	if err := sql.DB().Create(records).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

func GetSuspiciousTXList(req GetSuspiciousTXListRequest, paginator *modelpkg.Paginator, searcher modelpkg.Searcher) ([]SuspiciousTX, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]SuspiciousTX, 0)
	query := sql.DB().
		Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `u` ON `t`.`users_id` = `u`.`id`", users.TableName)).
		Scopes(modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt))

	if req.DedicatedReviewResult > 0 {
		query = query.Where("`t`.`dedicated_review_result` = ?", int32(req.DedicatedReviewResult))
	}
	if req.Type > 0 {
		query = query.Where("`t`.`type` = ?", int32(req.Type))
	}
	query = searcher.AddToQuery(query, []string{"`t`.`users_id`", "`t`.`order_id`"}).Session(&gorm.Session{}) // Session for Reuse query

	if err := query.Select(
		"`t`.*",
		"`u`.`last_name` AS `last_name`",
		"`u`.`first_name` AS `first_name`",
		"`u`.`account` AS `email`",
	).Scopes(modelpkg.WithPaginator(paginator)).
		Order("`t`.`id` DESC").
		Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetExportCSV(req ExportSuspiciousTXCSVRequest) ([]csv.Record, *errpkg.Error) {
	var records []SuspiciousTX

	query := sql.DB().Table(TableName).Scopes(modelpkg.WithStartDateAndEnd("created_at", true, req.StartAt, req.EndAt))
	if len(req.DedicatedReviewResults) > 0 {
		query = query.Where("`dedicated_review_result` IN ?", req.DedicatedReviewResults)
	}
	if len(req.Types) > 0 {
		query = query.Where("`type` IN ?", req.Types)
	}

	if err := query.Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}

	return output, nil
}

func GetSuspiciousTXDetail(id int64) (SuspiciousTXDetail, *errpkg.Error) {
	var record SuspiciousTXDetail

	if err := sql.DB().
		Select(
			"`t`.*",
			"`u`.`last_name` AS `last_name`",
			"`u`.`first_name` AS `first_name`",
			"`u`.`account` AS `email`",
			"`u`.`phone` AS `phone`",
			"`u`.`national_id` AS `national_id`",
			"`u`.`countries_code` AS `countries_code`",
			"`u`.`dual_nationality_code` AS `dual_nationality_code`",
			"`u`.`industrial_classifications_id` AS `industrial_classifications_id`",
			"`u`.`annual_income` AS `annual_income`",
			"`u`.`funds_sources` AS `funds_sources`",
			"`u`.`purpose_of_use` AS `purpose_of_use`",
			"`u`.`created_at` AS `register_at`",
			"JSON_VALUE(`u`.`extra`, '$.RIP') AS `register_ip`").
		Where("`t`.`id` = ?", id).
		Table(fmt.Sprintf("`%s` AS `t`", TableName)).
		Joins(fmt.Sprintf("LEFT JOIN `%s` AS `u` ON `t`.`users_id` = `u`.`id`", users.TableName)).
		Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return SuspiciousTXDetail{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return SuspiciousTXDetail{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func UpdateRecord(req UpdateRequest) *errpkg.Error {
	var record Model

	if err := sql.DB().Where("`id` = ?", req.ID).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	updateMap := make(map[string]any)
	switch req.Type {
	case UpdateTypeInformationReviewComment:
		if req.Comment != record.InformationReviewComment {
			updateMap["information_review_comment"] = req.Comment
		}
	case UpdateTypeRiskReviewResult:
		if req.RiskReviewResult != record.RiskReviewResult {
			updateMap["risk_review_result"] = req.RiskReviewResult
		}
	case UpdateTypeDedicatedReview:
		if req.DedicatedReviewResult != record.DedicatedReviewResult {
			updateMap["dedicated_review_result"] = req.DedicatedReviewResult
		}
		if req.Comment != record.DedicatedReviewComment {
			updateMap["dedicated_review_comment"] = req.Comment
		}
	case UpdateTypeReportMJIBAt:
		if !record.ReportMJIBAt.Valid || !req.ReportMJIBAt.Time.Equal(record.ReportMJIBAt.Time) {
			updateMap["report_mjib_at"] = req.ReportMJIBAt
		}
	default:
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("bad type")}
	}

	if len(updateMap) == 0 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeRecordNoChange}
	}

	if err := sql.DB().Model(&record).Updates(updateMap).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

func UpdateFiles(id int64, updateType UpdateFilesType, filename string) *errpkg.Error {
	var (
		record             Model
		column             string
		oldFiles, newFiles Files
	)

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("`id` = ?", id).Scopes(modelpkg.WithLock(true)).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		switch updateType {
		case UpdateFilesTypeInformationReview:
			column = "information_review_files"
			oldFiles = record.InformationReviewFiles
			newFiles = oldFiles.Add(filename)
		case UpdateFilesTypeRiskReview:
			column = "risk_review_files"
			oldFiles = record.RiskReviewFiles
			newFiles = oldFiles.Add(filename)
		case -UpdateFilesTypeInformationReview:
			column = "information_review_files"
			oldFiles = record.InformationReviewFiles
			newFiles = oldFiles.Remove(filename)
		case -UpdateFilesTypeRiskReview:
			column = "risk_review_files"
			oldFiles = record.RiskReviewFiles
			newFiles = oldFiles.Remove(filename)
		default:
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("bad update type")}
		}

		if newLen := len(newFiles); len(oldFiles) == newLen {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeRecordNoChange}
		} else if updateType > 0 && newLen > 10 {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeTooManyFiles}
		}

		if err := tx.Model(&record).Update(column, newFiles).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}
