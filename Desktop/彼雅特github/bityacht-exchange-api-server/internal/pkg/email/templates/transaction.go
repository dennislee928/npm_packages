package emailtemplates

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

// #nosec G101
const (
	transactionFailedTitle    = "[BitYacht] 交易失敗通知信"
	transactionFailedTemplate = "tx-failed.html.tmpl"

	transactionSuccessfulTitle    = "[BitYacht] 交易成功通知信"
	transactionSuccessfulTemplate = "tx-successful.html.tmpl"
)

type TransactionPayload struct {
	Time           string
	TransactionsID string
	BaseSymbol     string
	QuoteSymbol    string
	Side           string
	Quantity       string
}

func ExecTransactionFailed(payload TransactionPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(transactionFailedTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return transactionFailedTitle, content, nil
}

func ExecTransactionSuccessful(payload TransactionPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(transactionSuccessfulTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return transactionSuccessfulTitle, content, nil
}
