package errpkg

import (
	"bityacht-exchange-api-server/internal/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handle will send the error response to client if err is not nil, and return true.
func Handler(ctx *gin.Context, err *Error) bool {
	if err == nil {
		return false
	}

	if err.HttpStatus == 0 {
		err.HttpStatus = http.StatusInternalServerError
	}
	if err.Code == 0 {
		err.Code = CodeUnknown
	}

	ctx.JSON(err.HttpStatus, err)

	return true
}

// HandlerWithCode will send the error response to client if err is not nil, and return true.
func HandlerWithCode(ctx *gin.Context, httpStatus int, errorCode Code, err error) bool {
	if err == nil {
		return false
	}

	resp := Error{Code: errorCode}

	switch err := err.(type) {
	case *Error:
		if err == nil {
			warnLogger := logger.GetGinRequestLogger(ctx)
			warnLogger.Warn().Msg("HandlerWithCode get *Error(nil) but err == nil, Don't Let it Happen !!!")

			httpStatus = http.StatusInternalServerError
			resp.Code = CodeBadErrorType
		} else {
			resp.Data = err.Data
			resp.Err = err.Err
		}
	case Error:
		resp.Data = err.Data
		resp.Err = err.Err
	default:
		resp.Err = err
	}

	ctx.JSON(httpStatus, resp)
	return true
}
