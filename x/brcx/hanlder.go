package brcx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"

	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
	sdkerrors "github.com/brc20-collab/brczero/libs/cosmos-sdk/types/errors"
	"github.com/brc20-collab/brczero/x/brcx/types"
	"github.com/ethereum/go-ethereum/common"
)

var (
	contractJson = ``
)

// NewHandler creates a sdk.Handler for all the slashing type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx.SetEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgInscription:
			return handleInscription(ctx, msg, k)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}
	}
}

func handleInscription(ctx sdk.Context, msg MsgInscription, k Keeper) (*sdk.Result, error) {
	inscription := make(map[string]interface{})
	err := json.Unmarshal([]byte(msg.Inscription), &inscription)
	if err != nil {
		return &sdk.Result{}, ErrValidateInput("msg Inscription json marshal failed")
	}
	p, ok := inscription["p"]
	if !ok {
		return &sdk.Result{}, ErrValidateInput("can not anaylize protocol")
	}
	protocol, ok := p.(string)
	if !ok {
		return &sdk.Result{}, ErrValidateInput("the type of protocol must be string")
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			EventTypeBRCX,
			sdk.NewAttribute(AttributeProtocol, protocol),
		),
	)
	switch protocol {
	case ManageContractProtocolName:
		result, err := handleManageContract(ctx, msg, k)
		if err != nil {
			return result, err
		}
		result.Events = append(result.Events, ctx.EventManager().Events()...)
		return result, nil
	default:
		return handleBRCX(ctx, msg, protocol, k)
	}
}

func handleManageContract(ctx sdk.Context, msg MsgInscription, k Keeper) (*sdk.Result, error) {
	from, err := ConvertBTCAddress(msg.InscriptionContext.Sender)
	if err != nil {
		return nil, ErrValidateInput(fmt.Sprintf("InscriptionContext.Sender %s is not address: %s ", msg.InscriptionContext.Sender, err))
	}

	var manageContract ManageContract
	if err := json.Unmarshal([]byte(msg.Inscription), &manageContract); err != nil {
		return nil, ErrValidateInput(fmt.Sprintf("Inscription json unmarshal failed: %s ", err))
	}

	if err := manageContract.ValidateBasic(); err != nil {
		return nil, err
	}
	calldata, err := manageContract.GetCallData()
	if err != nil {
		return nil, ErrValidateInput(fmt.Sprintf("Inscription data is not hex: %s ", err))
	}
	manageContractEvent := sdk.NewEvent(EventTypeManageContract, sdk.NewAttribute(AttributeManageContractOperation, manageContract.Operation))
	var result sdk.Result
	switch manageContract.Operation {
	case ManageCreateContract:
		executeResult, contractResult, err := k.CallEvm(ctx, common.BytesToAddress(from[:]), nil, common.Big0, calldata)
		if err != nil {
			return nil, ErrExecute(fmt.Sprintf("create contract failed: %s", err))
		}
		result = *executeResult.Result
		k.InsertContractAddressWithName(ctx, manageContract.Name, contractResult.ContractAddress.Bytes())
		manageContractEvent = manageContractEvent.AppendAttributes(
			sdk.NewAttribute(AttributeManageContractAddress, contractResult.ContractAddress.Hex()),
			sdk.NewAttribute(AttributeManageOutput, hex.EncodeToString(contractResult.Ret)))
	case ManageCallContract:
		to := common.HexToAddress(manageContract.Contract)
		executeResult, contractResult, err := k.CallEvm(ctx, common.BytesToAddress(from[:]), &to, common.Big0, calldata)
		if err != nil {
			return nil, fmt.Errorf("create contract failed: %v", err)
		}
		manageContractEvent = manageContractEvent.AppendAttributes(
			sdk.NewAttribute(AttributeManageOutput, hex.EncodeToString(contractResult.Ret)),
		)
		result = *executeResult.Result
	default:
		return nil, ErrUnknownOperationOfManageContract(manageContract.Operation)
	}

	ctx.EventManager().EmitEvent(manageContractEvent)
	return &result, nil
}

func handleBRCX(ctx sdk.Context, msg MsgInscription, protocol string, k Keeper) (*sdk.Result, error) {
	to, err := k.GetContractAddrByProtocol(ctx, protocol)
	inscriptionBytes, err := json.Marshal(msg.Inscription)
	if err != nil {
		return &sdk.Result{}, err
	}

	input, err := types.GetEntryPointInput(msg.InscriptionContext, string(inscriptionBytes))
	executionResult, _, err := k.CallEvm(ctx, common.BytesToAddress(k.GetBRCXAddress().Bytes()), &to, big.NewInt(0), input)
	if err != nil {
		return nil, err
	}
	return executionResult.Result, nil
}

type CompiledContract struct {
	ABI abi.ABI
	Bin string
}
