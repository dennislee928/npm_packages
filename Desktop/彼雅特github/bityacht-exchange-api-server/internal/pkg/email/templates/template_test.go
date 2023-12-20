package emailtemplates

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"testing"
)

func Test_executeTemplate(t *testing.T) {
	type args struct {
		name string
		data any
	}
	tests := []struct {
		name        string
		args        args
		wantErrCode errpkg.Code
	}{
		// Account Change
		{"account change disable", args{accountChangeDisableTemplate, nil}, 0},
		{"account change forzen", args{accountChangeForzenTemplate, nil}, 0},
		// IDV
		{"idv failed", args{idvFailedTemplate, nil}, 0},
		{"idv on going", args{idvOnGoingTemplate, IDVOnGoingPayload{Time: "2023-07-25 09:57:17 (UTC +08:00)"}}, 0},
		{"idv pass", args{idvPassTemplate, nil}, 0},
		{"idv risk", args{idvRiskTemplate, nil}, 0},
		// Notify
		{"notify deposit fiat", args{notifyDepositFiatTemplate, NotifyDepositFiatPayload{Time: "2023-07-25 09:57:17 (UTC +08:00)", Value: "1000"}}, 0},
		{"notify deposit spot", args{notifyDepositSpotTemplate, NotifyDepositSpotPayload{Time: "2023-07-25 09:57:17 (UTC +08:00)", CurrenciesSymbol: "BTC", Mainnet: "BTC", Amount: "1"}}, 0},
		{"notify invoice", args{notifyInvoiceTemplate, NotifyInvoicePayload{OrderID: "1023G10549324", YearMonth: "112年07-08月", InvNo: "RK-82947670", InvoiceTime: "2023-07-29 23:59:59", RandomNumber: "1670", Total: "3", Seller: "53926705", HandlingCharge: "10*1", SalesAmount: "10"}}, 0},
		{"login notify", args{notifyLoginTemplate, NotifyLoginPayload{Account: "example@test.com", Device: "Mac", IP: "127.0.0.1", Browser: "Chrome", Location: "Taipei", Time: "2023-07-25 09:57:17 (UTC +08:00)"}}, 0},
		{"notify registered", args{notifyRegisteredTemplate, NotifyRegisteredPayload{Time: "2023-07-25 09:57:17 (UTC +08:00)"}}, 0},
		{"notify transaction failed", args{notifyTransactionFailedTemplate, NotifyTransactionFailedPayload{OrderID: "1023G10549324", TransactionPair: "BTC/USDT", TransactionDirection: "賣", Amount: "1"}}, 0},
		{"notify withdraw fiat", args{notifyWithdrawFiatTemplate, NotifyWithdrawFiatPayload{Time: "2023-07-25 09:57:17 (UTC +08:00)", BankAccount: "12340340940", BankName: "台灣銀行", Name: "張小明", Amount: "1000"}}, 0},
		{"notify withdraw spot", args{notifyWithdrawSpotTemplate, NotifyWithdrawSpotPayload{Time: "2023-07-25 09:57:17 (UTC +08:00)", CurrenciesSymbol: "BTC", Mainnet: "BTC", Address: "234dfewrw45f4wdsf", Amount: "1"}}, 0},
		// Verification
		{"verification login 2fa", args{verificationLogin2FATemplate, VerificationPayload{Code: "123456", CodeLifeTime: "5"}}, 0},
		{"verification register", args{verificationRegisterTemplate, VerificationPayload{Code: "123456", CodeLifeTime: "5"}}, 0},
		{"verification reset password", args{verificationResetPasswordTemplate, VerificationPayload{Code: "123456", CodeLifeTime: "5"}}, 0},
		// Transaction
		{"transaction failed", args{transactionFailedTemplate, TransactionPayload{TransactionsID: "1023G10549324", BaseSymbol: "BTC", QuoteSymbol: "USDT", Side: "賣", Quantity: "1"}}, 0},
		{"transaction successful", args{transactionSuccessfulTemplate, TransactionPayload{Time: "2023-07-25 09:57:17 (UTC +08:00)", TransactionsID: "1023G10549324", BaseSymbol: "BTC", QuoteSymbol: "USDT", Side: "買", Quantity: "1"}}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1 := executeTemplate(tt.args.name, tt.args.data)
			if !got1.CodeEqualTo(tt.wantErrCode) {
				t.Errorf("executeTemplate() got1 = %v, want %v", got1, tt.wantErrCode)
			}
		})
	}
}
