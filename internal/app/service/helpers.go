package service

import (
	"strconv"
	"strings"
)

func parsePrice(price string) int64 {
	price = strings.ReplaceAll(price, "₽", "")
	price = strings.ReplaceAll(price, " ", "")

	result, err := strconv.ParseInt(price, 10, 64)
	if err != nil {
		return 0
	}

	return result
}