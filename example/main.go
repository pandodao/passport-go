package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pandodao/passport-go/auth"
	"github.com/pandodao/passport-go/eip4361"
)

func main() {
	ctx := context.Background()
	authorizer := auth.New(
		[]string{"dapp_id"}, // issuers whitelist, oauth token must be issued by the issuers
		[]string{},          // domain whitelist
	)

	message, signature, err := makeMessage()
	if err != nil {
		panic(err)
	}

	user, err := authorizer.AuthorizeWithMvmValidator(
		ctx,
		&auth.AuthorizationParams{
			Method:           auth.AuthMethodMvm,
			MvmSignedMessage: message,
			MvmSignature:     signature,
		},
		func(ctx context.Context, message *eip4361.Message) error {
			// verify anything you want
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	bts, _ := json.MarshalIndent(user, "", "    ")
	fmt.Println(string(bts))
}

func makeMessage() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}

	var (
		address   = crypto.PubkeyToAddress(privateKey.PublicKey)
		domain    = "pando-apps.aspens.rocks"
		uri       = "https://pando-apps.aspens.rocks"
		nonce     = "oCxDubPgiNZdE8z71"
		issuedAt  = time.Now()
		expiredAt = issuedAt.Add(time.Hour)
		resources = "https://pando-apps.aspens.rocks"
	)

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
	return message, "0x" + hex.EncodeToString(signature), err
}
