package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auctiontypes "github.com/tharsis/ethermint/x/auction/types"
	"github.com/tharsis/ethermint/x/nameservice/helpers"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

func GetBlockChangeSetIndexKey(height int64) []byte {
	return append(PrefixBlockChangesetIndex, helpers.Int64ToBytes(height)...)
}

func getOrCreateBlockChangeset(store sdk.KVStore, legacyCodec codec.LegacyAmino, height int64) *types.BlockChangeset {
	bz := store.Get(GetBlockChangeSetIndexKey(height))

	if bz != nil {
		var changeset types.BlockChangeset
		legacyCodec.MustUnmarshal(bz, &changeset)
		return &changeset
	}

	return &types.BlockChangeset{
		Height:      height,
		Records:     []string{},
		Names:       []string{},
		Auctions:    []string{},
		AuctionBids: []auctiontypes.AuctionBidInfo{},
	}
}

func updateBlockChangeSetForAuction(ctx sdk.Context, k RecordKeeper, id string) {
	changeSet := getOrCreateBlockChangeset(ctx.KVStore(k.storeKey), k.legacyCodec, ctx.BlockHeight())

	found := false
	for _, elem := range changeSet.Auctions {
		if id == elem {
			found = true
			break
		}
	}

	if !found {
		changeSet.Auctions = append(changeSet.Auctions, id)
		saveBlockChangeSet(ctx.KVStore(k.storeKey), k.legacyCodec, changeSet)
	}
}

func saveBlockChangeSet(store sdk.KVStore, codec codec.LegacyAmino, changeset *types.BlockChangeset) {
	bz := codec.MustMarshal(*changeset)
	store.Set(GetBlockChangeSetIndexKey(changeset.Height), bz)
}

func (k Keeper) saveBlockChangeSet(ctx sdk.Context, changeSet *types.BlockChangeset) {
	saveBlockChangeSet(ctx.KVStore(k.storeKey), k.legacyCodec, changeSet)
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
		Height:      height,
		Records:     []string{},
		Names:       []string{},
		Auctions:    []string{},
		AuctionBids: []auctiontypes.AuctionBidInfo{},
	}
}

func updateBlockChangeSetForAuctionBid(ctx sdk.Context, k RecordKeeper, id, bidderAddress string) {
	changeSet := getOrCreateBlockChangeset(ctx.KVStore(k.storeKey), k.legacyCodec, ctx.BlockHeight())
	changeSet.AuctionBids = append(changeSet.AuctionBids, auctiontypes.AuctionBidInfo{AuctionID: id, BidderAddress: bidderAddress})
	saveBlockChangeSet(ctx.KVStore(k.storeKey), k.legacyCodec, changeSet)
}

func updateBlockChangeSetForNameAuthority(ctx sdk.Context, store sdk.KVStore, legacyCodec codec.LegacyAmino, name string) {
	changeset := getOrCreateBlockChangeset(store, legacyCodec, ctx.BlockHeight())
	changeset.NameAuthorities = append(changeset.NameAuthorities, name)
	saveBlockChangeSet(store, legacyCodec, changeset)
}
