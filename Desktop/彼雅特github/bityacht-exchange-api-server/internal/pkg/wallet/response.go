package wallet

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type QueryAPICodeStatusResp struct {
	Valid       APICodeStatusItem `json:"valid"`       // The activated API code
	Inactivated APICodeStatusItem `json:"inactivated"` // Not active API code
}

type APICodeStatusItem struct {
	APICode string `json:"api_code"` // API code for querying wallet
	Exp     int64  `json:"exp"`      // API code expiration unix time in UTC
}

type CreateDepositAddressesResp struct {
	Addresses []string `json:"addresses"`
	TXIDs     []string `json:"txids"`
}

type GetContractAddressItem struct {
	Address      string `json:"address"`
	Currency     uint64 `json:"currency"`
	Memo         string `json:"memo"`
	TokenAddress string `json:"token_address"`
}

type GetContractAddressesResp struct {
	Addresses json.RawMessage `json:"addresses"`
}

type WithdrawAssetsResp struct {
	Results map[string]int64 `json:"results"`
}

type VerifyAddressItem struct {
	Address      string `json:"address"`
	Valid        bool   `json:"valid"`
	MustNeedMemo bool   `json:"must_need_memo"` // Indicate whether the address needs a memo / destination tag when transferring cryptocurrency to that address
}

type VerifyAddressResp struct {
	Result []VerifyAddressItem `json:"result"` // Array of addresses' verification result
}

type AddWithdrawalWhitelistEntryResp struct {
	AddedItems []WithdrawalWhitelistEntryItem `json:"added_items"` // Array of added whitelist entries
}

type RemoveWithdrawalWhitelistEntryResp struct {
	RemovedItems []WithdrawalWhitelistEntryItem `json:"removed_items"` // Array of removed whitelist entries
}

type ErrResp struct {
	ErrorMsg  string `json:"error"`
	ErrorCode int    `json:"error_code"`
}

func (e ErrResp) Error() string {
	return fmt.Sprintf("%s: code: %d", e.ErrorMsg, e.ErrorCode)
}

func parseErrResp(r io.Reader) error {
	var resp ErrResp
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		return err
	}
	return resp
}

func checksumVerify(apiSecret, checksum string, r io.Reader) bool {
	sha := sha256.New()

	if _, err := io.Copy(sha, r); err != nil {
		return false
	}

	if _, err := sha.Write([]byte(apiSecret)); err != nil {
		return false
	}

	sum := sha.Sum(nil)
	calcHash := base64.URLEncoding.EncodeToString(sum)
	return calcHash == checksum
}
