package auth

import (
	"bytes"
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/passport-go/mvm"
	"github.com/golang-jwt/jwt/v4"
)

func (a *Authorizer) AuthorizeMixinToken(ctx context.Context, token string) (*User, error) {
	var claims struct {
		jwt.RegisteredClaims
		Scope string `json:"scp,omitempty"`
	}

	jwt.ParseWithClaims(token, &claims, nil)
	if claims.Scope != "FULL" && !govalidator.IsIn(claims.Issuer, a.issuers...) {
		return nil, ErrInvalidIssuer
	}

	user, err := mixin.NewFromAccessToken(token).UserMe(ctx)
	if err != nil {
		return nil, err
	}

	contractAddr, err := mvm.GetUserContract(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	emptyAddr := common.Address{}
	if bytes.Equal(contractAddr[:], emptyAddr[:]) {
		return nil, ErrBadLoginMethod
	}

	return &User{User: *user}, nil
}
