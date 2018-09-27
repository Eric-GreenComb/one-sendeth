package main

import (
	// "encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	nsq "github.com/nsqio/go-nsq"

	myconfig "github.com/Eric-GreenComb/one-sendeth/config"
	"github.com/Eric-GreenComb/one-sendeth/ethereum"
	"github.com/Eric-GreenComb/one-sendeth/handler"
	"github.com/Eric-GreenComb/one-sendeth/persist"
)

// MakeConsumer MakeConsumer
func MakeConsumer(topic, channel string, config *nsq.Config,
	handle func(message *nsq.Message) error) {
	consumer, _ := nsq.NewConsumer(topic, channel, config)
	consumer.AddHandler(nsq.HandlerFunc(handle))

	// 待深入了解
	// 連線到 NSQ 叢集，而不是單個 NSQ，這樣更安全與可靠。
	// err := q.ConnectToNSQLookupd("127.0.0.1:4161")

	err := consumer.ConnectToNSQD(myconfig.Nsq.Host)
	if err != nil {
		log.Panic("Could not connect")
	}
}

func main() {

	persist.InitDatabase()

	ethereum.Init()

	_nonce, err := ethereum.PendingNonce(myconfig.Ethereum.Address)
	if err != nil {
		log.Fatal(err)
	}
	handler.PendingNonce = _nonce
	fmt.Println(">>>>>>>>>  account nonce : ", _nonce)

	_config := nsq.NewConfig()
	_config.DefaultRequeueDelay = 0
	_config.MaxBackoffDuration = 20 * time.Millisecond
	_config.LookupdPollInterval = 1000 * time.Millisecond
	_config.RDYRedistributeInterval = 1000 * time.Millisecond
	_config.MaxInFlight = 2500

	MakeConsumer("send_oneyuan_1", "ch", _config, handler.OneYuanMessage)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	fmt.Println("exit")
}
