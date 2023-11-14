package rune

import (
	"encoding/hex"
	"fmt"

	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	bitcoinrpc "github.com/wyixin/runes-go/rpc"
)

var rpcClient bitcoinrpc.BTCCRPC

func init() {

	cwd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("error getting current working drectory, msg:%s", err))
	}

	envFilePath := filepath.Join(cwd, "../.env")
	viper.SetConfigFile(envFilePath)
	if err = viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("error on reading config, msg:%s", err))
	}

	rpcActive := viper.GetString("BITCOIN_CORE_RPC_ACTIVE")
	bitcoindHost := viper.GetString("BITCOIND_HOST")
	bitcoindUser := viper.GetString("BITCOIND_USER")
	bitcoindPass := viper.GetString("BITCOIND_PASS")
	quickNodeRPCURL := viper.GetString("QUICK_NODE_RPC_URL")

	rpcClient = bitcoinrpc.New(rpcActive, bitcoindHost, bitcoindUser, bitcoindPass, quickNodeRPCURL)

}

// rune deploy
// 6a01520b0001ff00752b7d000000000aff987806010000000012
func TestRUNEDeoloy(t *testing.T) {
	pubKey := "6a01520b0001ff00752b7d000000000aff987806010000000012"
	runeTx, err := ExtractRuneDataFromScriptPubKeyHexStr(pubKey)
	if err != nil {
		t.Fatalf("parser err, %v", err)
	}

	assign := &Assignment{
		ID:     0,
		Output: 1,
		Amount: 21000000,
	}
	expectedRune := &Transaction{
		Hash:      "",
		Transfers: make(Transfers, 1),
		Issuance: &Rune{
			Symbol:   "RUNE",
			Decimals: 18,
		},
	}

	expectedRune.Transfers = append(expectedRune.Transfers, assign)

	if runeTx.Issuance.Symbol != expectedRune.Issuance.Symbol {
		t.Fatalf("rune parser not worked as expected: %v ", runeTx)
	}
}

func TestDecodeTransfers(t *testing.T) {
	hashStr := "0001ff00752b7d000000000001ff00752b7d000000000001ff00752b7d00000000"

	hashHex, _ := hex.DecodeString(hashStr)

	expectedTransfers := make(Transfers, 0)
	expectedTransfers = append(expectedTransfers, &Assignment{
		ID:     0,
		Output: 1,
		Amount: 21000000,
	})

	trans, err := DecodeTransfers(hashHex)

	if err != nil {
		t.Fatalf("decode error, msg: %v", err)
	}

	for k, v := range expectedTransfers {

		t.Log(k, v, *trans[k])

		if v.ID != trans[k].ID || v.Output != trans[k].Output || v.Amount != trans[k].Amount {
			t.Fatal("decode err")
		}
	}
	t.Log("ok")
}

func TestDecodeIssuance(t *testing.T) {
	hashStr := "ff987806010000000012"
	hashHex, _ := hex.DecodeString(hashStr)

	experctedRume := &Rune{
		Symbol:   "RUNE",
		Decimals: 18,
	}

	rune, _ := DecodeIssuance(hashHex)

	if experctedRume.Symbol != rune.Symbol || experctedRume.Decimals != rune.Decimals {
		t.Fatalf("decode err: %s, %s, %d", rune.Symbol, experctedRume.Symbol, len(rune.Symbol))
	}
	t.Log("ok")
}
