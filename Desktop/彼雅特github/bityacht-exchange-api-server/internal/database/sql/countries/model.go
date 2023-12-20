package countries

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
)

// TableName of countries table
const TableName = "countries"

// Model of countries table
type Model struct {
	Code    string `gorm:"primaryKey"`
	Chinese string `gorm:"not null"`
	English string `gorm:"not null"`
	Locale  string `gorm:"not null;default:''"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetCountryListAndMap() ([]Country, map[string]Country, *errpkg.Error) {
	var records []Country

	if err := sql.DB().Table(TableName).Find(&records).Error; err != nil {
		return nil, nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	countryMap := make(map[string]Country, len(records))
	for _, country := range records {
		countryMap[country.Code] = country
	}

	return records, countryMap, nil
}
