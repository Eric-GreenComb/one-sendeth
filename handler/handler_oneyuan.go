package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	myhttp "github.com/Eric-GreenComb/contrib/net/http"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethcommon "github.com/ethereum/go-ethereum/common"
	nsq "github.com/nsqio/go-nsq"

	"github.com/Eric-GreenComb/one-sendeth/badger"
	"github.com/Eric-GreenComb/one-sendeth/bean"
	"github.com/Eric-GreenComb/one-sendeth/common"
	"github.com/Eric-GreenComb/one-sendeth/config"
	"github.com/Eric-GreenComb/one-sendeth/ethereum"
	"github.com/Eric-GreenComb/one-sendeth/persist"
)

// PendingNonce PendingNonce
var PendingNonce uint64

// OneYuanMessage OneYuanMessage
func OneYuanMessage(message *nsq.Message) error {

	fmt.Println("==== receive message ====")

	var _one bean.OneOrder

	err := json.Unmarshal(message.Body, &_one)
	if err != nil {
		return err
	}

	fmt.Println(">>>> order_code : ", _one.OrderCode)

	_txid, err := SendEthereumInfo(_one)
	if err != nil {
		fmt.Println("************** error : ", err.Error())
		return err
	}

	fmt.Println("==== txid : ", _txid)

	err = createOrder(_one, _txid)
	if err != nil {
		return err
	}

	_, err = PostJSON2Callback(_one.OrderCode, _txid)
	if err != nil {
		return err
	}
	fmt.Println("==== post callback ok")

	return nil
}

func createOrder(one bean.OneOrder, txid string) error {
	fmt.Println(">>>> save db")

	var _order bean.Order
	_order.OrderCode = one.OrderCode
	_order.GoodsID = one.GoodsID
	_order.GoodName = one.GoodName
	_order.Amount = one.Amount
	_order.BuyTime = one.BuyTime
	_order.UserName = one.UserName
	_order.Type = one.Type
	_order.Desc = one.Desc
	_order.TxID = txid

	err := persist.GetPersist().CreateOrder(_order)
	if err != nil {
		return err
	}

	return nil
}

// LoadAccountKey LoadAccountKey
func LoadAccountKey() (*keystore.Key, error) {

	_from := config.Ethereum.Address
	_pwd := string(common.FromHex(config.Ethereum.Passphrase))

	_value, err := badger.NewRead().Get(_from)
	if err != nil {
		return nil, err
	}

	var _keystore string
	_keystore = strings.Replace(string(_value), "\\\"", "\"", -1)
	_key, err := keystore.DecryptKey([]byte(_keystore), _pwd)
	if err != nil {
		return nil, err
	}

	fmt.Println(">>>>>>>>> account OK : ", _key.Address.String())

	return _key, nil
}

// SendEthereumInfo SendEthereumInfo
func SendEthereumInfo(order bean.OneOrder) (string, error) {

	_desc := GenEthereumInfo(order)

	_txID, err := SendEthereumCoin(_desc)
	if err != nil {
		return "", err
	}
	return _txID, nil
}

// GenEthereumInfo GenEthereumInfo
func GenEthereumInfo(order bean.OneOrder) string {
	var _buf bytes.Buffer

	switch order.Type {
	case 0:
		_buf.WriteString("订单号:")
		_buf.WriteString(order.OrderCode)
		_buf.WriteString(";商品ID:")
		_buf.WriteString(order.GoodsID)
		_buf.WriteString(";商品名称:")
		_buf.WriteString(order.GoodName)
		_buf.WriteString(";金额:")
		_buf.WriteString(order.Amount)
		_buf.WriteString(";时间:")
		_buf.WriteString(order.BuyTime)
		_buf.WriteString(";用户:")
		_buf.WriteString(order.UserName)
		_buf.WriteString(";购买编码:")
		_buf.WriteString(order.Desc)
		break
	case 1:
		_buf.WriteString("订单号:")
		_buf.WriteString(order.OrderCode)
		_buf.WriteString(";商品ID:")
		_buf.WriteString(order.GoodsID)
		_buf.WriteString(";商品名称:")
		_buf.WriteString(order.GoodName)
		_buf.WriteString(";时间:")
		_buf.WriteString(order.WinTime)
		_buf.WriteString(";备注:")
		_buf.WriteString(order.Desc)
		break
	}
	_desc := _buf.String()

	return _desc
}

// SendEthereumCoin SendEthereumCoin
func SendEthereumCoin(desc string) (string, error) {
	txID := &ethcommon.Hash{}

	_from := config.Ethereum.Address
	_to := _from

	_amountBigInt := ethereum.StringToWei("0.01", 18)
	_chainIDBigInt := big.NewInt(config.Ethereum.ChainID)

	_inputData := []byte(desc)

	_key, err := LoadAccountKey()
	if err != nil {
		fmt.Println(">>>>>>>>> load account key error : ", err.Error())
		return "", err
	}

	_txID, err := ethereum.SendEthCoins(_to, PendingNonce, _amountBigInt, _key.PrivateKey, _chainIDBigInt, _inputData)
	if err != nil {
		return txID.String(), err
	}

	PendingNonce++

	return _txID, nil
}

// PostJSON2Callback PostJSON2Callback
func PostJSON2Callback(orderCode, txID string) (string, error) {
	var _callback bean.Callback
	_callback.OrderCode = orderCode
	_callback.TxID = txID

	_postJSON, _ := json.Marshal(_callback)

	_res, err := myhttp.PostJSONString(config.Server.CallBack, string(_postJSON))
	if err != nil {
		return "", err
	}
	return _res, nil
}
