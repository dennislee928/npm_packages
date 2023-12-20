package migrate

import (
	"context"
	"errors"
	"fmt"
	"time"

	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/pkg/logger"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// SchemaMigration is the Model for schema_migrations table
type SchemaMigration struct {
	ID        uint64
	Command   string    `gorm:"not null"`
	Version   int64     `gorm:"not null"`
	Dirty     bool      `gorm:"not null"`
	Error     string    `gorm:"default:NULL"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP()"`
}

// TableName for gorm
func (SchemaMigration) TableName() string { return "schema_migrations" }

// IsClean will check migration is clean or dirty
func (sm *SchemaMigration) IsClean() error {
	if sm.Dirty {
		return errors.New("last migration is dirty")
	}

	return nil
}

func connectDB() (*gorm.DB, error) {
	conn, err := sql.Connect(configs.Config.Database.SQL.Name)
	var sqlErr *mysql.MySQLError
	if err != nil && errors.As(err, &sqlErr) && sqlErr.Number == 1049 {
		if !createDB {
			logger.Logger.Error().Err(err).Msg("if you want to create database, you can add --create to command and try again")
			return nil, err
		}

		if conn, err = sql.Connect(""); err == nil {
			if err = conn.Exec("CREATE DATABASE IF NOT EXISTS " + configs.Config.Database.SQL.Name + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;").Error; err == nil {
				logger.Logger.Info().Msg("Auto Create Database " + configs.Config.Database.SQL.Name + " successful!")
				conn, err = sql.Connect(configs.Config.Database.SQL.Name)
			}
		}
	}
	if err != nil {
		logger.Logger.Error().Err(err).Any("Database Config", configs.Config.Database.SQL)
		return nil, err
	}

	return conn, nil
}

func lockSchemaMigration(db *gorm.DB, createTable bool) error {
	var err error
	record := SchemaMigration{Command: "init", Version: 0, Dirty: false}
	lockQuery := fmt.Sprintf("LOCK TABLE `%s` WRITE NOWAIT", record.TableName())

	var sqlErr *mysql.MySQLError
	if err = db.Exec(lockQuery).Error; err == nil {
		return nil
	} else if errors.As(err, &sqlErr) && sqlErr.Number == 1146 {
		if !createTable {
			logger.Logger.Error().Err(err).Msg(fmt.Sprintf("if you are running goto/up command and want to create `%s` table, you can add --create to command and try again", record.TableName()))
			return err
		}

		if err = db.Transaction(func(tx *gorm.DB) error { // Nested Transaction
			if err := db.AutoMigrate(&record); err != nil {
				return err
			} else if err = db.Create(&record).Error; err != nil {
				return err
			}

			return nil
		}); err == nil {
			logger.Logger.Info().Msg(fmt.Sprintf("Auto Migrate `%s`.`%s` Table successful!", configs.Config.Database.SQL.Name, record.TableName()))
			return db.Exec(lockQuery).Error
		}
	}
	return err
}

func unlockSchemaMigration(db *gorm.DB) error {
	if err := db.Exec("UNLOCK TABLE").Error; err != nil {
		logger.Logger.Err(err).Msg("unlock schema_migrations table error")
		return err
	}

	return nil
}

func getSchemaMigration(db *gorm.DB) (SchemaMigration, error) {
	var record SchemaMigration
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.WithContext(ctx).Order("`id` DESC").First(&record).Error; err == nil {
		return record, nil
	} else if errors.Is(err, context.DeadlineExceeded) {
		logger.Logger.Info().Err(err).Msg("table may be locked")
	}
	return record, err
}
