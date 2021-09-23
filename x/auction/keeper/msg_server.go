package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tharsis/ethermint/x/auction/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (s msgServer) CreateAuction(c context.Context, msg *types.MsgCreateAuction) (*types.MsgCreateAuctionResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	resp, err := s.Keeper.CreateAuction(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAuction,
			sdk.NewAttribute(types.AttributeKeyCommitsDuration, strconv.FormatInt(msg.CommitsDuration, 10)),
			sdk.NewAttribute(types.AttributeKeyCommitFee, msg.CommitFee.String()),
			sdk.NewAttribute(types.AttributeKeyRevealFee, msg.RevealFee.String()),
			sdk.NewAttribute(types.AttributeKeyMinimumBid, msg.MinimumBid.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer.String()),
		),
	})

	return &types.MsgCreateAuctionResponse{Auction: resp}, nil
}

// CommitBid is the command for committing a bid
func (s msgServer) CommitBid(c context.Context, msg *types.MsgCommitBid) (*types.MsgCommitBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	resp, err := s.Keeper.CommitBid(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCommitBid,
			sdk.NewAttribute(types.AttributeKeyAuctionID, msg.AuctionId),
			sdk.NewAttribute(types.AttributeKeyCommitHash, msg.CommitHash),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer.String()),
		),
	})

	return &types.MsgCommitBidResponse{Auction: resp}, nil
}

//RevealBid is the command for revealing a bid
func (s msgServer) RevealBid(c context.Context, msg *types.MsgRevealBid) (*types.MsgRevealBidResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	resp, err := s.Keeper.RevealBid(ctx, *msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRevealBid,
			sdk.NewAttribute(types.AttributeKeyAuctionID, msg.AuctionId),
			sdk.NewAttribute(types.AttributeKeyReveal, msg.Reveal),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer.String()),
		),
	})

	return &types.MsgRevealBidResponse{Auction: resp}, nil
}
