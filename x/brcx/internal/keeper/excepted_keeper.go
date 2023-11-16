package keeper

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
	authexported "github.com/brc20-collab/brczero/libs/cosmos-sdk/x/auth/exported"
	evmtypes "github.com/brc20-collab/brczero/x/evm/types"
)

type EVMKeeper interface {
	GetChainConfig(ctx sdk.Context) (evmtypes.ChainConfig, bool)
	GenerateCSDBParams() evmtypes.CommitStateDBParams
	GetParams(ctx sdk.Context) evmtypes.Params
	GetCallToCM() vm.CallToWasmByPrecompile
	GetBlockHash() ethcmn.Hash
	AddInnerTx(...interface{})
	AddContract(...interface{})
}

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	SetAccount(ctx sdk.Context, acc authexported.Account)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
}

type BankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
