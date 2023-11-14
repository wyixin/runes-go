package bitcoinrpc

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
)

// bitcoin core rpc
type BTCCRPC interface {
	//	NewRPCClient(*rpcclient.ConnConfig) (*rpcclient.Client, error)
	//	NewRPCClient() (*rpcclient.Client, error)
	//	NewBTCRPCClient(*rpcclient.Client) (*BTCCRPC, error)
	GetBlockCount(context.Context) (uint32, error)
	GetBlock(context.Context, uint32) (*btcjson.GetBlockVerboseTxResult, error)
	GetTx(context.Context, string) (*btcjson.TxRawResult, error)
}
