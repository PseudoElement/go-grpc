package encryptorutils

import (
	"strconv"
)

func DecimalToHex(decimalStr string) (string, error) {
	decimalNum, err := strconv.Atoi(decimalStr)
	hex := strconv.FormatInt(int64(decimalNum), 16)
	return hex, err
}
