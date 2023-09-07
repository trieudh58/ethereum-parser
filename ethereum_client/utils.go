package ethereum_client

import (
	"strconv"
)

func baseDecimalToHexString(decimalNumber int64) string {
	return strconv.FormatInt(decimalNumber, 16)
}

func hexStringToDecimal(hexString string) (int64, error) {
	return strconv.ParseInt(hexString, 16, 64)
}
