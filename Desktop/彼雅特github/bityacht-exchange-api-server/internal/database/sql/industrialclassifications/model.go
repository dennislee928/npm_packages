package industrialclassifications

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
)

// TableName of industrial_classifications table
const TableName = "industrial_classifications"

// Model of industrial_classifications table
type Model struct {
	ID      int64
	Code    string `gorm:"not null;default:''"`
	Chinese string `gorm:"not null"`
	English string `gorm:"not null"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetICListAndMap() ([]IC, map[int64]IC, *errpkg.Error) {
	var records []IC

	if err := sql.DB().Table(TableName).Find(&records).Error; err != nil {
		return nil, nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	icMap := make(map[int64]IC, len(records))
	for _, ic := range records {
		icMap[ic.ID] = ic
	}

	return records, icMap, nil
}
