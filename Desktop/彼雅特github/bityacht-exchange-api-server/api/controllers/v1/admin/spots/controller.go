package spots

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	usersspottransfers "bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	csvpkg "bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/wallet"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary 	取得提入幣列表
// @Description 取得提入幣列表
// @Tags 		Admin-spots
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		searcher query modelpkg.Searcher false "searcher"
// @Param		query query usersspottransfers.GetListRequest false "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]usersspottransfers.Transfer}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/spots [get]
func GetListHandler(ctx *gin.Context) {
	var req usersspottransfers.GetListRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	searcher := modelpkg.GetSearcherFromQuery(ctx)

	var err *errpkg.Error
	if resp.Data, err = usersspottransfers.GetList(req, &resp.Paginator, searcher); errpkg.Handler(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出提入幣列表
// @Description 匯出提入幣列表
// @Tags 		Admin-spots
// @Security	BearerAuth
// @Accept 		json
// @Param 		query query usersspottransfers.GetExportRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/spots/export [get]
func ExportHandler(ctx *gin.Context) {
	var req usersspottransfers.GetExportRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := usersspottransfers.GetExport(req)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("spots_%s.csv", time.Now().Format("20060102"))
	csvpkg.ExportCSVFile(ctx, filename, usersspottransfers.GetTransferCSVHeaders(), data)
}

func newZipFileHeader(filenameFormat string, a ...any) *zip.FileHeader {
	return &zip.FileHeader{
		// Copy from zip/writer.go func (w *Writer) Create(name string) (io.Writer, error)
		Name:   fmt.Sprintf(filenameFormat, a...),
		Method: zip.Deflate,
		// Set the time
		Modified: time.Now(),
	}
}

func createAndWriteAegisImport(zipWriter *zip.Writer, timeStr string, csvFileIndex int, tempRecordOfAegisImport [][]string) error {
	if len(tempRecordOfAegisImport) == 0 {
		return nil
	}

	aegisImportWriter, err := zipWriter.CreateHeader(newZipFileHeader("aegisImport_%s_%d.csv", timeStr, csvFileIndex))
	if err != nil {
		return err
	}

	aegisImportCSV := csv.NewWriter(aegisImportWriter)
	if err = aegisImportCSV.Write([]string{"Address", "Amount"}); err != nil {
		return err
	}

	for _, tempRecord := range tempRecordOfAegisImport {
		if err = aegisImportCSV.Write(tempRecord); err != nil {
			return err
		}
	}

	aegisImportCSV.Flush()
	return nil
}

// @Summary 	Aegis 匯出
// @Description Aegis 匯出
// @Tags 		Admin-spots
// @Security	BearerAuth
// @Accept 		json
// @Param 		query query AegisExportRequest true "query"
// @Success 	200
// @Success 	204
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/spots/aegis-export [get]
func AegisExportHandler(ctx *gin.Context) {
	var req AegisExportRequest

	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	records, err := usersspottransfers.GetAegisExport(req.Mainnet, req.StartAt, req.EndAt)
	if errpkg.Handler(ctx, err) {
		return
	}

	csvRecordLimit := 50
	switch req.Mainnet {
	case wallet.MainnetBTC: // UTXO: 2400
		csvRecordLimit = 2400
	}

	tmpFile, rawErr := os.CreateTemp("", "aegis-export")
	if errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeWriteCSV, rawErr) {
		return
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			errLogger := logger.GetGinRequestLogger(ctx)
			errLogger.Err(err).Str("service", "aegis export").Msg("remove failed")
		}
	}()

	var (
		csvFileIndex            int
		bityachImportCSV        *csv.Writer
		zipWriter               = zip.NewWriter(tmpFile)
		timeStr                 = time.Now().Format("0601021503")
		addressMap              = make(map[string]struct{})
		tempRecordOfAegisImport = make([][]string, 0)
		csvFileRecordCount      = csvRecordLimit
		emptyZip                = true
	)

	for recordIndex, record := range records {
		if _, ok := addressMap[record.ToAddress]; ok {
			continue
		}
		addressMap[record.ToAddress] = struct{}{}

		if csvFileRecordCount+1 > csvRecordLimit {
			if bityachImportCSV != nil {
				bityachImportCSV.Flush()

				// Create and Write aegisImport
				if err := createAndWriteAegisImport(zipWriter, timeStr, csvFileIndex, tempRecordOfAegisImport); errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeWriteCSV, err) {
					return
				}
			}
			csvFileIndex++

			bityachImportWriter, rawErr := zipWriter.CreateHeader(newZipFileHeader("bityachImport_%s_%d.csv", timeStr, csvFileIndex))
			if errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeWriteCSV, rawErr) {
				return
			}

			if _, err := bityachImportWriter.Write(csvpkg.GetUTF8Bom()); err != nil {
				errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeWriteCSV, Err: err})
				return
			}

			bityachImportCSV = csv.NewWriter(bityachImportWriter)
			if rawErr = bityachImportCSV.Write(usersspottransfers.GetAegisTransferCSVHeaders()); errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeWriteCSV, rawErr) {
				return
			}

			emptyZip = false
			csvFileRecordCount = 0

			if remainRecordCount := len(records) - recordIndex; remainRecordCount < csvRecordLimit {
				tempRecordOfAegisImport = make([][]string, 0, remainRecordCount)
			} else {
				tempRecordOfAegisImport = make([][]string, 0, csvRecordLimit)
			}
		}

		if rawErr = bityachImportCSV.Write(record.ToCSV()); errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeWriteCSV, rawErr) {
			return
		}
		tempRecordOfAegisImport = append(tempRecordOfAegisImport, []string{record.ToAddress, record.Amount.String()})

		csvFileRecordCount++
	}

	if bityachImportCSV != nil {
		bityachImportCSV.Flush()

		if rawErr = createAndWriteAegisImport(zipWriter, timeStr, csvFileIndex, tempRecordOfAegisImport); errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeWriteCSV, rawErr) {
			return
		}
	}

	zipWriter.Close()

	if !emptyZip {
		ctx.FileAttachment(tmpFile.Name(), fmt.Sprintf("aegis_%s.zip", timeStr))
	} else {
		ctx.Status(http.StatusNoContent)
	}
}

// @Summary 	Aegis匯入
// @Description Aegis匯入
// @Tags 		Admin-spots
// @Security	BearerAuth
// @Accept 		mpfd
// @Param       csv formData file false "bityachImport.csv"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/spots/aegis-import [get]
func AegisImportHandler(ctx *gin.Context) {
	var req AegisImportRequest
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	csvFile, err := req.CSV.Open()
	if errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}
	defer func() {
		if err := csvFile.Close(); err != nil {
			logger.Logger.Err(err).Str("filename", req.CSV.Filename).Msg("AegisImportHandler close csvFile error")
		}
	}()

	if contentType := mime.TypeByExtension(path.Ext(req.CSV.Filename)); !strings.HasPrefix(contentType, "text/csv") {
		errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, errors.New("file type not allowed"))
		return
	}

	csvReader := csv.NewReader(csvFile)
	if headers, err := csvReader.Read(); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadCSVContent, err) {
		return
	} else if orignalHeaders := usersspottransfers.GetAegisTransferCSVHeaders(); len(headers) != len(orignalHeaders) {
		errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadCSVContent, errors.New("bad csv headers"))
		return
	}

	records := make([]usersspottransfers.AegisTransfer, 0, 2400)
	for {
		var record usersspottransfers.AegisTransfer

		rawRecord, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadCSVContent, err) {
			return
		} else if err := record.FromCSV(rawRecord); errpkg.Handler(ctx, err) {
			return
		}

		records = append(records, record)
	}

	if err := usersspottransfers.AegisImport(records); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}
