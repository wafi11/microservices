package pkg

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

type JwtToken struct {
	UserId string `json:"userId"`
	jwt.Claims
}

func VerifyToken(tokenString string) (*JwtToken, error) {
	secrets := []byte("f1404c536dbe0eb9b7ef959b09490eb0")

	tok, err := jwt.ParseSigned(tokenString, []jose.SignatureAlgorithm{jose.HS256})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %v", err)
	}

	claims := &JwtToken{}
	if err := tok.Claims(secrets, claims); err != nil {
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Baca dari Authorization header
		authHeader := c.GetHeader("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Fallback: baca dari cookie jika header tidak ada
		if tokenString == "" {
			cookie, err := c.Cookie("access_token")
			if err == nil {
				tokenString = cookie
			}
		}

		// 3. Tidak ada token sama sekali
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: no token provided",
			})
			return
		}

		// 4. Verifikasi token
		claims, err := VerifyToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: " + err.Error(),
			})
			return
		}

		// 5. Simpan claims ke context untuk dipakai di handler
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
