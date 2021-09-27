package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/nameservice/helpers"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

func GetBlockChangeSetIndexKey(height int64) []byte {
	return append(PrefixBlockChangesetIndex, helpers.Int64ToBytes(height)...)
}

func (k Keeper) saveBlockChangeSet(ctx sdk.Context, changeSet *types.BlockChangeset) {
	store := ctx.KVStore(k.storeKey)
	bz := k.legacyCodec.MustMarshal(changeSet)
	store.Set(GetBlockChangeSetIndexKey(changeSet.Height), bz)
}

func (k Keeper) updateBlockChangeSetForRecord(ctx sdk.Context, id string) {
	changeSet := k.getOrCreateBlockChangeSet(ctx, ctx.BlockHeight())
	changeSet.Records = append(changeSet.Records, id)
	k.saveBlockChangeSet(ctx, changeSet)
}

func (k Keeper) getOrCreateBlockChangeSet(ctx sdk.Context, height int64) *types.BlockChangeset {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(GetBlockChangeSetIndexKey(height))

	if bz != nil {
		var changeSet types.BlockChangeset
		err := k.legacyCodec.Unmarshal(bz, &changeSet)
		if err != nil {
			return nil
		}
		return &changeSet
	}

	return &types.BlockChangeset{
		Height:   height,
		Records:  []string{},
		Names:    []string{},
		Auctions: []string{},
		//AuctionBids: []auction.AuctionBidInfo{},
	}
}
