package passport

import "github.com/golang-jwt/jwt"

func ParseMVMToken(token string) (*Claim, error) {
	var c Claim
	if _, err := jwt.ParseWithClaims(token, &c, func(token *jwt.Token) (interface{}, error) {
		c := token.Claims.(*Claim)
		return c.Address(), nil
	}); err != nil {
		return nil, err
	}

	return &c, nil
}
