package wallet

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func testCybavoFixture() {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	configs.Config.Wallet.APIMode = int(APIModeTest)
	MustInit()
}

func Test_buildChecksum(t *testing.T) {
	type args struct {
		body      string
		urlValues url.Values
		secret    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "pass",
			args: args{
				body: `{"count":1}`,
				urlValues: url.Values{
					"t": []string{"1692781551"},
					"r": []string{"ifYkTUui"},
				},
				secret: "test",
			},
			want: "3e0d4a9abbdba8c29c01254228c9764f221f449a8d00662cb005cd9e4379a060",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildChecksum(tt.args.body, tt.args.urlValues, tt.args.secret); got != tt.want {
				t.Errorf("buildChecksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cybavo_CreateDepositAddress(t *testing.T) {
	testCybavoFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr, _, err := Cybavo.CreateDepositAddress(ctx, MainnetBTC)
	if err != nil {
		t.Errorf("cybavo.CreateDepositAddress() error = %v", err)
		return
	}

	t.Logf("BTC deposit addr: %v", addr)
}

func Test_cybavo_AddWithdrawalWhitelistEntry(t *testing.T) {
	testCybavoFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type args struct {
		m    Mainnet
		addr string
	}
	tests := []struct {
		name        string
		args        args
		wantErrCode errpkg.Code
	}{
		{"btc addr", args{MainnetBTC, "mohjSavDdQYHRYXcS3uS6ttaHP8amyvX78"}, 0},
		{"eth addr", args{MainnetETH, "0xBb977B2EE8a111D788B3477D242078d0B837E72b"}, 0},
		{"btc bad addr", args{MainnetBTC, "123"}, errpkg.CodeBadCryptocurrencyAddress},
		{"eth bad addr", args{MainnetETH, "456"}, errpkg.CodeBadCryptocurrencyAddress},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cybavo.AddWithdrawalWhitelistEntry(ctx, 0, tt.args.m, tt.args.addr); !got.CodeEqualTo(tt.wantErrCode) {
				t.Errorf("cybavo.AddWithdrawalWhitelistEntry() = %v, want %v", got, tt.wantErrCode)
			}
		})
	}
}

func Test_cybavo_RemoveWithdrawalWhitelistEntry(t *testing.T) {
	testCybavoFixture()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	type args struct {
		m    Mainnet
		addr string
	}
	tests := []struct {
		name        string
		args        args
		wantErrCode errpkg.Code
	}{
		{"btc addr 1", args{MainnetBTC, "mohjSavDdQYHRYXcS3uS6ttaHP8amyvX78"}, 0},
		{"eth addr 1", args{MainnetETH, "0xBb977B2EE8a111D788B3477D242078d0B837E72b"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cybavo.RemoveWithdrawalWhitelistEntry(ctx, 0, tt.args.m, tt.args.addr); !got.CodeEqualTo(tt.wantErrCode) {
				t.Errorf("cybavo.RemoveWithdrawalWhitelistEntry() = %v, want %v", got, tt.wantErrCode)
			}
		})
	}
}
