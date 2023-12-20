package email

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	"github.com/jordan-wright/email"
)

var sender ISender

func Init() {
	sender = newSMTPSender(configs.Config.Email)
}

//go:embed logo.png
var logoBytes []byte

// NewEmail with configs.Config.Email
func NewEmail(opts ...func(*email.Email)) *email.Email {
	mail := newEmail(configs.Config.Email)

	for _, opt := range opts {
		opt(mail)
	}

	return mail
}

// NewEmailWithConfig with specific config
func NewEmailWithConfig(emailConfig configs.EmailConfig, opts ...func(*email.Email)) *email.Email {
	mail := newEmail(emailConfig)

	for _, opt := range opts {
		opt(mail)
	}

	return mail
}

func WithLogo() func(*email.Email) {
	return func(mail *email.Email) {
		a, _ := mail.Attach(bytes.NewReader(logoBytes), "logo.png", "image/png")
		a.HTMLRelated = true
	}
}

func SendMail(mail *email.Email) *errpkg.Error {
	if !configs.Config.Email.Enable {
		return nil
	} else if sender == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNotInit, Err: errors.New("email sender is nil")}
	}

	return sender.SendMail(mail)
}

func getAddressOfAccount(account string) string {
	if strings.Contains(account, "@") {
		return account
	}

	return account + "@" + configs.Config.Email.Host
}

func newEmail(emailConfig configs.EmailConfig) *email.Email {
	e := email.NewEmail()

	if emailConfig.Nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", emailConfig.Nickname, getAddressOfAccount(emailConfig.Account))
	} else {
		e.From = getAddressOfAccount(emailConfig.Account)
	}

	return e
}

func IsDebug() bool {
	return !configs.Config.Email.Enable || configs.Config.Email.IsFake() // Disable or Is Fake
}
