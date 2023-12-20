package sms

// ErrorType is the common error type for Service Providers
type ErrorType int

const (
	ErrorTypeNone ErrorType = iota + 1
	ErrorTypeBadConfig
	ErrorTypeBadAccountStatus
	ErrorTypeBadRequest
	ErrorTypeProviderSystemError
	ErrorTypeTemporaryFailure
	ErrorTypeTimeout
	ErrorTypeUnknown
)
