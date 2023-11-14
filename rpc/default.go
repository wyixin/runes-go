package bitcoinrpc

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/wyixin/runes-go/pkg/log16"
	"github.com/wyixin/runes-go/rpc/btcd"
	"github.com/wyixin/runes-go/rpc/quicknode"
)

var log = log16.NewLogger("module", "bitcoinrpc")

var defaultRpcClient *DefaultRPCClient

type DefaultRPCClient struct {
	*rpcclient.Client
}

func New(rpcActive, bitcoindHost, bitcoindUser, bitcoindPass, quickNodeRPCURL string) BTCCRPC {

	var rpcclient *rpcclient.Client
	if rpcActive == "bitcoind" {
		rpcclient, _ = btcd.NewRPCClient(bitcoindHost, bitcoindPass, bitcoindUser)
	} else {
		rpcclient, _ = quicknode.NewRPCClient(quickNodeRPCURL)
	}

	return &DefaultRPCClient{Client: rpcclient}
}

func (b *DefaultRPCClient) GetBlockCount(ctx context.Context) (uint32, error) {

	blockCount, err := b.Client.GetBlockCount()
	log.Info("debug", "err", err, "cnt", blockCount)
	if err != nil {
		return 0, err
	}

	return uint32(blockCount), err
}

func (b *DefaultRPCClient) GetBlock(ctx context.Context, number uint32) (block *btcjson.GetBlockVerboseTxResult, err error) {

	blockHash, err := b.Client.GetBlockHash(int64(number))
	if err != nil {
		log.Error("rpc.getblockhash failed", "err", err)
		return
	}

	// Get the current block count.
	blockResult, err := b.Client.GetBlockVerboseTx(blockHash)

	if err != nil {
		log.Error("rpc.getblock  get block hash failed", "err", err)
		return
	}

	return blockResult, nil
}

func (b *DefaultRPCClient) GetTx(ctx context.Context, txHash string) (tx *btcjson.TxRawResult, err error) {

	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		log.Error("rpc.getTx hash string decode failed", "hash", txHash)
		return nil, err
	}

	tx, err = b.Client.GetRawTransactionVerbose(hash)
	if err != nil {
		log.Error("rpc.gettransaction rpc failed", "err", err)
		return nil, err
	}

	return tx, nil
}
