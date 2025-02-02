package service

import (
	"delivery/configs"
	"delivery/logger"
	"fmt"
	"os"

	"time"

	"github.com/dgrijalva/jwt-go"
)

// CustomClaims defines custom token fields
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Role     string `json:"role"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

// GenerateToken generates a JWT token with custom fields
func GenerateToken(userID uint, username, role, email string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Email:    email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(configs.AppSettings.AuthParams.JwtTtlMinutes)).Unix(),
			Issuer:    configs.AppSettings.AppParams.ServerName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

// ParseToken parses JWT token and returns custom fields
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Checking the token signature method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		logger.Error.Println("[service.ParseToken] cannot parse token. Error is:", err.Error())
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	logger.Error.Println("[service.ParseToken] invalid token")
	return nil, fmt.Errorf("invalid token")
}
