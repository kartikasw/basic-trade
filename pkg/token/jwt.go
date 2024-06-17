package token

import (
	"basic-trade/pkg/config"
	"crypto/rsa"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	AuthorizationHeader = "authorization"
	BearerScheme        = "bearer"
	JWTClaim            = "jwt_claim"
)

type JWT interface {
	CreateAccessToken(usrUUID uuid.UUID) (*JWTToken, error)
	CreateRefreshToken(usrUUID uuid.UUID) (*JWTToken, error) // Currently not used
	VerifyToken(token string, expectation ...Expectation) (*Claim, error)
}

type IJWT struct {
	cfg        config.Token
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func NewJWT(cfg config.Token) (JWT, error) {
	j := &IJWT{cfg: cfg}

	var err error
	j.privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(cfg.PrivateKey))
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse private key: %w", err)
	}

	j.publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.PublicKey))
	if err != nil {
		return nil, fmt.Errorf("Couldn't parse public key: %w", err)
	}

	return j, nil
}

func (j *IJWT) createJWTToken(claim Claim, duration time.Duration) (*JWTToken, error) {
	now := time.Now()
	expAt := now.Add(duration)
	exp := expAt.Unix()
	iat := now.Unix()

	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("Couldn't generate random UUID: %w", err)
	}

	claim.TokenID = uuid
	claim.StandardClaims = jwt.StandardClaims{
		ExpiresAt: exp,
		IssuedAt:  iat,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		return nil, fmt.Errorf("Couldn't sign token: %w", err)
	}

	jwtToken := &JWTToken{
		SignedToken: signedToken,
		Claim:       claim,
		ExpireAt:    expAt,
		Scheme:      BearerScheme,
	}

	return jwtToken, nil
}

func AccessTokenExpectation() Expectation {
	return Expectation(func(parsed Claim) error {
		if parsed.Scope != Scope(ScopeAccess) {
			return fmt.Errorf("Expect scope %s, have: %s", ScopeAccess, parsed.Scope)
		}
		return nil
	})
}

func (j *IJWT) CreateAccessToken(usrUUID uuid.UUID) (*JWTToken, error) {
	claim := Claim{
		UserID: usrUUID,
		Scope:  Scope(ScopeAccess),
	}

	return j.createJWTToken(claim, j.cfg.AccessTokenDuration)
}

func RefreshTokenExpectation() Expectation {
	return Expectation(func(parsed Claim) error {
		if parsed.Scope != ScopeRefresh {
			return fmt.Errorf("Expect scope %s, have: %s", ScopeRefresh, parsed.Scope)
		}
		return nil
	})
}

func (j *IJWT) CreateRefreshToken(usrUUID uuid.UUID) (*JWTToken, error) {
	claim := Claim{
		UserID: usrUUID,
		Scope:  Scope(ScopeRefresh),
	}

	return j.createJWTToken(claim, j.cfg.RefreshTokenDuration)
}

func (j *IJWT) VerifyToken(token string, expectations ...Expectation) (*Claim, error) {
	c := &Claim{}

	_, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return j.publicKey, nil
	})

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)
		if ok {
			if validationErr.Errors == jwt.ValidationErrorExpired {
				return nil, JWTExpirationError
			}
		}
		return nil, fmt.Errorf("Couldn't ParseWithClaims: %w", err)
	}

	if c.TokenID.String() == "" {
		return nil, fmt.Errorf("Empty token_id claim")
	}

	if strings.TrimSpace(c.UserID.String()) == "" {
		return nil, fmt.Errorf("Empty user_id claim")
	}

	for i, e := range expectations {
		err := e(*c)
		if err != nil {
			return nil, fmt.Errorf("Failed expectation: %d: %w", i, err)
		}
	}

	return c, nil
}
