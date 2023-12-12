package auth

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go/v2"
)

const (
	AuthMethodMixinToken AuthMethod = "mixin_token"
	AuthMethodMvm        AuthMethod = "mvm"
)

type (
	Authorizer struct {
		issuers []string
		domains []string
	}

	User struct {
		mixin.User
		MvmAddress common.Address
	}

	AuthMethod string

	AuthorizationParams struct {
		Method           AuthMethod `json:"method"`
		MixinToken       string     `json:"mixin_token"`
		MvmSignedMessage string     `json:"mvm_signed_message"`
		MvmSignature     string     `json:"mvm_signature"`
	}
)

func New(issuers, domains []string) *Authorizer {
	return &Authorizer{
		issuers: issuers,
		domains: domains,
	}
}

func (a *Authorizer) Authorize(ctx context.Context, params *AuthorizationParams) (*User, error) {
	var validator MvmValidator
	switch params.Method {
	case AuthMethodMvm:
		if len(a.domains) > 0 {
			validator = NewDomainsValidator(a.domains)
		}
	}
	return a.AuthorizeWithMvmValidator(ctx, params, validator)
}

func (a *Authorizer) AuthorizeWithMvmValidator(
	ctx context.Context,
	params *AuthorizationParams,
	validator MvmValidator,
) (*User, error) {
	switch params.Method {
	case AuthMethodMixinToken:
		return a.AuthorizeMixinToken(ctx, params.MixinToken)
	case AuthMethodMvm:
		return a.AuthorizeMvmMessage(ctx, params.MvmSignedMessage, params.MvmSignature, validator)
	default:
		return nil, NewBadLoginMethodError("unknown method " + string(params.Method))
	}
}
