package mainnets

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"errors"
	"net/http"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// TableName of mainnets table
const TableName = "mainnets"

// Model of mainnets table
type Model struct {
	CurrenciesSymbol         string          `gorm:"primaryKey"` // binance: coin
	Mainnet                  string          `gorm:"primaryKey"` // binance: network
	Name                     string          `gorm:"not null;default:''"`
	AddressRegex             string          `gorm:"not null;default:''"`
	WithdrawDecimalPrecision int32           `gorm:"not null;default:0"` // parse from binance: withdrawIntegerMultiple
	WithdrawFee              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	WithdrawMin              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	WithdrawMax              decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	Status                   Status          `gorm:"not null;default:0"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetMap() (map[string]map[string]Model, *errpkg.Error) {
	var records []Model

	if err := sql.DB().Where("`status` = ?", StatusEnable).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make(map[string]map[string]Model)

	for _, record := range records {
		if _, ok := output[record.CurrenciesSymbol]; !ok {
			output[record.CurrenciesSymbol] = make(map[string]Model)
		}

		output[record.CurrenciesSymbol][record.Mainnet] = record
	}

	return output, nil
}

func GetMainnetList(paginator *modelpkg.Paginator) ([]Mainnet, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	var (
		records []Mainnet
		query   = sql.DB().Table(TableName).Session(&gorm.Session{})
	)

	if err := query.Scopes(modelpkg.WithPaginator(paginator)).Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func UpdateMainnet(req UpdateRequest) *errpkg.Error {
	if req.WithdrawFee.LessThan(decimal.Zero) || req.WithdrawMin.LessThan(decimal.Zero) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody, Err: errors.New("withdraw fee & withdraw min cannot less than zero")}
	}

	var record Model
	if err := sql.DB().Table(TableName).Where("`currencies_symbol` = ? AND `mainnet` = ?", req.Currency, req.Mainnet).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
		}

		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	updateMap := map[string]any{
		"withdraw_fee": req.WithdrawFee,
		"withdraw_min": req.WithdrawMin,
	}

	if err := sql.DB().Model(&record).Updates(updateMap).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}
