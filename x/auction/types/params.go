package types

import (
	"bytes"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter namespace.
const (
	DefaultParamspace = ModuleName
)

var _ types.ParamSet = Params{}

func NewParams() Params {
	return Params{}
}

// ParamKeyTable - ParamTable for bond module.
func ParamKeyTable() types.KeyTable {
	return types.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs - implements params.ParamSet
func (p Params) ParamSetPairs() types.ParamSetPairs {
	return types.ParamSetPairs{}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		CommitsDuration: 24 * time.Hour,
		RevealsDuration: 24 * time.Hour,
		CommitFee: sdk.Coin{
			Amount: sdk.NewInt(10),
			Denom:  "aphoton",
		},
		RevealFee: sdk.Coin{
			Amount: sdk.NewInt(10),
			Denom:  "aphoton",
		},
		MinimumBid: sdk.Coin{
			Amount: sdk.NewInt(1000),
			Denom:  "aphoton",
		},
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	return sb.String()
}

// Validate a set of params.
func (p Params) Validate() error {
	return nil
}
