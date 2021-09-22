package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

// RecordKeeper exposes the bare minimal read-only API for other modules.
type RecordKeeper struct {
	//auctionKeeper auction.Keeper
	storeKey sdk.StoreKey      // Unexposed key to access store from sdk.Context
	cdc      codec.BinaryCodec // The wire codec for binary encoding/decoding.
}

// RemoveBondToRecordIndexEntry removes the Bond ID -> [Record] index entry.
func (k Keeper) RemoveBondToRecordIndexEntry(ctx sdk.Context, bondID string, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(getBondIDToRecordsIndexKey(bondID, id))
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

// ProcessRenewRecord renews a record.
func (k Keeper) ProcessRenewRecord(ctx sdk.Context, msg types.MsgRenewRecord) error {
	if !k.HasRecord(ctx, msg.RecordId) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Record not found.")
	}

	// Check if renewal is required (i.e. expired record marked as deleted).
	record := k.GetRecord(ctx, msg.RecordId)
	if !record.Deleted || record.ExpiryTime.After(ctx.BlockTime()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Renewal not required.")
	}

	recordType := record.ToRecordType()
	err := k.processRecord(ctx, &recordType, true)
	if err != nil {
		return err
	}

	return nil
}

// ProcessAssociateBond associates a record with a bond.
func (k Keeper) ProcessAssociateBond(ctx sdk.Context, msg types.MsgAssociateBond) error {

	if !k.HasRecord(ctx, msg.RecordId) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Record not found.")
	}
	//
	//if !k.bondKeeper.HasBond(ctx, msg.BondId) {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	//}

	// Check if already associated with a bond.
	record := k.GetRecord(ctx, msg.RecordId)
	if record.BondId != "" || len(record.BondId) != 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond already exists.")
	}

	// Only the bond owner can associate a record with the bond.
	//bond := k.bondKeeper.GetBond(ctx, msg.BondID)
	//if msg.Signer.String() != bond.Owner {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	//}

	record.BondId = msg.BondId
	k.PutRecord(ctx, record)
	k.AddBondToRecordIndexEntry(ctx, msg.BondId, msg.RecordId)

	// Required so that renewal is triggered (with new bond ID) for expired records.
	if record.Deleted {
		k.InsertRecordExpiryQueue(ctx, record)
	}

	return nil
}

// ProcessDissociateBond dissociates a record from its bond.
func (k Keeper) ProcessDissociateBond(ctx sdk.Context, msg types.MsgDissociateBond) error {
	if !k.HasRecord(ctx, msg.RecordId) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Record not found.")
	}

	// Check if associated with a bond.
	record := k.GetRecord(ctx, msg.RecordId)
	bondID := record.BondId
	if bondID == "" || len(bondID) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond not found.")
	}

	// Only the bond owner can dissociate a record from the bond.
	//bond := k.bondKeeper.GetBond(ctx, bondID)
	//if msg.Signer != bond.Owner {
	//	return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	//}

	// Clear bond ID.
	record.BondId = ""
	k.PutRecord(ctx, record)
	k.RemoveBondToRecordIndexEntry(ctx, bondID, record.Id)

	return nil
}

// ProcessDissociateRecords dissociates all records associated with a given bond.
func (k Keeper) ProcessDissociateRecords(ctx sdk.Context, msg types.MsgDissociateRecords) error {

	//if !k.bondKeeper.HasBond(ctx, msg.BondID) {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Bond not found.")
	//}

	// Only the bond owner can dissociate all records from the bond.
	//bond := k.bondKeeper.GetBond(ctx, msg.BondId)
	//if msg.Signer != bond.Owner {
	//	return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond owner mismatch.")
	//}

	// Dissociate all records from the bond.
	records := k.recordKeeper.QueryRecordsByBond(ctx, msg.BondId)
	for _, record := range records {
		// Clear bond ID.
		record.BondId = ""
		k.PutRecord(ctx, record)
		k.RemoveBondToRecordIndexEntry(ctx, msg.BondId, record.Id)
	}

	return nil
}

// ProcessReAssociateRecords switches records from and old to new bond.
func (k Keeper) ProcessReAssociateRecords(ctx sdk.Context, msg types.MsgReAssociateRecords) error {
	// TODO(gsk967): uncomment after bond module migration
	//if !k.bondKeeper.HasBond(ctx, msg.OldBondId) {
	//	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Old bond not found.")
	//}
	//
	//if !k.bondKeeper.HasBond(ctx, msg.NewBondId) {
	//	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "New bond not found.")
	//}
	//
	//// Only the bond owner can reassociate all records.
	//oldBond := k.bondKeeper.GetBond(ctx, msg.OldBondId)
	//if msg.Signer != oldBond.Owner {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Old bond owner mismatch.")
	//}
	//
	//newBond := k.bondKeeper.GetBond(ctx, msg.NewBondId)
	//if msg.Signer != newBond.Owner {
	//	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "New bond owner mismatch.")
	//}

	// Reassociate all records.
	records := k.recordKeeper.QueryRecordsByBond(ctx, msg.OldBondId)
	for _, record := range records {
		// Switch bond ID.
		record.BondId = msg.NewBondId
		k.PutRecord(ctx, record)

		k.RemoveBondToRecordIndexEntry(ctx, msg.OldBondId, record.Id)
		k.AddBondToRecordIndexEntry(ctx, msg.NewBondId, record.Id)

		// Required so that renewal is triggered (with new bond ID) for expired records.
		if record.Deleted {
			k.InsertRecordExpiryQueue(ctx, record)
		}
	}

	return nil
}
