package exchange

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

const maxAPIHost = "https://max-api.maicoin.com/api/v2"

type maxExchange struct {
	client *http.Client
}

func newMaxExchange() *maxExchange {
	exchange := new(maxExchange)

	// ref: https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 10
	transport.MaxIdleConnsPerHost = 10

	exchange.client = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return exchange
}

type MaxTrade struct {
	ID            int64           `json:"id"`               // trade id
	Price         decimal.Decimal `json:"price"`            // strike price
	Volume        decimal.Decimal `json:"volume"`           // traded volume
	Funds         decimal.Decimal `json:"funds"`            // total traded amount
	Market        string          `json:"market"`           // market id
	MarketName    string          `json:"market_name"`      // market name
	CreatedAt     int64           `json:"created_at"`       // DEPRECATED - created timestamp (second), use created_at_in_ms instead
	CreatedAtInMs int64           `json:"created_at_in_ms"` // created timestamp (millisecond)
	Side          string          `json:"side"`             // 'bid' or 'ask'; side of maker for public trades
}

func (me *maxExchange) GetLastTradePrice(ctx context.Context, symbol string) (decimal.Decimal, *errpkg.Error) {
	resp, err := me.GetTrades(ctx, symbol, 1)
	if err != nil {
		return decimal.Zero, err
	} else if len(resp) == 0 {
		return decimal.Zero, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallMaxAPI, Err: errors.New("bad resp length")}
	}

	return resp[0].Price, nil
}

// Ref: https://max.maicoin.com/documents/api_list/v2#!/public/getApiV2Trades
func (me *maxExchange) GetTrades(ctx context.Context, symbol string, limit int) ([]MaxTrade, *errpkg.Error) {
	const path = "/trades"

	url, err := url.Parse(maxAPIHost + path)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeParseURL, Err: err}
	}

	query := url.Query()
	query.Add("market", strings.ToLower(symbol)) // unique market id, check /api/v2/markets for available markets
	query.Add("limit", strconv.Itoa(limit))      // returned limit (1~1000, default 50)
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNewHTTPRequest, Err: err}
	}

	rawResp, err := me.client.Do(req)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeDoHTTPRequest, Err: err}
	}
	defer func() {
		if err := rawResp.Body.Close(); err != nil {
			logger.Logger.Err(err).Msg("max exchange read GetPriceFromMax close resp body error")
		}
	}()

	if rawResp.StatusCode >= http.StatusOK && rawResp.StatusCode < http.StatusMultipleChoices { // 200 Series
		var resp []MaxTrade
		if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
			return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONDecode, Err: err}
		}

		return resp, nil
	}

	rawRespBody, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeReadHTTPResponse, Err: err}
	}

	return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallMaxAPI, Err: fmt.Errorf("bad status code: %d, req uri: %q, resp body: %q", rawResp.StatusCode, url.String(), string(rawRespBody))}
}

func (me *maxExchange) GetHistoryPrice(ctx context.Context, symbol string) ([]decimal.Decimal, *errpkg.Error) {
	kLines, err := me.GetKLine(ctx, symbol, 168 /* = 7*24 */, 60)
	if err != nil {
		return nil, err
	}

	output := make([]decimal.Decimal, 0, len(kLines))
	for _, kLine := range kLines {
		output = append(output, kLine.High)
	}

	return output, nil
}

type MaxKLine struct {
	Timestamp int64
	Open      decimal.Decimal
	High      decimal.Decimal
	Low       decimal.Decimal
	Close     decimal.Decimal
	Volume    decimal.Decimal
}

func (mkl *MaxKLine) UnmarshalJSON(data []byte) error {
	var tmp []json.RawMessage
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	} else if len(tmp) != 6 {
		return fmt.Errorf("bad k line length: %q", string(data))
	} else if err = json.Unmarshal(tmp[0], &mkl.Timestamp); err != nil {
		return err
	} else if err = json.Unmarshal(tmp[1], &mkl.Open); err != nil {
		return err
	} else if err = json.Unmarshal(tmp[2], &mkl.High); err != nil {
		return err
	} else if err = json.Unmarshal(tmp[3], &mkl.Low); err != nil {
		return err
	} else if err = json.Unmarshal(tmp[4], &mkl.Close); err != nil {
		return err
	} else if err = json.Unmarshal(tmp[5], &mkl.Volume); err != nil {
		return err
	}

	return nil
}

// Ref: https://max.maicoin.com/documents/api_list/v2#!/public/getApiV2K
func (me *maxExchange) GetKLine(ctx context.Context, symbol string, limit int, period int) ([]MaxKLine, *errpkg.Error) {
	const path = "/k"

	url, err := url.Parse(maxAPIHost + path)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeParseURL, Err: err}
	}

	query := url.Query()
	query.Add("market", strings.ToLower(symbol)) // unique market id, check /api/v2/markets for available markets
	query.Add("limit", strconv.Itoa(limit))      // returned data points limit, default to 30
	query.Add("period", strconv.Itoa(period))    // time period of K line in minute, default to 1
	url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNewHTTPRequest, Err: err}
	}

	rawResp, err := me.client.Do(req)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeDoHTTPRequest, Err: err}
	}
	defer func() {
		if err := rawResp.Body.Close(); err != nil {
			logger.Logger.Err(err).Msg("max exchange read GetPriceFromMax close resp body error")
		}
	}()

	if rawResp.StatusCode >= http.StatusOK && rawResp.StatusCode < http.StatusMultipleChoices { // 200 Series
		var resp []MaxKLine
		if err = json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
			return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONDecode, Err: err}
		}

		return resp, nil
	}

	rawRespBody, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeReadHTTPResponse, Err: err}
	}

	return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallMaxAPI, Err: fmt.Errorf("bad status code: %d, req uri: %q, resp body: %q", rawResp.StatusCode, url.String(), string(rawRespBody))}
}

type MaxTicker struct {
	At       int64           `json:"at"`         // timestamp in seconds since Unix epoch ,
	Buy      decimal.Decimal `json:"buy"`        // highest buy price ,
	Sell     decimal.Decimal `json:"sell"`       // lowest sell price ,
	Open     decimal.Decimal `json:"open"`       // price before 24 hours ,
	Low      decimal.Decimal `json:"low"`        // lowest price within 24 hours ,
	High     decimal.Decimal `json:"high"`       // highest price within 24 hours ,
	Last     decimal.Decimal `json:"last"`       // last traded price ,
	Vol      decimal.Decimal `json:"vol"`        // traded volume within 24 hours ,
	VolInBTC decimal.Decimal `json:"vol_in_btc"` // traded volume within 24 hours in equal BTC
}

// Ref: https://max.maicoin.com/documents/api_list/v2#!/public/getApiV2TickersPathMarket
func (me *maxExchange) GetPrice(ctx context.Context, symbol string) (*BookTicker, *errpkg.Error) {
	const path = "/tickers/"

	url, err := url.Parse(maxAPIHost + path + strings.ToLower(symbol))
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeParseURL, Err: err}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNewHTTPRequest, Err: err}
	}

	rawResp, err := me.client.Do(req)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeDoHTTPRequest, Err: err}
	}
	defer func() {
		if err := rawResp.Body.Close(); err != nil {
			logger.Logger.Err(err).Msg("max exchange read GetPriceFromMax close resp body error")
		}
	}()

	if rawResp.StatusCode >= http.StatusOK && rawResp.StatusCode < http.StatusMultipleChoices { // 200 Series
		var resp MaxTicker
		if err = json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
			return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONDecode, Err: err}
		}

		output := new(BookTicker)
		output.fromMax(symbol, &resp)

		return output, nil
	}

	rawRespBody, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeReadHTTPResponse, Err: err}
	}

	return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallMaxAPI, Err: fmt.Errorf("bad status code: %d, req uri: %q, resp body: %q", rawResp.StatusCode, url.String(), string(rawRespBody))}
}
