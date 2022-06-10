package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2alpha1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/cosmos/cosmos-sdk/x/nft"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// UpgradeName defines the on-chain upgrade name for the sample simap upgrade from v045 to v046.
//
// NOTE: This upgrade defines a reference implementation of what an upgrade could look like
// when an application is migrating from Cosmos SDK version v0.45.x to v0.46.x.
const UpgradeName = "v045-to-v046"

func GetUpgradeStoreOption(keeper upgradekeeper.Keeper) baseapp.StoreOption {
	upgradeInfo, err := keeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeName && !keeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				group.ModuleName,
				nft.ModuleName,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		return upgradetypes.UpgradeStoreOption(uint64(upgradeInfo.Height), &storeUpgrades)
	}
	return nil
}

func (app *EthermintApp) registerUpgrade() {
	// mainnet upgrade handler
	app.UpgradeKeeper.SetUpgradeHandler(UpgradeName, func(ctx sdk.Context, plan upgradetypes.Plan, _ module.VersionMap) (module.VersionMap, error) {
		fromVM := map[string]uint64{
			"auth":         2,
			"authz":        1,
			"bank":         2,
			"capability":   1,
			"crisis":       1,
			"distribution": 2,
			"evidence":     1,
			"feegrant":     1,
			"gov":          2,
			"mint":         1,
			"params":       1,
			"slashing":     2,
			"staking":      2,
			"upgrade":      1,
			"vesting":      1,
			"ibc":          2,
			"genutil":      1,
			"transfer":     1,
		}

		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})
}
