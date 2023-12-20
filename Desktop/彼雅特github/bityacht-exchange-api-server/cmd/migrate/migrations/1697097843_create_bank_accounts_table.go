package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1697097843{})
}

// Migration1697097843 is a kind of IMigration, You can define schema here.
type Migration1697097843 struct {
	ID          int64
	UsersID     int64     `gorm:"not null"`
	BanksCode   string    `gorm:"not null"`
	BranchsCode string    `gorm:"not null"`
	Name        string    `gorm:"not null"`
	Account     string    `gorm:"not null"`
	CoverImage  []byte    `gorm:"type:mediumblob"`
	Status      int32     `gorm:"not null;default:1"`
	CreatedAt   time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	AuditTime   sql.NullTime

	Users       Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	BankBranchs Migration1697003675 `gorm:"foreignKey:BanksCode,BranchsCode;references:BanksCode,Code;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1697097843) TableName() string {
	return "bank_accounts"
}

// Version for IMigration
func (*Migration1697097843) Version() int64 {
	return 1697097843
}

func (m *Migration1697097843) getUserMigrationWithBankAccountsConstraint() interface{} {
	return &struct {
		Migration1689661893

		BankAccounts Migration1697097843 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	}{}
}

// Up for IMigration
func (m *Migration1697097843) Up(db *gorm.DB) error {
	if err := db.Migrator().CreateTable(m); err != nil {
		return err
	} else if err = db.Migrator().CreateConstraint(m.getUserMigrationWithBankAccountsConstraint(), "BankAccounts"); err != nil {
		return err
	}

	return nil
}

// Down for IMigration
func (m *Migration1697097843) Down(db *gorm.DB) error {
	if err := db.Migrator().DropConstraint(m.getUserMigrationWithBankAccountsConstraint(), "BankAccounts"); err != nil {
		return err
	} else if err := db.Migrator().DropTable(m); err != nil {
		return err
	}

	return nil
}
