package migrations

import (
	"bityacht-exchange-api-server/cmd/migrate"
	"database/sql"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

func init() {
	migrate.Register(&Migration1691389573{})
}

// Migration1691389573 is a kind of IMigration, You can define schema here.
type Migration1691389573 struct {
	ID              int64
	UsersID         int64           `gorm:"not null;default:0"`
	Type            int32           `gorm:"not null"`
	TaskID          string          `gorm:"not null;default:''"`
	SanctionMatched int32           `gorm:"not null;default:0"`
	PotentialRisk   int64           `gorm:"not null;default:0"`
	AuditAccepted   int32           `gorm:"not null;default:0"`
	Comment         string          `gorm:"not null;default:''"`
	Detail          json.RawMessage `gorm:"type:json;not null;default:'{}'"` // Save the original response, don't Use Detail as Type
	CreatedAt       time.Time       `gorm:"not null;default:UTC_TIMESTAMP()"`
	AuditTime       sql.NullTime

	Users Migration1689661893 `gorm:"constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}

// TableName for gorm
func (*Migration1691389573) TableName() string {
	return "due_diligences"
}

// Version for IMigration
func (*Migration1691389573) Version() int64 {
	return 1691389573
}

func getMigrationUserWithDDConstraint() interface{} {
	return &struct {
		Migration1689661893

		DueDiligences Migration1691389573 `gorm:"foreignKey:DueDiligencesID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	}{}
}

// Up for IMigration
func (m *Migration1691389573) Up(db *gorm.DB) error {
	if err := db.Migrator().CreateTable(m); err != nil {
		return err
	} else if err = db.Migrator().CreateConstraint(getMigrationUserWithDDConstraint(), "DueDiligences"); err != nil {
		return err
	}
	return nil
}

// Down for IMigration
func (m *Migration1691389573) Down(db *gorm.DB) error {
	if err := db.Migrator().DropConstraint(getMigrationUserWithDDConstraint(), "DueDiligences"); err != nil {
		return err
	} else if err = db.Migrator().DropTable(m); err != nil {
		return err
	}
	return nil
}
