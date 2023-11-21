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
	EventTypeBRCXProtocol            = types.EventTypeBRCXProtocol
	EventTypeManageContract          = types.EventTypeManageContract
	EventTypeEntryPoint              = types.EventTypeEntryPoint
	AttributeManageContractOperation = types.AttributeManageContractOperation

	AttributeManageContractAddress = types.AttributeManageContractAddress
	AttributeEvmOutput             = types.AttributeEvmOutput
	AttributeManageLog             = types.AttributeManageLog
	AttributeBTCTXID               = types.AttributeBTCTXID
)

var (
	RegisterCodec = types.RegisterCodec
	ModuleCdc     = types.ModuleCdc

	NewKeeper                           = keeper.NewKeeper
	NewQuerier                          = keeper.NewQuerier
	ErrUnknownOperationOfManageContract = types.ErrUnknownOperationOfManageContract
	ConvertBTCAddress                   = types.ConvertBTCAddress

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
