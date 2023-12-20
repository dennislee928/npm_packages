package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1697003675{})
}

// Migration1697003675 is a kind of IMigration, You can define schema here.
type Migration1697003675 struct {
	BanksCode string `gorm:"primaryKey"`
	Code      string `gorm:"primaryKey"`
	Chinese   string `gorm:"not null"`
	English   string `gorm:"not null;default:''"`

	Banks Migration1697001488 `gorm:"foreignKey:BanksCode;references:Code;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1697003675) TableName() string {
	return "bank_branchs"
}

// Version for IMigration
func (*Migration1697003675) Version() int64 {
	return 1697003675
}

// Up for IMigration
func (m *Migration1697003675) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1697003675) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
