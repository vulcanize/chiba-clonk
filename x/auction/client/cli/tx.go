package cli

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	utils "github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tharsis/ethermint/x/auction/types"

	wnsUtils "github.com/vulcanize/dxns/utils"
)

// GetTxCmd returns transaction commands for this module.
func GetTxCmd() *cobra.Command {
	auctionTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Auction transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// TODO(ashwin): Add Tx commands.
	auctionTxCmd.AddCommand(
		GetCmdCommitBid(),
		GetCmdRevealBid(),
	)

	return auctionTxCmd
}

// GetCmdCommitBid is the CLI command for committing a bid.
func GetCmdCommitBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-bid [auction-id] [bid-amount]",
		Short: "Commit sealed bid.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bidAmount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			mnemonic, err := wnsUtils.GenerateMnemonic()
			if err != nil {
				return err
			}

			chainID := viper.GetString("chain-id")
			auctionID := args[0]

			reveal := map[string]interface{}{
				"chainId":       chainID,
				"auctionId":     auctionID,
				"bidderAddress": clientCtx.GetFromAddress().String(),
				"bidAmount":     bidAmount.String(),
				"noise":         mnemonic,
			}

			commitHash, content, err := wnsUtils.GenerateHash(reveal)
			if err != nil {
				return err
			}

			// Save reveal file.
			ioutil.WriteFile(fmt.Sprintf("%s-%s.json", clientCtx.GetFromName(), commitHash), content, 0600)

			msg := types.NewMsgCommitBid(auctionID, commitHash, clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRevealBid is the CLI command for revealing a bid.
func GetCmdRevealBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reveal-bid [auction-id] [reveal-file-path]",
		Short: "Reveal bid.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionID := args[0]
			revealFilePath := args[1]

			revealBytes, err := ioutil.ReadFile(revealFilePath)
			if err != nil {
				return err
			}

			// TODO(ashwin): Before revealing, check if auction is in reveal phase.

			msg := types.NewMsgRevealBid(auctionID, hex.EncodeToString(revealBytes), clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
