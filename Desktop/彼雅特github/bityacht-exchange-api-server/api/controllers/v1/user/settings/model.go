package settings

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type UpdateMobileBarcodeRequest struct {
	MobileBarcode string `json:"mobileBarcode"`
}

type GetWithdraw2FAInfoRequest struct {
	VerificationCode string `form:"verificationCode" binding:"len=6"` // 驗證碼
}

type GetWithdraw2FAInfoResponse struct {
	Token  string `json:"token" binding:"required"`
	Secret string `json:"secret,omitempty"`
	QRCode []byte `json:"qrCode,omitempty"`
}

type UpdateWithdraw2FAInfoRequest struct {
	Token            string `json:"token"`
	VerificationCode string `form:"verificationCode" binding:"len=6"` // 驗證碼
}

//! Deprecated (Meeting at 2023/10/2)
// type UpdateLogin2FATypeRequest struct {
// 	Login2FAType users.Login2FAType `json:"login2FAType" binding:"gte=0,lte=1"`
// }
