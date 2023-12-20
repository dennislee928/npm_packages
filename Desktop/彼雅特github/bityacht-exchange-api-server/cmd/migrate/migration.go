package migrate

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"bityacht-exchange-api-server/internal/pkg/logger"

	"gorm.io/gorm"
)

var migrations []IMigration
var sortOnce sync.Once

// IMigration for migration file to implement
type IMigration interface {
	Version() int64
	TableName() string
	Up(db *gorm.DB) error
	Down(db *gorm.DB) error
}

// Register the migration
func Register(migration IMigration) {
	migrations = append(migrations, migration)
}

func getMigrations() []IMigration {
	sortOnce.Do(func() {
		sort.Slice(migrations, func(i, j int) bool {
			return migrations[i].Version() < migrations[j].Version()
		})
	})
	return migrations
}

type migrateAction int

const (
	migrateActionUp migrateAction = iota + 1
	migrateActionDown
	migrateActionGoto
)

func doMigration(action migrateAction, step int, gotoVersion int64) error {
	db, err := connectDB()
	if err != nil {
		return err
	} else if err = lockSchemaMigration(db, createDB); err != nil {
		return err
	}
	defer unlockSchemaMigration(db)

	if migrationRecord, err := getSchemaMigration(db); err != nil {
		logger.Logger.Err(err).Msg("Get Schema Migration Error")
		return err
	} else if err = migrationRecord.IsClean(); err != nil {
		logger.Logger.Err(err).Any("MigrationRecord", migrationRecord).Send()
		return err
	} else {
		var command string
		var actuallyRun int
		migrations := getMigrations()

		switch action {
		case migrateActionUp:
			command = fmt.Sprintf("up %d", step)
		case migrateActionDown:
			command = fmt.Sprintf("down %d", step)
		case migrateActionGoto: // translate to up/down
			found := false
			step = 1

			if gotoVersion == migrationRecord.Version {
				logger.Logger.Info().Msg("database's version is same as specified")
				return nil
			} else if gotoVersion == 0 {
				action = migrateActionDown
				step = -1
				found = true
			} else if gotoVersion > migrationRecord.Version {
				action = migrateActionUp
				for _, migration := range migrations {
					if migration.Version() <= migrationRecord.Version {
						continue
					} else if migration.Version() == gotoVersion {
						found = true
						break
					}
					step++
				}
			} else {
				action = migrateActionDown
				for index := len(migrations) - 1; index >= 0; index-- {
					if migrations[index].Version() >= migrationRecord.Version {
						continue
					} else if migrations[index].Version() == gotoVersion {
						found = true
						break
					}
					step++
				}
			}

			if !found {
				return errors.New("bad goto version")
			} else if action == migrateActionUp {
				command = fmt.Sprintf("goto %d [up %d]", gotoVersion, step)
			} else {
				command = fmt.Sprintf("goto %d [down %d]", gotoVersion, step)
			}
		default:
			return errors.New("bad action")
		}

		if action == migrateActionUp {
			for _, migration := range migrations {
				if migrationRecord.Version >= migration.Version() {
					continue
				} else if step >= 0 && actuallyRun == step {
					break
				}

				migrationRecord = SchemaMigration{Command: command, Version: migration.Version(), Dirty: true, Error: "init"}
				if err := db.Create(&migrationRecord).Error; err != nil {
					return err
				} else if err = migration.Up(db); err != nil {
					if errOfUpdate := db.Model(&migrationRecord).UpdateColumn("error", err.Error()).Error; errOfUpdate != nil {
						logger.Logger.Err(errOfUpdate).Str("TableName", migration.TableName()).Int64("Version", migration.Version()).Msg("update migration record error")
					}

					return err
				} else if err = db.Model(&migrationRecord).UpdateColumns(map[string]interface{}{"dirty": false, "error": gorm.Expr("NULL")}).Error; err != nil {
					logger.Logger.Err(err).Str("TableName", migration.TableName()).Int64("Version", migration.Version()).Msg("update migration record error")
					return err
				}

				logger.Logger.Info().Str("TableName", migration.TableName()).Int64("Version", migration.Version()).Msg("migration up successful!")
				actuallyRun++
			}

			if actuallyRun == 0 {
				logger.Logger.Info().Msg("database is up to date!")
			}
		} else { // down
			for index := len(migrations) - 1; index >= 0; index-- {
				migration := migrations[index]
				if migrationRecord.Version < migration.Version() {
					continue
				} else if step >= 0 && actuallyRun == step {
					break
				}

				migrationRecord = SchemaMigration{Command: command, Dirty: true, Error: "init"}
				if index > 0 {
					migrationRecord.Version = migrations[index-1].Version()
				}

				if err := db.Create(&migrationRecord).Error; err != nil {
					return err
				} else if err = migration.Down(db); err != nil {
					if errOfUpdate := db.Model(&migrationRecord).UpdateColumn("error", err.Error()).Error; errOfUpdate != nil {
						logger.Logger.Err(errOfUpdate).Msg("update migration record error")
					}

					return err
				} else if err = db.Model(&migrationRecord).UpdateColumns(map[string]interface{}{"dirty": false, "error": gorm.Expr("NULL")}).Error; err != nil {
					logger.Logger.Err(err).Msg("update migration record error")
					return err
				}

				logger.Logger.Info().Int64("Version", migration.Version()).Int64("NewVersion", migrationRecord.Version).Msg("migration down successful!")
				actuallyRun++
			}

			if actuallyRun == 0 {
				logger.Logger.Info().Msg("database is out of date!")
			}
		}

		return nil
	}
}
