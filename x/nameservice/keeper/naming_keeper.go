package keeper

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tharsis/ethermint/x/nameservice/helpers"
	"github.com/tharsis/ethermint/x/nameservice/types"
	"net/url"
)

func getAuthorityPubKey(pubKey cryptotypes.PubKey) string {
	if pubKey != nil {
		return helpers.BytesToBase64(pubKey.Bytes())
	}
	return ""
}

// GetNameAuthorityIndexKey Generates name -> NameAuthority index key.
func GetNameAuthorityIndexKey(name string) []byte {
	return append(PrefixNameAuthorityRecordIndex, []byte(name)...)
}

// GetNameRecordIndexKey Generates WRN -> NameRecord index key.
func GetNameRecordIndexKey(wrn string) []byte {
	return append(PrefixWRNToNameRecordIndex, []byte(wrn)...)
}

func GetCIDToNamesIndexKey(id string) []byte {
	return append(PrefixCIDToNamesIndex, []byte(id)...)
}

// GetNameAuthority - gets a name authority from the store.
func GetNameAuthority(store sdk.KVStore, codec codec.BinaryCodec, name string) *types.NameAuthority {
	authorityKey := GetNameAuthorityIndexKey(name)
	if !store.Has(authorityKey) {
		return nil
	}

	bz := store.Get(authorityKey)
	var obj types.NameAuthority
	codec.MustUnmarshal(bz, &obj)

	return &obj
}

func SetNameAuthority(ctx sdk.Context, store sdk.KVStore, codec codec.BinaryCodec, name string, authority types.NameAuthority) {
	store.Set(GetNameAuthorityIndexKey(name), codec.MustMarshal(&authority))
	//updateBlockChangesetForNameAuthority(ctx, store, codec, name)
}

// SetNameAuthority creates the NameAutority record.
func (k Keeper) SetNameAuthority(ctx sdk.Context, name string, authority types.NameAuthority) {
	SetNameAuthority(ctx, ctx.KVStore(k.storeKey), k.cdc, name, authority)
}

// GetNameAuthority - gets a name authority from the store.
func (k Keeper) GetNameAuthority(ctx sdk.Context, name string) *types.NameAuthority {
	return GetNameAuthority(ctx.KVStore(k.storeKey), k.cdc, name)
}

func (k Keeper) getAuthority(ctx sdk.Context, wrn string) (string, *url.URL, *types.NameAuthority, error) {
	parsedWRN, err := url.Parse(wrn)
	if err != nil {
		return "", nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid WRN.")
	}

	name := parsedWRN.Host
	authority := k.GetNameAuthority(ctx, name)
	if authority == nil {
		return name, nil, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Name authority not found.")
	}

	return name, parsedWRN, authority, nil
}

func (k Keeper) checkWRNAccess(ctx sdk.Context, signer sdk.AccAddress, wrn string) error {
	name, parsedWRN, authority, err := k.getAuthority(ctx, wrn)
	if err != nil {
		return err
	}

	formattedWRN := fmt.Sprintf("wrn://%s%s", name, parsedWRN.RequestURI())
	if formattedWRN != wrn {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Invalid WRN.")
	}

	if authority.OwnerAddress != signer.String() {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Access denied.")
	}

	if authority.Status != types.AuthorityActive {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Authority is not active.")
	}

	if authority.BondId == "" || len(authority.BondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Authority bond not found.")
	}

	if authority.OwnerPublicKey == "" {
		// Try to set owner public key if account has it available now.
		ownerAccount := k.accountKeeper.GetAccount(ctx, signer)
		pubKey := ownerAccount.GetPubKey()
		if pubKey != nil {
			// Update public key in authority record.
			authority.OwnerPublicKey = getAuthorityPubKey(pubKey)
			k.SetNameAuthority(ctx, name, *authority)
		}
	}

	return nil
}

// GetNameRecord - gets a name record from the store.
func GetNameRecord(store sdk.KVStore, codec codec.BinaryCodec, wrn string) *types.NameRecord {
	nameRecordKey := GetNameRecordIndexKey(wrn)
	if !store.Has(nameRecordKey) {
		return nil
	}

	bz := store.Get(nameRecordKey)
	var obj types.NameRecord
	codec.MustUnmarshal(bz, &obj)

	return &obj
}

// GetNameRecord - gets a name record from the store.
func (k Keeper) GetNameRecord(ctx sdk.Context, wrn string) *types.NameRecord {
	_, _, authority, err := k.getAuthority(ctx, wrn)
	if err != nil || authority.Status != types.AuthorityActive {
		// If authority is not active (or any other error), lookup fails.
		return nil
	}

	nameRecord := GetNameRecord(ctx.KVStore(k.storeKey), k.cdc, wrn)

	// Name record may not exist.
	if nameRecord == nil {
		return nil
	}

	// Name lookup should fail if the name record is stale.
	// i.e. authority was registered later than the name.
	if authority.Height > nameRecord.Latest.Height {
		return nil
	}

	return nameRecord
}

// RemoveRecordToNameMapping removes a name from the record ID -> []names index.
func RemoveRecordToNameMapping(store sdk.KVStore, codec codec.BinaryCodec, id string, wrn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	var names []string
	//codec.MustUnmarshalBinaryBare(store.Get(reverseNameIndexKey), &names)
	json.Unmarshal(store.Get(reverseNameIndexKey), &names)
	nameSet := helpers.SliceToSet(names)
	nameSet.Remove(wrn)

	if nameSet.Cardinality() == 0 {
		// Delete as storing empty slice throws error from baseapp.
		store.Delete(reverseNameIndexKey)
	} else {
		data, _ := json.Marshal(helpers.SetToSlice(nameSet))
		store.Set(reverseNameIndexKey, data)
	}
}

// AddRecordToNameMapping adds a name to the record ID -> []names index.
func AddRecordToNameMapping(store sdk.KVStore, codec codec.BinaryCodec, id string, wrn string) {
	reverseNameIndexKey := GetCIDToNamesIndexKey(id)

	var names []string
	if store.Has(reverseNameIndexKey) {
		//codec.MustUnmarshalBinaryBare(store.Get(reverseNameIndexKey), &names)
		json.Unmarshal(store.Get(reverseNameIndexKey), &names)
	}

	nameSet := helpers.SliceToSet(names)
	nameSet.Add(wrn)
	//store.Set(reverseNameIndexKey, codec.MustMarshalBinaryBare(helpers.SetToSlice(nameSet)))
	data, _ := json.Marshal(helpers.SetToSlice(nameSet))
	store.Set(reverseNameIndexKey, data)
}

// SetNameRecord - sets a name record.
func SetNameRecord(store sdk.KVStore, codec codec.BinaryCodec, wrn string, id string, height int64) {
	nameRecordIndexKey := GetNameRecordIndexKey(wrn)

	var nameRecord types.NameRecord
	fmt.Println("has ", store.Has(nameRecordIndexKey))
	if store.Has(nameRecordIndexKey) {
		bz := store.Get(nameRecordIndexKey)
		codec.MustUnmarshal(bz, &nameRecord)
		nameRecord.History = append(nameRecord.History, nameRecord.Latest)

		// Update old CID -> []Name index.
		if nameRecord.Latest.Id != "" || len(nameRecord.Latest.Id) != 0 {
			RemoveRecordToNameMapping(store, codec, nameRecord.Latest.Id, wrn)
		}
	}

	nameRecord.Latest = &types.NameRecordEntry{
		Id:     id,
		Height: uint64(height),
	}

	fmt.Println("name record ", nameRecord)
	store.Set(nameRecordIndexKey, codec.MustMarshal(&nameRecord))

	fmt.Println("name record stored")
	// Update new CID -> []Name index.
	if id != "" {
		AddRecordToNameMapping(store, codec, id, wrn)
	}
}

// SetNameRecord - sets a name record.
func (k Keeper) SetNameRecord(ctx sdk.Context, wrn string, id string) {
	SetNameRecord(ctx.KVStore(k.storeKey), k.cdc, wrn, id, ctx.BlockHeight())

	// Update changeset for name.
	//k.updateBlockChangesetForName(ctx, wrn)
}

// ProcessSetName creates a WRN -> Record ID mapping.
func (k Keeper) ProcessSetName(ctx sdk.Context, msg types.MsgSetName) error {
	signerAddress, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return err
	}
	err = k.checkWRNAccess(ctx, signerAddress, msg.Wrn)
	fmt.Println(" Err at checkWRNAccess ", err)
	if err != nil {
		return err
	}

	nameRecord := k.GetNameRecord(ctx, msg.Wrn)
	fmt.Println("nameRecord ", nameRecord)
	if nameRecord != nil && nameRecord.Latest.Id == msg.Cid {
		// Already pointing to same ID, no-op.
		return nil
	}

	k.SetNameRecord(ctx, msg.Wrn, msg.Cid)

	return nil
}

// ListNameRecords - get all name records.
func (k Keeper) ListNameRecords(ctx sdk.Context) []types.NameEntry {
	var nameEntries []types.NameEntry
	store := ctx.KVStore(k.storeKey)
	itr := sdk.KVStorePrefixIterator(store, PrefixWRNToNameRecordIndex)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		bz := store.Get(itr.Key())
		if bz != nil {
			var record types.NameRecord
			k.cdc.MustUnmarshal(bz, &record)
			nameEntries = append(nameEntries, types.NameEntry{
				Name:   string(itr.Key()[len(PrefixWRNToNameRecordIndex):]),
				Record: &record,
			})
		}
	}

	return nameEntries
}
