package auth

import (
	"context"

	"github.com/asaskevich/govalidator"
	"github.com/pandodao/passport-go/eip4361"
)

type (
	MvmValidator func(ctx context.Context, message *eip4361.Message) (bool, error)
)

func NewDomainsValidator(domains []string) MvmValidator {
	return func(ctx context.Context, message *eip4361.Message) (bool, error) {
		if !govalidator.IsIn(message.Domain, domains...) {
			return false, ErrInvalidDomain
		}
		return true, nil
	}
}
