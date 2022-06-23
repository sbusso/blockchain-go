package blockchain

import (
	"math/big"

	"github.com/shopspring/decimal"

	"encoding/hex"
	"strings"
)

// To Decimal
func StringToDecimal(str string) decimal.Decimal {
	if len(str) >= 2 && str[:2] == "0x" {
		b := new(big.Int)
		b.SetString(str[2:], 16)
		d := decimal.NewFromBigInt(b, 0)
		return d
	} else {
		v, err := decimal.NewFromString(str)
		if err != nil {
			panic(err)
		}
		return v
	}
}

func Bytes2Hex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func Hex2Bytes(str string) []byte {
	if strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X") {
		str = str[2:]
	}

	if len(str)%2 == 1 {
		str = "0" + str
	}

	h, _ := hex.DecodeString(str)
	return h
}

// with prefix '0x'
func Bytes2HexP(bytes []byte) string {
	return "0x" + hex.EncodeToString(bytes)
}
