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
	"github.com/gogo/protobuf/proto"
	"github.com/spf13/cobra"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmstate "github.com/tendermint/tendermint/proto/tendermint/state"
	"github.com/tendermint/tendermint/version"
	tmdb "github.com/tendermint/tm-db"
)

func StateMigrationFromIAVLtoSMT(keys map[string]*storetypes.KVStoreKey) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "iavl-to-smt [old-data-home-dir] [new-data-home-dir]",
		Short:   "State migraiton from iavl to smt",
		Example: fmt.Sprint("%s iavl-to-smt"),
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
	lcd, ss, err := v2multistore.MigrateFromV1(cms, ndb, opts)
	if err != nil {
		return err
	}

	fmt.Println("Commit hash", string(ss.Commit().Hash))
	stores := [...]string{
		"blockstore",
		"state",
		"peerstore",
		"tx_index",
		"evidence",
		"light",
	}

	for _, s := range stores {
		t, err := tmdb.NewGoLevelDB(s, dataDir)
		if err != nil {
			return err
		}
		// new badgerdb
		b, err := dbm.NewDB(s, dbm.BadgerDBBackend, newDataDir)
		if err != nil {
			return err
		}
		r := b.ReadWriter()
		iter, err := t.Iterator(nil, nil)
		if err != nil {
			return err
		}
		for ; iter.Valid(); iter.Next() {
			r.Set(iter.Key(), iter.Value())
		}
		err = r.Commit()
		if err != nil {
			return err
		}
		b.Close()
		t.Close()
	}

	fmt.Println("state tendetmint migration is done.")
	// access state db

	sdb, err := dbm.NewDB("state", dbm.BadgerDBBackend, newDataDir)
	if err != nil {
		return err
	}

	type Version struct {
		Consensus version.Consensus ` json:"consensus"`
		Software  string            ` json:"software"`
	}

	fmt.Println("state db init done")
	prefixState := int64(8)
	stateKey, err := orderedcode.Append(nil, prefixState)
	if err != nil {
		panic(err)
	}

	rat := sdb.ReadWriter()
	buf, err := rat.Get(stateKey)
	if err != nil {
		panic(err)
	}

	sp := new(tmstate.State)
	err = proto.Unmarshal(buf, sp)
	if err != nil {
		// DATA HAS BEEN CORRUPTED OR THE SPEC HAS CHANGED
		fmt.Sprintf(`LoadState: Data has been corrupted or its spec has changed:%v\n`, err)
	}

	fmt.Println("last hash form state ", sp.AppHash)
	fmt.Print("new hash from badger ", lcd.Hash)
	sp.AppHash = lcd.Hash
	b, err := proto.Marshal(sp)
	if err != nil {
		panic(err)
	}
	err = rat.Set(stateKey, b)
	if err != nil {
		panic(err)
	}
	err = rat.Commit()
	if err != nil {
		panic(err)
	}
	return nil
}
