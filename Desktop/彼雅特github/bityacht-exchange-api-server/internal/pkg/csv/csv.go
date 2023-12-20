package csv

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"encoding/csv"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const ContentType = "text/csv"

var utf8Bom = []byte{0xEF, 0xBB, 0xBF}

type Record interface {
	ToCSV() []string
}

func GetUTF8Bom() []byte {
	return utf8Bom
}

func ExportCSVFile(ctx *gin.Context, filename string, headers []string, records []Record) {
	tmpFile, err := os.CreateTemp("", "csv")
	if err != nil {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeWriteCSV, Err: err})
		return
	}
	defer func() {
		if err = os.Remove(tmpFile.Name()); err != nil {
			errLogger := logger.GetGinRequestLogger(ctx)
			errLogger.Err(err).Str("service", "export csv file").Msg("remove failed")
		}
	}()

	if _, err := tmpFile.Write(utf8Bom); err != nil {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeWriteCSV, Err: err})
		return
	}
	tmpFile.Sync()

	writer := csv.NewWriter(tmpFile)
	if err := writer.Write(headers); err != nil {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeWriteCSV, Err: err})
		return
	}

	for _, record := range records {
		if err := writer.Write(record.ToCSV()); err != nil {
			errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeWriteCSV, Err: err})
			return
		}
	}

	writer.Flush()

	ctx.FileAttachment(tmpFile.Name(), filename)
}
