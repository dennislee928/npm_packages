package walletsaddresses

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"errors"
	"fmt"
	"net/http"

	dbsql "database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TableName of wallets_address table
const TableName = "wallets_address"

// Model of wallets_address table
type Model struct {
	ID      int64
	UsersID int64            `gorm:"not null;uniqueIndex:idx_user_mainnet"`
	Mainnet string           `gorm:"not null;uniqueIndex:idx_user_mainnet;uniqueIndex:idx_mainnet_address"`
	Address dbsql.NullString `gorm:"uniqueIndex:idx_mainnet_address"`
	TxID    string           `gorm:"not null;default:''"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetByUserMainnet(userID int64, mainnet string) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().
		Where("`users_id` = ? AND `mainnet` = ?", userID, mainnet).
		Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Model{}, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeWalletAddressNotGen}
		}

		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func CreateAddress(usersID int64, mainnet, address, txID string) *errpkg.Error {
	var addr dbsql.NullString
	if address != "" {
		addr.String = address
		addr.Valid = true
	}

	if err := sql.DB().Create(&Model{
		UsersID: usersID,
		Address: addr,
		TxID:    txID,
		Mainnet: mainnet,
	}).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodeWalletAddressAlreadySet}
		}

		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

func Deposit(currenciesSymbol, mainnet, address string, val decimal.Decimal) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		walletRecord, wrapErr := takeWalletRecordWithLockingByAddress(tx, currenciesSymbol, mainnet, address)
		if wrapErr != nil {
			return wrapErr
		}

		if ok := walletRecord.Deposit(val); !ok {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadAmount}
		}

		if err := tx.Table(userswallets.TableName).Save(&walletRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

func WithdrawFailed(userID int64, currencyType wallet.CurrencyType, val decimal.Decimal) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		walletRecord, wrapErr := userswallets.GetOrCreateWithLock(tx, userID, currencyType.String())
		if wrapErr != nil {
			return wrapErr
		}

		if !walletRecord.UnlockAmount(val) {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeInsufficientBalance}
		}

		if err := tx.Table(userswallets.TableName).Save(&walletRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

func WithdrawDone(userID int64, currencyType wallet.CurrencyType, val decimal.Decimal) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		walletRecord, wrapErr := userswallets.GetOrCreateWithLock(tx, userID, currencyType.String())
		if wrapErr != nil {
			return wrapErr
		}

		if !walletRecord.Withdraw(val) {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeInsufficientBalance}
		}

		if err := tx.Table(userswallets.TableName).Save(&walletRecord).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

func GetUserByAddress(mainnet, address string) (users.Model, *errpkg.Error) {
	var userRecord users.Model

	if err := sql.DB().
		Select("`u`.*").
		Table(fmt.Sprintf("`%s` AS `u`", users.TableName)).
		Joins(fmt.Sprintf("INNER JOIN `%s` AS `t` ON `u`.`id` = `t`.`users_id`", TableName)).
		Where("`t`.`mainnet` = ? AND `t`.`address` = ?", mainnet, address).
		Take(&userRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userRecord, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}

		return userRecord, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return userRecord, nil
}

func UpdateAddress(model *Model, address string) *errpkg.Error {
	if model == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("model is nil")}
	}

	if err := sql.DB().
		Model(model).
		Update("address", address).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}

func takeWalletRecordWithLockingByAddress(tx *gorm.DB, currenciesSymbol, mainnet, address string) (*userswallets.Model, *errpkg.Error) {
	if tx == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("tx is nil")}
	}

	var record Model
	if err := tx.Where("`mainnet` = ? AND `address` = ?", mainnet, address).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeWalletAddressNotGen}
		}
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	walletRecord, wrapErr := userswallets.GetOrCreateWithLock(tx, record.UsersID, currenciesSymbol)
	if wrapErr != nil {
		return nil, wrapErr
	}

	return &walletRecord, nil
}
