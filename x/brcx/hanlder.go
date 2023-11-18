package brcx

import (
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
	err := json.Unmarshal(msg.Inscription, inscription)
	if err != nil {
		return &sdk.Result{}, err
	}
	p, ok := inscription["p"]
	if !ok {
		return &sdk.Result{}, fmt.Errorf("can not anaylize protocol")
	}
	protocol, ok := p.(string)
	if !ok {
		return &sdk.Result{}, fmt.Errorf("the type of protocol must be string")
	}

	switch protocol {
	case ManageContractProtocolName:
		return handleManageContract(ctx, msg, k)
	default:
		return handleBRCX(ctx, msg, protocol, k)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, ""),
			sdk.NewAttribute(sdk.AttributeKeySender, ""),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleManageContract(ctx sdk.Context, msg MsgInscription, k Keeper) (*sdk.Result, error) {
	if len(msg.InscriptionContext.CommitInput) < 1 {
		return &sdk.Result{}, fmt.Errorf("commit input length must be more than zero")
	}
	from, err := ConvertBTCPKScript([]byte(msg.InscriptionContext.Sender))
	if err != nil {
		return &sdk.Result{}, err
	}

	var manageContract ManageContract
	if err := json.Unmarshal(msg.Inscription, &manageContract); err != nil {
		return nil, err
	}

	if err := manageContract.ValidateBasic(); err != nil {
		return &sdk.Result{}, err
	}
	calldata, err := manageContract.GetCallData()
	if err != nil {
		return &sdk.Result{}, err
	}
	var result sdk.Result
	switch manageContract.Operation {
	case ManageCreateContract:
		executeResult, contractResult, err := k.CallEvm(ctx, common.BytesToAddress(from[:]), nil, common.Big0, calldata)
		if err != nil {
			return &sdk.Result{}, fmt.Errorf("create contract failed: %v", err)
		}
		result = *executeResult.Result
		k.InsertContractAddressWithName(ctx, manageContract.Name, contractResult.ContractAddress.Bytes())
	case ManageCallContract:
		to := common.HexToAddress(manageContract.Contract)
		executeResult, _, err := k.CallEvm(ctx, common.BytesToAddress(from[:]), &to, common.Big0, calldata)
		if err != nil {
			return &sdk.Result{}, fmt.Errorf("create contract failed: %v", err)
		}
		result = *executeResult.Result
	default:
		return &sdk.Result{}, ErrUnknownOperationOfManageContract(manageContract.Operation)
	}

	return &result, nil
}

func handleBRCX(ctx sdk.Context, msg MsgInscription, protocol string, k Keeper) (*sdk.Result, error) {
	to, err := k.GetContractAddrByProtocol(ctx, protocol)
	inscriptionBytes, err := msg.Inscription.MarshalJSON()
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
