package brcx

import (
	"encoding/json"
	"fmt"
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
	sdkerrors "github.com/brc20-collab/brczero/libs/cosmos-sdk/types/errors"
	"github.com/brc20-collab/brczero/x/brcx/internal/types"
)

// NewHandler creates an sdk.Handler for all the slashing type messages
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
	case "brczero":
		return handleManageContract(ctx, msg, protocol, inscription, k)
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

func handleManageContract(ctx sdk.Context, msg MsgInscription, protocol string, inscription map[string]interface{}, k Keeper) (*sdk.Result, error) {

	return nil, nil
}

func handleBRCX(ctx sdk.Context, msg MsgInscription, protocol string, k Keeper) (*sdk.Result, error) {
	return nil, nil
}
