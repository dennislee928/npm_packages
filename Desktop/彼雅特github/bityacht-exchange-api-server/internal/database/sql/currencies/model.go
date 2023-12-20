package currencies

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// TableName of currencies table
const TableName = "currencies"

// Model of currencies table
type Model struct {
	Symbol           string    `gorm:"primaryKey"`
	Name             string    `gorm:"not null;default:''"`
	Type             Type      `gorm:"not null;default:0"`
	DecimalPrecision int32     `gorm:"not null;default:0"` // Max: 9, if over max -> Need update all precision of decimal value in db.
	CreatedAt        time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	DeletedAt        *time.Time
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetBySymbol(symbol string) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`symbol` = ? AND `deleted_at` IS NULL", symbol).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("symbol")}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func GetAll() ([]Model, *errpkg.Error) {
	var records []Model

	if err := sql.DB().Scopes(modelpkg.WithNotDeleted()).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql}
	}

	return records, nil
}

func GetAllByType(currencyType Type) ([]Model, *errpkg.Error) {
	var records []Model

	if err := sql.DB().Scopes(modelpkg.WithNotDeleted()).Where("`type` = ?", currencyType).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql}
	}

	return records, nil
}
