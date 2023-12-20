package configs

import (
	"strconv"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/exp/constraints"
)

func getStrEnv(key string, defaultVal string) string {
	if val := viper.GetString(key); val != "" {
		return val
	}

	return defaultVal
}

func getIntEnv[T constraints.Signed](key string, defaultVal T) T {
	if val := viper.GetString(key); val != "" {
		if iVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			return T(iVal)
		}
	}

	return defaultVal
}

func getUintEnv[T constraints.Unsigned](key string, defaultVal T) T {
	if val := viper.GetString(key); val != "" {
		if uVal, err := strconv.ParseUint(val, 10, 64); err == nil {
			return T(uVal)
		}
	}

	return defaultVal
}

func getLogLevelEnv[T constraints.Signed](key string, defaultVal T, minVal T, maxVal T) T {
	if val := viper.GetString(key); val != "" {
		if iVal, err := strconv.ParseInt(val, 10, 64); err == nil {
			if iVal >= int64(minVal) && iVal <= int64(maxVal) {
				return T(iVal)
			}
		}
	}

	return defaultVal
}

func getBoolEnv(key string, defaultVal bool) bool {
	if val := viper.GetString(key); val != "" {
		if bVal, err := strconv.ParseBool(val); err == nil {
			return bVal
		}
	}

	return defaultVal
}

func getDurationEnv(key string, defaultVal time.Duration) time.Duration {
	if val := viper.GetString(key); val != "" {
		if dValue, err := time.ParseDuration(val); err == nil {
			return dValue
		}
	}

	return defaultVal
}

func getByteArrayEnv(key string, defaultVal []byte) []byte {
	if val := viper.GetString(key); val != "" {
		return []byte(val)
	}

	return defaultVal
}
