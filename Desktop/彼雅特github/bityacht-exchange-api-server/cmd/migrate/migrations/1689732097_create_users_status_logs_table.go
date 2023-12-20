package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

// ! Deprecated, Use Migration1697099650 (users_modify_logs) instead
func init() {
	migrate.Register(&Migration1689732097{})
}

// Migration1689732097 is a kind of IMigration, You can define schema here.
type Migration1689732097 struct {
	ID         int64
	ManagersID sql.NullInt64
	UsersID    int64     `gorm:"not null;index:idx_user_time"`
	Status     int32     `gorm:"not null;default:0"`
	Comment    string    `gorm:"type:text;not null;default:''"`
	CreatedAt  time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`

	Managers Migration1689661181 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Users    Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1689732097) TableName() string {
	return "users_status_logs"
}

// Version for IMigration
func (*Migration1689732097) Version() int64 {
	return 1689732097
}

// Up for IMigration
func (m *Migration1689732097) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689732097) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
