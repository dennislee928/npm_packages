package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689663148{})
}

// Migration1689663148 is a kind of IMigration, You can define schema here.
type Migration1689663148 struct {
	Symbol           string    `gorm:"primaryKey"`
	Name             string    `gorm:"not null;default:''"`
	Type             int32     `gorm:"not null;default:0"`
	DecimalPrecision int32     `gorm:"not null;default:0"` // Max: 9, if over max -> Need update all precision of decimal value in db.
	CreatedAt        time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	DeletedAt        *time.Time
}

// TableName for gorm
func (*Migration1689663148) TableName() string {
	return "currencies"
}

// Version for IMigration
func (*Migration1689663148) Version() int64 {
	return 1689663148
}

// Up for IMigration
func (m *Migration1689663148) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689663148) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
