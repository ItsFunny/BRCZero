package types

import ethcmn "github.com/ethereum/go-ethereum/common"

type BRCZeroRequestTx struct {
	HexRlpEncodeTx string `json:"hex_rlp_encode_tx"`
	BTCFee         uint64 `json:"btc_fee"`
}

type encodeTxData struct {
	AccountNonce uint64          `json:"nonce"`
	Price        string          `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *ethcmn.Address `json:"to" rlp:"nil"` // nil means contract creation
	Amount       string          `json:"value"`
	Payload      []byte          `json:"input"`

	// signature values
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`

	BTCFee string `json:"btc_fee"`
}
