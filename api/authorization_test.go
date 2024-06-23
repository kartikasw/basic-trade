package api

import (
	apiHelper "basic-trade/api/helper"
	"basic-trade/common"
	mockRepo "basic-trade/internal/repository/mock"
	"basic-trade/pkg/config"
	"basic-trade/pkg/token"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestProductAuthorization(t *testing.T) {
	usrUUID := uuid.New()
	prdUUID := uuid.New()

	testCases := []struct {
		name          string
		buildStubs    func(repo *mockRepo.MockAdminRepository)
		checkResponse func(t *testing.T, code int, data map[string]any)
	}{
		{
			name: "OK",
			buildStubs: func(repo *mockRepo.MockAdminRepository) {
				repo.EXPECT().CheckProductFromAdmin(gomock.Any(), usrUUID, prdUUID).Times(1).Return(true)
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusOK, code)
			},
		},
		{
			name: "UnauthorizedUser",
			buildStubs: func(repo *mockRepo.MockAdminRepository) {
				repo.EXPECT().CheckProductFromAdmin(gomock.Any(), usrUUID, prdUUID).Times(1).Return(false)
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusUnauthorized, code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server, repo := setUpTest(t, true)

			url := "/test/" + prdUUID.String()
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPut, url, nil)
			require.NoError(t, err)

			mockAuthorization(t, request, server.jwtImpl, usrUUID, token.BearerScheme)
			tc.buildStubs(repo)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder.Code, nil)
		})
	}
}

func TestVariantAuthorization(t *testing.T) {
	usrUUID := uuid.New()
	prdUUID := uuid.New()

	testCases := []struct {
		name          string
		buildStubs    func(repo *mockRepo.MockAdminRepository)
		checkResponse func(t *testing.T, code int, data map[string]any)
	}{
		{
			name: "OK",
			buildStubs: func(repo *mockRepo.MockAdminRepository) {
				repo.EXPECT().CheckVariantFromAdmin(gomock.Any(), usrUUID, prdUUID).Times(1).Return(true)
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusOK, code)
			},
		},
		{
			name: "UnauthorizedUser",
			buildStubs: func(repo *mockRepo.MockAdminRepository) {
				repo.EXPECT().CheckVariantFromAdmin(gomock.Any(), usrUUID, prdUUID).Times(1).Return(false)
			},
			checkResponse: func(t *testing.T, code int, data map[string]any) {
				require.Equal(t, http.StatusUnauthorized, code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server, repo := setUpTest(t, false)

			url := "/test/" + prdUUID.String()
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodPut, url, nil)
			require.NoError(t, err)

			mockAuthorization(t, request, server.jwtImpl, usrUUID, token.BearerScheme)
			tc.buildStubs(repo)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder.Code, nil)
		})
	}
}

func setUpTest(t *testing.T, product bool) (*Server, *mockRepo.MockAdminRepository) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	privateKey, publicKey := common.GenerateRSAKey(t)
	cfg := config.Config{
		Token: config.Token{
			AccessTokenDuration: 2 * time.Second,
			PrivateKey:          privateKey,
			PublicKey:           publicKey,
		},
		App: config.App{
			Timeout: 2 * time.Second,
		},
	}

	jwtImpl, err := token.NewJWT(cfg.Token)
	require.NoError(t, err)

	adminRepo := mockRepo.NewMockAdminRepository(ctrl)
	server := &Server{
		jwtImpl:       jwtImpl,
		authorization: *NewAuthorizationMiddleware(adminRepo),
	}

	router := gin.Default()

	authorization := server.authorization.ProductAuthorization()
	if !product {
		authorization = server.authorization.VariantAuthorization()
	}

	router.PUT(
		"/test/:uuid",
		Timeout(cfg.App.Timeout),
		Authentication(server.jwtImpl),
		authorization,
		func(ctx *gin.Context) {
			apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
				resChan <- apiHelper.ResponseData{
					StatusCode: http.StatusOK,
				}
			})
		},
	)

	server.router = router

	return server, adminRepo
}
