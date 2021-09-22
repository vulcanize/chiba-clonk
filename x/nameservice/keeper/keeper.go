package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	auth "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bank "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tharsis/ethermint/x/nameservice/helpers"
	"github.com/tharsis/ethermint/x/nameservice/types"
	"sort"
	"time"
)

var (

	// PrefixCIDToRecordIndex is the prefix for CID -> Record index.
	// Note: This is the primary index in the system.
	// Note: Golang doesn't support const arrays.
	PrefixCIDToRecordIndex = []byte{0x00}

	// PrefixNameAuthorityRecordIndex is the prefix for the name -> NameAuthority index.
	PrefixNameAuthorityRecordIndex = []byte{0x01}

	// PrefixWRNToNameRecordIndex is the prefix for the WRN -> NamingRecord index.
	PrefixWRNToNameRecordIndex = []byte{0x02}

	// PrefixBondIDToRecordsIndex is the prefix for the Bond ID -> [Record] index.
	PrefixBondIDToRecordsIndex = []byte{0x03}

	// PrefixBlockChangesetIndex is the prefix for the block changeset index.
	PrefixBlockChangesetIndex = []byte{0x04}

	// PrefixAuctionToAuthorityNameIndex is the prefix for the auction ID -> authority name index.
	PrefixAuctionToAuthorityNameIndex = []byte{0x05}

	// PrefixBondIDToAuthoritiesIndex is the prefix for the Bond ID -> [Authority] index.
	PrefixBondIDToAuthoritiesIndex = []byte{0x06}

	// PrefixExpiryTimeToRecordsIndex is the prefix for the Expiry Time -> [Record] index.
	PrefixExpiryTimeToRecordsIndex = []byte{0x10}

	// PrefixExpiryTimeToAuthoritiesIndex is the prefix for the Expiry Time -> [Authority] index.
	PrefixExpiryTimeToAuthoritiesIndex = []byte{0x11}

	// KeySyncStatus is the key for the sync status record.
	// Only used by WNS lite but defined here to prevent conflicts with existing prefixes.
	KeySyncStatus = []byte{0xff}

	// PrefixCIDToNamesIndex the the reverse index for naming, i.e. maps CID -> []Names.
	// TODO(ashwin): Move out of WNS once we have an indexing service.
	PrefixCIDToNamesIndex = []byte{0xe0}
)

// Keeper maintains the link to storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper auth.AccountKeeper
	bankKeeper    bank.Keeper
	recordKeeper  RecordKeeper
	//bondKeeper    bond.BondClientKeeper
	//auctionKeeper auction.Keeper

	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	cdc         codec.BinaryCodec // The wire codec for binary encoding/decoding.
	legacyCodec codec.LegacyAmino

	paramSubspace paramtypes.Subspace
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(cdc codec.BinaryCodec, legacyCdc codec.LegacyAmino, accountKeeper auth.AccountKeeper, bankKeeper bank.Keeper, recordKeeper RecordKeeper,
	storeKey sdk.StoreKey, ps paramtypes.Subspace) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		recordKeeper:  recordKeeper,
		//bondKeeper:    bondKeeper,
		//auctionKeeper: auctionKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
		legacyCodec:   legacyCdc,
		paramSubspace: ps,
	}
}

// GetRecordIndexKey Generates Bond ID -> Bond index key.
func GetRecordIndexKey(id string) []byte {
	return append(PrefixCIDToRecordIndex, []byte(id)...)
}

// HasRecord - checks if a record by the given ID exists.
func (k Keeper) HasRecord(ctx sdk.Context, id string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetRecordIndexKey(id))
}

// GetRecord - gets a record from the store.
func (k Keeper) GetRecord(ctx sdk.Context, id string) types.Record {
	store := ctx.KVStore(k.storeKey)
	result := store.Get(GetRecordIndexKey(id))
	var record types.Record
	err := k.cdc.Unmarshal(result, &record)
	if err != nil {
		return types.Record{}
	}
	return record
}

// ListRecords - get all records.
func (k Keeper) ListRecords(ctx sdk.Context) []types.Record {
	var records []types.Record

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixCIDToRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var obj types.Record
			k.cdc.MustUnmarshal(bz, &obj)
			//records = append(records, recordObjToRecord(store, k.cdc, obj))
			records = append(records, obj)
		}
	}

	return records
}

func (k Keeper) GetRecordExpiryQueue(ctx sdk.Context) (expired map[string][]string) {
	records := make(map[string][]string)

	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixExpiryTimeToRecordsIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		var record []string
		//k.cdc.MustUnmarshal(itr.Value(), &record)
		records[string(itr.Key()[len(PrefixExpiryTimeToRecordsIndex):])] = record
	}

	return records
}

// ProcessSetRecord creates a record.
func (k Keeper) ProcessSetRecord(ctx sdk.Context, msg types.MsgSetRecord) (*types.Record, error) {
	payload := msg.Payload.ToReadablePayload()
	record := types.RecordType{Attributes: payload.Record, BondId: msg.BondId}

	// Check signatures.
	resourceSignBytes, _ := record.GetSignBytes()
	cid, err := record.GetCID()
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid record JSON")
	}

	record.Id = cid

	if exists := k.HasRecord(ctx, record.Id); exists {
		// Immutable record already exists. No-op.
		//return record.ToRecordObj(), nil
		return nil, nil
	}

	record.Owners = []string{}
	for _, sig := range payload.Signatures {
		pubKey, err := legacy.PubKeyFromBytes(helpers.BytesFromBase64(sig.PubKey))
		if err != nil {
			fmt.Println("Error decoding pubKey from bytes: ", err)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Invalid public key.")
		}

		sigOK := pubKey.VerifySignature(resourceSignBytes, helpers.BytesFromBase64(sig.Sig))
		//sigOK := pubKey.VerifyBytes(resourceSignBytes, helpers.BytesFromBase64(sig.Sig))
		if !sigOK {
			fmt.Println("Signature mismatch: ", sig.PubKey)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Invalid signature.")
		}
		record.Owners = append(record.Owners, pubKey.Address().String())
	}

	// Sort owners list.
	sort.Strings(record.Owners)
	sdkErr := k.processRecord(ctx, &record, false)
	if sdkErr != nil {
		return nil, sdkErr
	}
	return nil, nil
}

func (k Keeper) processRecord(ctx sdk.Context, record *types.RecordType, isRenewal bool) error {
	params := k.GetParams(ctx)
	//rent := params.RecordRent

	// todo: remove this one once bond module implemented
	//err := k.bankKeeper.TransferCoinsToModuleAccount(ctx, r.BondId, types.RecordRentModuleAccountName, rent)
	//if err != nil{
	//	return err
	//}

	record.CreateTime = ctx.BlockHeader().Time
	record.ExpiryTime = ctx.BlockHeader().Time.Add(params.RecordRentDuration)
	record.Deleted = false

	k.PutRecord(ctx, record.ToRecordObj())
	//k.InsertRecordExpiryQueue(ctx, *record)

	// Renewal doesn't change the name and bond indexes.
	if !isRenewal {
		k.AddBondToRecordIndexEntry(ctx, record.BondId, record.Id)
	}

	return nil
}

// PutRecord - saves a record to the store and updates ID -> Record index.
func (k Keeper) PutRecord(ctx sdk.Context, record types.Record) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRecordIndexKey(record.Id), k.cdc.MustMarshal(&record))
	k.updateBlockChangesetForRecord(ctx, record.Id)
}

func (k Keeper) updateBlockChangesetForRecord(ctx sdk.Context, id string) {
	changeset := k.getOrCreateBlockChangeset(ctx, ctx.BlockHeight())
	changeset.Records = append(changeset.Records, id)
	k.saveBlockChangeset(ctx, changeset)
}

func (k Keeper) getOrCreateBlockChangeset(ctx sdk.Context, height int64) *types.BlockChangeset {
	return getOrCreateBlockChangeset(ctx, ctx.KVStore(k.storeKey), k.cdc, height)
}

func GetBlockChangesetIndexKey(height int64) []byte {
	return append(PrefixBlockChangesetIndex, helpers.Int64ToBytes(height)...)
}
func getOrCreateBlockChangeset(ctx sdk.Context, store sdk.KVStore, cdc codec.BinaryCodec, height int64) *types.BlockChangeset {

	bz := store.Get(GetBlockChangesetIndexKey(height))

	if bz != nil {
		var changeset types.BlockChangeset
		//codec.MustUnmarshalBinaryBare(bz, &changeset)
		err := json.Unmarshal(bz, &changeset)
		if err != nil {
			return nil
		}
		return &changeset
	}

	return &types.BlockChangeset{
		Height:   height,
		Records:  []string{},
		Names:    []string{},
		Auctions: []string{},
		//AuctionBids: []auction.AuctionBidInfo{},
	}
}

// AddBondToRecordIndexEntry adds the Bond ID -> [Record] index entry.
func (k Keeper) AddBondToRecordIndexEntry(ctx sdk.Context, bondID string, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getBondIDToRecordsIndexKey(bondID, id), []byte{})
}

// Generates Bond ID -> Records index key.
func getBondIDToRecordsIndexKey(bondID string, id string) []byte {
	return append(append(PrefixBondIDToRecordsIndex, []byte(bondID)...), []byte(id)...)
}

// InsertRecordExpiryQueue inserts a record CID to the appropriate timeslice in the record expiry queue.
func (k Keeper) InsertRecordExpiryQueue(ctx sdk.Context, val types.Record) {
	//timeSlice := k.GetRecordExpiryQueueTimeSlice(ctx, val.ExpiryTime)
	//timeSlice = append(timeSlice, val.ID)
	//k.SetRecordExpiryQueueTimeSlice(ctx, val.ExpiryTime, timeSlice)
}

// getRecordExpiryQueueTimeKey gets the prefix for the record expiry queue.
func getRecordExpiryQueueTimeKey(timestamp time.Time) []byte {
	timeBytes := sdk.FormatTimeBytes(timestamp)
	return append(PrefixExpiryTimeToRecordsIndex, timeBytes...)
}

// GetRecordExpiryQueueTimeSlice gets a specific record queue timeslice.
// A timeslice is a slice of CIDs corresponding to records that expire at a certain time.
func (k Keeper) GetRecordExpiryQueueTimeSlice(ctx sdk.Context, timestamp time.Time) (cids []string) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(getRecordExpiryQueueTimeKey(timestamp))
	if bz == nil {
		return []string{}
	}

	err := json.Unmarshal(bz, &cids)
	if err != nil {
		return nil
	}
	k.legacyCodec.MustUnmarshalLengthPrefixed(bz, &cids)
	return cids
}

// GetModuleBalances gets the nameservice module account(s) balances.
func (k Keeper) GetModuleBalances(ctx sdk.Context) []*types.AccountBalance {
	var balances []*types.AccountBalance
	accountNames := []string{types.RecordRentModuleAccountName, types.AuthorityRentModuleAccountName}

	for _, accountName := range accountNames {
		moduleAddress := k.accountKeeper.GetModuleAddress(accountName)
		moduleAccount := k.accountKeeper.GetAccount(ctx, moduleAddress)
		if moduleAccount != nil {
			accountBalance := k.bankKeeper.GetAllBalances(ctx, moduleAddress)
			balances = append(balances, &types.AccountBalance{
				AccountName: accountName,
				Balance:     accountBalance,
			})
		}
	}

	return balances
}
