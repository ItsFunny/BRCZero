package types

import (
	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec"
	codectypes "github.com/brc20-collab/brczero/libs/cosmos-sdk/codec/types"
	"github.com/brc20-collab/brczero/libs/ibc-go/modules/core/exported"
)

// RegisterInterfaces registers the tendermint concrete client-related
// implementations and interfaces.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*exported.ClientState)(nil),
		&ClientState{},
	)
	registry.RegisterImplementations(
		(*exported.ConsensusState)(nil),
		&ConsensusState{},
	)
	registry.RegisterImplementations(
		(*exported.Header)(nil),
		&Header{},
	)
	registry.RegisterImplementations(
		(*exported.Misbehaviour)(nil),
		&Misbehaviour{},
	)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&ClientState{}, "ibc.lightclients.tendermint.v1.ClientState", nil)
}