package ethereum

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"strings"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

const transferAbi string = "0xa9059cbb"
const lenTransferInput int = 138

var (
	gasDefaultLimit uint64 = 90000
)

// ToHexInt ToHexInt
func ToHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
}

// GenArgsAbi GenArgsAbi
func GenArgsAbi(args string) (string, bool) {
	abilen := 32 * 2
	argslen := len(args)
	if argslen > abilen {
		return "", false
	}
	str := strings.Repeat("0", abilen-argslen)
	str += args
	return str, true
}

// GenContractData GenContractData
func GenContractData(strs ...string) (string, bool) {
	str := transferAbi
	ret := true
	for _, s := range strs {
		ss, ok := GenArgsAbi(s)
		if !ok {
			ret = false
			return "", ret
		}
		str += ss
	}
	return str, ret
}

// GenTransferCode gen token transfer data
func GenTransferCode(toAddress string, amount *big.Int) ([]byte, error) {
	amountHexStr := ToHexInt(amount)
	to := toAddress
	if strings.HasPrefix(toAddress, "0x") || strings.HasPrefix(toAddress, "0X") {
		to = toAddress[2:] // 去掉目标地址0x
	}
	dataStr, ok := GenContractData(to, amountHexStr)
	if !ok {
		return nil, errors.New("param was too long")
	}

	data := ethcommon.FromHex(dataStr)

	return data, nil
}

// SendEthTokens SendEthTokens
func SendEthTokens(address, to string, nonce uint64, amount *big.Int, priv *ecdsa.PrivateKey, chainID *big.Int) (string, error) {
	client := Clients.Eth
	var err error

	data, err := GenTransferCode(to, amount)
	if err != nil {
		return "", err
	}

	gasPrice := GasPrice()

	tx := types.NewTransaction(nonce, ethcommon.HexToAddress(address), nil, gasDefaultLimit, gasPrice, data)

	var signed *types.Transaction
	if chainID != nil {
		signed, _ = types.SignTx(tx, types.NewEIP155Signer(chainID), priv)
	} else {
		signed, _ = types.SignTx(tx, types.HomesteadSigner{}, priv)
	}

	txData, err := rlp.EncodeToBytes(signed)
	if err != nil {
		return "", err
	}

	txID := &ethcommon.Hash{}
	err = client.Call(&txID, "eth_sendRawTransaction", ethcommon.ToHex(txData))
	if err != nil {
		return "", err
	}

	txid := txID.String()

	return txid, nil
}
