package infura

import (
	sdk "github.com/brc20-collab/brczero/libs/cosmos-sdk/types"
)

// BeginBlocker runs the logic of BeginBlocker with version 0.
// BeginBlocker resets keeper cache.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	if !k.stream.enable {
		return
	}
	k.stream.cache.Reset()
}
