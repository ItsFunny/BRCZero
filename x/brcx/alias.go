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
	EventTypeEntryPoint              = types.EventTypeEntryPoint
	AttributeManageContractOperation = types.AttributeManageContractOperation

	AttributeManageContractAddress = types.AttributeManageContractAddress
	AttributeEvmOutput             = types.AttributeEvmOutput
	AttributeManageLog             = types.AttributeManageLog
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewKeeper                           = keeper.NewKeeper
	NewQuerier                          = keeper.NewQuerier
	ErrUnknownOperationOfManageContract = types.ErrUnknownOperationOfManageContract
	ConvertBTCPKScript                  = types.ConvertBTCPKScript

	ErrInternal           = types.ErrInternal
	ErrValidateInput      = types.ErrValidateInput
	ErrExecute            = types.ErrExecute
	ErrGetContractAddress = types.ErrGetContractAddress
	ErrCallEntryPoint     = types.ErrCallEntryPoint
	ErrPackInput          = types.ErrPackInput
)

type (
	Keeper         = keeper.Keeper
	MsgInscription = types.MsgInscription
	ManageContract = types.ManageContract
)
