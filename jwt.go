package passport

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/storyicon/sigverify"
)

const SigningMethod = "ESMVM"

func init() {
	jwt.RegisterSigningMethod(SigningMethod, func() jwt.SigningMethod {
		return &MVMSigningMethod{}
	})
}

type MVMSigningMethod struct{}

func (sm *MVMSigningMethod) Verify(signingString, signature string, key interface{}) error {
	var pub common.Address
	switch k := key.(type) {
	case common.Address:
		pub = k
	default:
		return jwt.ErrInvalidKeyType
	}

	sig, err := jwt.DecodeSegment(signature)
	if err != nil {
		return err
	}

	addr, err := sigverify.EcRecoverEx([]byte(signingString), sig)
	if err != nil {
		return err
	}

	if addr != pub {
		return jwt.ErrSignatureInvalid
	}

	return nil
}

func (sm *MVMSigningMethod) Sign(signingString string, key interface{}) (string, error) {
	var signKey *ecdsa.PrivateKey
	switch k := key.(type) {
	case *ecdsa.PrivateKey:
		signKey = k
	default:
		return "", jwt.ErrInvalidKeyType
	}

	b, _ := accounts.TextAndHash([]byte(signingString))
	sig, err := crypto.Sign(b, signKey)
	if err != nil {
		return "", err
	}

	return jwt.EncodeSegment(sig), nil
}

func (sm *MVMSigningMethod) Alg() string {
	return SigningMethod
}
