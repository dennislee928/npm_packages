package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1690884579{})
}

// Migration1690884579 is a kind of IMigration, You can define schema here.
type Migration1690884579 struct {
	ID        int64
	Factor    string    `gorm:"not null;default:''"`
	SubFactor string    `gorm:"not null;default:''"`
	Detail    string    `gorm:"not null;default:''"`
	Score     int64     `gorm:"not null;default:0"`
	CreatedAt time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
}

// TableName for gorm
func (*Migration1690884579) TableName() string {
	return "risks"
}

// Version for IMigration
func (*Migration1690884579) Version() int64 {
	return 1690884579
}

// Up for IMigration
func (m *Migration1690884579) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1690884579) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
