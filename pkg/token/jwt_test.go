package token

import (
	"basic-trade/common"
	"basic-trade/pkg/config"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestExpiredJWTToken(t *testing.T) {
	private, public := common.GenerateRSAKey(t)

	cfg := config.Token{
		AccessTokenDuration: -time.Minute,
		PrivateKey:          private,
		PublicKey:           public,
	}

	jwtImpl, err := NewJWT(cfg)
	require.NoError(t, err)

	usrUUID := uuid.New()
	token, err := jwtImpl.CreateAccessToken(usrUUID)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	claim, err := jwtImpl.VerifyToken(token.SignedToken, AccessTokenExpectation())
	require.Error(t, err)
	require.EqualError(t, err, JWTExpirationError.Error())
	require.Nil(t, claim)
}

func TestInvalidJWTToken(t *testing.T) {
	testCases := []struct {
		name        string
		token       func(privateKey string) string
		checkResult func(claim *Claim, err error)
	}{
		{
			name: "NoSigningMethod",
			token: func(privateKey string) string {
				jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, Claim{})
				token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
				require.NoError(t, err)
				return token
			},
			checkResult: func(claim *Claim, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "Couldn't ParseWithClaims: Unexpected signing method: none")
				require.Nil(t, claim)
			},
		},
		{
			name: "InvalidIDFormat",
			token: func(privateKey string) string {
				key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
				require.NoError(t, err)

				jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, Claim{})
				token, err := jwtToken.SignedString(key)
				require.NoError(t, err)
				return token
			},
			checkResult: func(claim *Claim, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "Empty token_id claim")
				require.Nil(t, claim)
			},
		},
		{
			name: "InvalidScope",
			token: func(privateKey string) string {
				key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
				require.NoError(t, err)

				jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, Claim{
					TokenID: uuid.New(),
					UserID:  uuid.New(),
					Scope:   ScopeRefresh,
				})
				token, err := jwtToken.SignedString(key)
				require.NoError(t, err)
				return token
			},
			checkResult: func(claim *Claim, err error) {
				require.Error(t, err)
				require.EqualError(t, err, "Failed expectation: Scope access to have refresh")
				require.Nil(t, claim)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			private, public := common.GenerateRSAKey(t)

			cfg := config.Token{
				AccessTokenDuration: -time.Minute,
				PrivateKey:          private,
				PublicKey:           public,
			}

			jwtImpl, err := NewJWT(cfg)
			require.NoError(t, err)

			claim, err := jwtImpl.VerifyToken(tc.token(cfg.PrivateKey), AccessTokenExpectation())
			tc.checkResult(claim, err)
		})
	}
}
