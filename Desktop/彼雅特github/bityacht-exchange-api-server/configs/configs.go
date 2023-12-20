package configs

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	gormlogger "gorm.io/gorm/logger"
)

// ConfigFile is name of config file (without extension), default is 'config'
var ConfigFile string

// Config contain all config in this project, DONOT READ ME in init function
var Config AllConfig

func Init() {
	if ConfigFile == "" {
		ConfigFile = "config"
	}

	viper.SetConfigName(ConfigFile)
	viper.SetConfigType("yaml")    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")       // path to look for the config file in
	viper.AddConfigPath("configs") // path to look for the config file in

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // For Read environment variables
	viper.AutomaticEnv()                                   // For Read environment variables

	if err := viper.ReadInConfig(); err != nil && os.IsNotExist(err) {
		panic(err)
	}

	Config.Server = ServerConfig{
		Host:         getStrEnv("server.host", ""),
		Port:         getUintEnv("server.port", uint(80)),
		ReadTimeout:  getDurationEnv("server.readTimeout", 60*time.Second),
		WriteTimeout: getDurationEnv("server.writeTimeout", 60*time.Second),
		CertFile:     getStrEnv("server.certFile", ""),
		KeyFile:      getStrEnv("server.keyFile", ""),
		StoragePath:  getStrEnv("server.storagePath", "./"),
	}

	Config.JWT = JWTConfig{
		Key:                  getByteArrayEnv("jwt.key", []byte("bityacht-exchange-api-server")),
		AccessTokenLifetime:  getDurationEnv("jwt.accessTokenLifetime", 60*time.Minute),
		RefreshTokenLifetime: getDurationEnv("jwt.refreshTokenLifetime", 180*time.Minute),
	}

	Config.Database = DatabaseConfig{
		SQL: SQLConfig{
			Host:     getStrEnv("database.sql.host", "localhost"),
			Port:     getUintEnv("database.sql.port", uint(3306)),
			Name:     getStrEnv("database.sql.name", "bityacht_exchange"),
			User:     getStrEnv("database.sql.user", "user"),
			Password: getStrEnv("database.sql.password", "password"),
			LogLevel: getLogLevelEnv("database.sql.logLevel", gormlogger.Silent, gormlogger.Silent, gormlogger.Info),
		},
	}

	Config.Cache = CacheConfig{
		Redis: RedisConfig{
			Host:           getStrEnv("cache.redis.host", "localhost"),
			Port:           getUintEnv("cache.redis.port", uint(6379)),
			DB:             getIntEnv("cache.redis.db", 0),
			Username:       getStrEnv("cache.redis.username", ""),
			Password:       getStrEnv("cache.redis.password", ""),
			MaxConnections: getIntEnv("cache.redis.maxConnections", 0),
			ReadTimeout:    getDurationEnv("cache.redis.readTimeout", 3*time.Second),
			WriteTimeout:   getDurationEnv("cache.redis.writeTimeout", 3*time.Second),
		},
	}

	Config.Log = LogConfig{
		Level:      getLogLevelEnv("log.level", zerolog.InfoLevel, zerolog.TraceLevel, zerolog.Disabled),
		Filename:   getStrEnv("log.filename", ""),
		MaxSize:    getIntEnv("log.maxSize", 10),
		MaxBackups: getIntEnv("log.maxBackups", 10),
		MaxAge:     getIntEnv("log.maxAge", 30),
		Compress:   getBoolEnv("log.compress", true),
	}

	Config.Email = EmailConfig{
		Enable:              getBoolEnv("email.enable", false),
		Nickname:            getStrEnv("email.nickname", "example"),
		Account:             getStrEnv("email.account", "example@skycloud.com.tw"),
		Password:            getStrEnv("email.password", "example12345678"),
		Host:                getStrEnv("email.host", "smtp.gmail.com"),
		Port:                getStrEnv("email.port", "465"),
		SSL:                 getBoolEnv("email.ssl", false),
		TLS:                 getBoolEnv("email.tls", true),
		BitYachtFrontendURL: getStrEnv("email.bitYachtfrontendURL", "https://bityacht.io"),
	}

	Config.SMS = SMSConfig{
		Enable:    getBoolEnv("sms.enable", false),
		Provider:  getStrEnv("sms.provider", "mitake"),
		Host:      getStrEnv("sms.host", "https://smsapi.mitake.com.tw/api/mtk"),
		Username:  getStrEnv("sms.username", "username"),
		Password:  getStrEnv("sms.password", "password"),
		ConnCount: getIntEnv("sms.connCount", 15),
	}

	Config.Exchange = ExchangeConfig{
		UpdateTrendInterval: getDurationEnv("exchange.updateTrendInterval", 5*time.Second),
		Binance: BinanceConfig{
			Debug:            getBoolEnv("exchange.binance.debug", false),
			Apikey:           getStrEnv("exchange.binance.apikey", ""),
			SecretKey:        getStrEnv("exchange.binance.secretKey", ""),
			UseTestnet:       getBoolEnv("exchange.binance.useTestnet", false),
			TestnetApikey:    getStrEnv("exchange.binance.testnetApikey", ""),
			TestnetSecretKey: getStrEnv("exchange.binance.testnetSecretKey", ""),
		},
	}

	Config.Receipt = ReceiptConfig{
		APIMode: getIntEnv("receipt.api.mode", 0),
		AppCode: getStrEnv("receipt.app.code", ""),
		AppKey:  getStrEnv("receipt.app.key", ""),
		APIAcc:  getStrEnv("receipt.api.account", ""),
		APIPwd:  getStrEnv("receipt.api.password", ""),
	}

	Config.KYC = KYCConfig{
		APIToken:              getStrEnv("kyc.apiToken", ""),
		CallbackHost:          getStrEnv("kyc.callbackHost", ""),
		SuccessURL:            getStrEnv("kyc.successURL", ""),
		ErrorURL:              getStrEnv("kyc.errorURL", ""),
		DDTaskSearchSettingID: getIntEnv("kyc.ddTaskSearchSettingID", int64(0)),
	}

	if err := viper.UnmarshalKey("wallet", &Config.Wallet); err != nil {
		panic(err)
	}
}

type AllConfig struct {
	Server   ServerConfig
	JWT      JWTConfig
	Database DatabaseConfig
	Cache    CacheConfig
	Log      LogConfig
	Email    EmailConfig
	SMS      SMSConfig
	Exchange ExchangeConfig
	Receipt  ReceiptConfig
	KYC      KYCConfig
	Wallet   WalletConfig
}

type ServerConfig struct {
	Host         string
	Port         uint
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	CertFile     string
	KeyFile      string
	StoragePath  string
}

type JWTConfig struct {
	Key                  []byte
	AccessTokenLifetime  time.Duration
	RefreshTokenLifetime time.Duration
}

type DatabaseConfig struct {
	SQL SQLConfig
}

type SQLConfig struct {
	LogLevel gormlogger.LogLevel
	Host     string
	Port     uint
	Name     string // Database Name
	User     string
	Password string
}

type CacheConfig struct {
	Redis RedisConfig
}

type RedisConfig struct {
	Host           string
	Port           uint
	DB             int
	Username       string
	Password       string
	MaxConnections int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
}

type LogConfig struct {
	Level      zerolog.Level
	Filename   string // If Log File is empty, it will write to console
	MaxSize    int    // Megabytes
	MaxBackups int
	MaxAge     int // Days
	Compress   bool
}

type EmailConfig struct {
	Enable              bool
	Nickname            string
	Account             string
	Password            string
	Host                string
	Port                string
	SSL                 bool
	TLS                 bool
	BitYachtFrontendURL string // BitYacht Exchange Frontend
}

func (ec EmailConfig) IsFake() bool {
	return ec.Host == "localhost" && ec.Password == "fakeSMTP"
}

type SMSConfig struct {
	Enable    bool
	Provider  string
	Host      string
	Username  string
	Password  string
	ConnCount int
}

type ExchangeConfig struct {
	UpdateTrendInterval time.Duration

	Binance BinanceConfig
}

type BinanceConfig struct {
	Debug            bool
	Apikey           string
	SecretKey        string
	UseTestnet       bool
	TestnetApikey    string
	TestnetSecretKey string
}

type ReceiptConfig struct {
	APIMode int

	AppCode, AppKey string
	APIAcc, APIPwd  string
}

type KYCConfig struct {
	APIToken              string
	CallbackHost          string
	SuccessURL            string
	ErrorURL              string
	DDTaskSearchSettingID int64
}

type WalletConfig struct {
	APIMode int `mapstructure:"apiMode"`

	Wallets []struct {
		Currency      string `mapstructure:"currency"`
		Mainnet       string `mapstructure:"mainnet"`
		Flow          int    `mapstructure:"flow"`
		ID            string `mapstructure:"id"`
		Token         string `mapstructure:"token"`
		Secret        string `mapstructure:"secret"`
		RefreshToken  string `mapstructure:"refreshToken"` // TODO: need to refresh?
		OrderIDPrefix string `mapstructure:"orderIDPrefix"`
	} `mapstructure:"wallets"`
}
