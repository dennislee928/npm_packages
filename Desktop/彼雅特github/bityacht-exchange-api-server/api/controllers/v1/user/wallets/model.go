package wallets

import (
	"bityacht-exchange-api-server/internal/database/sql/levellimits"
	"bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	"bityacht-exchange-api-server/internal/database/sql/userswithdrawalwhitelist"
	"bityacht-exchange-api-server/internal/pkg/wallet"

	"github.com/shopspring/decimal"
)

type GetDepositAddressRequest struct {
	// Currency Type
	// * 1: BTC
	// * 2: ETH
	// * 3: USDC
	// * 4: USDT
	CurrencyType wallet.CurrencyType `form:"currencyType" json:"currencyType" binding:"required,gte=1,lte=4"`

	// Mainnet
	// * 1: BTC
	// * 2: ETC
	// * 3: ERC20
	// * 4: TRC20
	Mainnet wallet.Mainnet `form:"mainnet" json:"mainnet" binding:"required,gte=1,lte=4"`
}

type GetDepositAddressResponse struct {
	Address    string `json:"address"`
	QRCodeData []byte `json:"qrCodeData" swaggertype:"string"`
}

type WithdrawRequest struct {
	// Currency Type
	// * 1: BTC
	// * 2: ETH
	// * 3: USDC
	// * 4: USDT
	CurrencyType wallet.CurrencyType `json:"currencyType" binding:"required,gte=1,lte=4"`

	// Mainnet
	// * 1: BTC
	// * 2: ETC
	// * 3: ERC20
	// * 4: TRC20
	Mainnet wallet.Mainnet `json:"mainnet" binding:"required,gte=1,lte=4"`

	WhitelistID int64 `json:"whitelistID" binding:"required"`

	Amount decimal.Decimal `json:"amount" binding:"required"`
}

type Withdraw2FARequest struct {
	OnePassKey            string `json:"onePassKey" binding:"required"`
	EmailVerificationCode string `json:"emailVerificationCode"`
	SMSVerificationCode   string `json:"smsVerificationCode"`
	GAVerificationCode    string `json:"gaVerificationCode"`
}

type WithdrawResponse struct {
	OnePassKey string `json:"onePassKey,omitempty"`

	// 兩階段認證類型(Bitwise):
	// * 0: None
	// * 1: Email
	// * 2: SMS
	// * 4: Google Authenticator
	TwoFAType users.TwoFAType `json:"twoFAType,omitempty"`

	EmailVerificationCode string `json:"emailVerificationCode,omitempty"` // 驗證碼(Only in Debug Mode)
	SMSVerificationCode   string `json:"smsVerificationCode,omitempty"`   // 驗證碼(Only in Debug Mode)
}

type GetCoinInfoRequest struct {
	// Currency Type
	// * 1: BTC
	// * 2: ETH
	// * 3: USDC
	// * 4: USDT
	CurrencyType wallet.CurrencyType `form:"currencyType" json:"currencyType" binding:"required,gte=1,lte=4"`

	// Mainnet
	// * 1: BTC
	// * 2: ETC
	// * 3: ERC20
	// * 4: TRC20
	Mainnet wallet.Mainnet `form:"mainnet" json:"mainnet" binding:"required,gte=1,lte=4"`
}

type GetCoinInfoResponse struct {
	Name        string `json:"name"`
	Network     string `json:"network"`
	WithdrawFee string `json:"withdrawFee"`
	WithdrawMax string `json:"withdrawMax"`
	WithdrawMin string `json:"withdrawMin"`

	usersspottransfers.AccWithdrawValuation

	// 等級限制
	LevelLimit levellimits.LimitByAction `json:"levelLimit"`
}

type GetWithdrawalWhitelistRequest struct {
	// 主網
	// * 1: BTC
	// * 2: ETH
	// * 3: ERC20
	// * 4: TRC20
	Mainnet wallet.Mainnet `form:"mainnet" binding:"required,gte=1,lte=4"`
}

type CreateWithdrawalWhitelistRequest struct {
	// 主網
	// * 1: BTC
	// * 2: ETH
	// * 3: ERC20
	// * 4: TRC20
	Mainnet wallet.Mainnet `json:"mainnet" binding:"required,gte=1,lte=4"`

	// 地址
	Address string `json:"address" binding:"required"`

	// 其他資訊
	Extra userswithdrawalwhitelist.Extra `json:"extra"`
}

func (cwwr CreateWithdrawalWhitelistRequest) ToModel(usersID int64) userswithdrawalwhitelist.Model {
	return userswithdrawalwhitelist.Model{
		UsersID: usersID,
		Mainnet: cwwr.Mainnet.BinanceNetwork(),
		Address: cwwr.Address,
		Extra:   cwwr.Extra,
	}
}

type DeleteWithdrawalWhitelistRequest struct {
	ID int64 `uri:"ID" binding:"required"`
}
