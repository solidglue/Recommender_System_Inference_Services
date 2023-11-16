package utils

import (
	"fmt"
	"strconv"
)

func FloatRound(f float32, n int) float64 {
	format := "%." + strconv.Itoa(n) + "f"
	res, _ := strconv.ParseFloat(fmt.Sprintf(format, f), 64)
	return res
}
