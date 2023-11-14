package utils

import (
	"encoding/hex"
	"encoding/json"
	"math/big"

	"strings"

	"github.com/btcsuite/btcd/txscript"
)

func PrettyPrint(st interface{}) {
	stJson, err := json.Marshal(st)
	if err != nil {
		panic(err)
	}
	println(string(stJson))
}

func PrettyStr(st interface{}) string {
	stJson, err := json.Marshal(st)
	if err != nil {
		panic(err)
	}
	return string(stJson)
}

func Str2Hex(str string) []byte {
	res, _ := hex.DecodeString(str)
	return res
}

// Reverse the map that can be used to output a human-readable name by code.
func OpcodeByCode(c byte) string {
	//	originalMap := txscript.OpcodeByName
	//	codeNameMap = make(map[int8]string, len(originalMap))

	for k, v := range txscript.OpcodeByName {
		if v == c {
			return k
		}
	}

	return "OP_UNKNOWN"
}

func HexToBase26(hexStr string) string {
	hexNum := new(big.Int)
	hexNum.SetString(hexStr, 16)
	var base26StrBuilder strings.Builder

	base26Chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	for hexNum.Cmp(big.NewInt(0)) > 0 {
		remainder := new(big.Int)
		remainder.Mod(hexNum, big.NewInt(26))
		hexNum.Div(hexNum, big.NewInt(26))
		base26StrBuilder.WriteByte(base26Chars[remainder.Int64()])
	}

	return base26StrBuilder.String()
}
