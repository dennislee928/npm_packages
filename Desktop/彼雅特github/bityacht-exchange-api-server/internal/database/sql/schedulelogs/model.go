package schedulelogs

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
	"time"
)

// TableName of schedule_logs table
const TableName = "schedule_logs"

// Model of schedule_logs table
type Model struct {
	ID        int64
	Key       string    `gorm:"not null;uniqueIndex"`
	Tag       string    `gorm:"not null;default:''"`
	Status    Status    `gorm:"not null;default:0"`
	Result    Result    `gorm:"type:json;not null;default:'{}'"`
	CreatedAt time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func Create(key string, tag string) (Model, *errpkg.Error) {
	record := Model{
		Key:    key,
		Tag:    tag,
		Status: StatusRunning,
	}

	if err := sql.DB().Create(&record).Error; err != nil {
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func Update(record *Model) *errpkg.Error {
	updateMap := map[string]interface{}{
		"status": record.Status,
		"result": record.Result,
	}

	if err := sql.DB().Model(record).Updates(updateMap).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return nil
}
