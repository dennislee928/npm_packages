package admin

type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken        string `json:"accessToken" binding:"required"` // JWT Token
	NeedChangePassword bool   `json:"needChangePassword,omitempty"`
	//* Maybe RefreshToken string
}

type ForgotPasswordRequest struct {
	Account string `json:"account" binding:"required,email"` // Email
}

type AuthResponse struct {
	IsValid bool `json:"isValid" binding:"required"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password" binding:"required"`
}
