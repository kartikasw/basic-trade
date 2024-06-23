package api

import (
	"basic-trade/common"
	"basic-trade/pkg/config"
	"basic-trade/pkg/token"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestAuthentication(t *testing.T) {
	testCases := []struct {
		name           string
		overrideConfig bool
		setupAuth      func(t *testing.T, request *http.Request, jwtImpl token.JWT)
		checkResponse  func(t *testing.T, code int, data map[string]any)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, jwtImpl token.JWT) {
				mockAuthorization(t, request, jwtImpl, uuid.New(), token.BearerScheme)
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusOK, code)
			},
		},
		{
			name:      "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, jwtImpl token.JWT) {},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusUnauthorized, code)
				require.Equal(t, data["message"], "Authorization header is not found")
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, jwtImpl token.JWT) {
				mockAuthorization(t, request, jwtImpl, uuid.New(), "")
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusUnauthorized, code)
				require.Equal(t, data["message"], "Authorization header is not in Bearer scheme")
			},
		},
		{
			name: "InvalidAuthorizationScheme",
			setupAuth: func(t *testing.T, request *http.Request, jwtImpl token.JWT) {
				mockAuthorization(t, request, jwtImpl, uuid.New(), "unsupported")
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusUnauthorized, code)
				require.Equal(t, data["message"], "Unsupported authorization scheme unsupported")
			},
		},
		{
			name:           "InvalidToken",
			overrideConfig: true,
			setupAuth: func(t *testing.T, request *http.Request, jwtImpl token.JWT) {
				mockAuthorization(t, request, jwtImpl, uuid.New(), token.BearerScheme)
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusUnauthorized, code)
				require.Equal(t, data["message"], "Couldn't verify token: token is expired")
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newAuthenticationTestServer(t, tc.overrideConfig)

			url := "/test"
			server.router.GET(
				url,
				Authentication(server.jwtImpl),
				func(ctx *gin.Context) {
					ctx.JSON(http.StatusOK, gin.H{})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.jwtImpl)

			server.router.ServeHTTP(recorder, request)

			var data map[string]interface{}
			json.Unmarshal(recorder.Body.Bytes(), &data)
			tc.checkResponse(t, recorder.Code, data)
		})
	}
}

func mockAuthorization(
	t *testing.T,
	request *http.Request,
	jwtImpl token.JWT,
	usrUUID uuid.UUID,
	authorizationType string,
) {
	jwtToken, err := jwtImpl.CreateAccessToken(usrUUID)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, jwtToken.SignedToken)
	request.Header.Set(token.AuthorizationHeader, authorizationHeader)
}

func newAuthenticationTestServer(t *testing.T, overrideConfig bool) *Server {
	privateKey, publicKey := common.GenerateRSAKey(t)
	cfg := config.Token{
		AccessTokenDuration: 2 * time.Second,
		PrivateKey:          privateKey,
		PublicKey:           publicKey,
	}

	if overrideConfig {
		cfg.AccessTokenDuration = -cfg.AccessTokenDuration
	}

	jwtImpl, err := token.NewJWT(cfg)
	require.NoError(t, err)

	server := &Server{
		jwtImpl: jwtImpl,
	}

	server.router = gin.Default()
	return server
}
