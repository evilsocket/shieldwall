package api

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/evilsocket/shieldwall/database"
	"time"
)

var (
	ErrTokenClaims       = errors.New("can't extract claims from jwt token")
	ErrTokenInvalid      = errors.New("jwt token not valid")
	ErrTokenExpired      = errors.New("jwt token expired")
	ErrTokenIncomplete   = errors.New("jwt token is missing required fields")
	ErrTokenUnauthorized = errors.New("jwt token authorized field is false (?!)")
)

func (api *API) tokenFor(user *database.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user.ID
	claims["expires_at"] = time.Now().Add(time.Duration(api.config.TokenTTL) * time.Second).Format(time.RFC3339)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(api.config.Secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (api *API) validateToken(header string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(api.config.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrTokenClaims
	} else if !token.Valid {
		return nil, ErrTokenInvalid
	}

	required := []string{
		"expires_at",
		"authorized",
		"user_id",
	}
	for _, req := range required {
		if _, found := claims[req]; !found {
			return nil, ErrTokenIncomplete
		}
	}

	// log.Debug("%+v", claims)

	if expiresAt, err := time.Parse(time.RFC3339, claims["expires_at"].(string)); err != nil {
		return nil, ErrTokenExpired
	} else if expiresAt.Before(time.Now()) {
		return nil, ErrTokenExpired
	} else if claims["authorized"].(bool) != true {
		return nil, ErrTokenUnauthorized
	}
	return claims, err
}