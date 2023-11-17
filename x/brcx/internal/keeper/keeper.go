package keeper

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/brc20-collab/brczero/libs/cosmos-sdk/codec"
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
	"github.com/brc20-collab/brczero/libs/tendermint/libs/log"
	"github.com/brc20-collab/brczero/x/brcx/internal/types"
)

type Keeper struct {
	cdc *codec.CodecProxy

	storeKey sdk.StoreKey
	logger   log.Logger

	evmKeeper     EVMKeeper
	accountKeeper AccountKeeper
	bankKeeper    BankKeeper
	supplyKeeper  SupplyKeeper
}

func NewKeeper(cdc *codec.CodecProxy, storeKey sdk.StoreKey, logger log.Logger, evmKeeper EVMKeeper, accountKeeper AccountKeeper, bk BankKeeper) *Keeper {
	logger = logger.With("module", types.ModuleName)
	// ensure brcx module account is set
	if addr := sk.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	return &Keeper{cdc: cdc, logger: logger, evmKeeper: evmKeeper, accountKeeper: accountKeeper, bankKeeper: bk, supplyKeeper: sk}
}

func (k Keeper) Logger() log.Logger {
	return k.logger
}

func (k Keeper) getAminoCodec() *codec.Codec {
	return k.cdc.GetCdc()
}

func (k Keeper) GetProtoCodec() *codec.ProtoCodec {
	return k.cdc.GetProtocMarshal()
}

func (k Keeper) GetContractAddressByName(ctx sdk.Context, name string) []sdk.AccAddress {
	kvStore := ctx.KVStore(k.storeKey)
	value := kvStore.Get(types.GetContractAddressByName(name))
	var addrs []sdk.AccAddress
	k.cdc.GetCdc().MustUnmarshalBinaryLengthPrefixed(value, &addrs)

	return addrs
}

func (k Keeper) InsertContractAddressWithName(ctx sdk.Context, name string, address sdk.AccAddress) {
	kvStore := ctx.KVStore(k.storeKey)
	value := kvStore.Get(types.GetContractAddressByName(name))
	var addrs []sdk.AccAddress
	k.cdc.GetCdc().MustUnmarshalBinaryLengthPrefixed(value, &addrs)
	addrs = append(addrs, address)

	v := k.cdc.GetCdc().MustMarshalBinaryLengthPrefixed(addrs)
	kvStore.Set(types.GetContractAddressByName(name), v)
}

func (k Keeper) GetBRCXAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) GetBRCXAddress() sdk.AccAddress {
	return k.supplyKeeper.GetModuleAddress(types.ModuleName)
}

func (k Keeper) GetContractAddrByProtocol(protocol string) (common.Address, error) {
	//todo
	return [20]byte{}, nil
}
