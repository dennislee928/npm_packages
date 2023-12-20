package user

import (
	"bityacht-exchange-api-server/internal/database/sql/users"
	"time"
)

type LoginRequest struct {
	Account  string `json:"account" binding:"required,email"` // Email
	Password string `json:"password" binding:"required"`      // 密碼
}

type LoginResponse struct {
	// 兩階段認證類型(Bitwise):
	// * 0: None
	// * 1: Email
	// * 2: SMS
	// * 4: Google Authenticator
	Login2FAType users.TwoFAType `json:"login2FAType" binding:"required"`

	AccountNotVerified bool `json:"accountNotVerified,omitempty"` // 帳號尚未認證

	ID int64 `json:"id" binding:"required"` // Users ID

	OnePassKey string `json:"onePassKey,omitempty"`

	VerificationCode string `json:"verificationCode,omitempty"` // 驗證碼(Only in Debug Mode)

	TwoFactorLoginResponse
}

type TwoFactorLoginRequest struct {
	ID               int64  `json:"id" binding:"required"` // Users ID
	OnePassKey       string `json:"onePassKey" binding:"required"`
	VerificationCode string `json:"verificationCode" binding:"required"` // 驗證碼
}

type TwoFactorLoginResponse struct {
	AccessToken           string     `json:"accessToken,omitempty" binding:"required"` // JWT Token
	RefreshToken          string     `json:"refreshToken,omitempty" binding:"required"`
	RefreshTokenExpiredAt *time.Time `json:"refreshTokenExpiredAt,omitempty" binding:"required"`
}

type ForgotPasswordRequest struct {
	Account string `json:"account" binding:"required,email"` // Email
}

type VerifyResetPasswordRequest struct {
	Account          string `json:"account" binding:"required,email"`    // Email
	VerificationCode string `json:"verificationCode" binding:"required"` // 驗證碼
}

type VerifyResetPasswordResponse struct {
	Token string `json:"token" binding:"required"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required"` // 密碼
}

type RefreshTokenRequest struct {
	ID           int64  `json:"id" binding:"required"` // Users ID
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RefreshTokenResponse = LoginResponse

type RegisterRequest struct {
	LoginRequest

	InviteCode string `json:"inviteCode"`
}

type RegisterResponse struct {
	ID               int64  `json:"id" binding:"required"`      // Users ID
	VerificationCode string `json:"verificationCode,omitempty"` // 驗證碼(Only in Debug Mode)
}

type ResendEmailVerificationCodeRequest struct {
	ID int64 `json:"id" binding:"required"` // Users ID
}

type VerifyEmailRequest struct {
	ID               int64  `json:"id" binding:"required"`               // Users ID
	VerificationCode string `json:"verificationCode" binding:"required"` // 驗證碼
}
