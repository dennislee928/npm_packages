package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

// ! Deprecated, Use Migration1697099650 (users_modify_logs) instead
func init() {
	migrate.Register(&Migration1691127122{})
}

// Migration1691127122 is a kind of IMigration, You can define schema here.
type Migration1691127122 struct {
	ID         int64
	ManagersID int64     `gorm:"not null"`
	UsersID    int64     `gorm:"not null;index:idx_user_time"`
	Type       int32     `gorm:"not null;default:0"`
	Status     int32     `gorm:"not null;default:0"`
	Comment    string    `gorm:"type:text;not null;default:''"`
	CreatedAt  time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`

	Managers Migration1689661181 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Users    Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1691127122) TableName() string {
	return "users_reviews_logs"
}

// Version for IMigration
func (*Migration1691127122) Version() int64 {
	return 1691127122
}

// Up for IMigration
func (m *Migration1691127122) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1691127122) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
