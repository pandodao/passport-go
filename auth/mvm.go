package auth

import (
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/fox-one/passport-go/eip4361"
	"github.com/fox-one/passport-go/mvm"
)

func (a *Authorizer) AuthorizeMvmMessage(ctx context.Context, signedMessage, signature string) (*User, error) {
	message, err := eip4361.Parse(signedMessage)
	if err != nil {
		return nil, ErrBadMvmLoginMessage
	}

	if err := message.Validate(time.Now()); err != nil {
		return nil, ErrBadMvmLoginMessage
	}

	if err := eip4361.Verify(message, signature); err != nil {
		return nil, ErrBadMvmLoginSignature
	}

	if !govalidator.IsIn(message.Domain, a.Issuers...) {
		return nil, ErrInvalidIssuer
	}

	addr := common.HexToAddress(message.Address)
	mvmUser, err := mvm.GetBridgeUser(ctx, addr)
	if err != nil {
		return nil, err
	}

	cli, err := mixin.NewFromKeystore(&mvmUser.Key)
	if err != nil {
		return nil, err
	}

	user, err := cli.UserMe(ctx)
	if err != nil {
		return nil, err
	}

	return &User{User: *user, MvmAddress: addr}, nil
}
