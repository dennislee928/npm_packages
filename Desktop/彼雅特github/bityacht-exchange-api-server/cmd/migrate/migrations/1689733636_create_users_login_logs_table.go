package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689733636{})
}

// Migration1689733636 is a kind of IMigration, You can define schema here.
type Migration1689733636 struct {
	ID               int64
	UsersID          int64     `gorm:"not null;index:idx_user_time"`
	UserAgent        string    `gorm:"type:text;not null;default:''"`
	Browser          string    `gorm:"type:text;not null;default:''"`
	Device           string    `gorm:"type:text;not null;default:''"`
	Location         string    `gorm:"not null;default:''"`
	IP               string    `gorm:"not null;default:''"`
	IPRelteadHeaders string    `gorm:"type:text;not null;default:''"`
	CreatedAt        time.Time `gorm:"not null;index:idx_user_time;default:UTC_TIMESTAMP()"`

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1689733636) TableName() string {
	return "users_login_logs"
}

// Version for IMigration
func (*Migration1689733636) Version() int64 {
	return 1689733636
}

// Up for IMigration
func (m *Migration1689733636) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689733636) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
