package risks

import (
	"errors"
	"net/http"
	"time"

	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersrisks"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TableName of managers table
const TableName = "risks"

// Model of managers table
type Model struct {
	ID        int64
	Factor    string    `gorm:"not null;default:''"`
	SubFactor string    `gorm:"not null;default:''"`
	Detail    string    `gorm:"not null;default:''"`
	Score     int64     `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName of Model
func (Model) TableName() string {
	return TableName
}

func Get() ([]Risk, *errpkg.Error) {
	var records []Risk
	if err := sql.DB().Find(&records).Error; err != nil {
		return records, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return records, nil
}

func Create(req CreateRequest) *errpkg.Error {
	record := Model{
		Factor:    req.Factor,
		SubFactor: req.SubFactor,
		Detail:    req.Detail,
		Score:     req.Score,
	}
	if err := sql.DB().Create(&record).Error; err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
	}
	return nil
}

func Update(req UpdateRequest) *errpkg.Error {
	updateMap := make(map[string]interface{})
	if req.Factor != "" {
		updateMap["factor"] = req.Factor
	}
	if req.SubFactor != "" {
		updateMap["sub_factor"] = req.SubFactor
	}
	if req.Detail != "" {
		updateMap["detail"] = req.Detail
	}
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		// 取得本來的分數
		var record Model
		if err := tx.Table(TableName).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.ID).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		if record.Score != req.Score {
			// 取得所有有此風險的使用者
			var usersRecords []users.Model
			if err := tx.Table(users.TableName+" as u").Select("u.*").
				Joins("inner join "+usersrisks.TableName+" as ur on u.id = ur.users_id").
				Clauses(clause.Locking{Strength: "UPDATE"}).Where("ur.risks_id = ?", req.ID).
				Scan(&usersRecords).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
				}
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			if len(usersRecords) > 0 {
				// 更新所有有此風險的使用者分數
				for _, userRecord := range usersRecords {
					userRecord.InternalRisksTotal += (req.Score - record.Score)
				}
				if err := tx.Save(&usersRecords).Error; err != nil {
					return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
				}
			}
			updateMap["score"] = req.Score
		}
		if len(updateMap) == 0 {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeRecordNoChange}
		}

		if err := tx.Table(TableName).Where("id = ?", req.ID).Updates(updateMap).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

func Delete(id int64) *errpkg.Error {
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		// 取得本來的分數
		var record Model
		if err := tx.Table(TableName).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Take(&record).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		// 取得所有有此風險的使用者
		var usersRecords []users.Model
		if err := tx.Table(users.TableName+" as u").Select("u.*").
			Joins("inner join "+usersrisks.TableName+" as ur on u.id = ur.users_id").
			Clauses(clause.Locking{Strength: "UPDATE"}).Where("ur.risks_id = ?", id).
			Scan(&usersRecords).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}

		if len(usersRecords) > 0 {
			// 更新所有有此風險的使用者分數
			for _, userRecord := range usersRecords {
				userRecord.InternalRisksTotal -= record.Score
			}
			if err := tx.Save(&usersRecords).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

		}

		if tx.Delete(&usersrisks.Model{}, "risks_id = ?", id).Error != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql}
		}
		if tx.Delete(&Model{}, "id = ?", id).Error != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql}
		}
		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}

// UpdateRisks 更新某使用者風險
func UpdateRisks(req UpdateRisksRequest) *errpkg.Error {
	create := make([]*usersrisks.Model, 0, len(req.RisksIDs))
	for _, riskID := range req.RisksIDs {
		create = append(create, &usersrisks.Model{UsersID: req.ID, RisksID: riskID})
	}
	if err := sql.DB().Transaction(func(tx *gorm.DB) error {
		// 鎖定使用者
		var userRecord users.Model
		if err := tx.Table(users.TableName).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", req.ID).Take(&userRecord).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &errpkg.Error{HttpStatus: http.StatusNotFound, Code: errpkg.CodeRecordNotFound}
			}
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		// 清除所有該使用者的風險
		if err := tx.Table(usersrisks.TableName).Delete(&usersrisks.Model{}, "users_id = ?", req.ID).Error; err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
		}
		if len(create) > 0 {
			// 建立使用者風險
			if err := tx.Table(usersrisks.TableName).Create(&create).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}

			// 更新該使用者的風險分數
			var total int64
			if err := tx.Select("SUM(score)").Table(TableName).Where("id IN ?", req.RisksIDs).Scan(&total).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
			if err := tx.Table(users.TableName).Where("id = ?", req.ID).Update("internal_risks_total", total).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
		} else {
			// 更新該使用者的風險分數為0
			if err := tx.Table(users.TableName).Where("id = ?", req.ID).Update("internal_risks_total", 0).Error; err != nil {
				return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSql, Err: err}
			}
		}

		return nil
	}); err != nil {
		return err.(*errpkg.Error)
	}
	return nil
}
