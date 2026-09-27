package app

import (
	"fmt"
	"math"
	"strconv"
)

func providerUsageMilliseconds(value any) int64 {
	seconds, err := strconv.ParseFloat(fmt.Sprint(value), 64)
	if err != nil || math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds <= 0 || seconds > 3600 {
		return 0
	}
	return int64(math.Round(seconds * 1000))
}
