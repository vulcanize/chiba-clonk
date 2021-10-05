package testutil

import (
	"fmt"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	"github.com/tharsis/ethermint/x/nameservice/cli"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

func (s *IntegrationTestSuite) TestGetCmdQueryParams() {
	val := s.network.Validators[0]
	sr := s.Require()

	testCases := []struct {
		name string
		args []string
	}{
		{
			"params",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.GetQueryParamsCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			sr.NoError(err)
			var param types.QueryParamsResponse
			err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &param)
			sr.NoError(err)
			params := types.DefaultParams()
			sr.Equal(param.GetParams().String(), params.String())
		})
	}
}

func (s *IntegrationTestSuite) TestGetCmdQueryForRecords() {
	val := s.network.Validators[0]
	sr := s.Require()
	var recordID string

	testCases := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
	}{
		{
			"empty list",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			0,
		},
		{
			"get records list",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdList()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryListRecordsResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.Equal(tc.noOfRecords, response.GetRecords())
				recordID = response.GetRecords()[0].GetId()
			}
		})
	}

	s.T().Log("Test Cases for getting records by record-id")
	testCasesByRecordID := []struct {
		name   string
		args   []string
		expErr bool
	}{
		{
			"invalid request without record id",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
		},
		{
			"get records by id",
			[]string{recordID, fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
		},
	}

	for _, tc := range testCasesByRecordID {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdGetResource()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryRecordByIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.NotNil(response.GetRecord())
			}
		})
	}

	s.T().Log("Test Cases for getting records by bond-id")
	testCasesByRecordByBondID := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
	}{
		{
			"invalid request without bond-id",
			[]string{fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			true,
			0,
		},
		{
			"get records by bond-id",
			[]string{s.GetBondId(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
		},
	}

	for _, tc := range testCasesByRecordByBondID {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdQueryByBond()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.QueryRecordByBondIdResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.Equal(tc.noOfRecords, len(response.GetRecords()))
			}
		})
	}

	s.T().Log("Test Cases for getting nameservice module account balance")
	testCasesForNameServiceModuleBalance := []struct {
		name        string
		args        []string
		expErr      bool
		noOfRecords int
	}{
		{
			"get nameservice module accounts balance",
			[]string{s.GetBondId(), fmt.Sprintf("--%s=json", tmcli.OutputFlag)},
			false,
			1,
		},
	}

	for _, tc := range testCasesForNameServiceModuleBalance {
		s.Run(tc.name, func() {
			cmd := cli.GetCmdBalance()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expErr {
				sr.Error(err)
			} else {
				sr.NoError(err)
				var response types.GetNameServiceModuleBalanceResponse
				err = clientCtx.Codec.UnmarshalJSON(out.Bytes(), &response)
				sr.NoError(err)
				sr.Equal(tc.noOfRecords, len(response.GetBalances()))
				balance := response.GetBalances()[0]
				sr.Equal(balance.AccountName, types.RecordRentModuleAccountName)
			}
		})
	}
}
