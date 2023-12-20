package modelpkg

import (
	"fmt"
	"time"

	"bityacht-exchange-api-server/internal/pkg/rand"
)

type Type int32

const (
	TypeSwap Type = iota + 1
	TypeSpot
	TypeFlat
)

type Action int32

const (
	ActionDepositOrBuy Action = iota + 1
	ActionWithdrawOrSell
)

var monthMap = map[time.Month]string{
	time.January:   "A",
	time.February:  "B",
	time.March:     "C",
	time.April:     "D",
	time.May:       "E",
	time.June:      "F",
	time.July:      "G",
	time.August:    "H",
	time.September: "I",
	time.October:   "J",
	time.November:  "K",
	time.December:  "L",
}

func GetOrderID(t Type, a Action) string {
	now := time.Now()
	rng := rand.Intn[int](1000000)
	str := fmt.Sprintf("%1d%1d%02d%1s%02d%06d",
		t,
		a,
		now.Year()%100,
		monthMap[now.Month()],
		now.Day(),
		rng,
	)
	return str
}
