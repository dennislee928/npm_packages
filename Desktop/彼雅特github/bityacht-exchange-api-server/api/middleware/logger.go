package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"bityacht-exchange-api-server/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// SensitivePathKey for Check API is sensitive or not, if sensitive then the body will not be logged.
const SensitivePathKey = "sensitivePath"

// SkipLogPathKey for Check API is need to Skip Log or not.
const SkipLogPathKey = "skipLogPath"

func SetSensitive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(SensitivePathKey, struct{}{})
	}
}

func SetSkipLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(SkipLogPathKey, struct{}{})
	}
}

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()

		reqBody := new(bytes.Buffer)
		tee := io.TeeReader(ctx.Request.Body, reqBody)
		ctx.Request.Body = io.NopCloser(tee)

		ctx.Next()

		if _, ok := ctx.Get(SkipLogPathKey); ok {
			return
		}

		end := time.Now()
		status := ctx.Writer.Status()

		logLevel := zerolog.InfoLevel
		switch {
		case status >= http.StatusInternalServerError:
			logLevel = zerolog.ErrorLevel
		case status >= http.StatusBadRequest:
			logLevel = zerolog.WarnLevel
		}

		reqLogger := logger.GetGinRequestLogger(ctx)
		reqLog := reqLogger.WithLevel(logLevel).Str("ip", ctx.ClientIP()).Int("status", status).Dur("duration", end.Sub(start))
		if len(ctx.Errors) > 0 {
			reqLog = reqLog.Str("err", ctx.Errors.String())
		}
		if reqBody.Len() > 0 {
			compactBuffer := new(bytes.Buffer)

			if _, ok := ctx.Get(SensitivePathKey); !ok {
				if err := json.Compact(compactBuffer, reqBody.Bytes()); err == nil {
					reqBody = compactBuffer
				}

				reqLog = reqLog.Str("body", reqBody.String())
			}
		}

		reqLog.Send()
	}
}
