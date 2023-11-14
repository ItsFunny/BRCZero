package brcx

import (
	"github.com/brc20-collab/brczero/x/brcx/internal/keeper"
	"github.com/brc20-collab/brczero/x/brcx/internal/types"
)

const (
	ModuleName = types.ModuleName
	StoreKey   = types.StoreKey
	RouterKey  = types.RouterKey
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	Keeper = keeper.Keeper
)
