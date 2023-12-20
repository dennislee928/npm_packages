package seed

import (
	"errors"
	"sort"
	"sync"

	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/pkg/logger"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var seederMap map[string]ISeeder
var seedersOrder []string
var seeders []ISeeder
var sortOnce sync.Once
var ErrNotImplement = errors.New("not implement")

func init() {
	seederMap = make(map[string]ISeeder)
}

// ISeeder for seeder file to implement
type ISeeder interface {
	SeederName() string // For Setting the Seeder order
	TableName() string
	Default(db *gorm.DB) error
	Fake(db *gorm.DB) error
}

// Register the seeder
func Register(seeder ISeeder) {
	seederMap[seeder.SeederName()] = seeder
}

// SetSeedersOrder for make sure the order of running seed
func SetSeedersOrder(order []string) {
	seedersOrder = order
}

func getSeeders() []ISeeder {
	sortOnce.Do(func() {
		sortedMap := make(map[string]struct{})

		for _, seederName := range seedersOrder {
			if iSeeder, ok := seederMap[seederName]; !ok {
				logger.Logger.Error().Msg("seeder '" + seederName + "' not found in seederMap. Make sure it is registered.")
			} else if _, ok := sortedMap[seederName]; !ok {
				seeders = append(seeders, iSeeder)
				sortedMap[seederName] = struct{}{}
			} else {
				logger.Logger.Error().Msg("seeder '" + seederName + "' has been shown in seeders order more than once.")
			}
		}

		sortedIndex := len(seeders)
		for seederName, iSeeder := range seederMap {
			if _, ok := sortedMap[seederName]; !ok {
				seeders = append(seeders, iSeeder)
			}
		}

		sort.Slice(seeders[sortedIndex:], func(i, j int) bool {
			return seeders[i+sortedIndex].SeederName() < seeders[j+sortedIndex].SeederName()
		})
	})
	return seeders
}

func connectDB() (*gorm.DB, error) {
	conn, err := sql.Connect(configs.Config.Database.SQL.Name)
	if err != nil {
		logger.Logger.Error().Err(err).Any("Database Config", configs.Config.Database.SQL)
		return nil, err
	}

	return conn, nil
}

func RetryCreateWhenDuplicate(db *gorm.DB, record any) error {
	if db == nil {
		return errors.New("nil gorm DB")
	}

	for retry := 0; retry < 50; retry++ {
		err := db.Create(record).Error
		if err == nil {
			return nil
		}

		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			continue
		}

		return err
	}
	return errors.New("retry create failed")
}
