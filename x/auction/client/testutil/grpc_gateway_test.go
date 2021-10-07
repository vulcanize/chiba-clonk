package testutil

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/testutil"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/cosmos/cosmos-sdk/types/rest"
	auctiontypes "github.com/tharsis/ethermint/x/auction/types"
)

const (
	randomAuctionID     = "randomAuctionID"
	randomBidderAddress = "randomBidderAddress"
	randomOwnerAddress  = "randomOwnerAddress"
)

func (suite *IntegrationTestSuite) TestGetAllAuctionsGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		headers         map[string]string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctions", val.APIAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctions", val.APIAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var auctions auctiontypes.AuctionsResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &auctions)

			if tc.isErrorExpected {
				//sr.Empty(auctions.GetAuctions().Auctions)
				sr.NotNil(err)
			} else {
				sr.NoError(err)
				sr.Len(auctions.GetAuctions(), tc.respObjCount)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryParamsGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/params", val.APIAddress),
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := rest.GetRequest(tc.url)
			sr.NoError(err)

			var params auctiontypes.QueryParamsResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &params)

			if tc.isErrorExpected {
				//sr.Empty(params.Params)
				sr.NotNil(err)
			} else {
				sr.NoError(err)
				sr.Equal(auctiontypes.DefaultParams(), *(params.Params))
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetAuctionGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		headers         map[string]string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctions/%s", val.APIAddress, randomAuctionID),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctions/%s", val.APIAddress, randomAuctionID),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var auction auctiontypes.AuctionResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &auction)

			if tc.isErrorExpected {
				//sr.Empty(auction.Auction)
				sr.NotNil(err)
			} else {
				sr.NoError(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetBidsGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		headers         map[string]string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/bids/%s", val.APIAddress, randomAuctionID),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/bids/%s", val.APIAddress, randomAuctionID),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var bids auctiontypes.BidsResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bids)

			if tc.isErrorExpected {
				//sr.Empty(bids.GetBids())
				sr.NotNil(err)
			} else {
				sr.NoError(err)
				sr.Len(bids.GetBids(), tc.respObjCount)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetBidGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		headers         map[string]string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/bids/%s/%s", val.APIAddress, randomAuctionID, randomBidderAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/bids/%s/%s", val.APIAddress, randomAuctionID, randomBidderAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var bid auctiontypes.BidResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &bid)

			if tc.isErrorExpected {
				//sr.Empty(bid.GetBid())
				sr.NotNil(err)
			} else {
				sr.NoError(err)
				sr.Len(bid.GetBid(), tc.respObjCount)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetAuctionsByBidderGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		headers         map[string]string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctionsbybidder/%s", val.APIAddress, randomBidderAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctionsbybidder/%s", val.APIAddress, randomBidderAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var auctions auctiontypes.AuctionsByBidderResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &auctions)

			if tc.isErrorExpected {
				//sr.Empty(auctions.GetAuctions())
				sr.NotNil(err)
			} else {
				sr.NoError(err)
				sr.Len(auctions.GetAuctions(), tc.respObjCount)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGetAuctionsByOwnerGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		headers         map[string]string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctionsbyowner/%s", val.APIAddress, randomOwnerAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/auctionsbyowner/%s", val.APIAddress, randomOwnerAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var auctions auctiontypes.AuctionsByOwnerResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &auctions)

			if tc.isErrorExpected {
				sr.Empty(auctions.GetAuctions())
				sr.NotNil(err)
			} else {
				sr.NoError(err)
				sr.Len(auctions.GetAuctions(), tc.respObjCount)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestQueryBalanceGrpc() {
	val := suite.network.Validators[0]
	sr := suite.Require()
	testCases := []struct {
		msg             string
		url             string
		headers         map[string]string
		respObjCount    int
		isErrorExpected bool
	}{
		{
			"invalid request with headers",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/balance", val.APIAddress),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			0,
			true,
		},
		{
			"valid request",
			fmt.Sprintf("%s/vulcanize/auction/v1beta1/balance", val.APIAddress),
			map[string]string{},
			0,
			false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
			sr.NoError(err)

			var balance auctiontypes.BalanceResponse
			err = val.ClientCtx.Codec.UnmarshalJSON(resp, &balance)

			sr.Empty(balance.GetBalance())
			if tc.isErrorExpected {
				sr.NotNil(err)
			} else {
				sr.NoError(err)
			}
		})
	}
}
