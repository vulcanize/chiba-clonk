package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/store/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2alpha1"
	v2multistore "github.com/cosmos/cosmos-sdk/store/v2alpha1/multi"
	"github.com/google/orderedcode"

	dbm "github.com/cosmos/cosmos-sdk/db"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"
)

func StateMigrationFromIAVLtoSMT(keys map[string]*storetypes.KVStoreKey) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "iavl-to-smt [old-data-home-dir] [new-data-home-dir]",
		Short:   "State migraiton from iavl to smt",
		Example: fmt.Sprint("%s iavl-to-smt", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			home := args[0]
			if _, err := os.Stat(home); err != nil {
				log.Fatal(err)
			}
			dataDir := path.Join(home, "data")

			newDir := args[1]
			if _, err := os.Stat(home); err != nil {
				log.Fatal(err)
			}
			newdataDir := path.Join(newDir, "data")

			err := iavlToSmt(keys, dataDir, newdataDir)
			if err != nil {
				fmt.Println("[X] State migration failed from iavl to smt.")
				return err
			}
			fmt.Println("[!] State migration completed from iavl to smt.")
			return nil
		},
	}

	return cmd
}

func iavlToSmt(keys map[string]*storetypes.KVStoreKey, dataDir, newDataDir string) error {
	oldLevelDb, err := tmdb.NewGoLevelDB("application", dataDir)
	if err != nil {
		return err
	}
	cms := rootmulti.NewStore(oldLevelDb, tmlog.NewNopLogger())
	// mount the kvStore
	for _, key := range keys {
		cms.MountStoreWithDB(key, types.StoreTypeIAVL, oldLevelDb)
	}

	fmt.Println("v1store init is done.")
	// new store v2
	// new dir
	opts := v2multistore.DefaultStoreParams()

	ndb, err := dbm.NewDB("application", dbm.BadgerDBBackend, newDataDir)
	if err != nil {
		return err
	}
	fmt.Println("badgerdb initial succeed")
	v2multistore.MigrateFromV1(cms, ndb, opts)

	// state migration of tm-db
	stateDB, err := tmdb.NewGoLevelDB("state", dataDir)
	if err != nil {
		return err
	}

	// tm key
	prefixState := int64(8)
	tmStateKey, err := orderedcode.Append(nil, prefixState)
	if err != nil {
		panic(err)
	}
	ss, err := stateDB.Get(tmStateKey)
	if err != nil {
		panic(err)
	}

	bStateDB, err := tmdb.NewDB("state", tmdb.BackendType(dbm.BadgerDBBackend), newDataDir)
	if err != nil {
		return err
	}

	err = bStateDB.Set(tmStateKey, ss)
	if err != nil {
		return err
	}
	fmt.Println("state tendetmint migration is done.")
	return nil
}
