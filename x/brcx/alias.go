package brcx

import (
	"github.com/brc20-collab/brczero/x/brcx/internal/keeper"
	"github.com/brc20-collab/brczero/x/brcx/internal/types"
)

const (
	ModuleName   = types.ModuleName
	StoreKey     = types.StoreKey
	RouterKey    = types.RouterKey
	QuerierRoute = types.QuerierRoute
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewKeeper  = keeper.NewKeeper
	NewQuerier = keeper.NewQuerier
)

type (
	Keeper            = keeper.Keeper
	MsgCreateContract = types.MsgCreateContract
)
