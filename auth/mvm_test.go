package auth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pandodao/passport-go/eip4361"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	domain    = "localhost:6006"
	uri       = "https://service.invalid/login"
	nonce     = "12345678"
	issuedAt  = time.Now().Truncate(time.Second)
	expiredAt = issuedAt.Add(time.Second * 10)
	resources = "ipfs://bafybeiemxf5abjwjbikoz4mc3a3dla6ual3jsgpdr4cjr3oz3evfyavhwq/"
)

func makeMessage(t *testing.T, privateKey *ecdsa.PrivateKey, domain, uri, nonce, resources string, issuedAt, expiredAt time.Time) (string, string) {
	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	message := fmt.Sprintf(`%s wants you to sign in with your Ethereum account:
%s

I accept the ServiceOrg Terms of Service: https://service.invalid/tos

URI: %s
Version: 1
Chain ID: 1
Nonce: %s
Issued At: %s
Expiration Time: %s
Request ID: F369349D-9B66-4367-BAF2-AE9D83E0F9FA
Resources:
- %s`,
		domain,
		address,
		uri,
		nonce,
		issuedAt.Format(time.RFC3339),
		expiredAt.Format(time.RFC3339),
		resources,
	)
	signature, err := crypto.Sign(accounts.TextHash([]byte(message)), privateKey)
	require.NoError(t, err)

	return message, "0x" + hex.EncodeToString(signature)
}

func TestMvmAuth(t *testing.T) {
	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	auth := Authorizer{}

	t.Run("message", func(t *testing.T) {
		{
			// Domain
			message, signature := makeMessage(t, privateKey, domain, uri, nonce, resources, issuedAt, expiredAt)

			_, err := auth.AuthorizeMvmMessage(
				ctx,
				message,
				signature,
				NewDomainsValidator([]string{"localhost"}),
			)
			assert.True(t, IsErrBadDomain(err), fmt.Sprintf("expect %d, got %v", ErrCodeBadDomain, err))
		}

		{
			// expire at
			message, signature := makeMessage(t, privateKey, domain, uri, nonce, resources, issuedAt, expiredAt.AddDate(0, 0, -1))

			_, err := auth.AuthorizeMvmMessage(
				ctx,
				message,
				signature,
				nil,
			)
			assert.True(t, IsErrBadLoginMessage(err), fmt.Sprintf("expect %d, got %v", ErrCodeBadLoginMessage, err))
		}
	})

	t.Run("signature", func(t *testing.T) {
		message, signature := makeMessage(t, privateKey, domain, uri, nonce, resources, issuedAt, expiredAt)

		{
			_, err := auth.AuthorizeMvmMessage(ctx, message, "no signature", nil)
			assert.True(t, IsErrBadLoginSignature(err), fmt.Sprintf("expect %d, got %v", ErrCodeBadLoginSignature, err))
		}

		{
			_, err = auth.AuthorizeMvmMessage(
				ctx,
				message,
				signature,
				func(ctx context.Context, message *eip4361.Message) error {
					require.Equal(t, message.Domain, domain)
					require.Equal(t, message.IssuedAt, issuedAt.Format(time.RFC3339))
					require.Equal(t, message.Nonce, nonce)
					require.Equal(t, message.Resources[0], resources)
					return nil
				},
			)
			require.NoError(t, err)
		}
	})
}
