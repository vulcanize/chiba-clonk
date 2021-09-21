package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

// RecordKeeper exposes the bare minimal read-only API for other modules.
type RecordKeeper struct {
	//auctionKeeper auction.Keeper
	storeKey sdk.StoreKey      // Unexposed key to access store from sdk.Context
	cdc      codec.BinaryCodec // The wire codec for binary encoding/decoding.
}

// NewRecordKeeper creates new instances of the nameservice RecordKeeper
func NewRecordKeeper(storeKey sdk.StoreKey, cdc codec.BinaryCodec) RecordKeeper {
	return RecordKeeper{
		//auctionKeeper: auctionKeeper,
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// QueryRecordsByBond - get all records for the given bond.
func (k RecordKeeper) QueryRecordsByBond(ctx sdk.Context, bondID string) []types.Record {
	var records []types.Record

	bondIDPrefix := append(PrefixBondIDToRecordsIndex, []byte(bondID)...)
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, bondIDPrefix)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		cid := itr.Key()[len(bondIDPrefix):]
		bz := store.Get(append(PrefixCIDToRecordIndex, cid...))
		if bz != nil {
			var obj types.Record
			k.cdc.MustUnmarshal(bz, &obj)
			//records = append(records, recordObjToRecord(store, k.cdc, obj))
			records = append(records, obj)
		}
	}

	return records
}
