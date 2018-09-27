package ethereum

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PendingNonce TODO nonce 理论上非并发安全
func PendingNonce(address string) (uint64, error) {
	var result hexutil.Uint64
	err := Clients.Eth.Call(&result, "eth_getTransactionCount", address, "pending")
	return uint64(result), err
}

// GetBalance GetBalance
func GetBalance(address string) (string, error) {
	var result string
	err := Clients.Eth.Call(&result, "eth_getBalance", address, "pending")
	return result, err
}
