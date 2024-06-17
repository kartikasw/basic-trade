package token

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTError string

func (e JWTError) Error() string {
	return string(e)
}

const JWTExpirationError = JWTError("token is expired")

type Claim struct {
	TokenID uuid.UUID `json:"token_id"`
	UserID  uuid.UUID `json:"user_id"`
	Scope   Scope     `json:"scope"`
	jwt.StandardClaims
}

type Scope string

const (
	ScopeAccess  = Scope("access")
	ScopeRefresh = Scope("refresh")
)

type Expectation func(parsed Claim) error

type JWTToken struct {
	SignedToken string
	Claim       Claim
	ExpireAt    time.Time
	Scheme      string
}
