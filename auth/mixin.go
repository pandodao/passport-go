package auth

import (
	"bytes"
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pandodao/passport-go/mvm"
)

func (a *Authorizer) AuthorizeMixinToken(ctx context.Context, token string) (*User, error) {
	var claims struct {
		jwt.RegisteredClaims
		Scope string `json:"scp,omitempty"`
	}

	jwt.ParseWithClaims(token, &claims, nil)
	if claims.Scope != "FULL" && !govalidator.IsIn(claims.Issuer, a.issuers...) {
		return nil, NewBadIssuerError("")
	}

	user, err := mixin.NewFromAccessToken(token).UserMe(ctx)
	if err != nil {
		return nil, NewError(fmt.Sprintf("read user profile failed (%v)", err.Error()))
	}

	contractAddr, err := mvm.GetUserContract(ctx, user.UserID)
	if err != nil {
		return nil, NewError(fmt.Sprintf("read bridge user failed (%v)", err.Error()))
	}

	emptyAddr := common.Address{}
	if !bytes.Equal(contractAddr[:], emptyAddr[:]) {
		return nil, NewBadLoginMethodError("mvm user")
	}

	return &User{User: *user}, nil
}
