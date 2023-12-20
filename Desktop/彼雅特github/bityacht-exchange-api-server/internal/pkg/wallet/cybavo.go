package wallet

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type APIMode int

const (
	APIModeTest APIMode = iota
	APIModeProduction
)

var (
	_once  sync.Once
	Cybavo CybavoWallet
)

func MustInit() {
	_once.Do(func() {
		cybavo, err := newCybavo(configs.Config.Wallet)
		if err != nil {
			logger.Logger.Fatal().Err(err).Msg("cybavo init failed")
		}
		Cybavo = cybavo
	})
}

type CybavoWallet interface {
	WithdrawalCallback(checksum string, r io.Reader) error

	Callback(checksum string, r io.Reader) (*CallbackInfo, error)

	QueryAPICodeStatus(ctx context.Context, walletID string) error

	CreateDepositAddress(ctx context.Context, m Mainnet) (address string, txID string, wrapErr *errpkg.Error)

	GetContractDepositAddress(ctx context.Context, ct CurrencyType, m Mainnet, txID string) (string, *errpkg.Error)

	WithdrawAssets(
		ctx context.Context, ct CurrencyType, m Mainnet, payload WithdrawPayload,
	) *errpkg.Error

	VerifyAddress(ctx context.Context, ct CurrencyType, m Mainnet, addr string) (*VerifyAddressItem, *errpkg.Error)

	AddWithdrawalWhitelistEntry(ctx context.Context, usersID int64, m Mainnet, addr string) *errpkg.Error

	RemoveWithdrawalWhitelistEntry(ctx context.Context, usersID int64, m Mainnet, addr string) *errpkg.Error
}

type cybavo struct {
	client  *http.Client
	baseURL string

	// wallets is a map of walletID to wallet.
	// NOTE: wallets shouldn't be modified after initialized.
	wallets    map[string]*wallet
	walletList []*wallet // The order of list is same as config.yaml
}

var _ CybavoWallet = (*cybavo)(nil)

func newCybavo(cfg configs.WalletConfig) (CybavoWallet, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 10
	transport.MaxIdleConnsPerHost = 10

	if len(cfg.Wallets) == 0 {
		return nil, errors.New("cybavo api token or api secret is empty")
	}

	c := &cybavo{
		client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},

		wallets:    make(map[string]*wallet, len(cfg.Wallets)),
		walletList: make([]*wallet, 0, len(cfg.Wallets)),
	}

	for _, walletCfg := range cfg.Wallets {
		if walletCfg.ID == "" {
			logger.Logger.Warn().Msg("wallet id is empty")
			continue
		}

		if _, exists := c.wallets[walletCfg.ID]; exists {
			logger.Logger.Warn().Str("walletID", walletCfg.ID).Msg("duplicate wallet id, skip")
			continue
		}

		ct, ok := ParseCurrencyType(walletCfg.Currency)
		if !ok {
			logger.Logger.Warn().Str("currencyType", walletCfg.Currency).Msg("parse currency type failed")
			continue
		}

		m, ok := ParseMainnet(walletCfg.Mainnet)
		if !ok {
			logger.Logger.Warn().Str("mainnet", walletCfg.Mainnet).Msg("parse mainnet failed")
			continue
		}

		f, ok := IntToFlow(walletCfg.Flow)
		if !ok {
			logger.Logger.Warn().Int("flow", walletCfg.Flow).Msg("parse flow failed")
			continue
		}

		w := &wallet{
			ID:            walletCfg.ID,
			currencyType:  ct,
			mainnet:       m,
			flow:          f,
			apiToken:      walletCfg.Token,
			apiSecret:     walletCfg.Secret,
			refreshToken:  walletCfg.RefreshToken,
			orderIDPrefix: walletCfg.OrderIDPrefix,
		}
		c.wallets[walletCfg.ID] = w
		c.walletList = append(c.walletList, w)
	}

	switch APIMode(cfg.APIMode) {
	case APIModeTest:
		c.baseURL = testBaseURL
	case APIModeProduction:
		c.baseURL = prodBaseURL
	default:
		return nil, fmt.Errorf("unknown api mode: %d", cfg.APIMode)
	}

	wg := new(sync.WaitGroup)
	for _, w := range c.walletList {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			ctx, cancle := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancle()

			if err := c.QueryAPICodeStatus(ctx, id); err != nil {
				logger.Logger.Err(err).Str("service", "cybavo").Msg("query API Code Status Failed")
			}
		}(w.ID)
	}

	wg.Wait()

	return c, nil
}

// Ref: https://www.cybavo.com/zh-tw/developers/rest-api/common-api/query-api-code-status/
// For Activate the API Code
func (c *cybavo) QueryAPICodeStatus(ctx context.Context, walletID string) error {
	w, ok := c.getWalletByID(walletID)
	if !ok {
		return &errpkg.Error{
			HttpStatus: http.StatusBadRequest,
			Code:       errpkg.CodeCybavoWalletNotFound,
			Err:        fmt.Errorf("wallet %q not found", walletID),
		}
	}

	uri := fmt.Sprintf("/v1/sofa/wallets/%s/apisecret", walletID)
	resp, wrapErr := c.sendReq(ctx, w, http.MethodGet, uri, nil, nil)
	if wrapErr != nil {
		return wrapErr
	}

	var respBody QueryAPICodeStatusResp
	if err := json.NewDecoder(resp).Decode(&respBody); err != nil {
		return wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	logger.Logger.Info().Str("service", "cybavo").Any("API Code Status", respBody).Send()

	return nil
}

// WithdrawalCallback is the callback info from Cybavo Wallet.(When Withdrawing)
func (c *cybavo) WithdrawalCallback(checksum string, r io.Reader) error {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	var req WithdrawPayload
	if err := json.NewDecoder(tee).Decode(&req); err != nil {
		return err
	}
	if len(req.Requests) == 0 {
		return errors.New("bad length of requests")
	}

	wallet, ok := c.getWalletByOrderIDAndFlow(req.Requests[0].OrderID, FlowWithdraw)
	if !ok {
		return errors.New("wallet not found")
	}

	if !checksumVerify(wallet.apiSecret, checksum, &buf) {
		return ErrChecksumVerify
	}

	// TODO: Check the OrderID is create from our system.

	return nil
}

// CallbackInfo is the callback info from Cybavo Wallet.
// Checksum is the value of header `HeaderCbChecksum`
func (c *cybavo) Callback(checksum string, r io.Reader) (*CallbackInfo, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	var info *CallbackInfo
	if err := json.NewDecoder(tee).Decode(&info); err != nil {
		return nil, err
	}

	wallet, ok := c.getWalletByID(strconv.FormatInt(info.WalletID, 10))
	if !ok {
		return nil, fmt.Errorf("wallet %q not found", info.WalletID)
	}

	if !checksumVerify(wallet.apiSecret, checksum, &buf) {
		return nil, ErrChecksumVerify
	}

	info.wallet = wallet
	info.currencyType = wallet.currencyType
	info.mainnet = wallet.mainnet

	return info, nil
}

func (c *cybavo) CreateDepositAddress(ctx context.Context, m Mainnet) (address string, txID string, wrapErr *errpkg.Error) {
	w, ok := c.getWalletByBinanceNetworkFlow(m, FlowDeposit)
	if !ok {
		return "", "", &errpkg.Error{
			HttpStatus: http.StatusBadRequest,
			Code:       errpkg.CodeCybavoWalletNotFound,
			Err:        fmt.Errorf("wallet not found, mainnet: %+v", m),
		}
	}

	// we only need one address for a user.
	const num = 1

	uri := fmt.Sprintf("/v1/sofa/wallets/%s/addresses", w.ID)
	resp, wrapErr := c.sendReq(ctx, w, http.MethodPost, uri, CreateDepositAddressesPayload{Count: num}, nil)
	if wrapErr != nil {
		return "", "", wrapErr
	}

	var respBody CreateDepositAddressesResp
	if err := json.NewDecoder(resp).Decode(&respBody); err != nil {
		return "", "", wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	if addrNum := len(respBody.Addresses); addrNum == num {
		address = respBody.Addresses[0]
	} else if txNums := len(respBody.TXIDs); txNums == num {
		txID = respBody.TXIDs[0]
	} else {
		return "", "", wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected number of addresses & txids: %d, %d (expected %d)", addrNum, txNums, num))
	}

	return
}

func (c *cybavo) GetContractDepositAddress(ctx context.Context, ct CurrencyType, m Mainnet, txID string) (string, *errpkg.Error) {
	wallet, ok := c.getWalletByTypeFlow(ct, m, FlowDeposit)
	if !ok {
		return "", wrapWithErr(errpkg.CodeCybavoWalletNotFound, fmt.Errorf("wallet not found, currencyType: %d", ct))
	}

	uri := fmt.Sprintf("/v1/sofa/wallets/%s/addresses/contract_txid", wallet.ID)
	resp, wrapErr := c.sendReq(ctx, wallet, http.MethodGet, uri, nil, url.Values{"txids": []string{txID}})
	if wrapErr != nil {
		return "", wrapErr
	}

	var respBody GetContractAddressesResp
	if err := json.NewDecoder(resp).Decode(&respBody); err != nil {
		return "", wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	var items map[string]GetContractAddressItem
	if err := json.Unmarshal(respBody.Addresses, &items); err != nil {
		return "", wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	switch len(items) {
	case 0:
		return "", &errpkg.Error{
			HttpStatus: http.StatusAccepted,
			Code:       errpkg.CodeAddressDeploying,
			Err:        errors.New("address is deploying"),
		}
	case 1:
		return items[txID].Address, nil

	default:
		return "", wrapWithErr(errpkg.CodeCallCybavoAPI, errors.New("unexpected address number of items"))
	}
}

func (c *cybavo) WithdrawAssets(
	ctx context.Context, ct CurrencyType, m Mainnet, payload WithdrawPayload,
) *errpkg.Error {
	w, ok := c.getWalletByTypeFlow(ct, m, FlowWithdraw)
	if !ok {
		return wrapWithErr(errpkg.CodeCybavoWalletNotFound, fmt.Errorf("wallet not found, currencyType: %d", ct))
	}

	for i := range payload.Requests {
		payload.Requests[i].OrderID = w.orderIDPrefix + payload.Requests[i].OrderID
	}

	uri := fmt.Sprintf("/v1/sofa/wallets/%s/sender/transactions", w.ID)
	resp, wrapErr := c.sendReq(ctx, w, http.MethodPost, uri, payload, nil)
	if wrapErr != nil {
		return wrapErr
	}

	var respBody WithdrawAssetsResp
	if err := json.NewDecoder(resp).Decode(&respBody); err != nil {
		return wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	logger.Logger.Info().Any("respBody", respBody).Msg("withdraw assets resp")

	return nil
}

func (c *cybavo) sendReq(
	ctx context.Context, wallet wallet, method, targetURI string, payload any, urlValues url.Values,
) (io.Reader, *errpkg.Error) {
	fullURL, err := url.JoinPath(c.baseURL, targetURI)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeParseURL, err)
	}

	req, reqWrapErr := wallet.newReq(ctx, method, fullURL, urlValues, payload)
	if reqWrapErr != nil {
		return nil, reqWrapErr
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeCallCybavoAPI, err)
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	tee := io.TeeReader(resp.Body, &buf)

	if resp.StatusCode != http.StatusOK {
		if err := parseErrResp(resp.Body); err != nil {
			return nil, wrapWithErr(errpkg.CodeCallCybavoAPI, err)
		}
		return nil, wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected err resp body, status code: %d", resp.StatusCode))
	}

	if !checksumVerify(wallet.apiSecret, resp.Header.Get(HeaderChecksum), tee) {
		return nil, wrapWithErr(errpkg.CodeCallCybavoAPI, ErrChecksumVerify)
	}

	return &buf, nil
}

// getWalletByMainnetFlow returns the first wallet that matches the given
// mainnet and flow.
// NOTE: It only returns the first matching wallet
func (c *cybavo) getWalletByMainnetFlow(m Mainnet, flow Flow) (wallet, bool) {
	for _, w := range c.walletList {
		if w.mainnet == m &&
			(w.flow&flow > 0) {

			return w.copy(), true
		}
	}

	return wallet{}, false
}

// getWalletByBinanceNetworkFlow returns the first wallet that matches the given
// BinanceNetwork and flow.
// NOTE: It only returns the first matching wallet
func (c *cybavo) getWalletByBinanceNetworkFlow(m Mainnet, flow Flow) (wallet, bool) {
	for _, w := range c.walletList {
		if w.mainnet.BinanceNetwork() == m.BinanceNetwork() &&
			(w.flow&flow > 0) {

			return w.copy(), true
		}
	}

	return wallet{}, false
}

// getWalletByTypeFlow returns the first wallet that matches the given
// currency type, mainnet and flow.
// NOTE: It only returns the first matching wallet
func (c *cybavo) getWalletByTypeFlow(ct CurrencyType, m Mainnet, flow Flow) (wallet, bool) {
	for _, w := range c.walletList {
		if w.currencyType == ct &&
			w.mainnet == m &&
			(w.flow&flow > 0) {

			return w.copy(), true
		}
	}

	return wallet{}, false
}

func (c *cybavo) getWalletByID(id string) (wallet, bool) {
	w, ok := c.wallets[id]
	if !ok {
		return wallet{}, false
	}

	return w.copy(), true
}

func (c *cybavo) getWalletByOrderIDAndFlow(orderID string, flow Flow) (wallet, bool) {
	for _, w := range c.walletList {
		if (w.flow&flow > 0) && strings.HasPrefix(orderID, w.orderIDPrefix) {
			return w.copy(), true
		}
	}

	return wallet{}, false
}

func (c *cybavo) VerifyAddress(ctx context.Context, ct CurrencyType, m Mainnet, addr string) (*VerifyAddressItem, *errpkg.Error) {
	w, ok := c.getWalletByTypeFlow(ct, m, FlowDeposit)
	if !ok {
		return nil, wrapWithErr(errpkg.CodeCybavoWalletNotFound, fmt.Errorf("wallet not found, currencyType: %d", ct))
	}

	uri := fmt.Sprintf("/v1/sofa/wallets/%s/addresses/verify", w.ID)
	resp, wrapErr := c.sendReq(ctx, w, http.MethodPost, uri, VerifyAddressPayload{Addresses: []string{addr}}, nil)
	if wrapErr != nil {
		return nil, wrapErr
	}

	var respBody VerifyAddressResp
	if err := json.NewDecoder(resp).Decode(&respBody); err != nil {
		return nil, wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	if resultNums := len(respBody.Result); resultNums != 1 {
		return nil, wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected number of result: %d", resultNums))
	} else if respBody.Result[0].Address != addr {
		return nil, wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected address of result: %q != (expected) %q", respBody.Result[0].Address, addr))
	}

	return &respBody.Result[0], nil
}

func (c *cybavo) AddWithdrawalWhitelistEntry(ctx context.Context, usersID int64, m Mainnet, addr string) *errpkg.Error {
	w, ok := c.getWalletByMainnetFlow(m, FlowWithdraw)
	if !ok {
		return wrapWithErr(errpkg.CodeCybavoWalletNotFound, fmt.Errorf("wallet not found, mainnet: %d", m))
	}

	uri := fmt.Sprintf("/v1/sofa/wallets/%s/sender/whitelist", w.ID)
	resp, wrapErr := c.sendReq(ctx, w, http.MethodPost, uri, AddWithdrawalWhitelistEntryPayload{Items: []WithdrawalWhitelistEntryItem{{Address: addr, UserID: strconv.FormatInt(usersID, 10)}}}, nil)
	if wrapErr != nil {
		var errResp ErrResp
		if wrapErr.Err != nil && errors.As(wrapErr.Err, &errResp) {
			switch errResp.ErrorCode {
			case 399: // Duplicated entry
				return nil
			case 703: // Operation failed: one of error in ["invalid address", "invalid user id", "wallet does not support memo"]
				return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadCryptocurrencyAddress}
			}
		}

		return wrapErr
	}

	var respBody AddWithdrawalWhitelistEntryResp
	if err := json.NewDecoder(resp).Decode(&respBody); err != nil {
		return wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	if resultNums := len(respBody.AddedItems); resultNums != 1 {
		return wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected number of result: %d", resultNums))
	} else if respBody.AddedItems[0].Address != addr {
		return wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected address of result: %q != (expected) %q", respBody.AddedItems[0].Address, addr))
	}

	return nil
}

func (c *cybavo) RemoveWithdrawalWhitelistEntry(ctx context.Context, usersID int64, m Mainnet, addr string) *errpkg.Error {
	w, ok := c.getWalletByMainnetFlow(m, FlowWithdraw)
	if !ok {
		return wrapWithErr(errpkg.CodeCybavoWalletNotFound, fmt.Errorf("wallet not found, mainnet: %d", m))
	}

	uri := fmt.Sprintf("/v1/sofa/wallets/%s/sender/whitelist", w.ID)
	resp, wrapErr := c.sendReq(ctx, w, http.MethodDelete, uri, RemoveWithdrawalWhitelistEntryPayload{Items: []WithdrawalWhitelistEntryItem{{Address: addr, UserID: strconv.FormatInt(usersID, 10)}}}, nil)
	if wrapErr != nil {
		return wrapErr
	}

	//! Only Check the format of Resp
	var respBody RemoveWithdrawalWhitelistEntryResp
	if err := json.NewDecoder(resp).Decode(&respBody); err != nil {
		return wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	//! respBody.RemovedItems may be empty when addr removed before this api call, So there is no necessary to check RemovedItems
	// if resultNums := len(respBody.RemovedItems); resultNums != 1 {
	// 	return wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected number of result: %d", resultNums))
	// } else if respBody.RemovedItems[0].Address != addr {
	// 	return wrapWithErr(errpkg.CodeCallCybavoAPI, fmt.Errorf("unexpected address of result: %q != (expected) %q", respBody.RemovedItems[0].Address, addr))
	// }

	return nil
}

func buildChecksum(body string, urlValues url.Values, secret string) string {
	values := make([]string, 0, len(urlValues)+2) // 2 for body and secret

	if body != "" {
		values = append(values, body)
	}

	for k, v := range urlValues {
		values = append(values, fmt.Sprintf("%s=%s", k, v[0]))
	}
	sort.Strings(values)
	values = append(values, fmt.Sprintf("secret=%s", secret))

	return fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(values, "&"))))
}

func wrapWithErr(code errpkg.Code, err error) *errpkg.Error {
	return &errpkg.Error{
		HttpStatus: http.StatusInternalServerError,
		Code:       code,
		Err:        err,
	}
}
