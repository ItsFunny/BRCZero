package types

import (
	"fmt"
	sdkerrors "github.com/brc20-collab/brczero/libs/cosmos-sdk/types/errors"
)

const (
	manageContract = 10
)

var (
	ErrChainConfigNotFound = sdkerrors.Register(ModuleName, 2, "chain configuration not found")
)

func ErrUnknownOperationOfManageContract(operation string) *sdkerrors.Error {
	return sdkerrors.New(ModuleName, manageContract+1, fmt.Sprintf("%s is unknown operation of manage contract", operation))
}

func ErrValidateBasic(msg string) *sdkerrors.Error {
	return sdkerrors.New(ModuleName, manageContract+2, fmt.Sprintf("ManageContract validateBasic error : %s", msg))
}

func ErrValidateInput(msg string) *sdkerrors.Error {
	return sdkerrors.New(ModuleName, manageContract+3, msg)
}

func ErrExecute(msg string) *sdkerrors.Error {
	return sdkerrors.New(ModuleName, manageContract+4, msg)
}
