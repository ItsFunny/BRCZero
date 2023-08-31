package simulator

import (
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
)

type Simulator interface {
	Simulate([]sdk.Msg, sdk.CacheMultiStore) (*sdk.Result, error)
	Context() *sdk.Context
	Release()
}

var NewWasmSimulator func() Simulator
