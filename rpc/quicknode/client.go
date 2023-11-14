package quicknode

import (
	"github.com/btcsuite/btcd/rpcclient"
)

func NewRPCClient(quickNodeRPCURL string) (*rpcclient.Client, error) {

	connCfg := &rpcclient.ConnConfig{
		Host:         quickNodeRPCURL,
		User:         "user",
		Pass:         "pass",
		HTTPPostMode: true,
		DisableTLS:   false,
	}

	client, err := rpcclient.New(connCfg, nil)

	return client, err
}
