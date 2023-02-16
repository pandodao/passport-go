package mvm

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go"
)

type BridgeUser struct {
	UserID    string         `json:"user_id"`
	FullName  string         `json:"full_name"`
	CreatedAt time.Time      `json:"created_at"`
	Key       mixin.Keystore `json:"key"`
	Contract  common.Address `json:"contract"`
}

func GetBridgeUser(ctx context.Context, addr common.Address) (*BridgeUser, error) {
	b := map[string]interface{}{
		"public_key": addr.String(),
	}

	resp, err := client.R().SetContext(ctx).SetBody(b).Post(getBridgeUserURL)
	if err != nil {
		return nil, err
	}

	var body struct {
		User BridgeUser `json:"user"`
	}

	if err := unmarshalResponse(resp, &body); err != nil {
		return nil, err
	}

	return &body.User, nil
}
