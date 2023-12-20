package email

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/smtp"

	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	"github.com/jordan-wright/email"
)

// ISender is the interface of Sender
type ISender interface {
	SendMail(*email.Email) *errpkg.Error
}

// SenderType represent the type of Sender
type SenderType int

const (
	SenderTypeSMTP SenderType = iota + 1
	//* Maybe sendmail command
)

// SMTPSender for sending E-Mail
type SMTPSender struct {
	addr      string
	auth      smtp.Auth
	tlsConfig *tls.Config
}

func newSMTPSender(emailConfig configs.EmailConfig) *SMTPSender {
	s := &SMTPSender{
		addr: net.JoinHostPort(emailConfig.Host, emailConfig.Port),
		auth: smtp.PlainAuth("", emailConfig.Account, emailConfig.Password, emailConfig.Host),
	}

	//! Only For Developing.
	if emailConfig.IsFake() {
		s.auth = LoginAuth(emailConfig.Account, emailConfig.Password)
	}

	if emailConfig.SSL || emailConfig.TLS {
		s.tlsConfig = &tls.Config{
			ServerName: emailConfig.Host,
			MinVersion: tls.VersionTLS12,
		}
	}

	return s
}

func (s *SMTPSender) SendMail(mail *email.Email) *errpkg.Error {
	var err error

	if s.tlsConfig == nil {
		err = mail.Send(s.addr, s.auth)
	} else if configs.Config.Email.TLS {
		err = mail.SendWithStartTLS(s.addr, s.auth, s.tlsConfig)
	} else {
		err = mail.SendWithTLS(s.addr, s.auth, s.tlsConfig)
	}

	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSendEmail, Err: err}
	}

	return nil
}
