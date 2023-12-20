package banners

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"errors"
	"net/http"

	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of banners table
const TableName = "banners"

// Model of banners table
type Model struct {
	ID        int64
	WebImage  string    `gorm:"not null;default:''"`
	AppImage  string    `gorm:"not null;default:''"`
	Priority  int64     `gorm:"not null;default:0"`
	Title     string    `gorm:"not null;default:''"`
	SubTitle  string    `gorm:"not null;default:''"`
	Button    string    `gorm:"not null;default:''"`
	ButtonUrl string    `gorm:"not null;default:''"`
	Status    Status    `gorm:"not null;default:0"`
	StartAt   time.Time `gorm:"not null;default:'0001-01-01 00:00:00'"`
	EndAt     time.Time `gorm:"not null;default:'0001-01-01 00:00:00'"`
	CreatedAt time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func Create(m *Model) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(m).Error; err != nil {
			return errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		} else if err = tx.Model(m).Update("priority", m.ID).Error; err != nil {
			return errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func Get(id int64) (Banner, *errpkg.Error) {
	var record Banner
	if err := sql.DB().Table(TableName).Where("id = ?", id).Take(&record).Error; err != nil {
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return record, nil
}

func GetEnableList() ([]Model, *errpkg.Error) {
	var records []Model

	if err := sql.DB().Where("`status` = ?", StatusEnable).Order("`priority` ASC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	return records, nil
}

func GetList(paginator *modelpkg.Paginator, status int32) ([]Banner, *errpkg.Error) {
	records := make([]Banner, 0)
	query := sql.DB().Table(TableName)

	if status != -1 {
		query = query.Where("`status` = ?", status)
	}
	query = query.Session(&gorm.Session{})
	if paginator == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("paginator is nil")}
	}
	if err := query.Limit(paginator.PageSize).Offset(paginator.Offset()).Order("`priority` ASC").Find(&records).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	if err := query.Count(&paginator.TotalRecord).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil
}

func Update(req UpdateRequest, webImage, appImage string) *errpkg.Error {
	var record Model
	if err := sql.DB().Transaction(
		func(tx *gorm.DB) error {
			updateMap := make(map[string]interface{})
			query := tx.Table(TableName).Where("`id` = ?", req.ID).Session(&gorm.Session{}) // Session for Reuse query

			if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&record).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
				}
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
			if webImage != "" {
				updateMap["web_image"] = webImage
			}
			if appImage != "" {
				updateMap["app_image"] = appImage
			}
			if req.Title != "" && req.Title != record.Title {
				updateMap["title"] = req.Title
			}
			if req.SubTitle != "" && req.SubTitle != record.SubTitle {
				updateMap["sub_title"] = req.SubTitle
			}
			if req.Button != "" && req.Button != record.Button {
				updateMap["button"] = req.Button
			}
			if req.ButtonUrl != "" && req.ButtonUrl != record.ButtonUrl {
				updateMap["button_url"] = req.ButtonUrl
			}
			if req.Status != record.Status {
				updateMap["status"] = req.Status
			}

			if !req.StartAt.Time.Equal(record.StartAt) {
				updateMap["start_at"] = req.StartAt
			}
			if !req.EndAt.Time.Equal(record.EndAt) {
				updateMap["end_at"] = req.EndAt
			}

			if len(updateMap) == 0 {
				return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeRecordNoChange}
			}
			if err := query.Updates(updateMap).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
			if err := query.Take(&record).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
			return nil
		}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

func PriorityUpdate(req PriorityUpdateRequset) *errpkg.Error {
	if err := sql.DB().Transaction(
		func(tx *gorm.DB) error {
			for _, v := range req.Rows {
				var record Model
				if err := tx.Table(TableName).Where("`id` = ?", v.ID).Clauses(clause.Locking{Strength: "UPDATE"}).Take(&record).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
					}
					return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
				}
				if err := tx.Model(&record).Update("priority", v.Priority).Error; err != nil {
					return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
				}
			}
			return nil
		}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

func Delete(id int64) (*Model, *errpkg.Error) {
	var old *Model
	if err := sql.DB().Table(TableName).Where("id = ?", id).Take(&old).Error; err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	if err := sql.DB().Delete(Model{}, id).Error; err != nil {
		return old, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return old, nil
}
