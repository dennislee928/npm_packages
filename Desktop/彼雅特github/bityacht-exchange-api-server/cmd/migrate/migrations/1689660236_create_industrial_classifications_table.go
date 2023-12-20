package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689660236{})
}

// Migration1689660236 is a kind of IMigration, You can define schema here.
type Migration1689660236 struct {
	ID      int64
	Code    string `gorm:"not null;default:''"`
	Chinese string `gorm:"not null"`
	English string `gorm:"not null"`
}

// TableName for gorm
func (*Migration1689660236) TableName() string {
	return "industrial_classifications"
}

// Version for IMigration
func (*Migration1689660236) Version() int64 {
	return 1689660236
}

// Up for IMigration
func (m *Migration1689660236) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689660236) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
