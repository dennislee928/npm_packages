package wallet

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/rand"
	"bytes"
	"strconv"
	"time"

	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

const (
	BTC  = "BTC"
	ETH  = "ETH"
	USDC = "USDC"
	USDT = "USDT"
	TRX  = "TRX"
)

type CurrencyType int

const (
	CurrencyTypeBTC CurrencyType = iota + 1
	CurrencyTypeETH
	CurrencyTypeUSDC
	CurrencyTypeUSDT
	CurrencyTypeTRX
)

func (ct CurrencyType) String() string {
	switch ct {
	case CurrencyTypeBTC:
		return BTC
	case CurrencyTypeETH:
		return ETH
	case CurrencyTypeUSDC:
		return USDC
	case CurrencyTypeUSDT:
		return USDT
	case CurrencyTypeTRX:
		return TRX
	default:
		return ""
	}
}

func (ct CurrencyType) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(ct), 10)), nil
}

func (ct *CurrencyType) UnmarshalBinary(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*ct = CurrencyType(i)

	return nil
}

func ParseCurrencyType(ct string) (CurrencyType, bool) {
	switch strings.ToUpper(ct) {
	case BTC:
		return CurrencyTypeBTC, true
	case ETH:
		return CurrencyTypeETH, true
	case USDC:
		return CurrencyTypeUSDC, true
	case USDT:
		return CurrencyTypeUSDT, true
	case TRX:
		return CurrencyTypeTRX, true
	default:
		return 0, false
	}
}

type Mainnet int

const (
	MainnetBTC Mainnet = iota + 1
	MainnetETH
	MainnetERC20
	MainnetTRC20
)

func (m Mainnet) String() string {
	switch m {
	case MainnetBTC:
		return "BTC"
	case MainnetETH:
		return "ETH"
	case MainnetERC20:
		return "ERC20"
	case MainnetTRC20:
		return "TRC20"
	default:
		return ""
	}
}

func (m Mainnet) BinanceNetwork() string {
	switch m {
	case MainnetBTC:
		return "BTC"
	case MainnetETH, MainnetERC20:
		return "ETH"
	case MainnetTRC20:
		return TRX
	default:
		return ""
	}
}

func (m Mainnet) MarshalBinary() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(m), 10)), nil
}

func (m *Mainnet) UnmarshalBinary(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*m = Mainnet(i)

	return nil
}

func ParseMainnet(m string) (Mainnet, bool) {
	switch strings.ToUpper(m) {
	case BTC:
		return MainnetBTC, true
	case ETH:
		return MainnetETH, true
	case "ERC20":
		return MainnetERC20, true
	case "TRC20":
		return MainnetTRC20, true
	default:
		return 0, false
	}
}

type Flow int

const (
	FlowDeposit Flow = 1 << iota
	FlowWithdraw
	FlowBoth = FlowDeposit | FlowWithdraw
)

func IntToFlow(i int) (Flow, bool) {
	switch Flow(i) {
	case FlowDeposit:
		return FlowDeposit, true
	case FlowWithdraw:
		return FlowWithdraw, true
	case FlowBoth:
		return FlowBoth, true
	default:
		return 0, false
	}
}

type wallet struct {
	ID           string
	currencyType CurrencyType
	mainnet      Mainnet
	flow         Flow

	apiToken, apiSecret string
	refreshToken        string

	orderIDPrefix string
}

func (w *wallet) copy() wallet {
	return wallet{
		ID:            w.ID,
		currencyType:  w.currencyType,
		mainnet:       w.mainnet,
		flow:          w.flow,
		apiToken:      w.apiToken,
		apiSecret:     w.apiSecret,
		refreshToken:  w.refreshToken,
		orderIDPrefix: w.orderIDPrefix,
	}
}

func (w *wallet) newReq(
	ctx context.Context,
	method, fullURL string,
	urlValues url.Values,
	payload any,
) (*http.Request, *errpkg.Error) {
	var (
		req *http.Request
		err error
	)

	if urlValues == nil {
		urlValues = make(url.Values, 2)
	}

	urlValues.Set("t", strconv.FormatInt(time.Now().Unix(), 10))
	urlValues.Set("r", rand.LetterAndNumberString(8))

	var bodyStr string
	if payload == nil {
		req, err = http.NewRequestWithContext(ctx, method, fullURL, nil)
	} else {
		bs, mErr := json.Marshal(payload)
		if mErr != nil {
			return nil, wrapWithErr(errpkg.CodeJSONMarshal, err)
		}

		req, err = http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(bs))
		bodyStr = string(bs)
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeNewHTTPRequest, err)
	}

	checksum := buildChecksum(bodyStr, urlValues, w.apiSecret)
	req.URL.RawQuery = urlValues.Encode()
	req.Header.Set(HeaderAPICode, w.apiToken)
	req.Header.Set(HeaderChecksum, checksum)

	return req, nil
}
