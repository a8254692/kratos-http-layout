package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"gitlab.top.slotssprite.com/my/api-layout/party/util/xcolor"
	"log"
	"strconv"
)

// Round 四舍五入
// req: 传入小数
// n: 保留位数, n < 0，它将整数部分四舍五入到最接近的10^(-n)。
// Examples:
// Round(0.145,2) // output: 0.15
// Round(0.144,2) // output: 0.14
// Round(0.144,-2) // output: 0
// Round(123,-2) // output: 100
// Round(178,-2) // output: 200
func Round(req float64, n int32) float64 {
	float := decimal.NewFromFloat(req)
	round := float.Round(n)
	res := DecimalToFloat(round)
	return res
}

func DecimalToFloat(d decimal.Decimal) float64 {
	f, exact := d.Float64()
	if exact {
		return f
	}
	f, err := strconv.ParseFloat(d.String(), 64)
	if err != nil {
		log.Printf(xcolor.RED, fmt.Sprintf("DecimalToFloat exec error. decimal: %v, decimal.String(): %v, error: %v", d, d.String(), err))
	}
	return f
}
