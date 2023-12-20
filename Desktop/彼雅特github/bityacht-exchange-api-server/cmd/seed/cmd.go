package seed

import (
	"errors"

	"bityacht-exchange-api-server/internal/pkg/logger"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	Offset    int
	UsersID   int64
	UserCount int
)

func init() {
	fakeSeedCmd.PersistentFlags().IntVar(&Offset, "offset", 0, "The offset of the fake seeder.")
	fakeSeedCmd.PersistentFlags().Int64Var(&UsersID, "usersID", 0, "The usersID of the fake seeder.")
	fakeSeedCmd.PersistentFlags().IntVar(&UserCount, "userCount", 10, "The user count of the fake seeder.")

	Cmd.AddCommand(defaultSeedCmd)
	Cmd.AddCommand(fakeSeedCmd)
}

var Cmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed can create records into Database (SQL), including default and fake data.",
}

var defaultSeedCmd = &cobra.Command{
	Use:   "default [SeederName]...",
	Short: "Default Seeder will create basic records into Database.",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := connectDB()
		if err != nil {
			return err
		}

		if len(args) == 0 { // Run all seeders
			for _, seeder := range getSeeders() {
				if err = seeder.Default(db); err != nil {
					if !errors.Is(err, ErrNotImplement) {
						logger.Logger.Err(err).Str("Seeder(default)", seeder.SeederName()).Msg("run seeder error!")
						return err
					}
				} else {
					logger.Logger.Info().Str("Seeder(default)", seeder.SeederName()).Msg("run seeder successful!")
				}
			}
		} else {
			for _, seederName := range args {
				if seeder, ok := seederMap[seederName]; !ok {
					return errors.New("seeder '" + seederName + "' not found.")
				} else if err = seeder.Default(db); err != nil {
					if !errors.Is(err, ErrNotImplement) {
						logger.Logger.Err(err).Str("Seeder(default)", seeder.SeederName()).Msg("run seeder error!")
						return err
					}
				} else {
					logger.Logger.Info().Str("Seeder(default)", seeder.SeederName()).Msg("run seeder successful!")
				}
			}
		}

		return nil
	},
}

var fakeSeedCmd = &cobra.Command{
	Use:   "fake [SeederName]...",
	Short: "Fake Seeder will create fake data into Database. (Only for development!!)",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := connectDB()
		if err != nil {
			return err
		}

		return db.Transaction(func(tx *gorm.DB) error {
			if len(args) == 0 { // Run all seeders
				for _, seeder := range getSeeders() {
					if err = seeder.Fake(tx); err != nil {
						if !errors.Is(err, ErrNotImplement) {
							logger.Logger.Err(err).Str("Seeder(fake)", seeder.SeederName()).Msg("run seeder error!")
							return err
						}
					} else {
						logger.Logger.Info().Str("Seeder(fake)", seeder.SeederName()).Msg("run seeder successful!")
					}
				}
			} else {
				for _, seederName := range args {
					if seeder, ok := seederMap[seederName]; !ok {
						return errors.New("seeder '" + seederName + "' not found.")
					} else if err = seeder.Fake(tx); err != nil {
						if !errors.Is(err, ErrNotImplement) {
							logger.Logger.Err(err).Str("Seeder(fake)", seeder.SeederName()).Msg("run seeder error!")
							return err
						}
					} else {
						logger.Logger.Info().Str("Seeder(fake)", seeder.SeederName()).Msg("run seeder successful!")
					}
				}
			}

			return nil
		})
	},
}
