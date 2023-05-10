package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/fox-one/mixin-sdk-go"
	"github.com/pandodao/passport-go/eip4361"
	"github.com/pandodao/passport-go/mvm"
)

func (a *Authorizer) AuthorizeMvmMessage(
	ctx context.Context,
	signedMessage, signature string,
	validator MvmValidator,
) (*User, error) {
	message, err := eip4361.Parse(signedMessage)
	if err != nil {
		return nil, NewBadLoginMessageError("")
	}

	if err := message.Validate(time.Now()); err != nil {
		return nil, NewBadLoginMessageError(fmt.Sprintf("validate failed (%v)", err))
	}

	if err := eip4361.Verify(message, signature); err != nil {
		return nil, NewBadLoginSignatureError(err.Error())
	}

	if validator != nil {
		if err := validator(ctx, message); err != nil {
			return nil, NewBadLoginMessageError(fmt.Sprintf("custom validate failed (%v)", err))
		}
	}

	addr := common.HexToAddress(message.Address)
	mvmUser, err := mvm.GetBridgeUser(ctx, addr)
	if err != nil {
		return nil, NewError(fmt.Sprintf("read bridge user failed (%v)", err.Error()))
	}

	cli, err := mixin.NewFromKeystore(&mvmUser.Key)
	if err != nil {
		return nil, NewError(fmt.Sprintf("load bridge keystore failed (%v)", err))
	}

	user, err := cli.UserMe(ctx)
	if err != nil {
		return nil, NewError(fmt.Sprintf("read user profile failed (%v)", err))
	}

	return &User{User: *user, MvmAddress: addr}, nil
}
