package auth

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"time"
)

const (
	UserRoleCode      = 1
	AdminRoleCode     = 2
	DefaultExpireTime = time.Hour * 24
	Iss               = "IMDB-In-Golang"
	ContextKey        = "user"
)

var (
	ErrTokenExpired      = errors.New("token expired")
	ErrInvalidIssuedAt   = errors.New("token issue date is not valid")
	ErrInvalidClaimsType = errors.New("invalid claims type")
	ErrTokenNotFound     = errors.New("token not found")
)

type JwtClaim struct {
	Role int
	jwt.StandardClaims
}

func Validate(claim JwtClaim) error {
	if claim.ExpiresAt < time.Now().Unix() {
		return ErrTokenExpired
	}
	if claim.IssuedAt > time.Now().Unix() {
		return ErrInvalidIssuedAt
	}
	return nil
}

func ExtractJWT(c echo.Context) (*JwtClaim, error) {
	token, ok := c.Get(ContextKey).(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("failed to extract cliams: %w", ErrTokenNotFound)
	}
	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		return nil, fmt.Errorf("type conversion failed: %t, %w", token.Claims, ErrInvalidClaimsType)
	}
	return claims, nil
}

func GenerateJWT(secret, username string, isAdmin bool) (string, error) {
	var roleCode int
	if isAdmin {
		roleCode = AdminRoleCode
	} else {
		roleCode = UserRoleCode
	}
	claim := JwtClaim{
		roleCode,
		jwt.StandardClaims{
			Id:        username,
			Issuer:    Iss,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(DefaultExpireTime).Unix(),
		},
	}
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := rawToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("can not generate token: %w", err)
	}
	return token, nil
}
