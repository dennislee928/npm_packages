package assets

import (
	"bityacht-exchange-api-server/internal/pkg/storage"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得 Banner 圖片
// @Description 取得 Banner 圖片
// @Tags 		Assets
// @Param 		filename path string true "banners filename"
// @Success 	200
// @Failure 	400
// @Failure 	500
// @Router 		/assets/banners/{filename} [get]
func GetBannerHandler(ctx *gin.Context) {
	ctx.File(storage.GetBannerPath(ctx.Param("Filename")))
}
