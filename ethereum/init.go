package ethereum

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang/glog"

	"github.com/Eric-GreenComb/one-sendeth/config"
)

// Clients global client connection
var Clients = &Conn{}

// Conn 各种链的链接走这里
type Conn struct {
	Eth *rpc.Client
}

// Init 初始化各种链接
// TODO 维护多个连接池（每个链考虑多个节点)
func Init() error {
	fmt.Println(">>>>>>>>> start dial ethereum")
	eth, err := rpc.Dial(config.Ethereum.Host)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(">>>>>>>>> dialed ethereum")

	Clients.Eth = eth

	// initNonceAt()
	go Forever(suggestGasPrice, 3*time.Minute)

	return nil
}

// GetNonceAt 原子操作 address=0x...
func GetNonceAt(address string) uint64 {
	_nonce, err := PendingNonce(address)
	if err != nil {
		glog.Exitln(err.Error())
	}

	return atomic.AddUint64(&_nonce, 1)
}
