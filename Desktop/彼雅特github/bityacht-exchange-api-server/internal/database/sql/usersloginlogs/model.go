package usersloginlogs

import (
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/mmdb"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/useragent"
	"errors"
	"net"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// TableName of users_login_logs table
const TableName = "users_login_logs"

// Model of users_login_logs table
type Model struct {
	ID               int64
	UsersID          int64     `gorm:"not null;index:idx_user_time"`
	UserAgent        string    `gorm:"type:text;not null;default:''"`
	Browser          string    `gorm:"type:text;not null;default:''"`
	Device           string    `gorm:"type:text;not null;default:''"`
	Location         string    `gorm:"not null;default:''"`
	IP               string    `gorm:"not null;default:''"`
	IPRelteadHeaders Headers   `gorm:"type:text;not null;default:''"`
	CreatedAt        time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func Create(usersID int64, userAgent string, ip string, headers http.Header) (Model, *errpkg.Error) {
	ua := useragent.Parse(userAgent)
	record := Model{
		UsersID:          usersID,
		UserAgent:        userAgent,
		Browser:          ua.GetBrowserString(),
		Device:           ua.GetDeviceString(),
		IP:               ip,
		IPRelteadHeaders: make(Headers, len(headers)),
	}

	if cityResult, err := mmdb.LookupCity(ip); err != nil {
		return Model{}, err
	} else {
		record.Location = cityResult.String()
	}

	for headerKey, headerVals := range headers {
		for _, val := range headerVals {
			if net.ParseIP(val) != nil {
				record.IPRelteadHeaders[headerKey] = headerVals
				break
			}
		}
	}

	if err := sql.DB().Create(&record).Error; err != nil {
		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return record, nil
}

func GetLogsByUser(id int64, paginator *modelpkg.Paginator) ([]Log, *errpkg.Error) {
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}

	query := sql.DB().Table(TableName).Where("`users_id` = ?", id).Session(&gorm.Session{}) // Session
	records := make([]Log, 0)

	if err := query.Limit(paginator.PageSize).Offset(paginator.Offset()).Order("`created_at` DESC").Find(&records).Error; err != nil {
		return records, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err = query.Count(&paginator.TotalRecord).Error; err != nil {
		return records, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetExport(id int64, startAt, endAt modelpkg.Date) ([]csv.Record, *errpkg.Error) {
	var records []Log
	query := sql.DB().Table(TableName).Where("users_id = ?", id).Scopes(modelpkg.WithStartDateAndEnd("created_at", true, startAt, endAt))
	if err := query.Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	output := make([]csv.Record, len(records))
	for k, v := range records {
		output[k] = v
	}
	return output, nil
}
