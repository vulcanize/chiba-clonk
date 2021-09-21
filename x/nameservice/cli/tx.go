package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/tharsis/ethermint/server/flags"
	"github.com/tharsis/ethermint/x/nameservice/types"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"strings"
)

// NewTxCmd returns a root CLI command handler for all x/bond transaction commands.
func NewTxCmd() *cobra.Command {
	bondTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "nameservice transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bondTxCmd.AddCommand(
		GetCmdSetRecord(),
		GetCmdSetName(),
	)

	return bondTxCmd
}

// GetCmdSetRecord is the CLI command for creating/updating a record.
func GetCmdSetRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [payload file path] [bond-id]",
		Short: "Set record.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create a new record with payload and bond id.
Example:
$ %s tx set [payload file path] [bond-id]
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			payload, err := getPayloadFromFile(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgSetRecord(payload.ToPayload(), args[1], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlags(cmd)
	return cmd
}

// GetCmdSetName is the CLI command for mapping a name to a CID.
func GetCmdSetName() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-name [wrn] [cid]",
		Short: "Set WRN to CID mapping.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSetName(args[0], args[1], clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlags(cmd)
	return cmd
}

// Load payload object from YAML file.
func getPayloadFromFile(filePath string) (*types.PayloadType, error) {
	var payload types.PayloadType

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
