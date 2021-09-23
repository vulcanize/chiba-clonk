package testutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tharsis/ethermint/testutil/network"
	"github.com/tharsis/ethermint/x/bond/client/cli"
	"github.com/tharsis/ethermint/x/bond/types"
	"gopkg.in/yaml.v2"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestGetCmdQueryParams() {
	val := s.network.Validators[0]

	testCases := []struct {
		name       string
		args       []string
		outputType string
	}{
		{
			"json output",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"json",
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetQueryParamsCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)
			var param types.QueryParamsResponse
			if tc.outputType == "json" {
				err := clientCtx.Codec.UnmarshalJSON(out.Bytes(), &param)
				s.Require().NoError(err)
			} else {
				err := yaml.Unmarshal(out.Bytes(), &param)
				s.Require().NoError(err)
			}
			s.Require().Equal(param.Params.MaxBondAmount, types.DefaultParams().MaxBondAmount)
		})
	}
}

func (s *IntegrationTestSuite) TestGetQueryBondLists() {
	val := s.network.Validators[0]
	suiteRequire := s.Require()
	var accountName = "newAccount"

	testCases := []struct {
		name       string
		args       []string
		outputType string
		createBond bool
		noOfBonds  int
	}{
		{
			"Empty Bonds",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"json",
			false,
			0,
		},
		{
			"Crate and Get Bond",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			"json",
			true,
			1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			clientCtx := val.ClientCtx
			if tc.createBond {
				consPrivKey := ed25519.GenPrivKey()
				consPubKeyBz, err := s.cfg.Codec.MarshalInterfaceJSON(consPrivKey.PubKey())
				suiteRequire.NoError(err)
				suiteRequire.NotNil(consPubKeyBz)

				info, _, err := val.ClientCtx.Keyring.NewMnemonic(accountName, keyring.English, sdk.FullFundraiserPath, keyring.DefaultBIP39Passphrase, hd.Secp256k1)
				suiteRequire.NoError(err)

				newAddr := sdk.AccAddress(info.GetPubKey().Address())
				_, err = banktestutil.MsgSendExec(
					val.ClientCtx,
					val.Address,
					newAddr,
					sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
				)
				suiteRequire.NoError(err)

				createBondCmd := cli.NewCreateBondCmd()
				args := []string{
					fmt.Sprintf("10%s", s.cfg.BondDenom),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, createBondCmd, args)
				suiteRequire.NoError(err)
				var d sdk.TxResponse
				clientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				suiteRequire.Zero(d.Code)
				err = s.network.WaitForNextBlock()
				suiteRequire.NoError(err)
			}
			cmd := cli.GetQueryBondLists()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			suiteRequire.NoError(err)
			var queryResponse types.QueryGetBondsResponse
			if tc.outputType == "json" {
				err := clientCtx.Codec.UnmarshalJSON(out.Bytes(), &queryResponse)
				suiteRequire.NoError(err)
			} else {
				err := yaml.Unmarshal(out.Bytes(), &queryResponse)
				suiteRequire.NoError(err)
			}
			suiteRequire.Equal(tc.noOfBonds, len(queryResponse.GetBonds()))
		})
	}
}
