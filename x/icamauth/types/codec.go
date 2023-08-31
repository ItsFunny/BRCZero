package types

import (
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec"
	cdctypes "github.com/brc20-collab/brczero/libs/cosmos-sdk/codec/types"
	txmsg "github.com/brc20-collab/brczero/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/types/msgservice"
)

var (
	ModuleCdc = codec.New()
	Marshal   *codec.CodecProxy
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgRegisterAccount{}, "icamauth/MsgRegisterAccount", nil)
	cdc.RegisterConcrete(MsgSubmitTx{}, "icamauth/MsgSubmitTx", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations(
		(*txmsg.Msg)(nil),
		&MsgRegisterAccount{},
		&MsgSubmitTx{},
	)
}
