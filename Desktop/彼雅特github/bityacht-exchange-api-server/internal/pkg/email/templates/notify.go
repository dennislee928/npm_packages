package emailtemplates

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

const (
	notifyDepositFiatTitle    = "[BitYacht] 入金通知信"
	notifyDepositFiatTemplate = "deposit-fiat.html.tmpl"

	notifyDepositSpotTitle    = "[BitYacht] 入幣通知信"
	notifyDepositSpotTemplate = "deposit-spot.html.tmpl"

	notifyInvoiceTitle    = "[BitYacht] 電子發票開立通知"
	notifyInvoiceTemplate = "invoice.html.tmpl"

	notifyLoginTitle    = "[BitYacht] 新設備登入提示"
	notifyLoginTemplate = "login.html.tmpl"

	notifyRegisteredTitle    = "[BitYacht] 歡迎加入 BitYacht 交易所！"
	notifyRegisteredTemplate = "registered.html.tmpl"

	notifyTransactionFailedTitle    = "[BitYacht] 交易失敗通知信"
	notifyTransactionFailedTemplate = "transaction-failed.html.tmpl"

	notifyWithdrawFiatTitle    = "[BitYacht] 出金通知信"
	notifyWithdrawFiatTemplate = "withdraw-fiat.html.tmpl"

	notifyWithdrawSpotTitle    = "[BitYacht] 提幣通知信"
	notifyWithdrawSpotTemplate = "withdraw-spot.html.tmpl"

	notifyWithdrawSpotFailedTitle    = "[BitYacht] 提幣失敗通知信"
	notifyWithdrawSpotFailedTemplate = "withdraw-spot-failed.html.tmpl"
)

type NotifyDepositFiatPayload struct {
	Time  string
	Value string
}

func ExecNotifyDepositFiat(payload NotifyDepositFiatPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyDepositFiatTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyDepositFiatTitle, content, nil
}

type NotifyDepositSpotPayload struct {
	Time             string
	CurrenciesSymbol string
	Mainnet          string
	Amount           string
}

func ExecNotifyDepositSpot(payload NotifyDepositSpotPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyDepositSpotTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyDepositSpotTitle, content, nil
}

type NotifyInvoicePayload struct {
	OrderID        string
	YearMonth      string
	InvNo          string
	InvoiceTime    string
	RandomNumber   string
	Total          string
	Seller         string
	HandlingCharge string
	SalesAmount    string
}

func ExecNotifyInvoice(payload NotifyInvoicePayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyInvoiceTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyInvoiceTitle, content, nil
}

type NotifyLoginPayload struct {
	Account  string
	Device   string
	IP       string
	Browser  string
	Location string
	Time     string
}

func ExecNotifyLogin(payload NotifyLoginPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyLoginTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyLoginTitle, content, nil
}

type NotifyRegisteredPayload struct {
	Time string
	URL  string
}

func ExecNotifyRegistered(payload NotifyRegisteredPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyRegisteredTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyRegisteredTitle, content, nil
}

type NotifyTransactionFailedPayload struct {
	OrderID              string
	TransactionPair      string
	TransactionDirection string
	Amount               string
}

func ExecNotifyTransactionFailed(payload NotifyTransactionFailedPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyTransactionFailedTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyTransactionFailedTitle, content, nil
}

type NotifyWithdrawFiatPayload struct {
	Time        string
	BankAccount string
	BankName    string
	Name        string
	Amount      string
}

func ExecNotifyWithdrawFiat(payload NotifyWithdrawFiatPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyWithdrawFiatTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyWithdrawFiatTitle, content, nil
}

type NotifyWithdrawSpotPayload struct {
	Time             string
	CurrenciesSymbol string
	Mainnet          string
	Address          string
	Amount           string
}

func ExecNotifyWithdrawSpot(payload NotifyWithdrawSpotPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyWithdrawSpotTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyWithdrawSpotTitle, content, nil
}

type NotifyWithdrawSpotFailedPayload struct {
	Time             string
	CurrenciesSymbol string
	Mainnet          string
	Amount           string
}

func ExecNotifyWithdrawSpotFailed(payload NotifyWithdrawSpotFailedPayload) (string, []byte, *errpkg.Error) {
	content, err := executeTemplate(notifyWithdrawSpotFailedTemplate, payload)
	if err != nil {
		return "", nil, err
	}

	return notifyWithdrawSpotFailedTitle, content, nil
}
