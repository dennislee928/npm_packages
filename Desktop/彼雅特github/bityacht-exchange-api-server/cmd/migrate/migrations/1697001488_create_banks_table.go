package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1697001488{})
}

// Migration1697001488 is a kind of IMigration, You can define schema here.
type Migration1697001488 struct {
	Code    string `gorm:"primaryKey"`
	Chinese string `gorm:"not null"`
	English string `gorm:"not null;default:''"`
}

// TableName for gorm
func (*Migration1697001488) TableName() string {
	return "banks"
}

// Version for IMigration
func (*Migration1697001488) Version() int64 {
	return 1697001488
}

// Up for IMigration
func (m *Migration1697001488) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1697001488) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
