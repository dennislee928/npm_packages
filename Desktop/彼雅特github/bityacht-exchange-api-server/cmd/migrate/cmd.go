package migrate

import (
	"errors"
	"fmt"
	"strconv"

	"bityacht-exchange-api-server/internal/pkg/logger"

	"github.com/spf13/cobra"
)

var createDB bool
var downAll bool

func init() {
	gotoAndUpUsage := fmt.Sprintf("Auto create database and `%s` table if not exist.", SchemaMigration{}.TableName())
	gotoCmd.PersistentFlags().BoolVar(&createDB, "create", false, gotoAndUpUsage)
	upCmd.PersistentFlags().BoolVar(&createDB, "create", false, gotoAndUpUsage)
	downCmd.PersistentFlags().BoolVar(&downAll, "all", false, "Apply all down migrations.")

	Cmd.AddCommand(gotoCmd)
	Cmd.AddCommand(upCmd)
	Cmd.AddCommand(downCmd)
	Cmd.AddCommand(forceCmd)
	Cmd.AddCommand(versionCmd)
}

var Cmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration (SQL).",
}

var gotoCmd = &cobra.Command{
	Use:   "goto {V}",
	Short: "Migrate to version V.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetVersion, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		return doMigration(migrateActionGoto, 0, targetVersion)
	},
}

var upCmd = &cobra.Command{
	Use:   "up [N]",
	Short: "Apply all or N up migrations.",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		nStep := -1
		if len(args) == 1 {
			if nStep, err = strconv.Atoi(args[0]); err != nil {
				return err
			} else if nStep <= 0 {
				return errors.New("bad n, n should > 0 or just empty")
			}
		}

		return doMigration(migrateActionUp, nStep, 0)
	},
}

var downCmd = &cobra.Command{
	Use:   "down [N]",
	Short: "Apply all or N down migrations, or Use --all to apply all down migrations.",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		nStep := 0
		if downAll {
			nStep = -1
		} else {
			if len(args) == 1 {
				if nStep, err = strconv.Atoi(args[0]); err != nil {
					return err
				}
			}
			if nStep <= 0 {
				return errors.New("bad n, n should > 0 or Use --all to apply all down migrations")
			}
		}

		return doMigration(migrateActionDown, nStep, 0)
	},
}

var forceCmd = &cobra.Command{
	Use:   "force V",
	Short: "Set version V but don't run migration (ignores dirty state).",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetVersion, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}

		db, err := connectDB()
		if err != nil {
			return err
		} else if err = lockSchemaMigration(db, false); err != nil {
			return err
		}
		defer unlockSchemaMigration(db)

		if targetVersion == 0 {
			return db.Create(&SchemaMigration{Command: "force", Version: targetVersion, Dirty: false}).Error
		}
		for _, migration := range getMigrations() {
			if migration.Version() == targetVersion {
				return db.Create(&SchemaMigration{Command: "force", Version: targetVersion, Dirty: false}).Error
			}
		}

		return errors.New("v is not found in migrations")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current migration version.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if conn, err := connectDB(); err != nil {
			return
		} else if migrationRecord, err := getSchemaMigration(conn); err != nil {
			logger.Logger.Err(err).Send()
			return
		} else {
			logger.Logger.Info().Any("MigrationRecord", migrationRecord).Send()
		}
	},
}
