package exchange

var (
	Binance *binanceExchange
	Max     *maxExchange
)

func Init() {
	Binance = newBinaceExchange()
	Max = newMaxExchange()
}
