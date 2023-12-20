package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1691478857{})
}

// Migration1691478857 is a kind of IMigration, You can define schema here.
type Migration1691478857 struct {
	ID             int64
	UsersID        int64     `gorm:"not null;index"`
	Type           int32     `gorm:"not null"`
	TaskID         string    `gorm:"not null;default:''"`
	IDImage        []byte    `gorm:"type:mediumblob"`
	IDBackImage    []byte    `gorm:"type:mediumblob"`
	PassportImage  []byte    `gorm:"type:mediumblob"`
	FaceImage      []byte    `gorm:"type:mediumblob"`
	IDAndFaceImage []byte    `gorm:"type:mediumblob"`
	ResultImage    []byte    `gorm:"type:mediumblob"`
	State          int32     `gorm:"not null;default:0"`
	AuditStatus    int32     `gorm:"not null;default:0"`
	Detail         string    `gorm:"type:json;not null;default:'{}'"`
	CreatedAt      time.Time `gorm:"not null;default:UTC_TIMESTAMP()"`
	AuditTime      sql.NullTime

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1691478857) TableName() string {
	return "id_verifications"
}

// Version for IMigration
func (*Migration1691478857) Version() int64 {
	return 1691478857
}

func (m *Migration1691478857) GetUserMigrationWithIDVConstraint() interface{} {
	return &struct {
		Migration1689661893

		IDVerifications Migration1691478857 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	}{}
}

// Up for IMigration
func (m *Migration1691478857) Up(db *gorm.DB) error {
	if err := db.Migrator().CreateTable(m); err != nil {
		return err
	} else if err = db.Migrator().CreateConstraint(m.GetUserMigrationWithIDVConstraint(), "IDVerifications"); err != nil {
		return err
	}

	return nil
}

// Down for IMigration
func (m *Migration1691478857) Down(db *gorm.DB) error {
	if err := db.Migrator().DropConstraint(m.GetUserMigrationWithIDVConstraint(), "IDVerifications"); err != nil {
		return err
	} else if err = db.Migrator().DropTable(m); err != nil {
		return err
	}

	return nil
}
