package service

import (
	"bityacht-exchange-api-server/internal/database/sql/receipts"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bityacht-exchange-api-server/internal/pkg/receipt"
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"

	"github.com/sourcegraph/conc/iter"
)

// RestoreIssuingReceipts is a service to restore issuing receipts.
// This service should be called when app start.
func RestoreIssuingReceipts(ctx context.Context) {
	logger := logger.Logger.With().Str("service", "restore-issuing-receipts").Logger()

	orderIDs, wrapErr := receipts.GetIssuingReceiptOrderIDs()
	if wrapErr != nil {
		logger.Err(wrapErr.Err).Msg("get issuing receipt order ids failed")
		return
	}

	iter.ForEach(orderIDs, func(id *string) {
		if wrapErr := issue(ctx, *id); wrapErr != nil {
			if wrapErr.Err == nil {
				wrapErr.Err = errors.New("unknown error")
			}
			logger.Err(wrapErr.Err).Str("id", *id).Msg("issue receipt failed")
			return
		}
		logger.Info().Str("id", *id).Msg("issue receipt success")
	})
}

// IssueReceipt is a service to issue receipt.
func IssueReceipt(orderID string) *errpkg.Error {
	if wrapErr := receipts.SetReciptStatusByID(orderID, receipts.StatusIssuing); wrapErr != nil {
		return wrapErr
	}

	return issue(context.Background(), orderID)
}

func newChargeProd(sales int64) receipt.ProdItem {
	return receipt.ProdItem{
		Title:   "平台手續費",
		IncTax:  true,
		Sales:   float64(sales),
		Qty:     1,
		TaxType: 1,
	}
}

func issue(ctx context.Context, orderID string) *errpkg.Error {
	failed := true
	defer func() {
		if failed {
			if err := receipts.SetReciptStatusByID(orderID, receipts.StatusFailed); err != nil {
				logger.Logger.Err(err).Str("id", orderID).Msg("rollback receipt status failed")
			}
		}
	}()

	r, wrapErr := receipts.GetByID(orderID)
	if wrapErr != nil {
		return wrapErr
	}

	userID := r.UserID
	sales := r.InvoiceAmount
	barcode := r.Barcode

	invoice, wrapErr := receipt.EZ.FastCreateB2CInvoice(ctx, receipt.FastCreateB2CPayload{
		ID:         orderID,
		Title:      strconv.FormatInt(userID, 10),
		ProdItem:   newChargeProd(sales),
		CarrierNum: barcode,
	})
	if wrapErr != nil {
		return wrapErr
	}

	invoiceTime, err := time.ParseInLocation(time.DateTime, invoice.InvoiceTime, modelpkg.DefaultTimeLoc)
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadParam, Err: err}
	}

	if wrapErr := receipts.IssueReceiptByID(orderID, receipts.IssueReceiptParam{
		InvoiceID:   invoice.InvNo,
		CreatedAt:   invoiceTime,
		SalesAmount: int64(invoice.SalesAmount),
		Tax:         int64(invoice.TaxAmount),
	}); wrapErr != nil {
		return wrapErr
	}

	failed = false

	return nil
}
