package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1697445575{})
}

// Migration1697445575 is a kind of IMigration, You can define schema here.
type Migration1697445575 struct {
	ID      int64  `gorm:"primaryKey"`
	UsersID int64  `gorm:"not null;uniqueIndex:idx_user_mainnet_addr_unique"`
	Mainnet string `gorm:"not null;uniqueIndex:idx_user_mainnet_addr_unique"`
	Address string `gorm:"not null;uniqueIndex:idx_user_mainnet_addr_unique"`
	Extra   string `gorm:"type:json;not null;default:'{}'"`

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1697445575) TableName() string {
	return "users_withdrawal_whitelist"
}

// Version for IMigration
func (*Migration1697445575) Version() int64 {
	return 1697445575
}

// Up for IMigration
func (m *Migration1697445575) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1697445575) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
