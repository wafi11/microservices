package internal

import (
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type JwtToken struct {
	UserId string `json:"userId"`
	jwt.Claims
}

func NewJwtToken(userId string) *JwtToken {
	now := time.Now()
	return &JwtToken{
		UserId: userId,
		Claims: jwt.Claims{
			Issuer:    "wafiuddin",
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Expiry:    jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	}
}

func GenerateToken(userId string) (string, error) {
	secret := []byte("f1404c536dbe0eb9b7ef959b09490eb0")

	sig, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: secret},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return "", fmt.Errorf("could not create signer: %v", err)
	}

	token := NewJwtToken(userId)
	raw, err := jwt.Signed(sig).Claims(token).Serialize()
	if err != nil {
		return "", fmt.Errorf("could not serialize token: %v", err)
	}

	return raw, nil
}

func VerifyToken(tokenString string) (*JwtToken, error) {
	secret := []byte("f1404c536dbe0eb9b7ef959b09490eb0")

	tok, err := jwt.ParseSigned(tokenString, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	claims := &JwtToken{}
	if err := tok.Claims(secret, claims); err != nil {
		return nil, fmt.Errorf("could not verify token: %v", err)
	}

	if err := claims.Validate(jwt.Expected{
		Issuer: "wafiuddin",
		Time:   time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("token invalid: %v", err)
	}

	return claims, nil
}
