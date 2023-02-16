package mvm

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

const (
	ChainID = 73927

	RPCEndpoint    = "https://geth.mvm.dev"
	APIEndpoint    = "https://api.mvm.dev"
	BridgeEndpoint = "https://bridge.mvm.dev"
)

var (
	client = resty.New()

	getUserContractURL, _ = url.JoinPath(APIEndpoint, "/user_contract")
	getContractUserURL, _ = url.JoinPath(APIEndpoint, "/contract_user")
	getBridgeUserURL, _   = url.JoinPath(BridgeEndpoint, "/users")
)

func unmarshalResponse(r *resty.Response, v any) error {
	if r.IsError() {
		return fmt.Errorf("status: %d, body: %s", r.StatusCode(), r.Body())
	}

	return json.Unmarshal(r.Body(), v)
}
