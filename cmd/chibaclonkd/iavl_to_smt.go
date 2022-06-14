package main

import (
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/store/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2alpha1"
	v2multistore "github.com/cosmos/cosmos-sdk/store/v2alpha1/multi"
	"github.com/cosmos/cosmos-sdk/version"

	dbm "github.com/cosmos/cosmos-sdk/db"
	"github.com/spf13/cobra"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmdb "github.com/tendermint/tm-db"
)

var logger tmlog.Logger

func init() {
	logger, _ = tmlog.NewDefaultLogger("plain", "info", false)
	logger = logger.With("module", "state_migration")
}

// StateMigrationFromIAVLtoSMT cli cmd for migrate the data from iavl to smt and goleveldb to badgerdb
func StateMigrationFromIAVLtoSMT(keys map[string]*storetypes.KVStoreKey) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "iavl-to-smt [old-data-home-dir] [new-data-home-dir]",
		Short:   "State migraiton from iavl to smt",
		Example: fmt.Sprintf("%s iavl-to-smt ~/.chibaclonkd ~/.newchibaclonkd", version.AppName),
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			home := args[0]
			if _, err := os.Stat(home); err != nil {
				logger.Error(err.Error())
				return err
			}
			dataDir := path.Join(home, "data")

			newDir := args[1]
			if _, err := os.Stat(home); err != nil {
				logger.Error(err.Error())
				return err
			}
			newdataDir := path.Join(newDir, "data")

			err := iavlToSmt(keys, dataDir, newdataDir)
			if err != nil {
				logger.Error(fmt.Sprintf("State migration failed from iavl to smt. Err %v\n", err))
				return err
			}
			logger.Info("State migration is completed from iavl to smt and goleveldb to badgerdb.")
			return nil
		},
	}

	return cmd
}

// iavlToSmt migrate the data from iavl to smt
func iavlToSmt(keys map[string]*storetypes.KVStoreKey, dataDir, newDataDir string) error {
	v1LevelDB, err := tmdb.NewGoLevelDB("application", dataDir)
	if err != nil {
		return err
	}
	// v1 Store
	cms := rootmulti.NewStore(v1LevelDB, tmlog.NewNopLogger())
	// mount the kvStore
	for _, key := range keys {
		cms.MountStoreWithDB(key, types.StoreTypeIAVL, v1LevelDB)
	}
	logger.Info("v1store init is done.")

	//  v2 new store
	opts := v2multistore.DefaultStoreParams()
	ndb, err := dbm.NewDB("application", dbm.BadgerDBBackend, newDataDir)
	if err != nil {
		return err
	}
	logger.Info("badgerdb init done for application")

	// migrating the data from v1 to v2
	_, err = v2multistore.MigrateFromV1(cms, ndb, opts)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	// migrating the data from goleveldb to badgerdb
	stores := [...]string{
		"blockstore",
		"state",
		"peerstore",
		"tx_index",
		"evidence",
		"light",
	}

	for _, s := range stores {
		// goleveldb init
		goleveldb, err := tmdb.NewGoLevelDB(s, dataDir)
		if err != nil {
			return err
		}
		iter, err := goleveldb.Iterator(nil, nil)
		if err != nil {
			return err
		}

		// badgerdb init
		bagerDb, err := dbm.NewDB(s, dbm.BadgerDBBackend, newDataDir)
		if err != nil {
			return err
		}

		// read writer for badgerdb
		rw := bagerDb.ReadWriter()

		for ; iter.Valid(); iter.Next() {
			err = rw.Set(iter.Key(), iter.Value())
			if err != nil {
				return err
			}
		}

		// commit the data into badgerdb
		err = rw.Commit()
		if err != nil {
			return err
		}

		// close the badgerdb connection
		err = bagerDb.Close()
		if err != nil {
			return err
		}

		// close the goleveldb connection
		err = goleveldb.Close()
		if err != nil {
			return err
		}
	}

	logger.Info("state tendetmint migration is done.")
	return nil
}
