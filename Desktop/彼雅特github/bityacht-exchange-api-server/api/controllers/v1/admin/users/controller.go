package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	sqlcache "bityacht-exchange-api-server/internal/cache/memory/sql"
	redisusers "bityacht-exchange-api-server/internal/cache/redis/users"
	"bityacht-exchange-api-server/internal/database/sql"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/database/sql/userswallets"
	"bityacht-exchange-api-server/internal/database/sql/userswithdrawalwhitelist"
	"bityacht-exchange-api-server/internal/pkg/csv"
	"bityacht-exchange-api-server/internal/pkg/email"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	passwordpkg "bityacht-exchange-api-server/internal/pkg/password"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bityacht-exchange-api-server/internal/pkg/wallet"

	dbsql "database/sql"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// @Summary 	建立法人帳號
// @Description 建立法人帳號
// @Tags 		Admin-users
// @Security	BearerAuth
// @Accept 		json
// @Param 		body body CreateRequest true "Request Body"
// @Success 	201
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users [post]
func CreateHandler(ctx *gin.Context) {
	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	record := &users.Model{
		Account:                              req.Account,
		FirstName:                            req.Name,
		JuridicalPersonNature:                req.JuridicalPersonNature,
		Address:                              req.Address,
		BirthDate:                            req.BirthDate.ToDBTime(),
		Phone:                                req.Phone,
		FundsSources:                         req.JuridicalPersonFundsSources,
		JuridicalPersonCryptocurrencySources: req.JuridicalPersonCryptocurrencySources,
		AuthorizedPersonName:                 req.AuthorizedPersonName,
		AuthorizedPersonNationalID:           req.AuthorizedPersonNationalID,
		AuthorizedPersonPhone:                req.AuthorizedPersonPhone,
		Type:                                 users.TypeJuridicalPerson,
		Level:                                1,
		Status:                               usersmodifylogs.SLStatusEnable,
		Comment:                              req.Comment,
	}
	if req.NationalID != "" {
		record.NationalID = dbsql.NullString{
			String: req.NationalID,
			Valid:  true,
		}
	}
	if req.Country != "" {
		record.CountriesCode = dbsql.NullString{
			String: req.Country,
			Valid:  true,
		}
	}
	if req.IndustrialClassificationsID != 0 {
		record.IndustrialClassificationsID = dbsql.NullInt64{
			Int64: req.IndustrialClassificationsID,
			Valid: true,
		}
	}

	var (
		err         *errpkg.Error
		rawPassword = rand.LetterAndNumberString(8)
	)
	if record.Password, err = passwordpkg.Encrypt(rawPassword); errpkg.Handler(ctx, err) {
		return
	} else if err := sql.DB().Create(&record).Error; errpkg.HandlerWithCode(ctx, http.StatusInternalServerError, errpkg.CodeSql, err) {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			errpkg.HandlerWithCode(ctx, http.StatusConflict, errpkg.CodeAccountDuplicated, err)
			return
		}
		return
	}

	passwordMail := email.NewEmail()
	passwordMail.To = []string{record.Account}
	passwordMail.Subject = "BitYacht 兑幣所法人用戶密碼函"
	passwordMail.Text = []byte(fmt.Sprintf("法人用戶 %s 您好：%s 為您的預設密碼，請於首次登入後變更密碼。", record.FirstName, rawPassword))

	if err = email.SendMail(passwordMail); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusCreated)
}

// @Summary 	取得用戶列表
// @Description 取得用戶列表
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Param 		searcher query modelpkg.Searcher false "searcher"
// @Param		query query users.GetListRequest false "query"
// @Success 	200 {object} modelpkg.GetResponse{data=[]users.User}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users [get]
func GetListHandler(ctx *gin.Context) {
	var req users.GetListRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	resp := modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
	searcher := modelpkg.GetSearcherFromQuery(ctx)

	var err *errpkg.Error
	if resp.Data, err = users.GetList(req, &resp.Paginator, &searcher); errpkg.Handler(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	取得使用者相關選項
// @Description 取得使用者相關選項
// @Tags 		Admin-users
// @Security	BearerAuth
// @Produce		json
// @Success 	200 {object} sqlcache.IDVOptionsResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/options [get]
func GetOptionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, sqlcache.GetUserIDVOptionsResponse())
}

// @Summary 	匯出用戶列表
// @Description 匯出用戶列表
// @Tags 		Admin-users
// @Security	BearerAuth
// @Accept 		json
// @Param 		query query ExportRequest true "query"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/export [get]
func ExportHandler(ctx *gin.Context) {
	var req users.ExportRequest
	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	} else if err := modelpkg.ValidStartAndEndAt(req.StartAt.Time, req.EndAt.Time); errpkg.Handler(ctx, err) {
		return
	}

	data, err := users.GetExport(req)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("members_%s.csv", time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, users.GetUserCSVHeaders(), data)
}

// @Summary 	取得用戶提幣白名單
// @Description 取得用戶提幣白名單
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param		id path int true "user id"
// @Param		paginator query modelpkg.Paginator false "paginator"
// @Success 	200 {object} modelpkg.GetResponse{data=[]userswithdrawalwhitelist.Record}
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/withdrawal-whitelist [get]
func GetWithdrawalWhitelistHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	var (
		resp = modelpkg.GetResponse{Paginator: modelpkg.GetPaginatorFromQuery(ctx)}
		err  *errpkg.Error
	)

	if resp.Data, err = userswithdrawalwhitelist.GetRecordsByUser(req.ID, &resp.Paginator); errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	匯出用戶提幣白名單
// @Description 匯出用戶提幣白名單
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param		id path int true "user id"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/withdrawal-whitelist/export [get]
func ExportWithdrawalWhitelistHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	data, err := userswithdrawalwhitelist.GetExportByUser(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	filename := fmt.Sprintf("withdrawal-whitelist-%d_%s.csv", req.ID, time.Now().Format("20060102"))
	csv.ExportCSVFile(ctx, filename, userswithdrawalwhitelist.GetCSVHeaders(), data)
}

// @Summary 	刪除用戶提幣白名單
// @Description 刪除用戶提幣白名單
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param		id path int true "user id"
// @Param		whitelistID path int true "withdrawal whitelist id"
// @Success 	204
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/withdrawal-whitelist/{whitelistID} [delete]
func DeleteWithdrawalWhitelistHandler(ctx *gin.Context) {
	var req DeleteWithdrawalWhitelistRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	record, err := userswithdrawalwhitelist.Delete(req.ID, req.UsersID)
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusNoContent)

	if record != nil {
		errLogger := logger.GetGinRequestLogger(ctx)

		m, ok := wallet.ParseMainnet(record.Mainnet)
		if !ok {
			errLogger.Error().Any("record", *record).Msg("parse mainnet failed")
			return
		}

		if err := wallet.Cybavo.RemoveWithdrawalWhitelistEntry(context.Background(), req.UsersID, m, record.Address); err != nil {
			errLogger.Err(err.Err).Any("record", *record).Msg("remove withdrawal whitelist from cybavo failed")
		}
	}
}

// @Summary 	更新用戶資料等級
// @Description 更新用戶資料等級
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Param 		body body UpdateLevelRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/level [patch]
func UpdateLevelHandler(ctx *gin.Context) {
	var req UpdateLevelRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	} else if err := users.UpdateLevel(req.ID, req.Level); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	更新用戶資料狀態
// @Description 更新用戶資料狀態
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Accept 		json
// @Param 		body body UpdateStatusRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/status [patch]
func UpdateStatusHandler(ctx *gin.Context) {
	var req UpdateStatusRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	} else if err := users.UpdateStatus(claims.ManagerPayload.ID, req.ID, req.Status, req.Comment); errpkg.Handler(ctx, err) {
		return
	}

	if err := redisusers.ForceLogout(ctx, req.ID); err != nil {
		errLogger := logger.GetGinRequestLogger(ctx)
		errLogger.Err(err).Msg("force logout failed")
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	用戶資料帶資產列表
// @Description 用戶資料帶資產列表
// @Tags 		Admin-users
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Success 	200 {object} users.UserWithAsset
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id} [get]
func GetHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	user, err := users.GetUserByID(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	assets, err := userswallets.GetUserAsset(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	resp := &users.UserWithAsset{
		User:          user,
		MobileBarcode: user.Extra.MobileBarcode,
		Assets:        assets,
	}

	ctx.JSON(http.StatusOK, resp)
}
