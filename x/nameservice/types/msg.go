package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = &MsgSetRecord{}
	_ sdk.Msg = &MsgSetName{}
)

// NewMsgSetRecord is the constructor function for MsgSetRecord.
func NewMsgSetRecord(payload Payload, bondID string, signer sdk.AccAddress) MsgSetRecord {
	return MsgSetRecord{
		Payload: payload,
		BondId:  bondID,
		Signer:  signer.String(),
	}
}

func (msg MsgSetRecord) ValidateBasic() error {
	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Signer)
	}
	owners := msg.Payload.Record.Owners
	for _, owner := range owners {
		if owner == "" {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Record owner not set.")
		}
	}

	if len(msg.BondId) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Bond ID is required.")
	}
	return nil
}

func (msg MsgSetRecord) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

// GetSignBytes gets the sign bytes for the msg MsgCreateBond
func (msg MsgSetRecord) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// NewMsgSetName is the constructor function for MsgSetName.
func NewMsgSetName(wrn string, cid string, signer sdk.AccAddress) MsgSetName {
	return MsgSetName{
		Wrn:    wrn,
		Cid:    cid,
		Signer: signer.String(),
	}
}

// Route Implements Msg.
func (msg MsgSetName) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetName) Type() string { return "set-name" }

// ValidateBasic Implements Msg.
func (msg MsgSetName) ValidateBasic() error {

	if msg.Wrn == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "WRN is required.")
	}

	if msg.Cid == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "CID is required.")
	}

	if len(msg.Signer) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid signer")
	}

	return nil
}

// GetSignBytes gets the sign bytes for the msg MsgSetName
func (msg MsgSetName) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners Implements Msg.
func (msg MsgSetName) GetSigners() []sdk.AccAddress {
	accAddr, _ := sdk.AccAddressFromBech32(msg.Signer)
	return []sdk.AccAddress{accAddr}
}
