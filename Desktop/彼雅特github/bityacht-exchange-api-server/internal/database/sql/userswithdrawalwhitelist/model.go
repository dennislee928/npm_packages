package userswithdrawalwhitelist

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// TableName of users_withdrawal_whitelist table
const TableName = "users_withdrawal_whitelist"

// Model of users_withdrawal_whitelist table
type Model struct {
	ID      int64  `gorm:"primaryKey"`
	UsersID int64  `gorm:"not null;uniqueIndex:idx_user_mainnet_addr_unique"`
	Mainnet string `gorm:"not null;uniqueIndex:idx_user_mainnet_addr_unique"`
	Address string `gorm:"not null;uniqueIndex:idx_user_mainnet_addr_unique"`
	Extra   Extra  `gorm:"type:json;not null;default:'{}'"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func CountByUserAndMainnet(usersID int64, mainnet wallet.Mainnet) (int64, *errpkg.Error) {
	var count int64
	if err := sql.DB().Table(TableName).Where("`users_id` = ? AND `mainnet` = ?", usersID, mainnet.BinanceNetwork()).Count(&count).Error; err != nil {
		return 0, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return count, nil
}

func GetByIDAndUser(id int64, usersID int64) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`id` = ? AND `users_id` = ?", id, usersID).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Model{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func GetRecordsByUser(usersID int64, paginator *modelpkg.Paginator) ([]Record, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	query := sql.DB().Table(TableName).Where("`users_id` = ?", usersID).Session(&gorm.Session{})

	records := make([]Record, 0)
	if err := query.Scopes(modelpkg.WithPaginator(paginator)).Order("`id` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetExportByUser(usersID int64) ([]csv.Record, *errpkg.Error) {
	var records []Record

	if err := sql.DB().Table(TableName).Where("`users_id` = ?", usersID).Order("`id` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for i, v := range records {
		output[i] = v
	}

	return output, nil
}

func GetRecordsByUserAndMainnet(usersID int64, mainnet wallet.Mainnet) ([]Record, *errpkg.Error) {
	records := make([]Record, 0)
	if err := sql.DB().Table(TableName).Where("`users_id` = ? AND `mainnet` = ?", usersID, mainnet.BinanceNetwork()).Order("`id` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func Create(record *Model) *errpkg.Error {
	if record == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("record is nil")}
	}

	if err := sql.DB().Create(record).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodeWalletAddressAlreadySet}
		}

		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

// Delete the record and check the address is still exist or not, if not exist -> return the record
func Delete(id int64, usersID int64) (*Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`id` = ? AND `users_id` = ?", id, usersID).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	if err := sql.DB().Delete(&record).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return &record, nil
}
