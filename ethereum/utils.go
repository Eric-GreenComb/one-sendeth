package ethereum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/shopspring/decimal"
)

var gWei int64 = 1000000000
var gasPriceCache = big.NewInt(21 * gWei)
var cacheLock sync.RWMutex

// response of  https://ethgasstation.info/json/ethgasAPI.json
type gasAPI struct {
	Fastest float32 `json:"fastest"`
	Average float32 `json:"average"`
	SafeLow float32 `json:"safeLow"`
}

// GasPrice get gasPrice from Cache
func GasPrice() *big.Int {
	cacheLock.RLock()
	cached := gasPriceCache
	cacheLock.RUnlock()

	return cached
}

// SetGasPrice 后台管理需要的Set接口
func SetGasPrice(gasPrice *big.Int) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	gasPriceCache = gasPrice
}

func suggestGasPrice() {
	// var hex hexutil.Big
	// common.Clients.Eth.Call(&hex, "eth_gasPrice")

	// gasPrice := (*big.Int)(&hex)
	// glog.V(5).Info("next Ethereum gasPric is Wei ", gasPrice)

	// // TODO enable dymanic
	// SetGasPrice(gasPrice)
	// // 1Gwei 相当于1Ether可以发送1亿笔交易
	// // SetGasPrice(big.NewInt(1000000000))
	url := "https://ethgasstation.info/json/ethgasAPI.json"

	gasClient := http.Client{
		Timeout: time.Second * 20, // Maximum of 8 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		glog.Errorln("NewRequest" + err.Error())
		return
	}

	res, getErr := gasClient.Do(req)
	if getErr != nil {
		glog.Errorln("gasClient.Do" + getErr.Error())
		return
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		glog.Errorln("ReadAll" + readErr.Error())
		return
	}

	gas := gasAPI{}
	jsonErr := json.Unmarshal(body, &gas)
	if jsonErr != nil {
		glog.Errorln("Unmarshal" + jsonErr.Error())
		return
	}

	// Average = gWei ×　１０, if gas.Average < 1, will be 0
	gasPrice := big.NewInt(int64(gas.Average) * gWei / 10)

	// if gasPrice is broken, will set the last value.
	if gasPrice.Cmp(big.NewInt(0)) < 0 {
		glog.Errorln("Got the gasPirce is Broken!", gasPrice)
		return
	}
	glog.V(5).Info("next Ethereum gasPric is Wei ", gasPrice)
	SetGasPrice(gasPrice)
}

// FloatToWei 传入值和精度, 对于以太来说是返回amount Wei
func FloatToWei(val float64, decimals int) *big.Int {
	v1 := decimal.NewFromFloat(val)
	v2 := decimal.New(1, int32(decimals))
	v3 := v1.Mul(v2)
	// ToDo error ?
	result, _bool := new(big.Int).SetString(v3.String(), 10)
	if !_bool {
		fmt.Println("error")
	}
	fmt.Println(result)
	glog.V(7).Infoln("Sending Wei is", result)
	return result
}

// StringToWei 传入值和精度, 对于以太来说是返回amount Wei
func StringToWei(val string, decimals int) *big.Int {
	v1, _ := decimal.NewFromString(val)
	v2 := decimal.New(1, int32(decimals))
	v3 := v1.Mul(v2)
	result, _ := new(big.Int).SetString(v3.String(), 10)
	glog.V(7).Infoln("Sending Wei is", result)
	return result
}

// hex to uint
func hexToDecimal(s string) uint64 {
	bigInt := new(big.Int)
	bigInt.SetString(s, 0)
	// bigIntByte, _ := bigInt.MarshalText()
	// return string(bigIntByte)
	return bigInt.Uint64()
}
