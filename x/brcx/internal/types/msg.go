package types

import (
	"encoding/json"
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
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
