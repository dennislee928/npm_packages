package usersmodifylogs

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/managers"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"fmt"

	dbsql "database/sql"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// TableName of users_modify_logs table
const TableName = "users_modify_logs"

// Model of users_modify_logs table
type Model struct {
	ID         int64
	ManagersID dbsql.NullInt64
	UsersID    int64     `gorm:"not null;index:idx_user_time"`
	Type       Type      `gorm:"not null"`
	SubType    int32     `gorm:"not null;default:0"` //! SubType's meaning is depend on Type.
	Status     int32     `gorm:"not null;default:0"` //! Status's meaning is depend on Type.
	Comment    string    `gorm:"type:text;not null;default:''"`
	CreatedAt  time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func withDefaultSelect(db *gorm.DB) *gorm.DB {
	return db.Select("`t`.*, `m`.`name` AS `managers_name`").Order("`id` DESC")
}

func withDefaultQuery(t Type, usersID int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.
			Where("`t`.`type` = ? AND `t`.`users_id` = ?", t, usersID).
			Table(fmt.Sprintf("`%s` AS `t`", TableName)).
			Joins(fmt.Sprintf("LEFT JOIN `%s` AS `m` ON `t`.`managers_id` = `m`.`id`", managers.TableName))
	}
}

func GetStatusLogList(paginator *modelpkg.Paginator, usersID int64) ([]StatusLog, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]StatusLog, 0)
	query := sql.DB().Scopes(withDefaultQuery(TypeStatusLog, usersID)).Session(&gorm.Session{})

	if err := query.Scopes(withDefaultSelect, modelpkg.WithPaginator(paginator)).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil
}

func GetStatusLogExport(usersID int64) ([]csv.Record, *errpkg.Error) {
	var records []StatusLog

	if err := sql.DB().Scopes(withDefaultSelect, withDefaultQuery(TypeStatusLog, usersID)).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}

	return output, nil
}

func GetReviewLogList(paginator *modelpkg.Paginator, usersID int64) ([]ReviewLog, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]ReviewLog, 0)
	query := sql.DB().Scopes(withDefaultQuery(TypeReviewLog, usersID)).Session(&gorm.Session{})

	if err := query.Scopes(withDefaultSelect, modelpkg.WithPaginator(paginator)).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil
}

func GetReviewLogExport(usersID int64) ([]csv.Record, *errpkg.Error) {
	var records []ReviewLog

	if err := sql.DB().Scopes(withDefaultSelect, withDefaultQuery(TypeReviewLog, usersID)).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}
	return output, nil
}

func GetBankAccountLogList(paginator *modelpkg.Paginator, usersID int64) ([]BankAccountLog, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	records := make([]BankAccountLog, 0)
	query := sql.DB().Scopes(withDefaultQuery(TypeBankAccountLog, usersID)).Session(&gorm.Session{})

	if err := query.Scopes(withDefaultSelect, modelpkg.WithPaginator(paginator)).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil
}

func GetBankAccountLogExport(usersID int64) ([]csv.Record, *errpkg.Error) {
	var records []BankAccountLog

	if err := sql.DB().Scopes(withDefaultSelect, withDefaultQuery(TypeBankAccountLog, usersID)).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}
	return output, nil
}
