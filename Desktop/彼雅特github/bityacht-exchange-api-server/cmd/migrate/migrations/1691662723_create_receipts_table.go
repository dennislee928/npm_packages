package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1691662723{})
}

// Migration1691662723 is a kind of IMigration, You can define schema here.
type Migration1691662723 struct {
	ID              string
	Status          int32     `gorm:"not null;default:1"`
	UserID          int64     `gorm:"not null;index"`
	InvoiceAmount   int64     `gorm:"not null;default:0"`
	InvoiceID       string    `gorm:"not null;index;default:''"`
	SalesAmount     int64     `gorm:"not null;default:0"`
	Tax             int64     `gorm:"not null;default:0"`
	InvoiceIssuedAt time.Time `gorm:"not null;default:'0001-01-01 00:00:00'"`
	Barcode         string    `gorm:"not null;default:''"`
	CreatedAt       time.Time `gorm:"not null;index;default:UTC_TIMESTAMP()"`
	DeletedAt       *time.Time

	User Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1691662723) TableName() string {
	return "receipts"
}

// Version for IMigration
func (*Migration1691662723) Version() int64 {
	return 1691662723
}

// Up for IMigration
func (m *Migration1691662723) Up(db *gorm.DB) error {
	return db.Migrator().CreateTable(m)
}

// Down for IMigration
func (m *Migration1691662723) Down(db *gorm.DB) error {
	return db.Migrator().DropTable(m)
}
