package nameservice

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tharsis/ethermint/x/nameservice/keeper"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) []abci.ValidatorUpdate {
	keeper.SetParams(ctx, data.Params)

	//for _, record := range data.Records {
	//obj := record.ToRecord()
	//keeper.PutRecord(ctx, obj)
	//
	//// Add to record expiry queue if expiry time is in the future.
	//if obj.ExpiryTime.After(ctx.BlockTime()) {
	//	keeper.InsertRecordExpiryQueue(ctx, obj)
	//}
	//
	//// Note: Bond genesis runs first, so bonds will already be present.
	//if record.BondID != "" {
	//	keeper.AddBondToRecordIndexEntry(ctx, record.BondID, record.ID)
	//}
	//}

	//for _, authority := range data.Authorities {
	// Only import authorities that are marked active.
	//if authority.Entry.Status == types.AuthorityActive {
	//	keeper.SetNameAuthority(ctx, authority.Name, authority.Entry)
	//
	//	// Add authority name to expiry queue.
	//	keeper.InsertAuthorityExpiryQueue(ctx, authority.Name, authority.Entry.ExpiryTime)
	//
	//	// Note: Bond genesis runs first, so bonds will already be present.
	//	if authority.Entry.BondID != "" {
	//		keeper.AddBondToAuthorityIndexEntry(ctx, authority.Entry.BondID, authority.Name)
	//	}
	//}
	//}

	//for _, nameEntry := range data.Names {
	//keeper.SetNameRecord(ctx, nameEntry.Name, nameEntry.Entry.ID)
	//}

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)

	//records := keeper.ListRecords(ctx)
	//recordEntries := []types.RecordObj{}
	//for _, record := range records {
	//	recordEntries = append(recordEntries, record.ToRecordObj())
	//}
	//
	//authorities := keeper.ListNameAuthorityRecords(ctx)
	//authorityEntries := []AuthorityEntry{}
	//for name, record := range authorities {
	//	authorityEntries = append(authorityEntries, AuthorityEntry{
	//		Name:  name,
	//		Entry: record,
	//	})
	//}
	//
	//names := keeper.ListNameRecords(ctx)
	//nameEntries := []types.NameEntry{}
	//for name, record := range names {
	//	nameEntries = append(nameEntries, NameEntry{
	//		Name:  name,
	//		Entry: record,
	//	})
	//}

	return types.GenesisState{
		Params: params,
		//Records:     recordEntries,
		//Authorities: authorityEntries,
		//Names:       nameEntries,
	}
}
