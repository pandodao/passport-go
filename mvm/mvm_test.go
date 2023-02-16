package mvm

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI(t *testing.T) {
	ctx := context.Background()

	addr := common.HexToAddress("0xE2aD78Fdf6C29338f5E2434380740ac889457256")

	user, err := GetBridgeUser(ctx, addr)
	require.NoError(t, err)

	contract, err := GetUserContract(ctx, user.UserID)
	require.NoError(t, err)

	assert.Equal(t, user.Contract, contract)

	userID, err := GetContractUser(ctx, user.Contract)
	require.NoError(t, err)
	assert.Equal(t, user.UserID, userID)
}
