package types

type BRCZeroRequestTx struct {
	HexRlpEncodeTx string `json:"hex_rlp_encode_tx"`
	BTCFee         uint64 `json:"btc_fee"`
}
