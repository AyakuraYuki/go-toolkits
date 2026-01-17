package decimals

import "github.com/shopspring/decimal"

func ParseString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.NewFromInt(0)
	}
	return d
}
