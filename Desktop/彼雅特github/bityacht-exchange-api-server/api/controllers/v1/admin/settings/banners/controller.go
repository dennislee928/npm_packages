package banners

import (
	"net/http"
	"os"
	"path/filepath"

	"bityacht-exchange-api-server/internal/database/sql/banners"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

// Banner's allowed extensions
var allowedExtensions = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/gif":  {},
}

// @Summary     建立橫幅資料
// @Description 建立橫幅資料
// @Tags        Admin-settings
// @Security	BearerAuth
// @Accept      mpfd
// @Produce     json
// @Param       webImage formData file false "webImage"
// @Param       appImage formData file false "appImage"
// @Param       data formData banners.CreateRequest true "Data"
// @Success     201
// @Failure     400 {object} errpkg.JsonError
// @Failure     500 {object} errpkg.JsonError
// @Router      /admin/settings/banners [post]
func CreateHandler(ctx *gin.Context) {
	var req banners.CreateRequest
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	allCreated := false
	web, app := "", ""
	saved := make([]string, 0, 2)

	defer func() {
		if allCreated {
			return
		}
		for _, filePath := range saved {
			if err := os.Remove(filePath); err != nil {
				logger.Logger.Err(err).Msg("remove saved file error")
			}
		}
	}()

	if req.WebImage != nil {
		web = genFilenameUUID(req.WebImage.Filename)
		filePath := storage.GetBannerPath(web)
		if err := storage.CheckAndSaveUploadFile(allowedExtensions, req.WebImage, filePath, os.O_EXCL); errpkg.Handler(ctx, err) {
			return
		}
		saved = append(saved, filePath)
	}

	if req.AppImage != nil {
		app = genFilenameUUID(req.AppImage.Filename)
		filePath := storage.GetBannerPath(app)
		if err := storage.CheckAndSaveUploadFile(allowedExtensions, req.AppImage, filePath, os.O_EXCL); errpkg.Handler(ctx, err) {
			return
		}
		saved = append(saved, filePath)
	}

	recode := &banners.Model{
		WebImage:  web,
		AppImage:  app,
		Title:     req.Title,
		SubTitle:  req.SubTitle,
		Button:    req.Button,
		ButtonUrl: req.ButtonUrl,
		Status:    req.Status,
		StartAt:   req.StartAt.Time,
		EndAt:     req.EndAt.Time,
	}

	if err := banners.Create(recode); errpkg.Handler(ctx, err) {
		return
	}

	allCreated = true
	ctx.Status(http.StatusCreated)
}

// @Summary 	取得橫幅資料列表
// @Description 取得橫幅資料列表
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		status query int false "banners status, all = -1"
// @Success 	200 {object} modelpkg.GetResponse{data=[]banners.Banner}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/banners [get]
func GetListHandler(ctx *gin.Context) {
	status, err := modelpkg.GetIntFromQuery[int32](ctx, "status", -1)
	if errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	if resp.Data, err = banners.GetList(&resp.Paginator, status); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	編輯橫幅資料
// @Description 編輯橫幅資料
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param 		id path int true "banners id"
// @Param       webImage formData file false "webImage"
// @Param       appImage formData file false "appImage"
// @Param 		body body banners.UpdateRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/banners/{id} [patch]
func UpdateHandler(ctx *gin.Context) {
	var req banners.UpdateRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	var oldWeb, oldApp string
	if record, err := banners.Get(req.ID); errpkg.Handler(ctx, err) {
		return
	} else {
		oldWeb, oldApp = record.WebImage, record.AppImage
	}

	var (
		web, app    string
		isCompleted bool
		saved       = make([]string, 0, 2)
	)

	defer func() {
		if isCompleted {
			if web != "" {
				if err := os.Remove(storage.GetBannerPath(oldWeb)); err != nil {
					logger.Logger.Err(err).Str("oldFileDst", oldWeb).Msg("remove old temp file failed")
				}
			}
			if app != "" {
				if err := os.Remove(storage.GetBannerPath(oldApp)); err != nil {
					logger.Logger.Err(err).Str("oldFileDst", oldApp).Msg("remove old temp file failed")
				}
			}
			return
		}

		for _, filePath := range saved {
			if err := os.Remove(filePath); err != nil {
				logger.Logger.Err(err).Msg("remove saved file error")
			}
		}
	}()

	if req.WebImage != nil {
		web = genFilenameUUID(req.WebImage.Filename)
		filePath := storage.GetBannerPath(web)
		if err := storage.CheckAndSaveUploadFile(allowedExtensions, req.WebImage, filePath, os.O_EXCL); errpkg.Handler(ctx, err) {
			return
		}
		saved = append(saved, filePath)
	}
	if req.AppImage != nil {
		app = genFilenameUUID(req.AppImage.Filename)
		filePath := storage.GetBannerPath(app)
		if err := storage.CheckAndSaveUploadFile(allowedExtensions, req.AppImage, filePath, os.O_EXCL); errpkg.Handler(ctx, err) {
			return
		}
		saved = append(saved, filePath)
	}

	if err := banners.Update(req, web, app); errpkg.Handler(ctx, err) {
		return
	}

	isCompleted = true
	ctx.Status(http.StatusOK)
}

// @Summary 	更新橫幅資料優先度
// @Description 更新橫幅資料優先度
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body banners.PriorityUpdateRequset true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/banners/priority [post]
func PriorityHandler(ctx *gin.Context) {
	var req banners.PriorityUpdateRequset
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := banners.PriorityUpdate(req); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	刪除橫幅資料
// @Description 刪除橫幅資料
// @Tags 		Admin-settings
// @Security	BearerAuth
// @Param 		id path int true "banners id"
// @Success 	204
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/settings/banners/{id} [delete]
func DeleteHandler(ctx *gin.Context) {
	var req banners.DeleteRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	if m, err := banners.Delete(req.ID); errpkg.Handler(ctx, err) {
		return
	} else {
		if m.WebImage != "" {
			filePath := storage.GetBannerPath(m.WebImage)
			if err := os.Remove(filePath); err != nil {
				logger.Logger.Err(err).Str("filename", filePath).Msg("remove file error")
			}
		}
		if m.AppImage != "" {
			filePath := storage.GetBannerPath(m.AppImage)
			if err := os.Remove(filePath); err != nil {
				logger.Logger.Err(err).Str("filename", filePath).Msg("remove file error")
			}
		}
	}
	ctx.Status(http.StatusNoContent)
}

func genFilenameUUID(filename string) string {
	ext := filepath.Ext(filename)
	return uuid.NewString() + ext
}
