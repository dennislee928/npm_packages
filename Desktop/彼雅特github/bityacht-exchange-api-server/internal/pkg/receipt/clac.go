package receipt

import (
	"github.com/shopspring/decimal"
)

var one = decimal.NewFromInt(1)

func CalcSalesAndTaxFromTotal(total decimal.Decimal, taxRate decimal.Decimal) (sales int64, tax int64) {
	if taxRate.LessThanOrEqual(decimal.Zero) {
		return total.IntPart(), 0
	}

	// Ref: https://www.tron-island.com/blog/fapiao-rounding-issue
	// Law: https://law.moj.gov.tw/LawClass/LawSingle.aspx?pcode=G0340081&flno=32-1
	tax = total.Div(one.Add(taxRate)).Mul(taxRate).Round(0).IntPart()
	sales = total.IntPart() - tax

	return
}
