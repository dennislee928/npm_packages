package suspicioustxs

import (
	"bityacht-exchange-api-server/internal/database/sql/suspicioustransactions"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/storage"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary 	取得可疑交易列表
// @Description 取得可疑交易列表
// @Tags 		Admin-SuspiciousTxs
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		searcher query modelpkg.Searcher false "searcher"
// @Param		query query suspicioustransactions.GetSuspiciousTXListRequest false "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]suspicioustransactions.SuspiciousTX}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/suspicious-txs [get]
func GetListHandler(ctx *gin.Context) {
	var req suspicioustransactions.GetSuspiciousTXListRequest

	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	var (
		resp     = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
		searcher = modelpkg.GetSearcherFromQuery(ctx)
		err      *errpkg.Error
	)

	if resp.Data, err = suspicioustransactions.GetSuspiciousTXList(req, &resp.Paginator, searcher); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	取得可疑交易詳細資訊
// @Description 取得可疑交易詳細資訊
// @Tags 		Admin-SuspiciousTxs
// @Security	BearerAuth
// @Param 		id path string true "Suspicious Transactions ID"
// @Success 	200 {object} suspicioustransactions.SuspiciousTXDetail
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/suspicious-txs/:id [get]
func GetDetailHandler(ctx *gin.Context) {
	var req GetDetailRequest

	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	resp, err := suspicioustransactions.GetSuspiciousTXDetail(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	更新可疑交易資訊
// @Description 更新可疑交易資訊
// @Tags 		Admin-SuspiciousTxs
// @Security	BearerAuth
// @Param 		id path string true "Suspicious Transactions ID"
// @Param 		body body suspicioustransactions.UpdateRequest true "body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/suspicious-txs/:id [patch]
func PatchHandler(ctx *gin.Context) {
	var req suspicioustransactions.UpdateRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	if err := req.Validate(); errpkg.Handler(ctx, err) {
		return
	}

	if err := suspicioustransactions.UpdateRecord(req); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	匯出可疑交易列表
// @Description 匯出可疑交易列表
// @Tags 		Admin-SuspiciousTxs
// @Accept 		json
// @Param 		query query suspicioustransactions.ExportSuspiciousTXCSVRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router		/admin/suspicious-txs/export-csv [get]
func ExportCSVHandler(ctx *gin.Context) {
	var req suspicioustransactions.ExportSuspiciousTXCSVRequest

	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := suspicioustransactions.GetExportCSV(req)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("suspiciousList_%s.csv", time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, suspicioustransactions.GetSuspiciousTXCSVHeaders(), data)
}

func getFilePath(id int64, fileType suspicioustransactions.UpdateFilesType, filename string) string {
	return storage.GetSuspiciousTxPath(fmt.Sprintf("%d/%d", id, fileType), filename)
}

// @Summary 	上傳可疑交易相關檔案
// @Description 上傳可疑交易相關檔案
// @Tags 		Admin-SuspiciousTxs
// @Accept		mpfd
// @Param 		id path string true "Suspicious Transactions ID"
// @Param       file formData file false "file"
// @Param 		body formData UploadFileRequest true "body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router		/admin/suspicious-txs/:id/file [post]
func UploadFileHandler(ctx *gin.Context) {
	var (
		req         UploadFileRequest
		isCompleted bool
	)

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	filePath := getFilePath(req.ID, req.UploadType, req.File.Filename)
	if err := storage.CheckAndSaveUploadFile(nil, req.File, filePath, os.O_EXCL); errpkg.Handler(ctx, err) {
		return
	}

	defer func() {
		if !isCompleted {
			if err := os.Remove(filePath); err != nil {
				errLogger := logger.GetGinRequestLogger(ctx)
				errLogger.Err(err).Str("filePath", filePath).Msg("remove saved file error")
			}
		}
	}()

	if err := suspicioustransactions.UpdateFiles(req.ID, req.UploadType, req.File.Filename); err != nil && err.Code != errpkg.CodeRecordNoChange {
		errpkg.Handler(ctx, err)
		return
	}
	isCompleted = true

	ctx.Status(http.StatusOK)
}

// @Summary 	下載可疑交易相關檔案
// @Description 下載可疑交易相關檔案
// @Tags 		Admin-SuspiciousTxs
// @Param 		id path string true "Suspicious Transactions ID"
// @Param       query query DownloadFileRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router		/admin/suspicious-txs/:id/file [get]
func DownloadFileHandler(ctx *gin.Context) {
	var req DownloadFileRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	ctx.File(getFilePath(req.ID, req.FileType, req.Filename))
}

// @Summary 	刪除可疑交易相關檔案
// @Description 刪除可疑交易相關檔案
// @Tags 		Admin-SuspiciousTxs
// @Param 		id path string true "Suspicious Transactions ID"
// @Param       query query DeleteFileRequest true "query"
// @Success 	204
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router		/admin/suspicious-txs/:id/file [delete]
func DeleteFileHandler(ctx *gin.Context) {
	var req DeleteFileRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	if err := suspicioustransactions.UpdateFiles(req.ID, -req.FileType, req.Filename); errpkg.Handler(ctx, err) {
		return
	}

	filePath := getFilePath(req.ID, req.FileType, req.Filename)
	if err := os.Remove(filePath); err != nil {
		errLogger := logger.GetGinRequestLogger(ctx)
		errLogger.Err(err).Str("filePath", filePath).Msg("remove file error")
	}

	ctx.Status(http.StatusNoContent)
}
