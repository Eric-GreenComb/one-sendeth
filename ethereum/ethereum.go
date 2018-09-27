package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/pborman/uuid"
)

// Address eth地址结构
type Address struct {
	// Wif     string `json:"address"`
	ID       string `json:"id" `
	UserName string `json:"username" form:"username"`
	Label    string `json:"label" form:"label"`
	Private  string `json:"private"`
	Public   string `json:"public"`
}

// NewAddress NewAddress
func NewAddress() *Address {
	a := &Address{}
	key, _ := crypto.GenerateKey()
	a.ID = uuid.NewRandom().String()
	// a.UserName a.Label  from post body
	a.Private = hex.EncodeToString(crypto.FromECDSA(key))
	a.Public = crypto.PubkeyToAddress(key.PublicKey).String()
	return a
}

// SendEthCoins SendEthCoins
func SendEthCoins(to string, nonce uint64, amountWei *big.Int, priv *ecdsa.PrivateKey, chainID *big.Int, data []byte) (string, error) {

	fmt.Println("----------- send eth coin")

	client := Clients.Eth

	gasPrice := GasPrice()

	fmt.Println("gas price : ", gasPrice)
	fmt.Println("nonce ", nonce)

	tx := types.NewTransaction(nonce, ethcommon.HexToAddress(to), amountWei, gasDefaultLimit, gasPrice, data)

	var signed *types.Transaction
	if chainID != nil {
		signed, _ = types.SignTx(tx, types.NewEIP155Signer(chainID), priv)
	} else {
		signed, _ = types.SignTx(tx, types.HomesteadSigner{}, priv)
	}

	data, err := rlp.EncodeToBytes(signed)
	if err != nil {
		return "", err
	}

	txID := &ethcommon.Hash{}
	err = client.Call(&txID, "eth_sendRawTransaction", ethcommon.ToHex(data))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	return txID.String(), nil
}
