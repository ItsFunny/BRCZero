package types

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
)

type InscriptionContext struct {
	InscriptionId     string `json:"inscription_id" yaml:"inscription_id"`         // 铭文id
	InscriptionNumber int64  `json:"inscription_number" yaml:"inscription_number"` // 铭文编号
	IsTransfer        bool   `json:"is_transfer" yaml:"is_transfer"`               // 是否为转移铭文
	Txid              Hash   `json:"txid" yaml:"txid"`                             // 铭文的所在交易的txid
	Sender            string `json:"sender" yaml:"sender"`                         // 铭文的发送者
	Receiver          string `json:"receiver" yaml:"receiver"`                     // 铭文接收者
	//todo: recover type
	//CommitInput       []TxIn   `json:"commit_input" yaml:"commit_input"`             //提交交易的输入
	//RevealOutput      []TxOut  `json:"reveal_output" yaml:"reveal_output"`           //揭示交易的输出
	//OldSatPoint       SatPoint `json:"old_sat_point" yaml:"old_sat_point"`           // 转移前，铭文所在的聪
	//NewSatPoint       SatPoint `json:"new_sat_point" yaml:"new_sat_point"`           // 转移后，铭文所在的聪
	CommitInput  string `json:"commit_input" yaml:"commit_input"`   //提交交易的输入
	RevealOutput string `json:"reveal_output" yaml:"reveal_output"` //揭示交易的输出
	OldSatPoint  string `json:"old_sat_point" yaml:"old_sat_point"` // 转移前，铭文所在的聪
	NewSatPoint  string `json:"new_sat_point" yaml:"new_sat_point"` // 转移后，铭文所在的聪

	BlockHash   Hash   `json:"block_hash" yaml:"block_hash"`     // btc block hash
	BlockTime   uint32 `json:"block_time" yaml:"block_time"`     // btc block time
	BlockHeight uint64 `json:"block_height" yaml:"block_height"` // btc block height
}

// Sat
type SatPoint struct {
	Outpoint Outpoint `json:"outpoint" yaml:"outpoint"` // btc的utxo
	Offset   uint64   `json:"offset" yaml:"offset"`     // 聪数在utxo里的偏移量
}

// UTXO
type Outpoint struct {
	Txid Hash   `json:"txid" yaml:"txid"` //btc 交易id
	Vout uint32 `json:"vout" yaml:"vout"` // 第几个输出
}

// 交易输出
type TxOut struct {
	Txid     Hash   `json:"txid" yaml:"txid"`
	Vout     uint32 `json:"vout" yaml:"vout"`
	Value    uint64 `json:"value" yaml:"value"`
	PkScript []byte `json:"pk_script" json:"pkScript"`
}

type TxIn struct {
	PrevTxOut   TxOut  `json:"prev_tx_out" yaml:"prevTxOut"`
	CurrentTxid Hash   `json:"current_txid" yaml:"currentTxid"`
	IndexIn     uint32 `json:"index_in" yaml:"indexIn"`
}

type Hash []byte

func (h Hash) ValidateBasic() error {
	if len(h) != common.HashLength {
		return fmt.Errorf("%s error length 32", hex.EncodeToString(h))
	}
	return nil
}

type MsgInscriptionFromOrd struct {
	Inscription        string                `json:"inscription" yaml:"inscription"`
	InscriptionContext InscriptionContextOrd `json:"inscription_context" yaml:"inscription_context"`
}

type InscriptionContextOrd struct {
	Txid              string `json:"txid"`
	InscriptionId     string `json:"inscription_id"`
	InscriptionNumber int64  `json:"inscription_number"`
	OldSatpoint       string `json:"old_satpoint"`
	NewSatpoint       string `json:"new_satpoint"`
	From              struct {
		Address string `json:"Address"`
	} `json:"from"`
	To struct {
		Address string `json:"Address"`
	} `json:"to,omitempty"`
	SatInOutputs bool   `json:"sat_in_outputs"`
	IsTransfer   bool   `json:"is_transfer"`
	Blockheight  uint64 `json:"blockheight"`
	Blocktime    uint32 `json:"blocktime"`
	Blockhash    string `json:"blockhash"`
}

func (i InscriptionContextOrd) ConvertToInscriptionContext() (InscriptionContext, error) {
	txid, err := hex.DecodeString(i.Txid)
	if err != nil {
		return InscriptionContext{}, err
	}
	blockhash, err := hex.DecodeString(i.Blockhash)
	if err != nil {
		return InscriptionContext{}, err
	}
	return InscriptionContext{
		InscriptionId:     i.InscriptionId,
		InscriptionNumber: i.InscriptionNumber,
		IsTransfer:        i.IsTransfer,
		Txid:              txid,
		Sender:            i.From.Address,
		Receiver:          i.To.Address,
		CommitInput:       "",
		RevealOutput:      "",
		OldSatPoint:       i.OldSatpoint,
		NewSatPoint:       i.NewSatpoint,

		BlockHash:   blockhash,
		BlockTime:   i.Blocktime,
		BlockHeight: i.Blockheight,
	}, nil
}
