package logger

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"bityacht-exchange-api-server/configs"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger zerolog.Logger

const timeFormat = time.RFC3339Nano

func Init() {
	zerolog.SetGlobalLevel(configs.Config.Log.Level)
	zerolog.TimeFieldFormat = timeFormat
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		if cwd, err := os.Getwd(); err == nil {
			if rel, err := filepath.Rel(cwd, file); err == nil {
				file = rel
			}
		}

		return file + ":" + strconv.Itoa(line)
	}

	var logWriter io.Writer
	switch configs.Config.Log.Filename {
	case "":
		logWriter = os.Stderr
	case "console":
		logWriter = zerolog.NewConsoleWriter()
	default:
		logWriter = &lumberjack.Logger{
			Filename:   configs.Config.Log.Filename,
			MaxSize:    configs.Config.Log.MaxSize,
			MaxBackups: configs.Config.Log.MaxBackups,
			MaxAge:     configs.Config.Log.MaxAge,
			Compress:   configs.Config.Log.Compress,
		}
	}

	Logger = zerolog.New(logWriter).With().Timestamp().Logger()
}

func GetGinRequestLogger(ctx *gin.Context) zerolog.Logger {
	return Logger.With().Str("reqID", requestid.Get(ctx)).Str("uri", ctx.Request.Method+" "+ctx.Request.RequestURI).Logger()
}
