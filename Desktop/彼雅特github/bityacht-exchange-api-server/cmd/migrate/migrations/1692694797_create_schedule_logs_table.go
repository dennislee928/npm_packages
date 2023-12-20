package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1692694797{})
}

// Migration1692694797 is a kind of IMigration, You can define schema here.
type Migration1692694797 struct {
	ID        int64
	Key       string    `gorm:"not null;uniqueIndex"`
	Tag       string    `gorm:"not null;default:''"`
	Status    int32     `gorm:"not null;default:0"`
	Result    string    `gorm:"type:json;not null;default:'{}'"`
	CreatedAt time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName for gorm
func (*Migration1692694797) TableName() string {
	return "schedule_logs"
}

// Version for IMigration
func (*Migration1692694797) Version() int64 {
	return 1692694797
}

// Up for IMigration
func (m *Migration1692694797) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1692694797) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
