package brcx

import (
	"github.com/brc20-collab/brczero/x/brcx/internal/keeper"
	"github.com/brc20-collab/brczero/x/brcx/types"
)

const (
	ModuleName                 = types.ModuleName
	StoreKey                   = types.StoreKey
	RouterKey                  = types.RouterKey
	QuerierRoute               = types.QuerierRoute
	ManageCreateContract       = types.ManageCreateContract
	ManageCallContract         = types.ManageCallContract
	ManageContractProtocolName = types.ManageContractProtocolName
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewKeeper                           = keeper.NewKeeper
	NewQuerier                          = keeper.NewQuerier
	ErrUnknownOperationOfManageContract = types.ErrUnknownOperationOfManageContract
	ConvertBTCPKScript                  = types.ConvertBTCPKScript
)

type (
	Keeper         = keeper.Keeper
	MsgInscription = types.MsgInscription
	ManageContract = types.ManageContract
)
