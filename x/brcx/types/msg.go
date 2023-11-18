package types

import (
	"encoding/json"
	"fmt"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec"
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
	authtypes "github.com/brc20-collab/brczero/libs/cosmos-sdk/x/auth/types"
	"github.com/brc20-collab/brczero/libs/tendermint/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// verify interface at compile time
var _ sdk.Msg = &MsgInscription{}

// MsgInscription - struct for create contract
type MsgInscription struct {
	Inscription        json.RawMessage    `json:"inscription" yaml:"inscription"`
	InscriptionContext InscriptionContext `json:"inscription_context" yaml:"inscriptionContext"`
}

// NewMsgUnjail creates a new MsgUnjail instance
func NewMsgCreateContract(Inscription json.RawMessage, ctx InscriptionContext) MsgInscription {
	return MsgInscription{
		Inscription:        Inscription,
		InscriptionContext: ctx,
	}
}

// nolint
func (msg MsgInscription) Route() string { return RouterKey }
func (msg MsgInscription) Type() string  { return "inscription" }
func (msg MsgInscription) GetSigners() []sdk.AccAddress {
	return nil
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgInscription) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgInscription) ValidateBasic() error {

	return nil
}

type MsgInscriptionFromOrd struct {
	Inscription        string             `json:"inscription" yaml:"inscription"`
	InscriptionContext InscriptionContext `json:"inscription_context_1" yaml:"inscriptionContext"`
}

// Decoder Try to decode as MsgInscription by json
func Decoder(_ codec.CdcAbstraction, txBytes []byte) (tx sdk.Tx, err error) {
	var brczeroTx types.BRCZeroRequestTx

	if err = rlp.DecodeBytes(txBytes, &brczeroTx); err == nil {
		// TODO
		var msgOrd MsgInscriptionFromOrd
		if err = json.Unmarshal([]byte(brczeroTx.HexRlpEncodeTx), &msgOrd); err == nil {
			msgInscription := MsgInscription{
				Inscription:        json.RawMessage(msgOrd.Inscription),
				InscriptionContext: msgOrd.InscriptionContext,
			}

			// TODO fee
			fee := authtypes.NewStdFee(200000, sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDec(100))})
			return authtypes.NewStdTx([]sdk.Msg{msgInscription}, fee, nil, ""), nil
		}
	}
	fmt.Println("3333333---", err)
	return nil, err
}
