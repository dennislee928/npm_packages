package commissions

import (
	"bityacht-exchange-api-server/internal/database/sql/userscommissions"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得返佣與提領紀錄
// @Description 取得返佣與提領紀錄
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Produce		json
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} GetCommissionsResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/commissions [get]
func GetHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	resp := GetCommissionsResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	if resp.TotalReward, resp.NotWithdrew, err = userscommissions.GetRewardByUser(claims.ID()); errpkg.Handler(ctx, err) {
		return
	} else if resp.Data, err = userscommissions.GetCommissionsByUser(claims.ID(), 0, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	提領返佣
// @Description 提領返佣
// @Tags 		User-MemberCenter
// @Security	BearerAuth
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/commissions/withdraw [post]
func WithdrawHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err = userscommissions.Withdraw(claims.ID()); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusCreated)
}
