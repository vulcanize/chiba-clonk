package testutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	bondcli "github.com/tharsis/ethermint/x/bond/client/cli"
	"github.com/tharsis/ethermint/x/bond/types"
	"github.com/tharsis/ethermint/x/nameservice/cli"
	"os"
)

func (s *IntegrationTestSuite) TestTxCreateBond() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			"create bond",
			[]string{
				fmt.Sprintf("100000000000%s", s.cfg.BondDenom),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			cmd := bondcli.NewCreateBondCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			}
		})
	}
}

func (s *IntegrationTestSuite) GetBondId() string {
	cmd := bondcli.GetQueryBondLists()
	val := s.network.Validators[0]
	sr := s.Require()
	clientCtx := val.ClientCtx

	out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)})
	sr.NoError(err)
	var queryResponse types.QueryGetBondsResponse
	err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
	sr.NoError(err)

	// extract bond id from bonds list
	bond := queryResponse.GetBonds()[0]
	return bond.GetId()
}

func (s *IntegrationTestSuite) TestGetCmdSetRecord() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			"invalid request without bond id/without payload",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
		},
		{
			"success",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			if !tc.err {
				// create the bond
				s.TestTxCreateBond()
				dir, err := os.Getwd()
				sr.NoError(err)
				payloadPath := dir + "/example1.yml"
				// get the bond id from bond list
				bondId := s.GetBondId()
				tc.args = append([]string{payloadPath, bondId}, tc.args...)
			}
			clientCtx := val.ClientCtx
			cmd := cli.GetCmdSetRecord()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.Nil(d)
				sr.NoError(err)
				sr.Zero(d.Code)
			}
		})
	}
}
