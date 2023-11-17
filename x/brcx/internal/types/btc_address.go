package types

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/ethereum/go-ethereum/common"
)

func ConvertBTCPKScript(pkScript []byte) (common.Address, error) {
	from := make([]byte, 0)
	scriptClass, addrs, num, err := txscript.ExtractPkScriptAddrs(pkScript, &chaincfg.MainNetParams)
	if err != nil {
		return common.Address{}, fmt.Errorf("commit input from script is error:%v", err)
	} else if num != 1 {
		return common.Address{}, fmt.Errorf("commit input from script is error: num is %d must be 1", num)

	} else {
		switch scriptClass {
		case txscript.PubKeyTy:
			from = btcutil.Hash160(addrs[0].ScriptAddress())
		case txscript.PubKeyHashTy:
			from = addrs[0].ScriptAddress()
		case txscript.WitnessV0PubKeyHashTy:
			from = addrs[0].ScriptAddress()
		default:
			return common.Address{}, fmt.Errorf("%s can not support type of address", scriptClass.String())
		}
	}

	return common.BytesToAddress(from), nil
}
