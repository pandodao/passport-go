package auth

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go"
)

type (
	Authorizer struct {
		Issuers []string
	}

	User struct {
		mixin.User
		MvmAddress common.Address
	}

	AuthorizationParams struct {
		Method           string
		MixinToken       string
		MvmSignedMessage string
		MvmSignature     string
	}
)

func New(issuers []string) *Authorizer {
	return &Authorizer{
		Issuers: issuers,
	}
}

func (a *Authorizer) Authorize(ctx context.Context, params *AuthorizationParams) (*User, error) {
	switch params.Method {
	case "mixin_token":
		return a.AuthorizeMixinToken(ctx, params.MixinToken)
	case "mvm":
		return a.AuthorizeMvmMessage(ctx, params.MvmSignedMessage, params.MvmSignature)
	default:
		return nil, ErrBadLoginMethod
	}
}
