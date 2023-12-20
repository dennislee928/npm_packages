package suspicioustransactionspkg

import (
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	"bityacht-exchange-api-server/internal/database/sql/userstransactions"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Result struct {
	Type        Type
	Information Information
}

var _ sql.Scanner = (*Information)(nil)
var _ driver.Valuer = (*Information)(nil)

type Information struct {
	LoginAt       *time.Time                                `json:"loginAt,omitempty"`
	Transactions  []userstransactions.TransactionForManager `json:"transactions,omitempty"`
	SpotTransfers []usersspottransfers.Transfer             `json:"spotTransfers,omitempty"`
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (i *Information) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("suspicious_transactions.Information: bad type")
	}

	return json.Unmarshal(bytes, i)
}

// Value return json value, implement driver.Valuer interface
func (i Information) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type Type int32

const (
	// 【多筆提領態樣】
	// 用戶在一定時間內，申請轉出多筆法定貨幣或虛擬通貨。
	// 於1小時內申請轉出法定貨幣或虛擬通貨或兩者達3筆以上者。
	TypeMultipleWithdrawal Type = iota + 1

	// 【多筆同額態樣】
	// 用戶在一定時間內，申請轉出多筆數額相同之法定貨幣或虛擬通貨。
	// 於1小時內申請轉出相同額度之法定貨幣或虛擬通貨達3筆(含)以上者。
	TypeMultipleSameAmount

	// 【同一出金地址態樣】
	// 不同帳戶使用者於一定時間內，出金至同一錢包地址。
	// 於24小時內，2個（含）以上用戶出金至同一錢包地址。
	TypeMultipleSameWithdrawalAddress

	// 【迅速轉出態樣】
	// 用戶登入後或用戶存入法定貨幣或存入、買入或接收虛擬通貨後，迅速轉出虛擬通貨或法定貨幣，一定期間內合計達特定筆數或金額，且交易內容與用戶身分顯不相當或無合理原因者。
	// 用戶登入後10分鐘內，轉出大於相當新台幣 5 萬元之虛擬通貨或法定貨幣。
	// 存入法定貨幣或存入、買入、接收虛擬通貨後，於 1 小時內(含)轉出虛擬通貨或法定貨幣達 5 筆或累積金額相當於新台幣 300 萬元。
	TypeQuicklyWithdraw

	// 【迅速買賣態樣】
	// 存入法定貨幣或存入、買入或接收虛擬通貨後，迅速買賣虛擬通貨或法定貨幣，一定期間內合計達特定筆數或金額，且交易內容與用戶身分顯不相當或無合理原因者。
	// 存入法定貨幣或存入、買入、接收虛擬通貨後，於 1 小時內(含)買賣虛擬通貨達 20 筆或累積金額相當於新台幣 300 萬元。
	TypeQuicklyTrade

	// 【小額接收大額轉出態樣】
	// 用戶小額接收多種或多筆虛擬通貨後，將虛擬通貨整筆轉出或賣出取得法定貨幣款項，且無合理原因者。
	// 該用戶於24小時之內接收超過 50 次入金，並於其中第一次入金後 24 小時內賣出取得法定貨幣或轉出虛擬通貨達相當於新台幣 50 萬元。
	TypeSmallInBigOut
)

func (t Type) Chinese() string {
	switch t {
	case TypeMultipleWithdrawal:
		return "多筆提領態樣"
	case TypeMultipleSameAmount:
		return "多筆同額態樣"
	case TypeMultipleSameWithdrawalAddress:
		return "同一出金地址態樣"
	case TypeQuicklyWithdraw:
		return "迅速轉出態樣"
	case TypeQuicklyTrade:
		return "迅速買賣態樣"
	case TypeSmallInBigOut:
		return "小額接收大額轉出態樣"
	}

	return "未知錯誤"
}

type Action int32

const (
	ActionDepositCryptocurrency Action = iota + 1
	ActionWithdrawCryptocurrency

	ActionBuySpot
	ActionSellSpot

	ActionDepositFiat
	ActionWithdrawFiat
)
