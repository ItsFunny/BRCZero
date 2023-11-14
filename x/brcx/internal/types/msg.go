package types

import sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"

// verify interface at compile time
var _ sdk.Msg = &MsgCreateContract{}

// MsgCreateContract - struct for create contract
type MsgCreateContract struct {
	ValidatorAddr sdk.ValAddress `json:"address" yaml:"address"` // address of the validator operator
}

// NewMsgUnjail creates a new MsgUnjail instance
func NewMsgCreateContract(validatorAddr sdk.ValAddress) MsgCreateContract {
	return MsgCreateContract{
		ValidatorAddr: validatorAddr,
	}
}

// nolint
func (msg MsgCreateContract) Route() string { return RouterKey }
func (msg MsgCreateContract) Type() string  { return "create" }
func (msg MsgCreateContract) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ValidatorAddr)}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgCreateContract) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgCreateContract) ValidateBasic() error {

	return nil
}
