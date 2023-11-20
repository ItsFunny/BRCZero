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

	AttributeProtocol                = types.AttributeProtocol
	EventTypeBRCX                    = types.EventTypeBRCX
	EventTypeManageContract          = types.EventTypeManageContract
	AttributeManageContractOperation = types.AttributeManageContractOperation

	AttributeManageContractAddress = types.AttributeManageContractAddress
	AttributeManageOutput          = types.AttributeManageOutput
	AttributeManageLog             = types.AttributeManageLog
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewKeeper                           = keeper.NewKeeper
	NewQuerier                          = keeper.NewQuerier
	ErrUnknownOperationOfManageContract = types.ErrUnknownOperationOfManageContract
	ConvertBTCAddress                   = types.ConvertBTCAddress

	ErrValidateInput = types.ErrValidateInput
	ErrExecute       = types.ErrExecute
)

type (
	Keeper         = keeper.Keeper
	MsgInscription = types.MsgInscription
	ManageContract = types.ManageContract
)
