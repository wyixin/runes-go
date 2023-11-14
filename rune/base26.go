package rune

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const base = 0

func base26Encode(input []byte) (string, error) {

	toInt := new(big.Int)
	toInt.SetString(hex.EncodeToString(input), 16)
	toStr := fmt.Sprintf("%s", toInt)
	if len(toStr)%2 != 0 {
		errors.New("base26Encode - Invalid Char Amount, str:" + toStr)
	}

	var A byte = 'A'
	encoded := make([]byte, len(input))
	for i := 0; i < len(toStr)/2; i++ {
		cuts := toStr[i*2 : i*2+2]
		intCuts, _ := strconv.Atoi(cuts)
		encoded[i] = A + byte(intCuts)
	}

	return strings.TrimSpace(string(encoded)), nil
}
