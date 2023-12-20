package userswallets

import (
	"bityacht-exchange-api-server/internal/cache/memory/spottrend"
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/currencies"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"errors"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of users_wallets table
const TableName = "users_wallets"

// Model of users_wallets table
type Model struct {
	ID               int64
	UsersID          int64           `gorm:"not null;uniqueIndex:idx_users_currency,priority:1"`
	CurrenciesSymbol string          `gorm:"not null;uniqueIndex:idx_users_currency,priority:2"`
	Type             currencies.Type `gorm:"not null;default:0"`
	FreeAmount       decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	LockedAmount     decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	DeletedAt        *time.Time
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func (m *Model) Deposit(val decimal.Decimal) bool {
	if val.IsNegative() {
		return false
	}

	if val.IsZero() {
		return true
	}

	m.FreeAmount = m.FreeAmount.Add(val)

	return true
}

func (m *Model) Withdraw(val decimal.Decimal) bool {
	if val.IsNegative() || m.LockedAmount.LessThan(val) {
		return false
	}

	if val.IsZero() {
		return true
	}

	m.LockedAmount = m.LockedAmount.Sub(val)

	return true
}

func (m *Model) LockAmount(val decimal.Decimal) bool {
	if val.IsNegative() || m.FreeAmount.LessThan(val) {
		return false
	}

	if val.IsZero() {
		return true
	}

	m.FreeAmount = m.FreeAmount.Sub(val)
	m.LockedAmount = m.LockedAmount.Add(val)

	return true
}

func (m *Model) UnlockAmount(val decimal.Decimal) bool {
	if val.IsNegative() || m.LockedAmount.LessThan(val) {
		return false
	}

	if val.IsZero() {
		return true
	}

	m.FreeAmount = m.FreeAmount.Add(val)
	m.LockedAmount = m.LockedAmount.Sub(val)

	return true
}

func GetAllFreeAmountGreaterThanZero() ([]Model, *errpkg.Error) {
	var records []Model

	if err := sql.DB().Scopes(modelpkg.WithNotDeleted()).Where("`free_amount` > 0").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetUserAsset(userID int64) ([]Asset, *errpkg.Error) {
	var record []Asset
	if err := sql.DB().Table(TableName).Select("currencies_symbol, free_amount").Where("users_id = ?", userID).Find(&record).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return record, nil
}

func GetAssetForUserMapByUser(usersID int64) (map[string]AssetForUser, *errpkg.Error) {
	var records []AssetForUser

	if err := sql.DB().Table(TableName).Scopes(modelpkg.WithNotDeleted()).Where("`users_id` = ?", usersID).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make(map[string]AssetForUser, len(records))
	for _, record := range records {
		output[record.CurrenciesSymbol] = record
	}

	return output, nil
}

func GetByUserID(usersID int64, currenciesSymbol string) (*Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Scopes(modelpkg.WithNotDeleted()).Where("`users_id` = ? AND `currencies_symbol` = ?", usersID, currenciesSymbol).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: err}
		}

		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return &record, nil
}

func Create(usersID int64, currenciesSymbol string) (*Model, *errpkg.Error) {
	info, err := spottrend.GetCurrencyInfo(currenciesSymbol)
	if err != nil {
		return nil, err
	}

	record := Model{
		UsersID:          usersID,
		CurrenciesSymbol: currenciesSymbol,
		Type:             info.Type,
	}

	if err := sql.DB().Create(&record).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return &record, nil
}

func GetOrCreateWithLock(tx *gorm.DB, usersID int64, currenciesSymbol string) (Model, *errpkg.Error) {
	if tx == nil {
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("tx is nil")}
	}

	var record Model
	query := tx.Scopes(modelpkg.WithLock(true), modelpkg.WithNotDeleted()).Where("`users_id` = ? AND `currencies_symbol` = ?", usersID, currenciesSymbol).Session(&gorm.Session{})

	if err := query.Take(&record).Error; err == nil {
		return record, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	info, err := spottrend.GetCurrencyInfo(currenciesSymbol)
	if err != nil {
		return Model{}, err
	}

	record = Model{
		UsersID:          usersID,
		CurrenciesSymbol: currenciesSymbol,
		Type:             info.Type,
	}

	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&record).Error; err != nil {
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := query.Take(&record).Error; err != nil {
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func LockAmount(userID int64, currenciesSymbol string, val decimal.Decimal) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		var record Model

		if err := tx.Table(TableName).
			Where("`users_id` = ? AND `currencies_symbol` = ?", userID, currenciesSymbol).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: err}
			}

			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if !record.LockAmount(val) {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeInsufficientBalance}
		}

		if err := tx.Table(TableName).Save(&record).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func UnlockAmount(userID int64, currenciesSymbol string, val decimal.Decimal) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		var record Model

		if err := tx.Table(TableName).
			Where("`users_id` = ? AND `currencies_symbol` = ?", userID, currenciesSymbol).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: err}
			}

			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if !record.UnlockAmount(val) {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeInsufficientBalance}
		}

		if err := tx.Table(TableName).Save(&record).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}
