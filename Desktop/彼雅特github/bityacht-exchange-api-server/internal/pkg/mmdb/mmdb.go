package mmdb

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	_ "embed"
	"errors"
	"net"
	"net/http"

	"github.com/oschwald/maxminddb-golang"
)

var reader *maxminddb.Reader

//go:embed city.mmdb
var cityBytes []byte

func init() {
	var err error

	if reader, err = maxminddb.FromBytes(cityBytes); err != nil {
		panic(err)
	}
}

func LookupCity(ip string) (CityResult, *errpkg.Error) {
	var result CityResult

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return result, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeParseIP, Err: errors.New(ip)}
	}

	if err := reader.Lookup(parsedIP, &result); err != nil {
		return result, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeLookUpMMDB, Err: err}
	}

	return result, nil
}
