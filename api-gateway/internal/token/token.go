package token

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	SecretKey string
}

func NewJWTMaker(secret string) *JWTMaker {
	return &JWTMaker{SecretKey: secret}
}

type CustomClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (j *JWTMaker) VerifyToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Optional: make sure token uses HMAC SHA512
		// if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		// }
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (j *JWTMaker) JWTAuthMiddleware(jwtMaker *JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		claims, err := jwtMaker.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
