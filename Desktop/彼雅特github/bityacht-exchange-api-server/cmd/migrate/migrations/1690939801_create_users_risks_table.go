package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1690939801{})
}

// Migration1690939801 is a kind of IMigration, You can define schema here.
type Migration1690939801 struct {
	UsersID int64 `gorm:"primaryKey"`
	RisksID int64 `gorm:"primaryKey"`

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	Risks Migration1690884579 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1690939801) TableName() string {
	return "users_risks"
}

// Version for IMigration
func (*Migration1690939801) Version() int64 {
	return 1690939801
}

// Up for IMigration
func (m *Migration1690939801) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1690939801) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
