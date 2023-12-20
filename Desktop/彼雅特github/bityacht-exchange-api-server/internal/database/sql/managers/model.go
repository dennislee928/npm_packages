package managers

import (
	"bityacht-exchange-api-server/internal/database/sql"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	passwordpkg "bityacht-exchange-api-server/internal/pkg/password"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of managers table
const TableName = "managers"

// ! If modify adminAccount, make sure UpdateRequest.Account's binding Tag is modified too
const adminAccount = "admin"

// Model of managers table
type Model struct {
	ID              int64
	Account         string    `gorm:"not null;uniqueIndex"`
	ManagersRolesID int64     `gorm:"not null"`
	Password        string    `gorm:"not null"`
	Name            string    `gorm:"not null;default:''"`
	Extra           Extra     `gorm:"type:json;not null;default:'{}'"`
	Status          Status    `gorm:"not null;default:0"`
	CreatedAt       time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	DeletedAt       *time.Time
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

// GetAvailableIDs is only for buffer now, so it use the native error.
func GetAvailableIDs() ([]int64, error) {
	var ids []int64

	if err := sql.DB().Table(TableName).Select("`id`").Where("`account` != ? AND `deleted_at` IS NULL", adminAccount).Pluck("id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func Login(account string, password string) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Where("`account` = ? AND `deleted_at` IS NULL AND `status` = ? ", account, StatusEnable).Take(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return record, &errpkg.Error{HttpStatus: http.StatusUnauthorized, Code: errpkg.CodeUnauthorized}
		}
		return record, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	} else if err := passwordpkg.Validate(record.Password, password); err != nil {
		return record, err
	}

	return record, nil
}

func Create(ctx *gin.Context, account string, name string, managersRolesID int64) (Model, *errpkg.Error) {
	var err *errpkg.Error
	rawPassword := rand.LetterAndNumberString(8)
	record := Model{
		Account:         account,
		Name:            name,
		ManagersRolesID: managersRolesID,
		Status:          1,
		Extra: Extra{
			NeedChangePassword: true,
		},
	}

	if record.Password, err = passwordpkg.Encrypt(rawPassword); err != nil {
		return Model{}, err
	} else if err := sql.DB().Create(&record).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return Model{}, &errpkg.Error{HttpStatus: http.StatusConflict, Code: errpkg.CodeAccountDuplicated}
		}

		return Model{}, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}

	record.Password = rawPassword
	return record, nil
}

// GetManagerList will return all not deleted managers in db, if id > 0 will retrieve specific record only
func GetManagerList(paginator *modelpkg.Paginator) ([]Manager, *errpkg.Error) {
	records := make([]Manager, 0)
	query := sql.DB().Table(TableName).Where("`deleted_at` IS NULL").Order("`created_at` DESC").Session(&gorm.Session{}) // Session for Reuse query
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

func Update(ctx *gin.Context, req UpdateRequest) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Transaction(

		func(tx *gorm.DB) error {
			updateMap := make(map[string]interface{})
			query := tx.Table(TableName).Where("`id` = ?", req.ID).Session(&gorm.Session{}) // Session for Reuse query

			if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Where("`deleted_at` IS NULL").Take(&record).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
				}
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			if req.Account != "" && req.Account != record.Account {
				updateMap["account"] = req.Account
			}
			if req.Name != "" && req.Name != record.Name {
				updateMap["name"] = req.Name
			}
			if req.Status != record.Status {
				updateMap["status"] = req.Status
			}
			if req.ManagersRolesID != 0 && req.ManagersRolesID != record.ManagersRolesID {
				updateMap["managers_roles_id"] = req.ManagersRolesID
			}

			if req.Password != "" {
				var err *errpkg.Error
				if err = passwordpkg.StrengthValidate(req.Password); err != nil {
					return err
				} else if updateMap["password"], err = passwordpkg.Encrypt(req.Password); err != nil {
					return err
				}
			}
			// 如果欄位沒有異動的話，就不更新
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
		return Model{}, err.(*errpkg.Error)
	}

	if req.Password != "" {
		record.Password = req.Password
	} else {
		record.Password = ""
	}

	return record, nil
}

func Delete(id int64) *errpkg.Error {
	if err := sql.DB().Transaction(
		func(tx *gorm.DB) error {
			var record Model
			query := tx.Table(TableName).Where("`id` = ?", id).Session(&gorm.Session{}) // Session for Reuse query

			if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Where("`deleted_at` IS NULL").Take(&record).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
				}
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			timeStr := strconv.FormatInt(time.Now().Unix(), 36)
			updateMap := map[string]interface{}{
				"deleted_at": tx.Raw("UTC_TIMESTAMP()"),
			}

			if len(timeStr)+len(record.Account) <= 254 { // len( timeStr[record.Account] ) Sholud less than or equal to 256
				updateMap["account"] = fmt.Sprintf("%s[%s]", timeStr, record.Account)
			} else {
				updateMap["account"] = fmt.Sprintf("%s[%s...]", timeStr, record.Account[:251-len(timeStr)])
			}

			if err := query.Updates(updateMap).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			return nil
		}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}

func ResetPasswordByEamil(account string) (Model, *errpkg.Error) {
	var record Model

	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		var err *errpkg.Error

		query := tx.Table(TableName).Session(&gorm.Session{}) // Session for Reuse query
		if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Where("`account` = ? AND `deleted_at` IS NULL", account).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		record.Extra.NeedChangePassword = true
		record.Password = rand.LetterAndNumberString(8)

		updateMap := make(map[string]interface{}, 2)
		updateMap["extra"] = record.Extra

		if updateMap["password"], err = passwordpkg.Encrypt(record.Password); err != nil {
			return err
		}
		if err := query.Where("id = ?", record.ID).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		return nil
	}); err != nil {
		return record, err.(*errpkg.Error)
	}

	return record, nil
}

func UpdatePassword(id int64, password string) *errpkg.Error {
	password, err := passwordpkg.Encrypt(password)
	if err != nil {
		return err
	}
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		var record Model
		query := tx.Table(TableName).Session(&gorm.Session{}) // Session for Reuse query
		if err := query.Clauses(clause.Locking{Strength: "UPDATE"}).Where("`id` = ? AND `deleted_at` IS NULL", id).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		record.Password = password
		if err := query.Save(record).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}

	return nil
}
