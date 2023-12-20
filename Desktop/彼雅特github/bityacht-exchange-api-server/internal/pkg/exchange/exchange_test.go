package exchange

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"context"
	"testing"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
)

func TestGetPrice(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if bookTicker, err := Binance.GetPrice(ctx, "ETHUSDT"); err != nil {
		t.Error(err)
	} else {
		t.Logf("bookTicker: %+v\n", bookTicker)
	}

	if bookTicker, err := Max.GetPrice(ctx, "USDTTWD"); err != nil {
		t.Error(err)
	} else {
		t.Logf("bookTicker: %+v\n", bookTicker)
	}
}

func TestGetFiatPrice(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	defer cancel()

	if price, err := Max.GetLastTradePrice(ctx, "USDTTWD"); err != nil {
		t.Error(err)
	} else {
		t.Logf("price: %+v\n", price)
	}
}

func TestGetHistoryPrice(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if histories, err := Binance.GetHistoryPrice(ctx, "ETHUSDT"); err != nil {
		t.Error(err)
	} else {
		t.Logf("histories: %+v\n", histories)
	}

	if histories, err := Max.GetHistoryPrice(ctx, "USDCTWD"); err != nil {
		t.Error(err)
	} else {
		t.Logf("histories: %+v\n", histories)
	}
}

func TestBinance(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	const symbol = "ETHUSDT"

	if resp, err := Binance.GetPrice(ctx, symbol); err != nil {
		t.Error(err)
	} else {
		t.Logf("GetLastTradePrice: %+v\n", resp)
	}

	if resp, err := Binance.GetHistoryPrice(ctx, symbol); err != nil {
		t.Error(err)
	} else {
		t.Logf("GetLastTradePrice: %+v\n", resp)
	}

	// Binance.CreateOrder
}

func TestMax(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	const symbol = "USDTTWD"

	if resp, err := Max.GetLastTradePrice(ctx, symbol); err != nil {
		t.Error(err)
	} else {
		t.Logf("GetLastTradePrice: %+v\n", resp)
	}

	if resp, err := Max.GetTrades(ctx, symbol, 5); err != nil {
		t.Error(err)
	} else {
		t.Logf("GetTrades: %+v\n", resp)
	}

	if resp, err := Max.GetHistoryPrice(ctx, symbol); err != nil {
		t.Error(err)
	} else {
		t.Logf("GetHistoryPrice: %+v\n", resp)
	}

	if resp, err := Max.GetKLine(ctx, symbol, 10, 60); err != nil {
		t.Error(err)
	} else {
		t.Logf("GetKLine: %+v\n", resp)
	}

	if resp, err := Max.GetPrice(ctx, symbol); err != nil {
		t.Error(err)
	} else {
		t.Logf("GetPrice: %+v\n", resp)
	}
}

func Test_binanceExchange_CreateOrder(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	type args struct {
		symbol    string
		side      binance.SideType
		orderType binance.OrderType
		quantity  decimal.Decimal
		price     decimal.Decimal
	}
	tests := []struct {
		name        string
		args        args
		wantErrCode errpkg.Code
	}{
		{"test1", args{"BTCUSDT", binance.SideTypeBuy, binance.OrderTypeLimit, decimal.NewFromFloat(0.0004), decimal.NewFromFloat(26000.0)}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			got, got1 := Binance.CreateOrder(ctx, tt.args.symbol, tt.args.side, tt.args.orderType, tt.args.quantity, tt.args.price)
			if !got1.CodeEqualTo(tt.wantErrCode) {
				t.Errorf("binanceExchange.CreateOrder() got err = %v, want %v", got1, tt.wantErrCode)
			} else {
				t.Log(got)
			}
		})
	}
}

func Test_binanceExchange_GetExchangeInfo(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	type args struct {
		symbol string
	}
	tests := []struct {
		name        string
		args        args
		wantErrCode errpkg.Code
	}{
		{"BTCETH", args{"BTCETH"}, errpkg.CodeCallBinanceAPI},
		{"ETHBTC", args{"ETHBTC"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			got, got1 := Binance.GetExchangeInfo(ctx, tt.args.symbol)
			if !got1.CodeEqualTo(tt.wantErrCode) {
				t.Errorf("binanceExchange.GetExchangeInfo() got err = %v, want %v", got1, tt.wantErrCode)
			} else {
				t.Log(got)
			}
		})
	}
}

func Test_binanceExchange_GetAllCoinsInfo(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	Init()

	t.Run("GetAllCoinsInfo", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		got, got1 := Binance.GetAllCoinsInfo(ctx)
		if got1 != nil {
			t.Errorf("binanceExchange.GetAllCoinsInfo() got err = %v", got1)
		} else {
			for _, v := range got {
				t.Logf("%+v\n", v)
			}
		}
	})
}
