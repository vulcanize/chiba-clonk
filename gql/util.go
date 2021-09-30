package gql

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
	"strconv"
)

// OwnerAttributeName denotes the owner attribute name for a bond.
const OwnerAttributeName = "owner"

// BondIDAttributeName denotes the record bond ID.
const BondIDAttributeName = "bondId"

// ExpiryTimeAttributeName denotes the record expiry time.
const ExpiryTimeAttributeName = "expiryTime"

func getGQLCoin(coin sdk.Coin) *Coin {
	gqlCoin := Coin{
		Type:     coin.Denom,
		Quantity: strconv.FormatInt(coin.Amount.Int64(), 10),
	}

	return &gqlCoin
}

func getGQLCoins(coins sdk.Coins) []*Coin {
	gqlCoins := make([]*Coin, len(coins))
	for index, coin := range coins {
		gqlCoins[index] = getGQLCoin(coin)
	}

	return gqlCoins
}

func getGQLBond(bondObj *bondtypes.Bond) (*Bond, error) {
	// Nil record.
	if bondObj == nil {
		return nil, nil
	}

	return &Bond{
		ID:      bondObj.Id,
		Owner:   bondObj.Owner,
		Balance: getGQLCoins(bondObj.Balance),
	}, nil
}

func matchBondOnAttributes(bondObj *bondtypes.Bond, attributes []*KeyValueInput) bool {
	for _, attr := range attributes {
		switch attr.Key {
		case OwnerAttributeName:
			{
				if attr.Value.String == nil || bondObj.Owner != *attr.Value.String {
					return false
				}
			}
		default:
			{
				// Only attributes explicitly listed in the switch are queryable.
				return false
			}
		}
	}

	return true
}
