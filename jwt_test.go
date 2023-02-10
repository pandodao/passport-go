package passport

import (
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/storyicon/sigverify"
)

var (
	msg  = []byte("abc")
	addr = "0xf2d5c11730d20cdbb2d54294fcd9980631788f18"
	sig  = hexutil.MustDecode("0x583b40db49b5a43306b8e90daa9b4a55149f626ca6754e6c00ff5518e1fcb8a0715cc1906eb61cce841b54df2e271432b8d95a69d898869efa0cfb9e3ba119fa1c")
)

func TestHash(t *testing.T) {
	r, err := sigverify.EcRecoverEx(msg, sig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(r.String())
}

func TestGenerateToken(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	var c Claim
	c.Issuer = addr.String()
	c.Id = uuid.NewString()
	c.ExpiresAt = time.Now().Add(time.Hour).Unix()
	c.Audience = "a753e0eb-3010-4c4a-a7b2-a7bda4063f62"
	c.Nonce = 1
	c.Version = 1
	token := jwt.NewWithClaims(&MVMSigningMethod{}, c)
	tokenString, err := token.SignedString(key)
	if err != nil {
		t.Fatal("sign failed", err)
	}

	t.Log(tokenString)
	parsedClaim, err := ParseMVMToken(tokenString)
	if err != nil {
		t.Fatal("parse failed", err)
	}

	if parsedClaim.Issuer != addr.String() {
		t.Fatal("issuer not match")
	}
}
