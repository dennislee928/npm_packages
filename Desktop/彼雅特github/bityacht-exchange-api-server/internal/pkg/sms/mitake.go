package sms

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
)

var _ ISMSProvider = (*mitake)(nil)

// Some description(comment) or naming in this package come from: https://github.com/minchao/go-mitake
type mitake struct {
	client          *http.Client
	smSendURL       *url.URL
	smBulkSendURL   *url.URL
	smCancelURL     *url.URL
	timeFormat      string
	sendRespPattern *regexp.Regexp
	nilClientID     string
}

type mitakeSendResponse struct {
	ClientID     string       `json:"clientid"`
	MsgID        string       `json:"msgid"`
	StatusCode   mitakeStatus `json:"statuscode"`
	AccountPoint int64        `json:"AccountPoint"`
	Duplicate    bool         `json:"Duplicate"`
	SmsPoint     int64        `json:"smsPoint"`
	Error        string       `json:"Error"`
}

type mitakeCancelResponse struct {
	MsgID      string       `json:"msgid"`
	StatusCode mitakeStatus `json:"statuscode"`
	StatusTime time.Time    `json:"statustime"`
	// SmsPoint   int64        `json:"smsPoint"`
}

type mitakeStatus string

const (
	mitakeStatusSyetemError                  mitakeStatus = "*" // 系統發生錯誤，請聯絡三竹資訊窗口人員
	mitakeStatusTemporarilyUnavailable       mitakeStatus = "a" // 簡訊發送功能暫時停止服務，請稍候再試
	mitakeStatusTemporarilyUnavailable2      mitakeStatus = "b" // 簡訊發送功能暫時停止服務，請稍候再試
	mitakeStatusUsernameRequired             mitakeStatus = "c" // 請輸入帳號
	mitakeStatusPasswordRequired             mitakeStatus = "d" // 請輸入密碼
	mitakeStatusBadUsernameOrPassword        mitakeStatus = "e" // 帳號、密碼錯誤
	mitakeStatusAccountExpired               mitakeStatus = "f" // 帳號已過期
	mitakeStatusAccountDisabled              mitakeStatus = "h" // 帳號已被停用
	mitakeStatusInvalidConnectionAddress     mitakeStatus = "k" // 無效的連線位址
	mitakeStatusOverConnectionLimit          mitakeStatus = "l" // 帳號已達到同時連線數上限
	mitakeStatusChangePasswordRequired       mitakeStatus = "m" // 必須變更密碼，在變更密碼前，無法使用簡訊發送服務
	mitakeStatusPasswordExpired              mitakeStatus = "n" // 密碼已逾期，在變更密碼前，將無法使用簡訊發送服務
	mitakeStatusPermissionDenied             mitakeStatus = "p" // 沒有權限使用外部Http程式
	mitakeStatusSystemTemporarilyUnavailable mitakeStatus = "r" // 系統暫停服務，請稍後再試
	mitakeStatusAccountingFailure            mitakeStatus = "s" // 帳務處理失敗，無法發送簡訊
	mitakeStatusSMSExpired                   mitakeStatus = "t" // 簡訊已過期
	mitakeStatusSMSBodyEmpty                 mitakeStatus = "u" // 簡訊內容不得為空白
	mitakeStatusInvalidPhoneNumber           mitakeStatus = "v" // 無效的手機號碼
	mitakeStatusOverQueryLimit               mitakeStatus = "w" // 查詢筆數超過上限
	mitakeStatusDataTooBig                   mitakeStatus = "x" // 發送檔案過大，無法發送簡訊
	mitakeStatusBadArguments                 mitakeStatus = "y" // 參數錯誤
	mitakeStatusNotFound                     mitakeStatus = "z" // 查無資料
	mitakeStatusReservationForDelivery       mitakeStatus = "0" // 預約傳送中
	mitakeStatusDeliveredToCarrier           mitakeStatus = "1" // 已送達業者
	mitakeStatusDeliveredToCarrier2          mitakeStatus = "2" // 已送達業者
	mitakeStatusDelivered                    mitakeStatus = "4" // 已送達手機
	mitakeStatusContentError                 mitakeStatus = "5" // 內容有錯誤
	mitakeStatusPhoneNumberError             mitakeStatus = "6" // 門號有錯誤
	mitakeStatusSMSDisable                   mitakeStatus = "7" // 簡訊已停用
	mitakeStatusDeliveryTimeout              mitakeStatus = "8" // 逾時無送達
	mitakeStatusReservationCanceled          mitakeStatus = "9" // 預約已取消
)

func (ms *mitakeStatus) ToErrorType() ErrorType {
	switch *ms {
	case mitakeStatusReservationForDelivery, mitakeStatusDeliveredToCarrier, mitakeStatusDeliveredToCarrier2, mitakeStatusDelivered, mitakeStatusReservationCanceled:
		return ErrorTypeNone
	case mitakeStatusSyetemError:
		return ErrorTypeProviderSystemError
	case mitakeStatusUsernameRequired, mitakeStatusPasswordRequired, mitakeStatusBadUsernameOrPassword, mitakeStatusInvalidConnectionAddress, mitakeStatusOverConnectionLimit:
		return ErrorTypeBadConfig
	case mitakeStatusAccountExpired, mitakeStatusAccountDisabled, mitakeStatusChangePasswordRequired, mitakeStatusPasswordExpired, mitakeStatusPermissionDenied, mitakeStatusAccountingFailure, mitakeStatusSMSDisable:
		return ErrorTypeBadAccountStatus
	case mitakeStatusSMSBodyEmpty, mitakeStatusInvalidPhoneNumber, mitakeStatusOverQueryLimit, mitakeStatusDataTooBig, mitakeStatusBadArguments, mitakeStatusContentError, mitakeStatusPhoneNumberError, mitakeStatusNotFound:
		return ErrorTypeBadRequest
	case mitakeStatusTemporarilyUnavailable, mitakeStatusTemporarilyUnavailable2, mitakeStatusSystemTemporarilyUnavailable:
		return ErrorTypeTemporaryFailure
	case mitakeStatusSMSExpired, mitakeStatusDeliveryTimeout:
		return ErrorTypeTimeout
	}

	return ErrorTypeUnknown
}

func newMitake() *mitake {
	m := new(mitake)

	// ref: https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if configs.Config.SMS.ConnCount > 0 && configs.Config.SMS.ConnCount <= 15 {
		transport.MaxIdleConns = configs.Config.SMS.ConnCount
		transport.MaxIdleConnsPerHost = configs.Config.SMS.ConnCount
	} else {
		transport.MaxIdleConns = 15        // maximum
		transport.MaxIdleConnsPerHost = 15 // maximum
	}

	m.client = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	hostURL, err := url.Parse(configs.Config.SMS.Host)
	if err != nil {
		logger.Logger.Panic().Err(err).Msg("bad sms host")
	}
	m.smSendURL = hostURL.JoinPath("/SmSend")
	query := m.smSendURL.Query()
	query.Set("CharsetURL", "UTF-8")
	m.smSendURL.RawQuery = query.Encode()

	m.smBulkSendURL = hostURL.JoinPath("/SmBulkSend")
	m.smCancelURL = hostURL.JoinPath("/SmCancel")
	m.timeFormat = "20060102150405" // YYYYMMDDHHMMSS
	m.nilClientID = "nil-id"
	m.sendRespPattern = regexp.MustCompile(`^\[.*\]$`)

	return m
}

func (*mitake) GetBatchLimit() int {
	return 500
}

func (m *mitake) Send(ctx context.Context, message Message) *errpkg.Error {
	// Request Parameters (in document)
	body := url.Values{}
	body.Set("username", configs.Config.SMS.Username)
	body.Set("password", configs.Config.SMS.Password)
	body.Set("dstaddr", message.Phone)
	body.Set("smbody", message.Message)
	// body.Set("objectID", "")
	// body.Set("smsPointFlag", "1") // if set 1, it will return smsPoint of this message cost.

	if message.ReceiverName != "" {
		body.Set("destname", message.ReceiverName)
	}
	if message.CallbackURL != "" {
		body.Set("response", message.CallbackURL)
	}
	if message.DeliveryTime != nil && time.Now().Before(*message.DeliveryTime) {
		body.Set("dlvtime", message.DeliveryTime.Format(m.timeFormat))
	}
	if message.ValidateTime != nil {
		if diff := time.Until(*message.ValidateTime); diff > 0 && diff < 24*time.Hour {
			body.Set("vldtime", message.DeliveryTime.Format(m.timeFormat))
		}
	}
	if message.ID != "" {
		body.Set("clientid", message.ID)
	}

	// Sending Request
	rawReqBody := body.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.smSendURL.String(), strings.NewReader(rawReqBody))
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNewHTTPRequest, Err: err}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if resp, err := m.client.Do(req); err != nil { // nolint:bodyclose // it will close by readSendResponse
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeDoHTTPRequest, Err: err}
	} else if _, err := m.readSendResponse(rawReqBody, resp); err != nil { //* Maybe need to do something by the response
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSendSMS, Err: err}
	}

	return nil
}

func (m *mitake) BatchSend(ctx context.Context, batchID int64, messages []Message) *errpkg.Error {
	// Query String Parameters (in document)
	query := url.Values{}
	query.Set("username", configs.Config.SMS.Username)
	query.Set("password", configs.Config.SMS.Password)
	query.Set("Encoding_PostIn", "UTF-8")
	query.Set("objectID", strconv.FormatInt(batchID, 10))
	// query.Set("smsPointFlag", "1") // if set 1, it will return smsPoint of this message cost.

	u := *m.smBulkSendURL
	u.RawQuery = query.Encode()

	// Request Body Data (in document)
	strMessages := make([]string, len(messages))
	for index, message := range messages {
		strMessages[index] = fmt.Sprintf("%s$$%s$$", message.ID, message.Phone)

		if message.DeliveryTime != nil {
			strMessages[index] += message.DeliveryTime.Format(m.timeFormat)
		}
		strMessages[index] += "$$"

		if message.ValidateTime != nil {
			strMessages[index] += message.ValidateTime.Format(m.timeFormat)
		}
		strMessages[index] += "$$"

		strMessages[index] += fmt.Sprintf("%s$$%s$$%s", message.ReceiverName, message.CallbackURL, strings.ReplaceAll(message.Message, "\n", string(byte(6))))
	}

	// Sending Request
	rawReqBody := strings.Join(strMessages, "\r\n")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(rawReqBody))
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNewHTTPRequest, Err: err}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if resp, err := m.client.Do(req); err != nil { // nolint:bodyclose // it will close by readSendResponse
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeDoHTTPRequest, Err: err}
	} else if _, err = m.readSendResponse(rawReqBody, resp); err != nil { //* Maybe need to do something by the response
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSendSMS, Err: err}
	}

	return nil
}

func (m *mitake) readSendResponse(rawReqBody string, resp *http.Response) ([]mitakeSendResponse, error) { // nolint:unparam // response may be use in the future
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Err(err).Msg("sms mitake readSendResponse close resp body error!")
		}
	}()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices { // 200 Series
		var responses []mitakeSendResponse
		var respPtr *mitakeSendResponse
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())

			if len(text) == 0 {
				continue
			} else if m.sendRespPattern.MatchString(text) {
				responses = append(responses, mitakeSendResponse{ClientID: text[1 : len(text)-1]})
				respPtr = &responses[len(responses)-1]
			} else {
				if respPtr == nil { // should not be happen!
					responses = append(responses, mitakeSendResponse{ClientID: m.nilClientID})
					respPtr = &responses[len(responses)-1]
				}

				splitedText := strings.Split(text, "=")
				if len(splitedText) == 2 {
					switch splitedText[0] {
					case "msgid":
						respPtr.MsgID = splitedText[1]
					case "statuscode":
						respPtr.StatusCode = mitakeStatus(splitedText[1])
					case "AccountPoint":
						respPtr.AccountPoint, _ = strconv.ParseInt(splitedText[1], 10, 64) //* Maybe catch the error
					case "Duplicate":
						respPtr.Duplicate = strings.ToUpper(splitedText[1]) == "Y"
					case "smsPoint":
						respPtr.SmsPoint, _ = strconv.ParseInt(splitedText[1], 10, 64) //* Maybe catch the error
					case "Error":
						respPtr.Error = splitedText[1]
					}
				} //* Maybe need to log the bad message?
			}
		}

		var errorList []error
		var logedErrorType int
		for _, response := range responses {
			switch errType := response.StatusCode.ToErrorType(); errType {
			case ErrorTypeBadConfig, ErrorTypeBadAccountStatus, ErrorTypeProviderSystemError, ErrorTypeTemporaryFailure:
				if code := 1 << int(errType); logedErrorType&code == 0 { // Reduce the same error
					logedErrorType |= code
					errorList = append(errorList, fmt.Errorf("status: '%+v', error: %q", response.StatusCode, response.Error))
				}
			case ErrorTypeBadRequest, ErrorTypeTimeout:
				errorList = append(errorList, fmt.Errorf("id: [%s], status: '%+v', error: %q", response.ClientID, response.StatusCode, response.Error))
			case ErrorTypeUnknown:
				errorList = append(errorList, fmt.Errorf("id: [%s], status: '%+v', error: %q", response.ClientID, response.StatusCode, response.Error))
			default: // ErrorTypeNone
				continue
			}
		}
		return responses, errors.Join(errorList...)
	} else if respBody, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("bad status code: %d, req body: %q, read resp body error: %q", resp.StatusCode, rawReqBody, err.Error())
	} else {
		return nil, fmt.Errorf("bad status code: %d, req body: %q, resp body: %q", resp.StatusCode, rawReqBody, string(respBody))
	}
}

func (m *mitake) Cancel(ctx context.Context, msgIDList []string) error {
	// Request Parameters (in document)
	body := url.Values{}
	body.Set("username", configs.Config.SMS.Username)
	body.Set("password", configs.Config.SMS.Password)

	var errorList []error
	for start := 0; start < len(msgIDList); {
		end := start + 100 // limit by provider
		if end > len(msgIDList) {
			end = len(msgIDList)
		}

		body.Set("msgid", strings.Join(msgIDList[start:end], ","))
		start = end

		// Sending Request
		rawReqBody := body.Encode()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.smCancelURL.String(), strings.NewReader(rawReqBody))
		if err != nil {
			errorList = append(errorList, err)
			continue
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		if resp, err := m.client.Do(req); err != nil { // nolint:bodyclose // it will close by readCancelResponse
			errorList = append(errorList, err)
			continue
		} else if _, err = m.readCancelResponse(rawReqBody, resp); err != nil { //* Maybe need to do something by the response
			errorList = append(errorList, err)
			continue
		}
	}

	return errors.Join(errorList...)
}

func (m *mitake) readCancelResponse(rawReqBody string, resp *http.Response) ([]mitakeCancelResponse, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.Logger.Err(err).Msg("sms mitake readCancelResponse close resp body error!")
		}
	}()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices { // 200 Series
		var responses []mitakeCancelResponse
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			text := strings.TrimSpace(scanner.Text())

			if len(text) == 0 {
				continue
			} else if splitedText := strings.Split(text, "\t"); len(splitedText) == 3 {

			} //* Maybe need to log the bad message?
		}

		return responses, nil
	} else if respBody, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("bad status code: %d, req body: %q, read resp body error: %q", resp.StatusCode, rawReqBody, err.Error())
	} else {
		return nil, fmt.Errorf("bad status code: %d, req body: %q, resp body: %q", resp.StatusCode, rawReqBody, string(respBody))
	}
}
