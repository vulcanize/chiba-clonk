package testutil

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	bondcli "github.com/tharsis/ethermint/x/bond/client/cli"
	"github.com/tharsis/ethermint/x/bond/types"
	"github.com/tharsis/ethermint/x/nameservice/client/cli"
	nstypes "github.com/tharsis/ethermint/x/nameservice/types"
	"os"
)

func CreateBond(s *IntegrationTestSuite) {
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

func GetBondId(s *IntegrationTestSuite) string {
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
				CreateBond(s)
				// get the bond id from bond list
				bondId := GetBondId(s)
				dir, err := os.Getwd()
				sr.NoError(err)
				payloadPath := dir + "/example1.yml"

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
				sr.NoError(err)
				sr.Zero(d.Code)
			}
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdReserveName() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "testtest"
	testCases := []struct {
		name string
		args []string
		err  bool
	}{
		{
			"invalid request without name",
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
			"success for parent name",
			[]string{
				authorityName,
				fmt.Sprintf("--owner=%s", accountAddress),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
		},
		{
			"success for sub domains",
			[]string{
				"sub." + authorityName,
				fmt.Sprintf("--owner=%s", accountAddress),
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
			cmd := cli.GetCmdReserveName()

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

func (s *IntegrationTestSuite) TestGetCmdSetName() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "TestGetCmdSetName"
	testCases := []struct {
		name   string
		args   []string
		err    bool
		preRun func(authorityName string)
	}{
		{
			"invalid request without name",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			func(authorityName string) {

			},
		},
		{
			"success",
			[]string{
				fmt.Sprintf("wrn://%s/", authorityName),
				"test_hello_cid",
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			func(authorityName string) {
				// reserving the name
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdReserveName()
				args := []string{
					authorityName,
					fmt.Sprintf("--owner=%s", accountAddress),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)

				// creating the bond
				CreateBond(s)

				// Get the bond-id
				bondId := GetBondId(s)

				// adding bond-id to name authority
				args = []string{
					authorityName, bondId,
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				cmd = cli.GetCmdSetAuthorityBond()

				out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			if !tc.err {
				tc.preRun(authorityName)
			}

			clientCtx := val.ClientCtx
			cmd := cli.GetCmdSetName()

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

func (s *IntegrationTestSuite) TestGetCmdSetAuthorityBond() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "TestGetCmdSetAuthorityBond"

	testCases := []struct {
		name   string
		args   []string
		err    bool
		preRun func(authorityName string)
	}{
		{
			"invalid request without name",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			func(authorityName string) {

			},
		},
		{
			"success with name and bond-id",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			func(authorityName string) {
				// reserving the name
				clientCtx := val.ClientCtx
				cmd := cli.GetCmdReserveName()
				args := []string{
					authorityName,
					fmt.Sprintf("--owner=%s", accountAddress),
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			if !tc.err {
				// reserve the name
				tc.preRun(authorityName)
				// creating the  bond
				CreateBond(s)
				// getting the bond-id
				bondId := GetBondId(s)
				tc.args = append([]string{authorityName, bondId}, tc.args...)
			}
			clientCtx := val.ClientCtx
			cmd := cli.GetCmdSetAuthorityBond()

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

func (s *IntegrationTestSuite) TestGetCmdDeleteName() {
	val := s.network.Validators[0]
	sr := s.Require()
	var authorityName = "TestGetCmdDeleteName"
	testCasesForDeletingName := []struct {
		name   string
		args   []string
		err    bool
		preRun func(authorityName string, s *IntegrationTestSuite)
	}{
		{
			"invalid request without wrn",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			func(authorityName string, s *IntegrationTestSuite) {

			},
		},
		{
			"successfully delete name",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			func(authorityName string, s *IntegrationTestSuite) {
				setupForSetRecord(authorityName, s)
			},
		},
	}

	for _, tc := range testCasesForDeletingName {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			if !tc.err {
				tc.preRun(authorityName, s)
				tc.args = append([]string{fmt.Sprintf("wrn://%s/", authorityName)}, tc.args...)
			}
			clientCtx := val.ClientCtx
			cmd := cli.GetCmdDeleteName()

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

func (s *IntegrationTestSuite) TestGetCmdDissociateBond() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCasesForDeletingName := []struct {
		name    string
		args    []string
		err     bool
		preRun  func(s *IntegrationTestSuite) string
		postRun func(recordId string, s *IntegrationTestSuite)
	}{
		{
			"invalid request without wrn",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			func(s *IntegrationTestSuite) string {
				return ""
			},
			func(recordId string, s *IntegrationTestSuite) {

			},
		},
		{
			"successfully dissociate bond-id from record ",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			func(s *IntegrationTestSuite) string {
				// create the bond
				CreateBond(s)
				// get the bond id from bond list
				bondId := GetBondId(s)
				dir, err := os.Getwd()
				sr.NoError(err)
				payloadPath := dir + "/example1.yml"

				args := []string{
					payloadPath, bondId,
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}

				clientCtx := val.ClientCtx
				cmd := cli.GetCmdSetRecord()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)

				// retrieving the record-id
				args = []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
				cmd = cli.GetCmdList()
				out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var records []nstypes.RecordType
				err = json.Unmarshal(out.Bytes(), &records)
				sr.NoError(err)
				return records[0].Id
			},
			func(recordId string, s *IntegrationTestSuite) {
				// checking the bond-id removed or not
				clientCtx := val.ClientCtx
				args := []string{recordId, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
				cmd := cli.GetCmdGetResource()
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var response nstypes.QueryRecordByIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				record := response.GetRecord()
				sr.NotNil(record)
				sr.Zero(len(record.GetBondId()))
			},
		},
	}

	for _, tc := range testCasesForDeletingName {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			var recordId string
			if !tc.err {
				recordId = tc.preRun(s)
				tc.args = append([]string{recordId}, tc.args...)
			}
			clientCtx := val.ClientCtx
			cmd := cli.GetCmdDissociateBond()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
				// post-run
				tc.postRun(recordId, s)
			}
		})
	}
}

//
//func (s *IntegrationTestSuite) TestGetCmdDissociateRecords() {
//	val := s.network.Validators[0]
//	sr := s.Require()
//	testCasesForDeletingName := []struct {
//		name    string
//		args    []string
//		err     bool
//		preRun  func(s *IntegrationTestSuite) (string, string)
//		postRun func(recordId string, s *IntegrationTestSuite)
//	}{
//		{
//			"invalid request without wrn",
//			[]string{
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
//			},
//			true,
//			func(s *IntegrationTestSuite) (string, string) {
//				return "", ""
//			},
//			func(recordId string, s *IntegrationTestSuite) {
//
//			},
//		},
//		{
//			"successfully dissociate records from bond-id",
//			[]string{
//				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
//				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
//			},
//			false,
//			func(s *IntegrationTestSuite) (string, string) {
//				// create the bond
//				CreateBond(s)
//				// get the bond id from bond list
//				bondId := GetBondId(s)
//				dir, err := os.Getwd()
//				sr.NoError(err)
//				payloadPath := dir + "/example1.yml"
//
//				args := []string{
//					payloadPath, bondId,
//					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
//					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
//					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
//					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
//					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
//				}
//
//				clientCtx := val.ClientCtx
//				cmd := cli.GetCmdSetRecord()
//
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
//				sr.NoError(err)
//				var d sdk.TxResponse
//				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
//				sr.NoError(err)
//				sr.Zero(d.Code)
//
//				// retrieving the record-id
//				args = []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
//				cmd = cli.GetCmdList()
//				out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
//				sr.NoError(err)
//				var records []nstypes.RecordType
//				err = json.Unmarshal(out.Bytes(), &records)
//				sr.NoError(err)
//				for _, record := range records {
//					if len(record.BondId) != 0 {
//						return record.Id, record.BondId
//					}
//				}
//				return records[0].Id, records[0].BondId
//			},
//			func(recordId string, s *IntegrationTestSuite) {
//				// checking the bond-id removed or not
//				clientCtx := val.ClientCtx
//				args := []string{recordId, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
//				cmd := cli.GetCmdGetResource()
//				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
//				sr.NoError(err)
//				var response nstypes.QueryRecordByIdResponse
//				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
//				sr.NoError(err)
//				record := response.GetRecord()
//				sr.NotNil(record)
//				sr.Zero(len(record.GetBondId()))
//			},
//		},
//	}
//
//	for _, tc := range testCasesForDeletingName {
//		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
//			var bondId string
//			var recordId string
//			if !tc.err {
//				recordId, bondId = tc.preRun(s)
//				tc.args = append([]string{bondId}, tc.args...)
//			}
//			clientCtx := val.ClientCtx
//			cmd := cli.GetCmdDissociateRecords()
//
//			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
//			if tc.err {
//				sr.Error(err)
//			} else {
//				sr.NoError(err)
//				var d sdk.TxResponse
//				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
//				sr.NoError(err)
//				sr.Zero(d.Code)
//				// post-run
//				tc.postRun(recordId, s)
//			}
//		})
//	}
//}

func (s *IntegrationTestSuite) TestGetCmdAssociateBond() {
	val := s.network.Validators[0]
	sr := s.Require()
	testCasesForDeletingName := []struct {
		name    string
		args    []string
		err     bool
		preRun  func(s *IntegrationTestSuite) (string, string)
		postRun func(recordId, bondId string, s *IntegrationTestSuite)
	}{
		{
			"invalid request without wrn",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			true,
			func(s *IntegrationTestSuite) (string, string) {
				return "", ""
			},
			func(recordId, bondId string, s *IntegrationTestSuite) {

			},
		},
		{
			"successfully dissociate records from bond-id",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
			},
			false,
			func(s *IntegrationTestSuite) (string, string) {
				// create the bond
				CreateBond(s)
				// get the bond id from bond list
				bondId := GetBondId(s)
				dir, err := os.Getwd()
				sr.NoError(err)
				payloadPath := dir + "/example1.yml"

				txArgs := []string{
					payloadPath, bondId,
					fmt.Sprintf("--%s=%s", flags.FlagFrom, accountName),
					fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
					fmt.Sprintf("--%s=json", tmcli.OutputFlag),
					fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
					fmt.Sprintf("--%s=%s", flags.FlagFees, fmt.Sprintf("3%s", s.cfg.BondDenom)),
				}

				clientCtx := val.ClientCtx
				cmd := cli.GetCmdSetRecord()

				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, txArgs)
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)

				// retrieving the record-id
				args := []string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
				cmd = cli.GetCmdList()
				out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var records []nstypes.RecordType
				err = json.Unmarshal(out.Bytes(), &records)
				sr.NoError(err)

				// GetCmdDissociateBond bond
				cmd = cli.GetCmdDissociateBond()
				txArgs = append([]string{records[0].Id}, txArgs[2:]...)
				out, err = clitestutil.ExecTestCLICmd(clientCtx, cmd, txArgs)
				sr.NoError(err)
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)

				return records[0].Id, records[0].BondId
			},
			func(recordId, bondId string, s *IntegrationTestSuite) {
				// checking the bond-id removed or not
				clientCtx := val.ClientCtx
				args := []string{recordId, fmt.Sprintf("--%s=json", tmcli.OutputFlag)}
				cmd := cli.GetCmdGetResource()
				out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
				sr.NoError(err)
				var response nstypes.QueryRecordByIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				record := response.GetRecord()
				sr.NotNil(record)
				sr.Equal(record.GetBondId(), bondId)
			},
		},
	}

	for _, tc := range testCasesForDeletingName {
		s.Run(fmt.Sprintf("Case %s", tc.name), func() {
			var recordId string
			var bondId string
			if !tc.err {
				recordId, bondId = tc.preRun(s)
				tc.args = append([]string{recordId, bondId}, tc.args...)
			}
			clientCtx := val.ClientCtx
			cmd := cli.GetCmdAssociateBond()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.err {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var d sdk.TxResponse
				err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &d)
				sr.NoError(err)
				sr.Zero(d.Code)
				// post-run
				tc.postRun(recordId, bondId, s)
			}
		})
	}
}
