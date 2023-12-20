package kryptogo

import (
	"bityacht-exchange-api-server/internal/database/sql/duediligences"
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/pkg/email"
	emailtemplates "bityacht-exchange-api-server/internal/pkg/email/templates"
	"bityacht-exchange-api-server/internal/pkg/kyc"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func readReq(ctx *gin.Context, errorLogger zerolog.Logger, uriReq any, compactedRawReq *bytes.Buffer, req any) bool {
	if err := ctx.ShouldBindUri(uriReq); err != nil {
		errorLogger.Err(err).Msg("ShouldBindUri error")
		ctx.Status(http.StatusBadRequest)
		return true
	}

	if compactedRawReq == nil {
		if err := ctx.ShouldBindJSON(req); err != nil {
			errorLogger.Err(err).Msg("ShouldBindJSON error")
			ctx.Status(http.StatusBadRequest)
			return true
		}

		return false
	}

	// Read Raw Resp
	if rawReq, err := io.ReadAll(ctx.Request.Body); err != nil {
		errorLogger.Err(err).Msg("ReadAll error")
		ctx.Status(http.StatusBadRequest)
		return true
	} else if err = json.Compact(compactedRawReq, rawReq); err != nil {
		errorLogger.Err(err).Msg("Compact json error")
		ctx.Status(http.StatusBadRequest)
		return true
	} else if err = json.Unmarshal(compactedRawReq.Bytes(), &req); err != nil {
		errorLogger.Err(err).Msg("Unmarshal error")
		ctx.Status(http.StatusBadRequest)
		return true
	}

	return false
}

func IDVCallbackHandler(ctx *gin.Context) {
	var (
		req             kyc.IDVCallbackRequest
		compactedRawReq = new(bytes.Buffer)
		errorLogger     = logger.GetGinRequestLogger(ctx)

		uriReq struct {
			UsersID int64 `uri:"UsersID" binding:"gt=0"`
			IDVsID  int64 `uri:"IDVsID" binding:"gt=0"`
		}
	)

	if readReq(ctx, errorLogger, &uriReq, compactedRawReq, &req) {
		return
	} else if req.State == kyc.IDVStateInitial && len(req.RejectReasons) > 0 { // Special case, mean the idv is not created.
		req.State = kyc.IDVStateReject
	}

	record := idvCallbackReqToModel(req)
	record.ID = uriReq.IDVsID
	record.UsersID = uriReq.UsersID
	record.Detail = json.RawMessage(compactedRawReq.Bytes())

	//! Deprecated (Message from TG 2023/11/28)
	// if req.AuditStatus == kyc.AuditStatusAccepted {
	// 	// Update FirstName, LastName, NationalID, Gender, BirthDate to users table
	// 	updateArgs, err := idvCallbackReqToUpdateArgs(req)
	// 	if err != nil {
	// 		errLogger := logger.GetGinRequestLogger(ctx)
	// 		errLogger.Err(err.Err).Any("req", req).Msg("idvCallbackReqToUpdateArgs error")
	// 		//! No need to return, Continue the flow.
	// 	}

	// 	if err := users.UpdateFromIDVCallback(uriReq.UsersID, uriReq.IDVsID, updateArgs); err != nil {
	// 		errLogger := logger.GetGinRequestLogger(ctx)
	// 		errLogger.Err(err.Err).Any("args", updateArgs).Any("errorCode", err.Code).Msg("UpdateFromIDVCallback error")
	// 		//! No need to return, Continue the flow.
	// 	}
	// }

	needSendNotify, err := idverifications.UpdateFromCallback(ctx, idverifications.TypeKryptoGO, record, idverifications.UpdateImagesByURLRequest{
		IDImage:        req.IDImageURL,
		IDBackImage:    req.IDBackImageURL,
		FaceImage:      req.FaceImageURL,
		IDAndFaceImage: req.IDAndFaceImageURL,
	})
	if err != nil {
		errorLogger.Err(err.Err).Int64("errorCode", int64(err.Code)).Msg("UpdateFromCallback error")
		ctx.Status(err.HttpStatus)
		return
	}

	ctx.Status(http.StatusOK)

	if needSendNotify {
		userRecord, err := users.GetByID(uriReq.UsersID)
		if err != nil {
			errorLogger.Err(err.Err).Int64("errorCode", int64(err.Code)).Msg("[users]GetByID error")
		} else if userRecord.IDVerificationsID.Int64 != record.ID {
			return
		}

		idvMail := email.NewEmail(email.WithLogo())
		idvMail.To = []string{userRecord.Account}

		switch req.State {
		case kyc.IDVStateReject:
			idvMail.Subject, idvMail.HTML, err = emailtemplates.ExecIDVFailed()
		default:
			idvMail.Subject, idvMail.HTML, err = emailtemplates.ExecIDVOnGoing(emailtemplates.IDVOnGoingPayload{Time: emailtemplates.FormatTime(time.Now())})
		}

		if err != nil {
			errorLogger.Err(err.Err).Int64("errorCode", int64(err.Code)).Msg("ExecuteIDVOnGoingTemplate error")
		} else if err = email.SendMail(idvMail); err != nil {
			errorLogger.Err(err.Err).Int64("errorCode", int64(err.Code)).Msg("SendMail error")
		}
	}
}

func DDCallbackHandler(ctx *gin.Context) {
	var (
		req             kyc.DDCallbackRequest
		compactedRawReq = new(bytes.Buffer)
		errorLogger     = logger.GetGinRequestLogger(ctx)

		uriReq struct {
			UsersID int64 `uri:"UsersID" binding:"gt=0"`
			DDsID   int64 `uri:"DDsID" binding:"gt=0"`
		}
	)

	if readReq(ctx, errorLogger, &uriReq, compactedRawReq, &req) {
		return
	}

	record, err := duediligences.NewFromKryptoGOTaskSummary(req)
	if err != nil {
		errorLogger.Err(err).Str("rawReq", compactedRawReq.String()).Send()
	}

	record.ID = uriReq.DDsID
	record.UsersID = uriReq.UsersID
	record.Detail = json.RawMessage(compactedRawReq.Bytes())

	if err := duediligences.UpdateFromCallback(record); err != nil {
		errorLogger.Err(err.Err).Int64("errorCode", int64(err.Code)).Msg("UpdateFromCallback error")
		ctx.Status(err.HttpStatus)
		return
	}

	ctx.Status(http.StatusOK)
}
