package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1692933978{})
}

// Migration1692933978 is a kind of IMigration, You can define schema here.
type Migration1692933978 struct {
	ID      int64
	UsersID int64          `gorm:"not null;uniqueIndex:idx_user_mainnet"`
	Mainnet string         `gorm:"not null;uniqueIndex:idx_user_mainnet;uniqueIndex:idx_mainnet_address"`
	Address sql.NullString `gorm:"uniqueIndex:idx_mainnet_address"`
	TxID    string         `gorm:"not null;default:''"`

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1692933978) TableName() string {
	return "wallets_address"
}

// Version for IMigration
func (*Migration1692933978) Version() int64 {
	return 1692933978
}

// Up for IMigration
func (m *Migration1692933978) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1692933978) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
