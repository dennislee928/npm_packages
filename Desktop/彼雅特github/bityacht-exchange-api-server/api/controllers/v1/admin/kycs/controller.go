package kycs

import (
	"fmt"
	"net/http"
	"time"

	"bityacht-exchange-api-server/internal/database/sql/duediligences"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得審查列表
// @Description 取得審查列表
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		searcher query modelpkg.Searcher false "searcher"
// @Param		query query duediligences.GetWithDDListRequest false "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]duediligences.Review}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs [get]
func GetListHandler(ctx *gin.Context) {
	var req duediligences.GetWithDDListRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	searcher := modelpkg.GetSearcherFromQuery(ctx)

	var err *errpkg.Error
	if resp.Data, err = duediligences.GetWithDDList(&resp.Paginator, req, searcher); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出審查列表
// @Description 匯出審查列表
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Accept 		json
// @Param 		query query duediligences.ExportWithDDRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs/export [get]
func ExportHandler(ctx *gin.Context) {
	var req duediligences.ExportWithDDRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := duediligences.GetWithDDExport(req)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("kyc_%s.csv", time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, duediligences.GetReviewCSVHeaders(), data)
}

// @Summary 	取得年度審查列表
// @Description 取得年度審查列表
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		searcher query modelpkg.Searcher false "searcher"
// @Param		query query duediligences.GetAnnualWithDDListRequest false "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]duediligences.Review}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs/annual [get]
func GetAnnualListHandler(ctx *gin.Context) {
	var req duediligences.GetAnnualWithDDListRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	searcher := modelpkg.GetSearcherFromQuery(ctx)

	var err *errpkg.Error
	if resp.Data, err = duediligences.GetAnnualWithDDList(&resp.Paginator, req, searcher); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出年度審查列表
// @Description 匯出年度審查列表
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Accept 		json
// @Param 		query query duediligences.ExportAnnualWithDDRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/kycs/annual/export [get]
func ExportAnnualHandler(ctx *gin.Context) {
	var req duediligences.ExportAnnualWithDDRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := duediligences.GetAnnualWithDDExport(req)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("kyc_%s.csv", time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, duediligences.GetReviewCSVHeaders(), data)
}
