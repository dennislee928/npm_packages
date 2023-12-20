package userskycs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"bityacht-exchange-api-server/internal/database/sql/duediligences"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/kyc"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bityacht-exchange-api-server/internal/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// Name Check PDF's allowed extensions
var allowedExtensions = map[string]struct{}{
	"application/pdf": {},
}

// @Summary 	取得使用者帶審查資料
// @Description 取得使用者帶審查資料
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "users id"
// @Success 	200 {object} duediligences.UserWithDD
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs [get]
func GetWithDDHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	}

	resp, err := duediligences.GetWithDD(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary 	更新認證結果(外籍人士)
// @Description 更新認證結果(外籍人士)
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Accept		json
// @Param 		id path int true "users id"
// @Param 		body body UpdateIDVAuditStatusRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/idv-audit-status [patch]
func UpdateIDVAuditStatusHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	var req UpdateIDVAuditStatusRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	if err := idverifications.UpdateAuditStatusByUser(claims.ID(), req.ID, req.AuditStatus, req.Comment); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	編輯 KryptoGO 單號
// @Description 編輯 KryptoGO 單號
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "users id"
// @Param 		body body UpdateKryptoGOTaskIDRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/task-id [patch]
func UpdateKryptoGOTaskIDHandler(ctx *gin.Context) {
	var req UpdateKryptoGOTaskIDRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBindBody, err) {
		return
	}

	rawResp := new(bytes.Buffer)
	taskSummary, err := kyc.KryptoGO.GetTaskSummary(ctx, req.TaskID, rawResp)
	if errpkg.Handler(ctx, err) {
		return
	}

	newDDRecord, err := duediligences.NewFromKryptoGOTaskSummary(taskSummary)
	if errpkg.Handler(ctx, err) {
		return
	}

	newDDRecord.UsersID = req.ID
	newDDRecord.Type = duediligences.TypeManualSet
	newDDRecord.Detail = json.RawMessage(rawResp.Bytes())

	if err := duediligences.CreateForUser(&newDDRecord); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	重送用戶KryptoGo審查
// @Description 重送用戶KryptoGo審查
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/resent-krypto [post]
func ResentToKryptoGoHandler(ctx *gin.Context) {
	var req users.IDURIRequest

	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, err) {
		return
	} else if createTasksParams, ddRecord, err := duediligences.ResentKryptoGo(ctx, req.ID); errpkg.Handler(ctx, err) {
		return
	} else if resp, err := kyc.KryptoGO.CreateTasks(ctx, []kyc.CreateTasksParams{createTasksParams}); errpkg.Handler(ctx, err) {
		return
	} else if err := duediligences.UpdateTaskID(ddRecord.ID, strconv.FormatInt(resp[0].TaskID, 10)); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	更新用戶KryptoGo審查結果
// @Description 更新用戶KryptoGo審查結果
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Param 		body body UpdateKryptoReviewRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/krypto-review [patch]
func UpdateKryptoReviewHandler(ctx *gin.Context) {
	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	var req UpdateKryptoReviewRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	ddRecord, err := duediligences.GetLatestByUsersID(req.ID)
	if errpkg.Handler(ctx, err) {
		return
	} else if ddRecord.TaskID == "" {
		errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, errors.New("tasks id is empty"))
		return
	} else if ddRecord.AuditAccepted != duediligences.BoolUnknown {
		errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadParam, errors.New("krypto audit accepted is updated, can not update again"))
		return
	}

	var (
		managersID *int64
		accepted   bool
		comment    string
	)

	switch ddRecord.Type {
	case duediligences.TypeCreateByIDV, duediligences.TypeManualResend:
		mID := claims.ID()
		managersID = &mID

		resp, err := kyc.KryptoGO.UpdateTaskStatus(ctx, ddRecord.TaskID, req.Result == usersmodifylogs.RLStatusApproved, req.Comment)
		if errpkg.Handler(ctx, err) {
			return
		}

		accepted = resp.Accepted
		comment = resp.Comment
	case duediligences.TypeManualSet:
		updateResp, updateErr := kyc.KryptoGO.UpdateTaskStatus(ctx, ddRecord.TaskID, req.Result == usersmodifylogs.RLStatusApproved, req.Comment)
		if updateErr == nil {
			mID := claims.ID()
			managersID = &mID

			accepted = updateResp.Accepted
			comment = updateResp.Comment
			break
		}

		taskSummary, err := kyc.KryptoGO.GetTaskSummary(ctx, ddRecord.TaskID, nil)
		if errpkg.Handler(ctx, err) {
			return
		}

		newDDRecord, err := duediligences.NewFromKryptoGOTaskSummary(taskSummary)
		if errpkg.Handler(ctx, err) {
			return
		}

		if newDDRecord.AuditAccepted == duediligences.BoolUnknown {
			errpkg.Handler(ctx, updateErr)
			return
		}

		accepted = newDDRecord.AuditAccepted == duediligences.BoolTrue
		comment = newDDRecord.Comment
	}

	if err := duediligences.UpdateAuditInfo(managersID, &ddRecord, accepted, comment); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	更新用戶姓名檢核排除評估
// @Description 更新用戶姓名檢核排除評估
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Accept		mpfd
// @Param 		id path int true "user id"
// @Param       file formData file false "name check pdf"
// @Param 		body formData UpdateNameCheckRequest true "body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/name-check [patch]
func UpdateNameCheckHandler(ctx *gin.Context) {
	var req UpdateNameCheckRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	var (
		allCreated bool
		saved      string
	)
	defer func() {
		if allCreated {
			return
		}
		if saved != "" {
			if err := os.Remove(saved); err != nil {
				logger.Logger.Err(err).Msg("remove saved file error")
			}
		}
	}()

	var uid string
	if req.File != nil {
		uid = fmt.Sprintf("%d.pdf", req.ID)
		path := storage.GetNameCheckPdfPath(uid)
		if err := storage.CheckAndSaveUploadFile(allowedExtensions, req.File, path, os.O_TRUNC); errpkg.Handler(ctx, err) {
			return
		}
		saved = path
	}

	if _, err := users.UpdateReview(users.UpdateReviewParams{
		ReviewType: usersmodifylogs.RLTypeNameCheckUploadPDF,
		ManagerID:  claims.ManagerPayload.ID,
		UserID:     req.ID,
		Result:     req.Result,
		Notice:     req.File.Filename,
		Comment:    uid,
	}); errpkg.Handler(ctx, err) {
		return
	}

	allCreated = true
	ctx.Status(http.StatusOK)
}

// @Summary 	取得用戶姓名檢核PDF檔案
// @Description 取得用戶姓名檢核PDF檔案
// @Tags 		Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/name-check/pdf [get]
func GetNameCheckPdfHandler(ctx *gin.Context) {
	var req users.IDURIRequest
	if err := ctx.ShouldBindUri(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	ctx.File(storage.GetNameCheckPdfPath(fmt.Sprintf("%d.pdf", req.ID)))
}

// @Summary 	更新用戶法遵審查結果
// @Description 更新用戶法遵審查結果
// @Tags 	    Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Param 		body body UpdateComplianceReviewRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/compliance-review [patch]
func UpdateComplianceReviewHandler(ctx *gin.Context) {
	var req UpdateComplianceReviewRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	if _, err := users.UpdateReview(users.UpdateReviewParams{
		ReviewType: usersmodifylogs.RLTypeComplianceReview,
		ManagerID:  claims.ManagerPayload.ID,
		UserID:     req.ID,
		Result:     req.Result,
		Notice:     "",
		Comment:    req.Comment,
	}); errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)
}

// @Summary 	更新用戶最終審查結果
// @Description 更新用戶最終審查結果
// @Tags 	    Admin-kycs
// @Security	BearerAuth
// @Param 		id path int true "user id"
// @Param 		body body UpdateFinalReviewRequest true "Request Body"
// @Success 	200
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/admin/users/{id}/kycs/final-review [patch]
func UpdateFinalReviewHandler(ctx *gin.Context) {
	var req UpdateFinalReviewRequest

	ctx.ShouldBindUri(&req)
	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadBody, err) {
		return
	}

	claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
	if errpkg.Handler(ctx, err) {
		return
	}

	userRecord, err := users.UpdateReview(users.UpdateReviewParams{
		ReviewType: usersmodifylogs.RLTypeFinalReview,
		ManagerID:  claims.ManagerPayload.ID,
		UserID:     req.ID,
		Result:     req.Result,
		Comment:    req.Comment,
	})
	if errpkg.Handler(ctx, err) {
		return
	}

	ctx.Status(http.StatusOK)

	resultMail := email.NewEmail(email.WithLogo())
	resultMail.To = []string{userRecord.Account}

	switch req.Result {
	case usersmodifylogs.RLStatusApproved:
		resultMail.Subject, resultMail.HTML, err = emailtemplates.ExecIDVPass()
	case usersmodifylogs.RLStatusRejected:
		resultMail.Subject, resultMail.HTML, err = emailtemplates.ExecIDVRisk()
	case usersmodifylogs.RLStatusIDVRejected:
		resultMail.Subject, resultMail.HTML, err = emailtemplates.ExecIDVFailed()
	default:
		logger.Logger.Error().Any("result", req.Result).Msg("bad result of final review")
	}

	if err != nil {
		logger.Logger.Err(err).Msg("failed to execute template")
	} else if err = email.SendMail(resultMail); err != nil {
		logger.Logger.Err(err).Msg("failed to send result mail")
	}
}
