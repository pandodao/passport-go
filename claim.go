package passport

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/passport-go/mvm"
	"github.com/golang-jwt/jwt"
)

type Claim struct {
	jwt.StandardClaims
	Nonce   uint64 `json:"non,omitempty"`
	Version int64  `json:"ver,omitempty"`
	Scope   string `json:"scp,omitempty"`
}

func (c *Claim) Address() common.Address {
	return common.HexToAddress(c.Issuer)
}

func (c Claim) Valid() error {
	if mvm.IsNullAddress(c.Address()) {
		return jwt.NewValidationError("invalid issuer", jwt.ValidationErrorIssuer)
	}

	return c.StandardClaims.Valid()
}
