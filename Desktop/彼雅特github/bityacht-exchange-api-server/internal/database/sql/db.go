package sql

import (
	"fmt"
	"sync"
	"time"

	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var conn *gorm.DB
var lock sync.RWMutex

// DB will return the gorm.DB
func DB() *gorm.DB {
	lock.RLock()
	if conn != nil {
		defer lock.RUnlock()

		return conn
	}
	lock.RUnlock()

	lock.Lock()
	defer lock.Unlock()

	if conn != nil {
		return conn
	}

	var err error
	if conn, err = Connect(configs.Config.Database.SQL.Name); err != nil {
		logger.Logger.Fatal().Err(err).Any("DatabaseConfig", configs.Config.Database).Msg("fail to connect SQL Server")
	}

	return conn
}

// Connect to DB by dbName and config of this project
func Connect(dbName string) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		DSN:                       fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=UTC&readTimeout=30s", configs.Config.Database.SQL.User, configs.Config.Database.SQL.Password, configs.Config.Database.SQL.Host, configs.Config.Database.SQL.Port, dbName),
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{
		Logger: gormlogger.New(&logger.Logger, gormlogger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			LogLevel:                  configs.Config.Database.SQL.LogLevel,
		}),
	})
}
