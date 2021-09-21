package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tharsis/ethermint/x/nameservice/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bond MsgServer interface for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) SetRecord(c context.Context, msg *types.MsgSetRecord) (*types.MsgSetRecordResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}

	_, err = m.Keeper.ProcessSetRecord(ctx, types.MsgSetRecord{
		BondId:  msg.GetBondId(),
		Signer:  msg.GetSigner(),
		Payload: msg.GetPayload(),
	})
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetRecord,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})

	return &types.MsgSetRecordResponse{}, nil
}

func (m msgServer) SetName(c context.Context, msg *types.MsgSetName) (*types.MsgSetNameResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	err = m.Keeper.ProcessSetName(ctx, *msg)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSetRecord,
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
			sdk.NewAttribute(types.AttributeKeyWRN, msg.Wrn),
			sdk.NewAttribute(types.AttributeKeyCID, msg.Cid),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySigner, msg.Signer),
		),
	})
	return &types.MsgSetNameResponse{}, nil
}
