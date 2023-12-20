package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1697099650{})
}

// Migration1697099650 is a kind of IMigration, You can define schema here.
type Migration1697099650 struct {
	ID         int64
	ManagersID sql.NullInt64
	UsersID    int64     `gorm:"not null;index:idx_user_time"`
	Type       int32     `gorm:"not null"`
	SubType    int32     `gorm:"not null;default:0"`
	Status     int32     `gorm:"not null;default:0"`
	Comment    string    `gorm:"type:text;not null;default:''"`
	CreatedAt  time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`

	Managers Migration1689661181 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Users    Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1697099650) TableName() string {
	return "users_modify_logs"
}

// Version for IMigration
func (*Migration1697099650) Version() int64 {
	return 1697099650
}

// Up for IMigration
func (m *Migration1697099650) Up(db *gorm.DB) error {
	if err := db.Migrator().CreateTable(m); err != nil {
		return err
	} else if err = m.insertFromUsersStatusLogs(db); err != nil {
		return err
	} else if err = m.insertFromUsersReviewsLogs(db); err != nil {
		return err
	} else if err = db.Migrator().DropTable(m.getUsersStatusLogsModel()); err != nil { // users_status_logs
		return err
	} else if err = db.Migrator().DropTable(m.getUsersReviewsLogsModel()); err != nil { // users_reviews_logs
		return err
	}

	return nil
}

// Down for IMigration
func (m *Migration1697099650) Down(db *gorm.DB) error {
	if err := db.Migrator().CreateTable(m.getUsersReviewsLogsModel()); err != nil { // users_reviews_logs
		return err
	} else if err = db.Migrator().CreateTable(m.getUsersStatusLogsModel()); err != nil { // users_status_logs
		return err
	} else if err = m.insertToUsersReviewsLogs(db); err != nil {
		return err
	} else if err = m.insertToUsersStatusLogs(db); err != nil {
		return err
	} else if err = db.Migrator().DropTable(m); err != nil {
		return err
	}

	return nil
}

func (*Migration1697099650) getUsersStatusLogsModel() interface{} {
	return &Migration1689732097{}
}

func (m *Migration1697099650) insertFromUsersStatusLogs(db *gorm.DB) error {
	var oldRecords []Migration1689732097

	if err := db.Order("`id` ASC").Find(&oldRecords).Error; err != nil {
		return err
	}

	newRecords := make([]Migration1697099650, len(oldRecords))
	for i, oldRecord := range oldRecords {
		newRecords[i] = Migration1697099650{
			ManagersID: oldRecord.ManagersID,
			UsersID:    oldRecord.UsersID,
			Type:       int32(usersmodifylogs.TypeStatusLog),
			Status:     oldRecord.Status,
			Comment:    oldRecord.Comment,
			CreatedAt:  oldRecord.CreatedAt,
		}
	}

	if len(newRecords) == 0 {
		return nil
	}

	return db.Create(&newRecords).Error
}

func (m *Migration1697099650) insertToUsersStatusLogs(db *gorm.DB) error {
	var oldRecords []Migration1697099650

	if err := db.Where("`type` = ?", int32(usersmodifylogs.TypeStatusLog)).Order("`id` ASC").Find(&oldRecords).Error; err != nil {
		return err
	}

	newRecords := make([]Migration1689732097, len(oldRecords))
	for i, oldRecord := range oldRecords {
		newRecords[i] = Migration1689732097{
			ManagersID: oldRecord.ManagersID,
			UsersID:    oldRecord.UsersID,
			Status:     oldRecord.Status,
			Comment:    oldRecord.Comment,
			CreatedAt:  oldRecord.CreatedAt,
		}
	}

	if len(newRecords) == 0 {
		return nil
	}

	return db.Create(&newRecords).Error
}

func (*Migration1697099650) getUsersReviewsLogsModel() interface{} {
	return &Migration1691127122{}
}

func (m *Migration1697099650) insertFromUsersReviewsLogs(db *gorm.DB) error {
	var oldRecords []Migration1691127122

	if err := db.Order("`id` ASC").Find(&oldRecords).Error; err != nil {
		return err
	}

	newRecords := make([]Migration1697099650, len(oldRecords))
	for i, oldRecord := range oldRecords {
		newRecord := Migration1697099650{
			ManagersID: sql.NullInt64{Int64: oldRecord.ManagersID, Valid: true},
			UsersID:    oldRecord.UsersID,
			Type:       int32(usersmodifylogs.TypeReviewLog),
			SubType:    oldRecord.Type,
			Status:     oldRecord.Status,
			Comment:    oldRecord.Comment,
			CreatedAt:  oldRecord.CreatedAt,
		}

		newRecords[i] = newRecord
	}

	if len(newRecords) == 0 {
		return nil
	}

	return db.Create(&newRecords).Error
}

func (m *Migration1697099650) insertToUsersReviewsLogs(db *gorm.DB) error {
	var oldRecords []Migration1697099650

	if err := db.Where("`type` = ?", int32(usersmodifylogs.TypeReviewLog)).Order("`id` ASC").Find(&oldRecords).Error; err != nil {
		return err
	}

	newRecords := make([]Migration1691127122, len(oldRecords))
	for i, oldRecord := range oldRecords {
		newRecord := Migration1691127122{
			ManagersID: oldRecord.ManagersID.Int64,
			UsersID:    oldRecord.UsersID,
			Type:       oldRecord.SubType,
			Status:     oldRecord.Status,
			Comment:    oldRecord.Comment,
			CreatedAt:  oldRecord.CreatedAt,
		}

		if newRecord.ManagersID == 0 {
			return errors.New("bad managers_id of users_reviews_logs")
		}
		newRecords[i] = newRecord
	}

	if len(newRecords) == 0 {
		return nil
	}

	return db.Create(&newRecords).Error
}
