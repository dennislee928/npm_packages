package wallet

import (
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

const (
	testBaseURL = "https://sofatest.sandbox.cybavo.com"
	prodBaseURL = "https://multisig.aegiscustody.com"
)

const (
	// HeaderChecksum is the header name of checksum
	HeaderChecksum = "X-CHECKSUM"

	// HeaderAPICode is the header name of API code
	HeaderAPICode = "X-API-CODE"
)

type CallbackType int

const (
	CallbackTypeUnknown CallbackType = iota
	CallbackTypeDeposit
	CallbackTypeWithdraw
	CallbackTypeCollect
	CallbackTypeAirdrop
)

type CallbackState int64

const (
	CallbackStateHolding       CallbackState = 1  // Processing batch in KMS (1)
	CallbackStateInPool        CallbackState = 2  // KMS process done, TXID created (2)
	CallbackStateInChain       CallbackState = 3  // TXID in chain (3)
	CallbackStateFailed        CallbackState = 5  // Failed (5)
	CallbackStateCancelled     CallbackState = 8  // cancelled
	CallbackStateDropped       CallbackState = 10 // Dropped
	CallbackStateInChainFailed CallbackState = 11 // Transaction Failed (11)
)

type ProcessingState int8

const (
	ProcessingStateFailed ProcessingState = -1 + iota
	ProcessingStateInPool
	ProcessingStateInChain
	ProcessingStateDone
)

type CallbackInfo struct {
	Type              CallbackType           `json:"type"`
	Serial            int64                  `json:"serial"`
	OrderIDWithPrefix string                 `json:"order_id"`
	Currency          string                 `json:"currency"`
	TXID              string                 `json:"txid"`
	BlockHeight       int64                  `json:"block_height"`
	TIndex            int                    `json:"tindex"`
	VOutIndex         int                    `json:"vout_index"`
	Amount            string                 `json:"amount"`
	Fees              string                 `json:"fees"`
	Memo              string                 `json:"memo"`
	BroadcastAt       int64                  `json:"broadcast_at"`
	ChainAt           int64                  `json:"chain_at"`
	FromAddress       string                 `json:"from_address"`
	ToAddress         string                 `json:"to_address"`
	WalletID          int64                  `json:"wallet_id"`
	State             CallbackState          `json:"state"`
	ConfirmBlocks     int64                  `json:"confirm_blocks"`
	ProcessingState   ProcessingState        `json:"processing_state"`
	Addon             map[string]interface{} `json:"addon"`
	Decimals          int32                  `json:"decimal"`

	wallet       wallet
	currencyType CurrencyType
	mainnet      Mainnet
}

func (c CallbackInfo) AmountDecimal() (decimal.Decimal, error) {
	if c.Amount == "" {
		return decimal.Zero, nil
	}

	val, err := decimal.NewFromString(c.Amount)
	if err != nil {
		return decimal.Zero, err
	}

	if c.Decimals == 0 {
		return val, nil
	}

	return val.Shift(-c.Decimals), nil
}

func (c CallbackInfo) FeeDecmial() (decimal.Decimal, error) {
	if c.Fees == "" {
		return decimal.Zero, nil
	}

	val, err := decimal.NewFromString(c.Fees)
	if err != nil {
		return decimal.Zero, err
	}

	if c.Addon == nil {
		return decimal.Zero, errors.New("addon is nil")
	}

	feeDecimalVals, ok := c.Addon["fee_decimal"]
	if !ok {
		return decimal.Zero, errors.New("fee_decimal not found")
	}

	decimalsFloat, ok := feeDecimalVals.(float64)
	if !ok {
		return decimal.Zero, errors.New("fee_decimal is not int32")
	}

	return val.Shift(-int32(decimalsFloat)), nil
}

func (c CallbackInfo) UserID() (string, error) {
	if c.Addon == nil {
		return "", errors.New("addon is nil")
	}

	userID, ok := c.Addon["user_id"]
	if !ok {
		return "", errors.New("user_id not found")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", errors.New("user_id is not string")
	}

	return userIDStr, nil
}

func (c CallbackInfo) CurrencyType() CurrencyType {
	return c.currencyType
}

func (c CallbackInfo) Mainnet() Mainnet {
	return c.mainnet
}

func (c CallbackInfo) OrderID() string {
	return strings.TrimPrefix(c.OrderIDWithPrefix, c.wallet.orderIDPrefix)
}

type CreateDepositAddressesPayload struct {
	Count  int      `json:"count"`
	Labels []string `json:"labels,omitempty"`
}

type WithdrawReq struct {
	// OrderID: Specify an unique ID, order ID must be prefixed, max 255 chars
	OrderID string `json:"order_id"`

	// Address: Outgoing address (address must be a contract address, if
	// the contract_abi is not empty)
	Address string `json:"address"`

	// Amount: 	Withdrawal amount
	Amount string `json:"amount"`

	// ContractAbi: Specify the ABI method and the parameters,
	// in the format ABI_method:parameters
	//
	// required, if calls contract ABI
	ContractAbi string `json:"contract_abi,omitempty"`

	// Memo: Memo on blockchain (This memo will be sent to blockchain).
	// Refer to [Memo Requirement]
	//
	// [Memo Requirement]: https://www.cybavo.com/zh-tw/developers/appendix/memo-requirement/
	Memo string `json:"memo,omitempty"`

	// UserID: 	Specify certain user
	UserID string `json:"user_id,omitempty"`

	// Message: Message (This message only savced on CYBAVO, not sent
	// to blockchain)
	Message string `json:"message,omitempty"`

	// BlockAverageFee: Use average blockchain fee within latest N blocks.
	// This option does not work for XRP, XLM, BNB, DOGE, EOS, TRX, ADA,
	// DOT and SOL cryptocurrencies.
	//
	// range 1~100
	BlockAverageFee int `json:"block_average_fee,omitempty"`

	// ManualFee: Specify blockchain fee in smallest unit of wallet currency
	// (For ETH/BSC/HECO/OKT/OP/ARB/CELO/FTM/PALM, the unit is gwei.
	// The unit returned by the Query Average Fee API is wei, divided by
	// 1000000000 to get the correct unit..
	// This option does not work for XRP, XLM, BNB, DOGE, EOS, TRX, ADA,
	// DOT and SOL cryptocurrencies.
	//
	// range 1~2000
	ManualFee int `json:"manual_fee,omitempty"`

	// TokenID: Specify the token ID to be transferred
	TokenID string `json:"token_id,omitempty"`
}

type WithdrawPayload struct {
	Requests []WithdrawReq `json:"requests"`

	// IgnoreBlackList: After setting, the address check will not be performed.
	// Apply to all orders.
	IgnoreBlackList bool `json:"ignore_black_list"`
}

type VerifyAddressPayload struct {
	Addresses []string `json:"addresses"` // Specify the address for verification
}

type AddWithdrawalWhitelistEntryPayload struct {
	Items []WithdrawalWhitelistEntryItem `json:"items"` // Specify the whitelist entries
}

type WithdrawalWhitelistEntryItem struct {
	Address string `json:"address"` // requried	The outgoing address
	Memo    string `json:"memo"`    // optional	The memo of the outgoing address
	UserID  string `json:"user_id"` // optional, max length 255	The custom user ID of the outgoing address
}

type RemoveWithdrawalWhitelistEntryPayload AddWithdrawalWhitelistEntryPayload
