package mvm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-resty/resty/v2"
)

var (
	client = resty.New().SetBaseURL("https://api.mvm.dev")

	nullAddress = common.Address{}
)

func UserToAddress(ctx context.Context, userID string) (common.Address, error) {
	resp, err := client.R().SetContext(ctx).SetQueryParam("user", userID).Get("/user_contract")
	if err != nil {
		return common.Address{}, err
	}

	if resp.IsError() {
		return common.Address{}, fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	var body struct {
		UserContract common.Address `json:"user_contract"`
	}

	_ = json.Unmarshal(resp.Body(), &body)
	return body.UserContract, nil
}

func IsNullAddress(addr common.Address) bool {
	return addr == nullAddress
}

func AddressToUser(ctx context.Context, addr common.Address) (string, error) {
	resp, err := client.R().SetContext(ctx).SetQueryParam("contract", addr.String()).Get("/contract_user")
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	var body struct {
		Users     []string `json:"users"`
		Threshold int      `json:"threshold"`
	}

	_ = json.Unmarshal(resp.Body(), &body)

	if len(body.Users) != 1 {
		return "", fmt.Errorf("expected 1 user, got %d", len(body.Users))
	}

	return body.Users[0], nil
}
