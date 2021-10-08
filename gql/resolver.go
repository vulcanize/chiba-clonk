package gql

import (
	"context"
	"encoding/base64"
	"github.com/cosmos/cosmos-sdk/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	bondtypes "github.com/tharsis/ethermint/x/bond/types"
	"strconv"
)

// DefaultLogNumLines is the number of log lines to tail by default.
const DefaultLogNumLines = 50

// MaxLogNumLines is the max number of log lines that can be tailed.
const MaxLogNumLines = 1000

type Resolver struct {
	ctx     client.Context
	logFile string
}

// Query is the entry point to query execution.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type queryResolver struct{ *Resolver }

func (q queryResolver) GetStatus(ctx context.Context) (*Status, error) {
	nodeInfo, syncInfo, validatorInfo, err := getStatusInfo(q.ctx)
	if err != nil {
		return nil, err
	}

	numPeers, peers, err := getNetInfo(q.ctx)
	if err != nil {
		return nil, err
	}

	validatorSet, err := getValidatorSet(q.ctx)
	if err != nil {
		return nil, err
	}

	diskUsage, err := GetDiskUsage(NodeDataPath)
	if err != nil {
		return nil, err
	}

	return &Status{
		Version:    NameServiceVersion,
		Node:       nodeInfo,
		Sync:       syncInfo,
		Validator:  validatorInfo,
		Validators: validatorSet,
		NumPeers:   numPeers,
		Peers:      peers,
		DiskUsage:  diskUsage,
	}, nil
}

func (q queryResolver) GetAccounts(ctx context.Context, addresses []string) ([]*Account, error) {
	accounts := make([]*Account, len(addresses))
	for index, address := range addresses {
		account, err := q.GetAccount(ctx, address)
		if err != nil {
			return nil, err
		}
		accounts[index] = account
	}
	return accounts, nil
}

func (q queryResolver) GetAccount(ctx context.Context, address string) (*Account, error) {
	authQueryClient := authtypes.NewQueryClient(q.ctx)
	accountResponse, err := authQueryClient.Account(ctx, &authtypes.QueryAccountRequest{Address: address})
	if err != nil {
		return nil, err
	}
	var account authtypes.AccountI
	err = q.ctx.Codec.UnpackAny(accountResponse.GetAccount(), &account)
	if err != nil {
		return nil, err
	}
	var pubKey *string
	if account.GetPubKey() != nil {
		pubKeyStr := base64.StdEncoding.EncodeToString(account.GetPubKey().Bytes())
		pubKey = &pubKeyStr
	}

	// Get the account balance
	bankQueryClient := banktypes.NewQueryClient(q.ctx)
	balance, err := bankQueryClient.AllBalances(ctx, &banktypes.QueryAllBalancesRequest{Address: address})

	accNum := strconv.FormatUint(account.GetAccountNumber(), 10)
	seq := strconv.FormatUint(account.GetSequence(), 10)

	return &Account{
		Address:  address,
		Number:   accNum,
		Sequence: seq,
		PubKey:   pubKey,
		Balance:  getGQLCoins(balance.GetBalances()),
	}, nil
}

func (q queryResolver) GetBondsByIds(ctx context.Context, ids []string) ([]*Bond, error) {
	bonds := make([]*Bond, len(ids))
	for index, id := range ids {
		bondObj, err := q.GetBond(ctx, id)
		if err != nil {
			return nil, err
		}
		bonds[index] = bondObj
	}

	return bonds, nil
}

func (q *queryResolver) GetBond(ctx context.Context, id string) (*Bond, error) {
	bondQueryClient := bondtypes.NewQueryClient(q.ctx)
	bondResp, err := bondQueryClient.GetBondById(context.Background(), &bondtypes.QueryGetBondByIdRequest{Id: id})
	if err != nil {
		return nil, err
	}

	bond := bondResp.GetBond()
	if bond == nil {
		return nil, nil
	}
	return getGQLBond(bondResp.GetBond())
}

func (q queryResolver) QueryBonds(ctx context.Context, attributes []*KeyValueInput) ([]*Bond, error) {
	bondQueryClient := bondtypes.NewQueryClient(q.ctx)
	bonds, err := bondQueryClient.Bonds(context.Background(), &bondtypes.QueryGetBondsRequest{})
	if err != nil {
		return nil, err
	}

	gqlResponse := make([]*Bond, len(bonds.GetBonds()))
	for i, bondObj := range bonds.GetBonds() {
		gqlBond, err := getGQLBond(bondObj)
		if err != nil {
			return nil, err
		}
		gqlResponse[i] = gqlBond
	}

	return gqlResponse, nil
}

// QueryBondsByOwner will return bonds by owner
func (q queryResolver) QueryBondsByOwner(ctx context.Context, ownerAddresses []string) ([]*OwnerBonds, error) {
	ownerBonds := make([]*OwnerBonds, len(ownerAddresses))
	for index, ownerAddress := range ownerAddresses {
		bondsObj, err := q.GetBondsByOwner(ctx, ownerAddress)
		if err != nil {
			return nil, err
		}
		ownerBonds[index] = bondsObj
	}

	return ownerBonds, nil
}

func (q queryResolver) GetBondsByOwner(ctx context.Context, address string) (*OwnerBonds, error) {
	bondQueryClient := bondtypes.NewQueryClient(q.ctx)
	bondResp, err := bondQueryClient.GetBondsByOwner(context.Background(), &bondtypes.QueryGetBondsByOwnerRequest{Owner: address})
	if err != nil {
		return nil, err
	}

	ownerBonds := make([]*Bond, len(bondResp.GetBonds()))
	for i, bond := range bondResp.GetBonds() {
		bondObj, err := getGQLBond(&bond)
		if err != nil {
			return nil, err
		}
		ownerBonds[i] = bondObj
	}

	return &OwnerBonds{Bonds: ownerBonds, Owner: address}, nil
}
