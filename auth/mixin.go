package auth

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pandodao/passport-go/mvm"
)

func init() {
	// bind process name to mixin client
	client := mixin.GetRestyClient()
	client.SetHeader("User-Agent", getUserAgent())
}

func getUserAgent() string {
	process := path.Base(os.Args[0])
	return fmt.Sprintf("mixin-sdk-go/v2/%s", process)
}

func (a *Authorizer) AuthorizeMixinToken(ctx context.Context, token string) (*User, error) {
	if err := validateMixinToken(token, a.issuers); err != nil {
		return nil, err
	}

	user, err := mixin.NewFromAccessToken(token).UserMe(ctx)
	if err != nil {
		return nil, NewError(fmt.Sprintf("read user profile failed (%v)", err.Error()))
	}

	if user.IdentityNumber == "" || user.IdentityNumber == "0" {
		contractAddr, err := mvm.GetUserContract(ctx, user.UserID)
		if err != nil {
			return nil, NewError(fmt.Sprintf("read bridge user failed (%v)", err.Error()))
		}

		emptyAddr := common.Address{}
		if !bytes.Equal(contractAddr[:], emptyAddr[:]) {
			return nil, NewBadLoginMethodError("mvm user")
		}
	}

	return &User{User: *user}, nil
}

func validateMixinToken(token string, issuers []string) error {
	var claims struct {
		jwt.RegisteredClaims
		Scope     string `json:"scp,omitempty"`
		UserID    string `json:"uid,omitempty"`
		Signature string `json:"sig,omitempty"`
	}

	jwt.ParseWithClaims(token, &claims, nil)

	if !claims.VerifyExpiresAt(time.Now(), true) {
		return NewBadMixinTokenError("expired")
	}

	if claims.Scope != "FULL" {
		if !govalidator.IsIn(claims.UserID, issuers...) {
			return NewBadIssuerError("")
		}
	} else {
		// user id is required
		if uuid.FromStringOrNil(claims.UserID).IsNil() {
			return NewBadMixinTokenError("invalid user id")
		}

		// sig is required
		if _, err := hex.DecodeString(claims.Signature); err != nil {
			return NewBadMixinTokenError("invalid signature")
		}
	}

	return nil
}
