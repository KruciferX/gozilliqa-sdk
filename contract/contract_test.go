/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package contract

import (
	"encoding/json"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/Zilliqa/gozilliqa-sdk/account"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
)

func TestContract_Deploy(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	host := "https://dev-api.zilliqa.com/"
	privateKey := "e19d05c5452598e24caad4a0d85a49146f7be089515c905ae6a19e8a578a6930"
	chainID := 333
	msgVersion := 1

	publickKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	address := keytools.GetAddressFromPublic(publickKey)
	pubkey := util.EncodeHex(publickKey)
	provider := provider2.NewProvider(host)

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(privateKey)

	code, _ := ioutil.ReadFile("./fungible.scilla")
	init := []Value{
		{
			"_scilla_version",
			"Uint32",
			"0",
		},
		{
			"owner",
			"ByStr20",
			"0x" + address,
		},
		{
			"total_tokens",
			"Uint128",
			"1000000000",
		},
		{
			"decimals",
			"Uint32",
			"0",
		},
		{
			"name",
			"String",
			"BobCoin",
		},
		{
			"symbol",
			"String",
			"BOB",
		},
	}
	contract := Contract{
		Code:     string(code),
		Init:     init,
		Signer:   wallet,
		Provider: provider,
	}

	result, _ := provider.GetBalance(address)

	nonce, _ := result.Result.(map[string]interface{})["nonce"].(json.Number).Int64()

	result, _ = provider.GetMinimumGasPrice()
	gasPrice := result.Result.(string)

	deployParams := DeployParams{
		Version:      strconv.FormatInt(int64(util.Pack(chainID, msgVersion)), 10),
		Nonce:        strconv.FormatInt(nonce+1, 10),
		GasPrice:     gasPrice,
		GasLimit:     "10000",
		SenderPubKey: pubkey,
	}

	tx, err := contract.Deploy(deployParams)
	assert.Nil(t, err, err)
	tx.Confirm(tx.ID, 1000, 10, provider)
	assert.True(t, tx.Status == transaction.Confirmed)
}

func TestContract_Call(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	host := "https://dev-api.zilliqa.com/"
	privateKey := "e19d05c5452598e24caad4a0d85a49146f7be089515c905ae6a19e8a578a6930"
	chainID := 333
	msgVersion := 1

	publickKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	address := keytools.GetAddressFromPublic(publickKey)
	pubkey := util.EncodeHex(publickKey)
	provider := provider2.NewProvider(host)

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(privateKey)

	contract := Contract{
		Address:  "bd7198209529dC42320db4bC8508880BcD22a9f2",
		Signer:   wallet,
		Provider: provider,
	}

	args := []Value{
		{
			"to",
			"ByStr20",
			"0x" + address,
		},
		{
			"tokens",
			"Uint128",
			"10",
		},
	}

	res, err := provider.GetBalance("9bfec715a6bd658fcb62b0f8cc9bfa2ade71434a")
	assert.Nil(t, err, err)
	nonce, _ := res.Result.(map[string]interface{})["nonce"].(json.Number).Int64()
	n := nonce + 1
	result, _ := provider.GetMinimumGasPrice()
	gasPrice := result.Result.(string)

	params := CallParams{
		Nonce:        strconv.FormatInt(n, 10),
		Version:      strconv.FormatInt(int64(util.Pack(chainID, msgVersion)), 10),
		GasPrice:     gasPrice,
		GasLimit:     "1000",
		SenderPubKey: pubkey,
		Amount:       "0",
	}

	tx, err2 := contract.Call("Transfer", args, params, true)
	assert.Nil(t, err2, err2)
	tx.Confirm(tx.ID, 1000, 3, provider)
	assert.True(t, tx.Status == transaction.Confirmed)
}

func TestContract_Sign(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	host := "https://dev-api.zilliqa.com/"
	privateKey := "e19d05c5452598e24caad4a0d85a49146f7be089515c905ae6a19e8a578a6930"
	chainID := 333
	msgVersion := 1

	publickKey := keytools.GetPublicKeyFromPrivateKey(util.DecodeHex(privateKey), true)
	pubkey := util.EncodeHex(publickKey)
	provider := provider2.NewProvider(host)

	wallet := account.NewWallet()
	wallet.AddByPrivateKey(privateKey)

	contract := Contract{
		Address:  "84eb5C96Bec8d29eDdFBe36865E9B7F26b816f0F",
		Signer:   wallet,
		Provider: provider,
	}

	args := []Value{
		{
			"proxyTokenContract",
			"ByStr20",
			"0x39550ab45d74cce5fef70e857c1326b2d9bee096",
		},
		{
			"to",
			"ByStr20",
			"0x39550ab45d74cce5fef70e857c1326b2d9bee096",
		},
		{
			"value",
			"Uint128",
			"10000000",
		},
	}

	result, _ := provider.GetBalance("9bfec715a6bd658fcb62b0f8cc9bfa2ade71434a")

	nonce, _ := result.Result.(map[string]interface{})["nonce"].(json.Number).Int64()
	n := nonce + 1

	result, _ = provider.GetMinimumGasPrice()
	gasPrice := result.Result.(string)

	params := CallParams{
		Nonce:        strconv.FormatInt(n, 10),
		Version:      strconv.FormatInt(int64(util.Pack(chainID, msgVersion)), 10),
		GasPrice:     gasPrice,
		GasLimit:     "1000",
		SenderPubKey: pubkey,
		Amount:       "0",
	}

	err, tx := contract.Sign("SubmitCustomMintTransaction", args, params, true)
	assert.Nil(t, err, err)


	pl := tx.ToTransactionPayload()
	j, _ := pl.ToJson()

	_, err2 := provider2.NewFromJson(j)
	assert.Nil(t, err2, err2)

}
