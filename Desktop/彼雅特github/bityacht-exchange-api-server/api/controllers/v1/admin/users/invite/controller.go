package invite

import (
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/userscommissions"
	"bityacht-exchange-api-server/internal/pkg/csv"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得邀請獎勵資訊
// @Description 取得邀請獎勵資訊
// @Tags 		Admin-users
// @Security	BearerAuth
// @Produce		json
// @Param 		id path int true "users id"
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} GetInfoResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/invite-info [get]
func GetInfoHandler(ctx *gin.Context) {
	var req GetRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp = GetInfoResponse{InviteCode: users.GetInviteCodeByID(req.ID)}
		err  *errpkg.Error
	)

	if resp.TotalInvited, resp.TotalSucceed, err = users.GetInviteCount(req.ID); errpkg.Handler(ctx, err) {
		return
	} else if resp.TotalReward, resp.NotWithdrew, err = userscommissions.GetRewardByUser(req.ID); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	取得邀請清單
// @Description 取得邀請清單
// @Tags 		Admin-users
// @Security	BearerAuth
// @Produce		json
// @Param 		id path int true "users id"
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=users.Invitee}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/invite [get]
func GetHandler(ctx *gin.Context) {
	var req GetRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
		err  *errpkg.Error
	)
	if resp.Data, err = users.GetInviteeList(req.ID, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	取得獎勵清單
// @Description 取得獎勵清單
// @Tags 		Admin-users
// @Security	BearerAuth
// @Produce		json
// @Param 		id path int true "users id"
// @Param 		query query GetRewardsRequest true "query"
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]userscommissions.Commission}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/invite-rewards [get]
func GetRewardsHandler(ctx *gin.Context) {
	var req GetRewardsRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
		err  *errpkg.Error
	)

	if resp.Data, err = userscommissions.GetCommissionsByUser(req.ID, req.Action, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出獎勵清單
// @Description 匯出獎勵清單
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param 		id path int true "users id"
// @Param 		query query GetExportRewardsHandler true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/invite-rewards/export [get]
func ExportRewardsHandler(ctx *gin.Context) {
	var req GetExportRewardsHandler

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	data, err := userscommissions.ExportCommissionsByUser(req.ID, req.Action)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("reward_%d.csv", req.ID)
	csv.ExportCSVFile(ctx, filename, userscommissions.GetCommissionCSVHeaders(), data)
}
