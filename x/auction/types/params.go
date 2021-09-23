package types

import (
	"bytes"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Default parameter namespace.
const (
	DefaultParamspace = ModuleName
)

var _ subspace.ParamSet = Params{}

func NewParams() Params {
	return Params{}
}

// ParamKeyTable - ParamTable for bond module.
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs - implements params.ParamSet
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := ModuleCdc.MustMarshalLengthPrefixed(&p)
	bz2 := ModuleCdc.MustMarshalLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{}
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
