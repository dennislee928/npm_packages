package userscommissions

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	dbsql "database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TableName of users_commissions table
const TableName = "users_commissions"

// Model of users_commissions table
type Model struct {
	ID             int64
	UsersID        int64 `gorm:"not null"`
	FromUsersID    dbsql.NullInt64
	TransactionsID dbsql.NullString
	Action         Action          `gorm:"not null;default:0"`
	Amount         decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	CreatedAt      time.Time       `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetRewardByUser(usersID int64) (total decimal.Decimal, notWithdrew decimal.Decimal, err *errpkg.Error) {
	// Select total and withdrew
	if sqlErr := sql.DB().Table(fmt.Sprintf("`%s` AS `t`", TableName)).Where("`t`.`users_id` = ?", usersID).Select("COALESCE(SUM(CASE `t`.`action` WHEN ? THEN `t`.`amount` ELSE 0 END), 0) AS `total`, COALESCE(SUM(CASE `t`.`action` WHEN ? THEN `t`.`amount` ELSE 0 END), 0) AS `withdrew`", int32(ActionDeposit), int32(ActionWithdraw)).Row().Scan(&total, &notWithdrew); sqlErr != nil {
		err = &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: sqlErr}
		return
	}

	// notWithdrew = total - withdrew
	notWithdrew = total.Sub(notWithdrew)
	return
}

func GetCommissionsByUser(usersID int64, action Action, paginator *modelpkg.Paginator) ([]Commission, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	query := sql.DB().Table(fmt.Sprintf("`%s` AS `t`", TableName)).Where("`t`.`users_id` = ?", usersID)
	records := make([]Commission, 0)

	if action > 0 {
		query = query.Where("`action` = ?", action)
	}

	query = query.Session(&gorm.Session{})
	if err := query.Select("`t`.*, `u`.`account`").Joins(fmt.Sprintf("LEFT JOIN `%s` AS `u` ON `t`.`from_users_id` = `u`.`id`", users.TableName)).Order("`id` DESC").Limit(paginator.PageSize).Offset(paginator.Offset()).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func ExportCommissionsByUser(usersID int64, action Action) ([]csv.Record, *errpkg.Error) {
	var records []Commission

	query := sql.DB().Select("`t`.*, `u`.`account`").Table(fmt.Sprintf("`%s` AS `t`", TableName)).Joins(fmt.Sprintf("LEFT JOIN `%s` AS `u` ON `t`.`from_users_id` = `u`.`id`", users.TableName)).Where("`t`.`users_id` = ?", usersID).Order("`id` DESC")
	if action > 0 {
		query = query.Where("`action` = ?", action)
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

func Withdraw(usersID int64) *errpkg.Error {
	newRecord := Model{UsersID: usersID, Action: ActionWithdraw}

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Table(TableName).Scopes(modelpkg.WithLock(true)).Select("SUM(CASE `action` WHEN ? THEN `amount` WHEN ? THEN -`amount` ELSE 0 END) AS `total_amount`", int32(ActionDeposit), int32(ActionWithdraw)).Where("`users_id` = ?", usersID).Row().Scan(&newRecord.Amount); err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if newRecord.Amount.LessThan(decimal.NewFromInt(10)) {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAmount}
		}

		wallet, err := userswallets.GetOrCreateWithLock(tx, usersID, "USDT")
		if err != nil {
			return err
		}
		wallet.Deposit(newRecord.Amount)

		if err := tx.Save(&wallet).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAmount, Err: err}
		} else if err := tx.Create(&newRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAmount, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}
