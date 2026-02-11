package token

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	secretKey string
}

func NewJwtMaker(secretKey string) (*JWT, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("secret key is empty")
	}
	return &JWT{secretKey: secretKey}, nil
}

type CustomClaims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (j *JWT) CreateToken(userID int64, duration time.Duration) (string, error) {
	fmt.Println("Creating token for userID:", userID, "and duration:", duration)
	claims := CustomClaims{
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWT) VerifyToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Optional: make sure token uses HMAC SHA512
		// if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		// }
		if t.Method != jwt.SigningMethodHS512 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(j.secretKey), nil
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

func (j *JWT) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		claims, err := j.VerifyToken(cookie.Value)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
