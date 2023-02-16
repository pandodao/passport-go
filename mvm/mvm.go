package mvm

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

// GetUserContract get the contract address of the user
// zero address means the user has no contract
func GetUserContract(ctx context.Context, userID string) (common.Address, error) {
	resp, err := client.R().SetContext(ctx).SetQueryParam("user", userID).Get(getUserContractURL)
	if err != nil {
		return common.Address{}, err
	}

	var body struct {
		UserContract common.Address `json:"user_contract"`
	}

	if err := unmarshalResponse(resp, &body); err != nil {
		return common.Address{}, err
	}

	return body.UserContract, nil
}

func GetContractUser(ctx context.Context, addr common.Address) (string, error) {
	resp, err := client.R().SetContext(ctx).SetQueryParam("contract", addr.String()).Get(getContractUserURL)
	if err != nil {
		return "", err
	}

	var body struct {
		Users     []string `json:"users"`
		Threshold int      `json:"threshold"`
	}

	if err := unmarshalResponse(resp, &body); err != nil {
		return "", err
	}

	if len(body.Users) != 1 {
		return "", fmt.Errorf("expected 1 user, got %d", len(body.Users))
	}

	return body.Users[0], nil
}
