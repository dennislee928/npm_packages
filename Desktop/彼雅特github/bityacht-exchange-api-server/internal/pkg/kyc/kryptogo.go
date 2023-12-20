package kyc

import (
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/database/sql/countries"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type kryptoGO struct {
	client  *http.Client
	baseURL string
}

func newKryptoGO() *kryptoGO {
	kg := &kryptoGO{
		baseURL: "https://external-api.kryptogo.com",
	}

	// ref: https://www.loginradius.com/blog/engineering/tune-the-go-http-client-for-high-performance/
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 10
	transport.MaxIdleConnsPerHost = 10

	kg.client = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return kg
}

func (kg *kryptoGO) callAPI(ctx context.Context, method string, targetURI string, body any, resp any, compactedRawResp *bytes.Buffer) *errpkg.Error {
	apiURL, err := url.JoinPath(kg.baseURL, targetURI)
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeParseURL, Err: err}
	}

	bodyBuffer := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(bodyBuffer).Encode(body); err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONMarshal, Err: err}
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL, bodyBuffer)
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNewHTTPRequest, Err: err}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("GOFACT-API-TOKEN", configs.Config.KYC.APIToken)

	rawResp, err := kg.client.Do(req)
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeDoHTTPRequest, Err: err}
	}
	defer func() {
		if err := rawResp.Body.Close(); err != nil {
			logger.Logger.Err(err).Msg("kyc kryptoGo close resp body error!")
		}
	}()
	errLogger := logger.Logger.Err(err).Str("GOFACT-API-TOKEN", req.Header.Get("GOFACT-API-TOKEN")).Str("api", fmt.Sprintf("%s %q", method, apiURL)).Int("httpStatus", rawResp.StatusCode)

	if rawResp.StatusCode < http.StatusOK || rawResp.StatusCode >= http.StatusMultipleChoices { // Not 200 Series
		err = errors.New("bad status code")

		if bytes, err := io.ReadAll(rawResp.Body); err != nil {
			errLogger.Str("respBody", err.Error())
		} else {
			errLogger.Str("respBody", string(bytes))
		}

		errLogger.Send()
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallKryptoGOAPI, Err: err}
	}

	// Read Resp
	if compactedRawResp == nil {
		if err = json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
			return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONUnmarshal, Err: err}
		}

		return nil
	}

	// Read Raw Resp
	if rawRespBody, err := io.ReadAll(rawResp.Body); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeReadHTTPResponse, Err: err}
	} else if err = json.Compact(compactedRawResp, rawRespBody); err != nil {
		errLogger.Str("rawRespBody", string(rawRespBody)).Send()
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONUnmarshal, Err: err}
	} else if err = json.Unmarshal(compactedRawResp.Bytes(), &resp); err != nil {
		errLogger.Str("rawRespBody", compactedRawResp.String()).Send()
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeJSONUnmarshal, Err: err}
	}

	return nil
}

// // GET /kya/{symbol}/address/{address}
// // Get blockchain address risk information
// // https://www.kryptogo.com/docs/kyabc#tag/DD-KYA/paths/~1kya~1%7Bsymbol%7D~1address~1%7Baddress%7D/get

type CreateTasksParams struct {
	Name              string
	Birthday          time.Time
	Citizenship       string
	CallBackURL       string
	CustomerReference string
	FromIDVID         int64
}

// POST /task
// Create search tasks
// https://www.kryptogo.com/docs/kyabc#tag/DD-KYB-and-KYC/paths/~1task/post
func (kg *kryptoGO) CreateTasks(ctx context.Context, tasks []CreateTasksParams) (CreateTasksResponse, *errpkg.Error) {
	const method = http.MethodPost
	const path = "/task"

	var resp CreateTasksResponse

	req := make(CreateTasksRequest, 0, len(tasks))
	for _, t := range tasks {
		task := Task{
			SearchSettingID: configs.Config.KYC.DDTaskSearchSettingID,
			Target: searchTaskKycTarget{
				Type:        1,
				Name:        t.Name,
				Birthday:    t.Birthday.In(modelpkg.DefaultTimeLoc).Format(time.DateOnly),
				Citizenship: t.Citizenship,
			},
			SearchSource:      []int32{0, 1, 2, 3, 4},
			CallBackURL:       t.CallBackURL,
			CustomerReference: t.CustomerReference,
			FromIDVID:         t.FromIDVID,
		}
		req = append(req, task)
	}

	if err := kg.callAPI(ctx, method, path, req, &resp, nil); err != nil {
		return resp, err
	} else if len(resp) != len(req) {
		return resp, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeCallKryptoGOAPI, Err: errors.New("bad resp length")}
	}
	return resp, nil
}

// GET /task/{task_id}/summary
// Get information of a search task
// https://www.kryptogo.com/docs/kyabc#tag/DD-KYB-and-KYC/paths/~1task~1%7Btask_id%7D~1summary/get
func (kg *kryptoGO) GetTaskSummary(ctx context.Context, taskID string, rawResp *bytes.Buffer) (TaskSummaryResponse, *errpkg.Error) {
	const method = http.MethodGet

	var (
		path = "/task/" + taskID + "/summary"
		resp TaskSummaryResponse
	)

	if err := kg.callAPI(ctx, method, path, nil, &resp, rawResp); err != nil {
		return resp, err
	}
	return resp, nil
}

// POST /task/{task_id}/accepted
// Update a search task's accepted status
// https://www.kryptogo.com/docs/kyabc#tag/DD-KYB-and-KYC/paths/~1task~1%7Btask_id%7D~1accepted/post
func (kg *kryptoGO) UpdateTaskStatus(ctx context.Context, taskID string, accepted bool, comment string) (UpdateTaskStatusResponse, *errpkg.Error) {
	const method = http.MethodPost

	var (
		path = "/task/" + taskID + "/accepted"
		req  = UpdateTaskStatusRequest{Accepted: accepted, Comment: comment}
		resp UpdateTaskStatusResponse
	)
	if err := kg.callAPI(ctx, method, path, req, &resp, nil); err != nil {
		return resp, err
	}
	return resp, nil
}

// // PUT /task/{task_id}/metadata
// // Set metadata to the search task
// // https://www.kryptogo.com/docs/kyabc#tag/DD-KYB-and-KYC/paths/~1task~1%7Btask_id%7D~1metadata/put

// // GET /tasks
// // Get search task history
// // https://www.kryptogo.com/docs/kyabc#tag/DD-KYB-and-KYC/paths/~1tasks/get

// // POST /idv
// // Create an IDV task
// // https://www.kryptogo.com/docs/kyabc#tag/ID-Verification/paths/~1idv/post

// // GET /idv/{idv_id}
// // Get IDV task detail
// // https://www.kryptogo.com/docs/kyabc#tag/ID-Verification/paths/~1idv~1%7Bidv_id%7D/get

// POST /idv/init
// Initiating an IDV transaction for Web Verification
// https://www.kryptogo.com/docs/kyabc#tag/ID-Verification/paths/~1idv~1init/post
func (kg *kryptoGO) InitIDV(ctx context.Context, country countries.Country, usersID int64, idvsID int64, ddsID int64, name string, birthDate modelpkg.Date, nationalID string, rawResp *bytes.Buffer) (InitIDVResponse, *errpkg.Error) {
	const method = http.MethodPost
	const path = "/idv/init"

	var (
		resp InitIDVResponse
		req  = InitIDVRequest{
			IDType:                IDTypeIDCard,
			Locale:                country.Locale,
			WorkflowID:            WorkflowIDAndIdentityCameraAndUpload,
			SuccessURL:            configs.Config.KYC.SuccessURL,
			ErrorURL:              configs.Config.KYC.ErrorURL,
			Country:               country.Code,
			ExpectedName:          name,
			ExpectedBirthday:      birthDate.Format(time.DateOnly),
			ExpectedIDNumber:      nationalID,
			CallbackURL:           GetKryptoGOIDVCallbackURL(usersID, idvsID),
			CustomerReference:     fmt.Sprintf("idvs_%020d", idvsID),
			AutoCreateDDTask:      true,
			DDTaskCallbackURL:     GetKryptoGODDCallbackURL(usersID, ddsID),
			DDTaskSearchSettingID: configs.Config.KYC.DDTaskSearchSettingID,
		}
	)

	if err := kg.callAPI(ctx, method, path, req, &resp, rawResp); err != nil {
		return resp, err
	}
	return resp, nil
}

func GetKryptoGOIDVCallbackURL(usersID int64, idvsID int64) string {
	if configs.Config.KYC.CallbackHost == "" {
		return ""
	}

	u, err := url.JoinPath(configs.Config.KYC.CallbackHost, "/api/v1/krypto-go/idv-callback/", strconv.FormatInt(usersID, 10), "/", strconv.FormatInt(idvsID, 10))
	if err != nil {
		logger.Logger.Err(err).Str("service", "KryptoGO").Msg("get idv callback url failed")
		return ""
	}

	return u
}

func GetKryptoGODDCallbackURL(usersID int64, ddsID int64) string {
	if configs.Config.KYC.CallbackHost == "" {
		return ""
	}

	u, err := url.JoinPath(configs.Config.KYC.CallbackHost, "/api/v1/krypto-go/dd-callback/", strconv.FormatInt(usersID, 10), "/", strconv.FormatInt(ddsID, 10))
	if err != nil {
		logger.Logger.Err(err).Str("service", "KryptoGO").Msg("get dd callback url failed")
		return ""
	}

	return u
}

func ParseKryptoGOUTCTime(utcTime string) (time.Time, error) {
	return time.Parse(time.RFC3339, utcTime)
}

func ParseKryptoGOTimestamp(unixTime int64) time.Time {
	if unixTime == 0 {
		return time.Time{}
	}
	return time.Unix(unixTime, 0)
}
