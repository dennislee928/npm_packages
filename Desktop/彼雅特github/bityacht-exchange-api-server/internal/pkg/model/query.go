package modelpkg

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/constraints"
)

func GetIntFromQuery[T constraints.Signed](ctx *gin.Context, key string, defaultValue T) (T, *errpkg.Error) {
	if val := ctx.Query(key); val == "" {
		return defaultValue, nil
	} else if iVal, err := strconv.ParseInt(val, 10, 64); err != nil {
		return defaultValue, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadQuery, Err: err}
	} else {
		return T(iVal), nil
	}
}

func GetBoolFromQuery(ctx *gin.Context, key string, defaultValue bool) (bool, *errpkg.Error) {
	if val := ctx.Query(key); val == "" {
		return defaultValue, nil
	} else if bVal, err := strconv.ParseBool(val); err != nil {
		return defaultValue, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadQuery, Err: err}
	} else {
		return bVal, nil
	}
}

func GetDateFromQuery(ctx *gin.Context, key string, defaultValue Date) (Date, *errpkg.Error) {
	val := ctx.Query(key)
	if val == "" {
		return defaultValue, nil
	}

	var dVal Date
	if err := dVal.Parse(JSONDateFormat, val); err != nil {
		return defaultValue, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadQuery, Err: err.Err}
	}

	return dVal, nil
}

func GetDateTimeFromQuery(ctx *gin.Context, key string, defaultValue DateTime) (DateTime, *errpkg.Error) {
	val := ctx.Query(key)
	if val == "" {
		return defaultValue, nil
	}

	var dtVal DateTime
	if err := dtVal.Parse(JSONDateTimeFormat, val); err != nil {
		return defaultValue, &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadQuery, Err: err.Err}
	}

	return dtVal, nil
}
