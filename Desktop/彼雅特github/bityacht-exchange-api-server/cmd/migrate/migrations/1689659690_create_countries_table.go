package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1689659690{})
}

// Migration1689659690 is a kind of IMigration, You can define schema here.
type Migration1689659690 struct {
	Code    string `gorm:"primaryKey"`
	Chinese string `gorm:"not null"`
	English string `gorm:"not null"`
	Locale  string `gorm:"not null;default:''"`
}

// TableName for gorm
func (*Migration1689659690) TableName() string {
	return "countries"
}

// Version for IMigration
func (*Migration1689659690) Version() int64 {
	return 1689659690
}

// Up for IMigration
func (m *Migration1689659690) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1689659690) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
