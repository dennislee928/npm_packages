package sms

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/pkg/email"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
)

var provider ISMSProvider

// ISMSProvider is the interface of SMS Provider
type ISMSProvider interface {
	GetBatchLimit() int
	Send(context.Context, Message) *errpkg.Error               //* Maybe return Common Response to do something
	BatchSend(context.Context, int64, []Message) *errpkg.Error //* Maybe return Common Response to do something
}

// Supported SMS Provider List
const (
	SMSProviderMitake = "mitake"
)

// InitSMS will init the sms service by provider in configs
func Init() {
	switch strings.ToLower(configs.Config.SMS.Provider) {
	case SMSProviderMitake:
		provider = newMitake()
	default:
		panic("bad sms provider")
	}
}

// Message for Send
type Message struct {
	ID           string // An unique identifier from client to identify SMS message
	Phone        string
	Message      string
	ReceiverName string     // Optional
	CallbackURL  string     // Optional, Callback URL to receive the delivery receipt of the message.
	DeliveryTime *time.Time // Optional
	ValidateTime *time.Time // Optional
}

// Send will send message to phone
func Send(ctx context.Context, message Message) *errpkg.Error {
	if provider == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNotInit, Err: errors.New("sms provider is nil")}
	}

	//* Maybe Create SMS History to DB
	message.ID = strconv.FormatInt(fakeMessageID.Add(1), 36)

	// Send to Service Provider
	//* Maybe Save the result to DB
	if configs.Config.SMS.Enable {
		return provider.Send(ctx, message)
	} else if email.IsDebug() {
		smsDebugMail := email.NewEmail()
		smsDebugMail.To = []string{"SMS@sms.debug"}
		smsDebugMail.Subject = "[SMS Debug]" + message.Phone
		smsDebugMail.Text = []byte(message.Message)
		email.SendMail(smsDebugMail)

		return nil
	}

	return nil
}

// SendBatch will send messages to phones, len(phones) == len(messages) is required.
func SendBatch(ctx context.Context, messages []Message) *errpkg.Error {
	if provider == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNotInit, Err: errors.New("sms provider is nil")}
	}

	batchLimit := provider.GetBatchLimit()

	//* Maybe Create SMS History to DB
	for index := 0; index < len(messages); index++ {
		messages[index].ID = strconv.FormatInt(fakeMessageID.Add(1), 36)
	}

	var errList []error
	// Send to Service Provider
	//* Maybe Save the result to DB
	errorChan := make(chan error)
	var routineCount int
	for start := 0; start < len(messages); {
		// TODO: get real batch id
		batchID := fakeBatchID.Add(1)

		end := start + batchLimit
		if end > len(messages) {
			end = len(messages)
		}

		if configs.Config.SMS.Enable {
			//* Maybe use ants to conrtol the total routine amount
			routineCount++
			go func() {
				errorChan <- provider.BatchSend(ctx, batchID, messages[start:end])
			}()
		}

		start = end
	}

	for ; routineCount > 0; routineCount-- {
		if err := <-errorChan; err != nil {
			errList = append(errList, err)
		}
	}

	if err := errors.Join(errList...); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSendSMS, Err: err}
	}

	return nil
}
