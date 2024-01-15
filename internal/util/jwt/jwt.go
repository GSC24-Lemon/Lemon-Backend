package jwt

import (
	"errors"
	"fmt"
	"lemon_be/internal/controller/http/errorWrapper"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const minSecretKeySize = 32

// JWTMaker JWT maker
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (JwtTokenMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken membuat jwt token baru berdurasi utk user
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(maker.secretKey))
	return token, payload, err
}

// VerifyToken cek jika token valid atau tidak
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errorWrapper.NewHTTPError(ErrInvalidToken, 401, "token is invalid")
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, errorWrapper.NewHTTPError(ErrExpiredToken, 401, "token has expired")
		}
		return nil, errorWrapper.NewHTTPError(ErrInvalidToken, 401, "token is invalid")
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errorWrapper.NewHTTPError(ErrInvalidToken, 401, "token is invalid")
	}

	return payload, nil
}
