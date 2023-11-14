package btcd

import (
	"github.com/btcsuite/btcd/rpcclient"
)

func NewRPCClient(bitcoindHost, bitcoindPass, bitcoindUser string) (*rpcclient.Client, error) {

	connCfg := &rpcclient.ConnConfig{
		Host:         bitcoindHost,
		User:         bitcoindUser,
		Pass:         bitcoindPass,
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	//	logger := btclog.NewBackend(os.Stdout).Logger("RPC")
	//	logger.SetLevel(btclog.LevelTrace)
	//	rpcclient.UseLogger(logger)

	client, err := rpcclient.New(connCfg, nil)

	return client, err
}
