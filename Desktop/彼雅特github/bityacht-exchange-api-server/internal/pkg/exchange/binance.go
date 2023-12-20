package exchange

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/adshao/go-binance/v2"
	binancecommon "github.com/adshao/go-binance/v2/common"
	"github.com/shopspring/decimal"
)

type binanceExchange struct {
	Client     *binance.Client
	SpotClient *binance.Client
}

func newBinaceExchange() *binanceExchange {
	exchange := &binanceExchange{
		Client: binance.NewClient(configs.Config.Exchange.Binance.Apikey, configs.Config.Exchange.Binance.SecretKey),
	}

	exchange.Client.Debug = configs.Config.Exchange.Binance.Debug
	exchange.SpotClient = exchange.Client

	if configs.Config.Exchange.Binance.UseTestnet {
		binance.UseTestnet = true
		exchange.SpotClient = binance.NewClient(configs.Config.Exchange.Binance.TestnetApikey, configs.Config.Exchange.Binance.TestnetSecretKey)
		exchange.SpotClient.Debug = configs.Config.Exchange.Binance.Debug
	}

	return exchange
}

func (be binanceExchange) GetAllCoinsInfo(ctx context.Context) ([]*binance.CoinInfo, *errpkg.Error) {
	// https://testnet.binance.vision/ F.A.Q.:
	// Can I use the /sapi endpoints on the Spot Test Network?
	// No, only the /api endpoints are available on the Spot Test Network:
	allCoinsInfo, err := be.Client.NewGetAllCoinsInfoService().Do(ctx)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: err}
	}

	return allCoinsInfo, nil
}

func (be *binanceExchange) GetExchangeInfo(ctx context.Context, symbol string) (*ExchangeInfo, *errpkg.Error) {
	exchangeInfo, err := be.SpotClient.NewExchangeInfoService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: err}
	} else if len(exchangeInfo.Symbols) != 1 {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: fmt.Errorf("bad ExchangeInfo resp: '%+v'", exchangeInfo)}
	}

	var output ExchangeInfo
	if err := output.fromBinance(exchangeInfo.Symbols[0]); err != nil {
		return nil, err
	}

	return &output, nil
}

// Ref: https://binance-docs.github.io/apidocs/spot/en/#symbol-order-book-ticker
func (be *binanceExchange) GetPrice(ctx context.Context, symbol string) (*BookTicker, *errpkg.Error) {
	bookTickerList, err := be.SpotClient.NewListBookTickersService().Symbol(symbol).Do(ctx)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: err}
	} else if len(bookTickerList) != 1 {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: fmt.Errorf("bad ListBookTickers resp: '%+v'", bookTickerList)}
	}

	var output BookTicker
	if err := output.fromBinance(bookTickerList[0]); err != nil {
		return nil, err
	}

	return &output, nil
}

// Ref: https://binance-docs.github.io/apidocs/spot/en/#kline-candlestick-data
func (be *binanceExchange) GetHistoryPrice(ctx context.Context, symbol string) ([]decimal.Decimal, *errpkg.Error) {
	// Interval: "1m", "3m", "5m", "15m", "30m", "1h", "2h", "4h", "6h", "8h", "12h", "1d", "3d", "1w", "1M"
	kLines, err := be.SpotClient.NewKlinesService().Symbol(symbol).Interval("1h").Limit(168 /* = 7*24 */).Do(ctx)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: err}
	} else if len(kLines) == 0 {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: fmt.Errorf("bad Klines resp: '%+v'", kLines)}
	}

	output := make([]decimal.Decimal, 0, len(kLines))
	for i, kLine := range kLines {
		high, err := decimal.NewFromString(kLine.High)
		if err != nil {
			return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: fmt.Errorf("bad Klines[%d]: '%+v'", i, *kLine)}
		}

		output = append(output, high)
	}

	return output, nil
}

// Ref: https://binance-docs.github.io/apidocs/spot/en/#new-order-trade
func (be *binanceExchange) CreateOrder(ctx context.Context, symbol string, side binance.SideType, orderType binance.OrderType, quantity decimal.Decimal, price decimal.Decimal) (*CreateOrderResponse, *errpkg.Error) {
	resp, err := be.SpotClient.NewCreateOrderService().
		Symbol(symbol).               // Mandatory
		Side(side).                   // Mandatory
		Type(binance.OrderTypeLimit). // Mandatory
		TimeInForce(binance.TimeInForceTypeFOK).
		Quantity(quantity.String()).
		Price(price.String()).
		Do(ctx)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallBinanceAPI, Err: err}
	}

	var output CreateOrderResponse
	if err := output.fromBinance(resp); err != nil {
		return nil, err
	}

	return &output, nil
}

func (binanceExchange) IsInvalidSymbol(err error) bool {
	if err == nil {
		return false
	}

	var binanceErr *binancecommon.APIError
	if errors.As(err, &binanceErr) && binanceErr.Code == -1121 {
		return true
	}

	return false
}
