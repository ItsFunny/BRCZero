package types

import sdkerrors "github.com/brc20-collab/brczero/libs/cosmos-sdk/types/errors"

var (
	ErrNoValidatorForAddress = sdkerrors.Register(ModuleName, 1, "address is not associated with any known validator")
	ErrChainConfigNotFound   = sdkerrors.Register(ModuleName, 2, "chain configuration not found")
)
