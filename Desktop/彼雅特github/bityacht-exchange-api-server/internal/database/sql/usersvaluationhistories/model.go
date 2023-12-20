package usersvaluationhistories

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// TableName of users_valuation_histories table
const TableName = "users_valuation_histories"

// Model of users_valuation_histories table
type Model struct {
	ID        int64
	UsersID   int64           `gorm:"not null;uniqueIndex:idx_user_date"`
	Date      time.Time       `gorm:"type:date;not null;uniqueIndex:idx_user_date"`
	Valuation decimal.Decimal `gorm:"type:decimal(27,9);not null;default:0"`
	CreatedAt time.Time       `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func GetHistoriesByUser(usersID int64, startAt time.Time, limit int) ([]History, *errpkg.Error) {
	var records []History

	if err := sql.DB().Table(TableName).Where("`users_id` = ? AND `created_at` > ?", usersID, startAt).Limit(limit).Order("`date` DESC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}
