package testdata

import (
	amino "github.com/tendermint/go-amino"

	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec/types"
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
	txmsg "github.com/brc20-collab/brczero/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/types/msgservice"
)

func NewTestInterfaceRegistry() types.InterfaceRegistry {
	registry := types.NewInterfaceRegistry()
	RegisterInterfaces(registry)
	return registry
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &TestMsg{})

	registry.RegisterInterface("Animal", (*Animal)(nil))
	registry.RegisterImplementations(
		(*Animal)(nil),
		&Dog{},
		&Cat{},
	)
	registry.RegisterImplementations(
		(*HasAnimalI)(nil),
		&HasAnimal{},
	)
	registry.RegisterImplementations(
		(*HasHasAnimalI)(nil),
		&HasHasAnimal{},
	)
	registry.RegisterImplementations(
		(*txmsg.Msg)(nil),
		&MsgCreateDog{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func NewTestAmino() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*Animal)(nil), nil)
	cdc.RegisterConcrete(&Dog{}, "testdata/Dog", nil)
	cdc.RegisterConcrete(&Cat{}, "testdata/Cat", nil)

	cdc.RegisterInterface((*HasAnimalI)(nil), nil)
	cdc.RegisterConcrete(&HasAnimal{}, "testdata/HasAnimal", nil)

	cdc.RegisterInterface((*HasHasAnimalI)(nil), nil)
	cdc.RegisterConcrete(&HasHasAnimal{}, "testdata/HasHasAnimal", nil)

	return cdc
}
