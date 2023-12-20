package transactionpairs

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of transaction_pairs table
const TableName = "transaction_pairs"

// Model of transaction_pairs table
type Model struct {
	BaseCurrenciesSymbol            string          `gorm:"primaryKey"`
	QuoteCurrenciesSymbol           string          `gorm:"primaryKey"`
	SpreadsOfBuy                    decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	SpreadsOfSell                   decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	HandlingChargeRate              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	BaseCurrenciesDecimalPrecision  int32           `gorm:"not null;default:0"`
	QuoteCurrenciesDecimalPrecision int32           `gorm:"not null;default:0"`
	Status                          Status          `gorm:"not null;default:0"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetBySymbol(base string, quote string) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`base_currencies_symbol` = ? AND `quote_currencies_symbol` = ?", base, quote).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound, Err: errors.New("symbol")}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func GetList(paginator *modelpkg.Paginator) ([]TransactionPair, *errpkg.Error) {
	records := make([]TransactionPair, 0)
	query := sql.DB().Table(TableName).Session(&gorm.Session{}) // Session for Reuse query
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}
	if err := query.Limit(paginator.PageSize).Offset(paginator.Offset()).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil

}

func Update(ctx *gin.Context, req UpdateRequest) *errpkg.Error {

	var record Model

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		query := tx.Table(TableName).Where("`base_currencies_symbol` = ? AND `quote_currencies_symbol` = ?", req.BaseCurrenciesSymbol, req.QuoteCurrenciesSymbol).Session(&gorm.Session{}) // Session for Reuse query
		if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		updateMap := make(map[string]interface{})

		if !req.SpreadsOfBuy.Equal(decimal.Zero) {
			updateMap["spreads_of_buy"] = req.SpreadsOfBuy
		}
		if !req.SpreadsOfSell.Equal(decimal.Zero) {
			updateMap["spreads_of_sell"] = req.SpreadsOfSell
		}
		if !req.HandlingChargeRate.Equal(decimal.Zero) {
			updateMap["handling_charge_rate"] = req.HandlingChargeRate
		}
		updateMap["status"] = req.Status

		if err := query.Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func GetAll() ([]Model, *errpkg.Error) {
	var records []Model

	if err := sql.DB().Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}
