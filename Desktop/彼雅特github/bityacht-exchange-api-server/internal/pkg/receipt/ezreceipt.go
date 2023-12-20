package receipt

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bytes"
	"context"
	"crypto/sha1" // #nosec G505: the api of ezreceipt requires sha1
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type APIMode int

const (
	APIModeTest APIMode = iota
	APIModeProduction
)

var (
	_once sync.Once
	EZ    EZReceipt
)

// MustInit initializes the EZReceipt.
// It panics if the EZReceipt fails to initialize.
func MustInit() {
	_once.Do(func() {
		ez, err := newEZReceipt(configs.Config.Receipt)
		if err != nil {
			logger.Logger.Fatal().Err(err).Msg("ezreceipt init error")
		}
		EZ = ez
	})
}

type EZReceipt interface {
	// Login 使用者登入。若登入成功，將會設定 token 及 token 有效期限
	Login(ctx context.Context) *errpkg.Error

	// CheckMobileCode check the mobile barcode is correct or not.
	//
	// 呼叫財政部電子發票營業人應用API檢查所輸入的手機條碼是否存在於財政部電子發票整合服務平台中
	CheckMobileCode(ctx context.Context, mobileCode string) *errpkg.Error

	// InvNumberList lists the invoice numbers.
	//
	// 條列字軌號碼的分段清單，也可以查詢過去的字軌配號
	// 本系統上可對原始字軌做分段管理，所以回傳的清單是分段後的結果 (除非沒有做任何分段的處理)
	InvNumberList(ctx context.Context) (*ListResp[InvNumberRespItem], *errpkg.Error)

	// ConfirmOrder confirms the order.
	//
	// 確認訂單，表示將草稿轉為正式訂單
	// 若該訂單已是正式訂單或作廢，將回報錯誤
	// 訂單中未包含任何銷售品項，若要進行確認也會回報錯誤
	// 當店家有提供折扣或優惠時，可利用 discounts 參數來輸入折扣原因與折扣金額
	// 系統會依據折扣內容自動調整訂單金額
	//
	// 若店家不使用金流服務時，可設定 paidTime 參數來表示該訂單已完成付款
	// 而且是以現金完成交易。若有啟用金流的店家，請以金流相關的 API 完成付款作業，勿使用 paidTime 參數
	ConfirmOrder(ctx context.Context, orderID uint64, payload ConfirmOrderPayload) (*ConfirmOrderRespItem, *errpkg.Error)

	// CreateOrder creates an order.
	//
	// 請用此服務節點產生訂單草稿。購物車視同訂單草稿
	// 若呼叫此節點時已可確認訂單，可將 `confirm` 參數設為 true
	// 否則後續要再呼叫 `/sales/order/confirm` 以便將草稿轉為正式訂單
	//
	// 訂單的開立對象若是公司（以後要開 B2B 發票者），`buyerID` 將是必要參數。
	// 這意味每一個 B2B 的買家，都需要有一個使用者帳號。先建立了 B2B 買家帳號，才能做後續的訂單開立動作
	// 退換貨作業時，請輸入原訂單號碼 (`voidOrderNo`)，以便系統作廢原訂單並對新舊訂單作連結
	// 新舊訂單的連結可以讓系統提供退換貨的資料比對
	//
	// 關於商品號碼：訂單明細中如果有提供商品號碼，系統將依據該商品號碼找出是否有該商品資訊。
	// 找得到就會以該商品資訊做為預設資料
	// 如果找不到，看用戶是否設定了「自動訂新增商品」的選項。如果有設定該選項，則系統會在商品管理中
	// 自動為用戶新增該商品；否則將回報錯誤。如果商品明細中未設定商品號碼，系統將檢查是否輸入商品
	// 名稱與價格，決定商品明細是否成立。在自動新增商品與未輸入商品號碼的狀況下，商品名稱與價格
	// 將是必要的輸入資料
	CreateOrder(ctx context.Context, payload CreateOrderPayload) (*CreateOrderRespItem, *errpkg.Error)

	// SetCarrier sets the carrier of the order.
	// 設定訂單未來要開立發票時的載具資訊。若是要開立 B2C 發票，開立發票前一定要先呼叫此服務節點
	SetCarrier(ctx context.Context, orderID uint64, payload SetCarrierPayload) (*SetCarrierRespItem, *errpkg.Error)

	// CreateB2CInvoice creates a B2C invoice.
	//
	// 依據現有訂單開立 B2C 電子發票。開立發票前，必須先成立訂單
	// 此外，還必須先呼叫 `/sales/order/setCarrier` 去完成發票載具的設定，否則系統將回報錯誤
	// 如果是發票已註銷再重開，`invNo` 和 `invoiceTime` 這二個參數都將被
	// 忽略（沿用原來的發票號碼或開立時間）
	// 如果訂單是發生退換貨後產生的新訂單，系統會自動將退換貨前的舊發票作廢
	CreateB2CInvoice(ctx context.Context, orderID uint64, payload CreateB2CInvoicePayload) (*CreateB2CInvoiceRespItem, *errpkg.Error)

	// FastCreateB2CInvoice creates a B2C invoice in one step. In this step, it will:
	// 	1. CreateOrder -> 2. SetCarrier -> 3. CreateB2CInvoice
	FastCreateB2CInvoice(ctx context.Context, payload FastCreateB2CPayload) (*CreateB2CInvoiceRespItem, *errpkg.Error)

	// VoidOrder voids the order.
	//
	// 將指定的訂單作廢。若是尚未成立的訂單（草稿訂單）將會直接刪除
	// 已付費的訂單需完成退費才能作廢
	VoidOrder(ctx context.Context, orderID uint64) *errpkg.Error

	// CleanupExistsOrder cleans up the exists order.
	//
	// 清除已存在的訂單
	CleanupExistsOrder(ctx context.Context, orderNo string) *errpkg.Error
}

type ezToken struct {
	mu sync.Mutex

	token   string
	validTo time.Time
}

type ezreceipt struct {
	client  *http.Client
	baseURL string

	appCode, appKey string
	apiAcc, apiPwd  string

	token ezToken
}

var _ EZReceipt = (*ezreceipt)(nil)

func newEZReceipt(cfg configs.ReceiptConfig) (EZReceipt, error) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 10
	transport.MaxIdleConnsPerHost = 10

	if cfg.AppCode == "" || cfg.AppKey == "" {
		return nil, errors.New("app code or app key is empty")
	}

	if cfg.APIAcc == "" || cfg.APIPwd == "" {
		return nil, errors.New("api acc or api pwd is empty")
	}

	ez := &ezreceipt{
		client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},

		appCode: cfg.AppCode,
		appKey:  cfg.AppKey,
		apiAcc:  cfg.APIAcc,
		apiPwd:  cfg.APIPwd,

		token: ezToken{validTo: time.Now()},
	}

	switch APIMode(cfg.APIMode) {
	case APIModeTest:
		ez.baseURL = testBaseURL
	case APIModeProduction:
		ez.baseURL = prodBaseURL
	default:
		return nil, fmt.Errorf("unknown api mode: %d", cfg.APIMode)
	}

	return ez, nil
}

func (ez *ezreceipt) newReq(ctx context.Context, targetURI string, payload any) (*http.Request, *errpkg.Error) {
	const apiMethod = http.MethodPost

	fullURL, err := url.JoinPath(ez.baseURL, targetURI)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeParseURL, err)
	}

	// The ezreceipt API requires an empty JSON if the payload is nil.
	if payload == nil {
		payload = struct{}{}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, wrapWithErr(errpkg.CodeJSONMarshal, err)
	}

	req, err := http.NewRequestWithContext(ctx, apiMethod, fullURL, &buf)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeNewHTTPRequest, err)
	}

	// login if token is invalid, but not login uri.
	isLoginURI := targetURI == loginURI
	if !isLoginURI {
		if err := ez.Login(ctx); err != nil {
			return nil, err
		}
	}
	ez.setEZHeader(req, !isLoginURI)

	return req, nil
}

// ezResp is the response of EZReceipt API.
// Use the `ezResp.Error()` to check if the response is success.
type ezResp struct {
	Code    ezErrCode       `json:"code"`
	Message string          `json:"message"`
	Value   json.RawMessage `json:"value"`
}

func (e ezResp) Err() error {
	if err := e.Code.Err(); err != nil {
		return fmt.Errorf("bad ez resp: %w: %s", err, e.Message)
	}

	return nil
}

func (ez *ezreceipt) setEZHeader(req *http.Request, withToken bool) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(headerAppCode, ez.appCode)
	req.Header.Set(headerAppKey, ez.appKey)
	if withToken {
		req.Header.Set(headerToken, ez.Token())
	}
}

const loginURI = "admin/user/login"

// Login logins to EZReceipt API.
// If success, it will set the token and token valid until time.
func (ez *ezreceipt) Login(ctx context.Context) *errpkg.Error {
	ez.token.mu.Lock()
	defer ez.token.mu.Unlock()

	if ez.isTokenValid(time.Now()) {
		return nil
	}

	if ez.apiAcc == "" || ez.apiPwd == "" {
		return wrapWithErr(errpkg.CodeBadParamEZReceiptAPI, errors.New("api acc or api pwd is empty"))
	}

	apiPwd, err := loginPwd(ez.apiAcc, ez.apiPwd)
	if err != nil {
		return wrapWithErr(errpkg.CodeBadParam, fmt.Errorf("login pwd: %w", err))
	}

	payload := struct {
		AccName string `json:"accName"`
		Passwd  string `json:"passwd"`
	}{
		AccName: ez.apiAcc,
		Passwd:  apiPwd,
	}
	req, reqWrapErr := ez.newReq(ctx, loginURI, payload)
	if err != nil {
		return reqWrapErr
	}

	rawResp, err := ez.client.Do(req)
	if err != nil {
		return wrapWithErr(errpkg.CodeCallEZReceiptAPI, err)
	}
	defer func() {
		if err := rawResp.Body.Close(); err != nil {
			logger.Logger.Err(err).Msg("ezreceipt login close resp body error")
		}
	}()

	var resp struct {
		ezResp
		Token struct {
			Token   string `json:"token"`
			ValidTo int64  `json:"validTo"`
		} `json:"token"`
	}
	if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
		return wrapWithErr(errpkg.CodeJSONDecode, err)
	}

	if reqErr := resp.ezResp.Err(); reqErr != nil {
		return wrapWithErr(errpkg.CodeCallEZReceiptAPI, err)
	}

	if resp.Token.Token != "" {
		ez.SetToken(resp.Token.Token, resp.Token.ValidTo)
	}

	return nil
}

func (ez *ezreceipt) CheckMobileCode(ctx context.Context, mobileCode string) *errpkg.Error {
	const uri = "openTax/carrier/checkMobileCode"

	resp, err := sendReqWithResp[CheckMobileCodeResp](ctx, ez, uri, CheckMobileCodePayload{MobileCode: mobileCode})
	if err != nil {
		return err
	} else if resp.IsExist != 1 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeMobileBarcodeNotExist, Data: mobileCode}
	}

	return nil
}

func (ez *ezreceipt) InvNumberList(ctx context.Context) (*ListResp[InvNumberRespItem], *errpkg.Error) {
	const uri = "eInvoice/invNumber/list"

	var payload struct{}

	return sendReqWithResp[ListResp[InvNumberRespItem]](ctx, ez, uri, payload)
}

func (ez *ezreceipt) ConfirmOrder(ctx context.Context, orderID uint64, payload ConfirmOrderPayload) (*ConfirmOrderRespItem, *errpkg.Error) {
	const uri = "sales/order/confirm"

	uriWithOrderID, err := makeURIWithID(uri, orderID)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeParseURL, err)
	}

	return sendReqWithResp[ConfirmOrderRespItem](ctx, ez, uriWithOrderID, payload)
}

func (ez *ezreceipt) CreateOrder(ctx context.Context, payload CreateOrderPayload) (*CreateOrderRespItem, *errpkg.Error) {
	const uri = "sales/order/create"

	return sendReqWithResp[CreateOrderRespItem](ctx, ez, uri, payload)
}

func (ez *ezreceipt) SetCarrier(ctx context.Context, orderID uint64, payload SetCarrierPayload) (*SetCarrierRespItem, *errpkg.Error) {
	const uri = "sales/order/setCarrier"

	uriWithOrderID, err := makeURIWithID(uri, orderID)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeParseURL, err)
	}

	return sendReqWithResp[SetCarrierRespItem](ctx, ez, uriWithOrderID, payload)
}

func (ez *ezreceipt) CreateB2CInvoice(ctx context.Context, orderID uint64, payload CreateB2CInvoicePayload) (*CreateB2CInvoiceRespItem, *errpkg.Error) {
	const uri = "eInvoice/invoice/createB2C"

	uriWithOrderID, err := makeURIWithID(uri, orderID)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeParseURL, err)
	}

	return sendReqWithResp[CreateB2CInvoiceRespItem](ctx, ez, uriWithOrderID, payload)
}

func (ez *ezreceipt) FastCreateB2CInvoice(ctx context.Context, payload FastCreateB2CPayload) (*CreateB2CInvoiceRespItem, *errpkg.Error) {
	if err := payload.Validate(); err != nil {
		return nil, wrapWithErr(errpkg.CodeBadParam, err)
	}

	if err := ez.CleanupExistsOrder(ctx, payload.ID); err != nil {
		return nil, err
	}

	order, wrapErr := ez.CreateOrder(ctx, CreateOrderPayload{
		OrderNo:  payload.ID,
		ProdList: []ProdItem{payload.ProdItem},
	})
	if wrapErr != nil {
		return nil, wrapErr
	}

	needRollback := false
	defer func() {
		if !needRollback {
			return
		}

		if err := ez.VoidOrder(ctx, order.OrderID); err != nil {
			logger.Logger.Err(err).Msg("ezreceipt fast create b2c invoice void order error")
		}
	}()

	if _, wrapErr := ez.ConfirmOrder(ctx, order.OrderID, ConfirmOrderPayload{}); wrapErr != nil {
		needRollback = true
		return nil, wrapErr
	}

	setCarrierPayload := SetCarrierPayload{
		CarrierType: 10,
	}

	if payload.CarrierNum != "" {
		setCarrierPayload.CarrierType = 2
		setCarrierPayload.CarrierInfo = payload.CarrierNum
	}

	if _, wrapErr := ez.SetCarrier(ctx, order.OrderID, setCarrierPayload); wrapErr != nil {
		needRollback = true
		return nil, wrapErr
	}

	invoice, wrapErr := ez.CreateB2CInvoice(ctx, order.OrderID, CreateB2CInvoicePayload{
		AutoInvNo: true,
		Title:     payload.Title,
	})
	if wrapErr != nil {
		needRollback = true
		return nil, wrapErr
	}

	return invoice, nil
}

func (ez *ezreceipt) VoidOrder(ctx context.Context, orderID uint64) *errpkg.Error {
	const uri = "sales/order/void"

	uriWithOrderID, err := makeURIWithID(uri, orderID)
	if err != nil {
		return wrapWithErr(errpkg.CodeParseURL, err)
	}

	if _, reqWrapErr := ez.sendReq(ctx, uriWithOrderID, nil); reqWrapErr != nil {
		logger.Logger.Err(reqWrapErr).Str("service", "receipt").Msg("send req err")
		return reqWrapErr
	}

	return nil
}

func (ez *ezreceipt) CleanupExistsOrder(ctx context.Context, orderNo string) *errpkg.Error {
	const uri = "sales/order/list"

	payload := struct {
		Prop      string `json:"prop"`
		PropValue string `json:"propValue"`
	}{
		Prop:      "orderNo",
		PropValue: orderNo,
	}

	items, wrapErr := sendReqWithResp[ListResp[OrderListItem]](ctx, ez, uri, payload)
	if wrapErr != nil {
		return wrapErr
	}

	if items.Entries == 0 {
		return nil
	}

	if items.Entries > 1 {
		return wrapWithErr(errpkg.CodeBadParam, fmt.Errorf("more than one order with order no: %s", orderNo))
	}

	return ez.VoidOrder(ctx, items.List[0].OrderID)
}

func (ez *ezreceipt) SetToken(token string, validTo int64) {
	ez.token.token = token
	ez.token.validTo = time.UnixMilli(validTo)
}

func (ez *ezreceipt) Token() string {
	return ez.token.token
}

func (ez *ezreceipt) isTokenValid(t time.Time) bool {
	return ez.token.token != "" && ez.token.validTo.After(t.Add(1*time.Minute))
}

// loginPwd returns the login password used to login to EZReceipt API.
//
//	雜湊規則為 `sha1(sha1($accName) + $password)`，其中
//
//	`$accName` 代表使用者帳號，`$password` 代表使用者密碼的明文
//
// Ref: ezreceipt 開發者資源 > 串接說明 > 串接說明
//
// #nosec G401: the login api of ezreceipt requires sha1
func loginPwd(accName, password string) (string, error) {
	if accName == "" || password == "" {
		return "", errors.New("app code or app key is empty")
	}

	hashAcc := fmt.Sprintf("%x", sha1.Sum([]byte(accName)))
	hashPwd := sha1.Sum(append([]byte(hashAcc), []byte(password)...))

	return fmt.Sprintf("%x", hashPwd), nil
}

func wrapWithErr(code errpkg.Code, err error) *errpkg.Error {
	return &errpkg.Error{
		HttpStatus: http.StatusInternalServerError,
		Code:       code,
		Err:        err,
	}
}

func (ez *ezreceipt) sendReq(ctx context.Context, uri string, payload any) (*ezResp, *errpkg.Error) {
	req, reqWrapErr := ez.newReq(ctx, uri, payload)
	if reqWrapErr != nil {
		return nil, reqWrapErr
	}

	rawResp, err := ez.client.Do(req)
	if err != nil {
		return nil, wrapWithErr(errpkg.CodeCallEZReceiptAPI, err)
	}
	defer func() {
		if err := rawResp.Body.Close(); err != nil {
			logger.Logger.Err(err).
				Str("uri", uri).
				Msg("ezreceipt close resp body error")
		}
	}()

	var resp *ezResp
	if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
		return nil, wrapWithErr(errpkg.CodeJSONDecode, err)
	}

	if err := resp.Err(); err != nil {
		return nil, wrapWithErr(errpkg.CodeCallEZReceiptAPI, err)
	}

	return resp, nil
}

// sendReqWithResp sends a request to EZReceipt API and
// returns the response value.
// The T is the type of the response value.
func sendReqWithResp[T any](ctx context.Context, ez *ezreceipt, uri string, payload any) (*T, *errpkg.Error) {
	resp, reqWrapErr := ez.sendReq(ctx, uri, payload)
	if reqWrapErr != nil {
		logger.Logger.Err(reqWrapErr).Str("service", "receipt").Msg("send req err")
		return nil, reqWrapErr
	}

	var value *T
	if err := json.Unmarshal(resp.Value, &value); err != nil {
		logger.Logger.Err(err).Str("service", "receipt").Any("resp", value).Msg("unmarshal json err")
		return nil, wrapWithErr(errpkg.CodeJSONUnmarshal, err)
	}

	return value, nil
}

func makeURIWithID(uri string, id uint64) (string, error) {
	return url.JoinPath(uri, strconv.FormatUint(id, 10))
}
