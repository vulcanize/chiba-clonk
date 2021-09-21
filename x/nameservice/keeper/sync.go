package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

func saveBlockChangeset(ctx sdk.Context, store sdk.KVStore, codec codec.BinaryCodec, changeset *types.BlockChangeset) {
	//bz := codec.MustMarshal(*changeset)
	//store.Set(GetBlockChangesetIndexKey(changeset.Height), bz)
}

func (k Keeper) saveBlockChangeset(ctx sdk.Context, changeset *types.BlockChangeset) {
	saveBlockChangeset(ctx, ctx.KVStore(k.storeKey), k.cdc, changeset)
}
