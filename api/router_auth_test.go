package api

import (
	"basic-trade/common"
	"basic-trade/internal/entity"
	"basic-trade/internal/handler"
	mockService "basic-trade/internal/service/mock"
	"basic-trade/pkg/config"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	admin := randomAdmin(true)

	testCases := []struct {
		name          string
		body          []testFormData
		buildStubs    func(service *mockService.MockAuthService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: []testFormData{
				{Fieldname: "name", Content: []byte(admin.Name)},
				{Fieldname: "email", Content: []byte(admin.Email)},
				{Fieldname: "password", Content: []byte(admin.Password)},
			},
			buildStubs: func(service *mockService.MockAuthService) {
				service.EXPECT().Register(gomock.Any(), admin).Times(1).Return(admin, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "EmailAlreadyRegistered",
			body: []testFormData{
				{Fieldname: "name", Content: []byte(admin.Name)},
				{Fieldname: "email", Content: []byte(admin.Email)},
				{Fieldname: "password", Content: []byte(admin.Password)},
			},
			buildStubs: func(service *mockService.MockAuthService) {
				service.
					EXPECT().
					Register(gomock.Any(), admin).Times(1).
					Return(entity.Admin{}, common.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server, service := setUpAuthTest(t)

			recorder := httptest.NewRecorder()
			url := "/auth/register"
			request := mockRequest(t, tc.body, url)

			tc.buildStubs(service)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestLogin(t *testing.T) {
	admin := randomAdmin(false)

	testCases := []struct {
		name          string
		body          []testFormData
		buildStubs    func(service *mockService.MockAuthService)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: []testFormData{
				{Fieldname: "email", Content: []byte(admin.Email)},
				{Fieldname: "password", Content: []byte(admin.Password)},
			},
			buildStubs: func(service *mockService.MockAuthService) {
				service.EXPECT().Login(gomock.Any(), admin).Times(1).Return(admin, "", nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "EmailNotRegistered",
			body: []testFormData{
				{Fieldname: "email", Content: []byte(admin.Email)},
				{Fieldname: "password", Content: []byte(admin.Password)},
			},
			buildStubs: func(service *mockService.MockAuthService) {
				service.
					EXPECT().
					Login(gomock.Any(), admin).Times(1).
					Return(entity.Admin{}, "", common.ErrRecordNotFound)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: []testFormData{
				{Fieldname: "email", Content: []byte(admin.Email)},
				{Fieldname: "password", Content: []byte(admin.Password)},
			},
			buildStubs: func(service *mockService.MockAuthService) {
				service.
					EXPECT().
					Login(gomock.Any(), admin).Times(1).
					Return(entity.Admin{}, "", errors.New("error"))
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server, service := setUpAuthTest(t)

			recorder := httptest.NewRecorder()
			url := "/auth/login"
			request := mockRequest(t, tc.body, url)

			tc.buildStubs(service)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func randomAdmin(name bool) entity.Admin {
	admin := entity.Admin{
		Email:    common.RandomEmail(),
		Password: common.RandomString(10),
	}

	if name {
		admin.Name = common.RandomName()

	}

	return admin
}

func setUpAuthTest(t *testing.T) (*Server, *mockService.MockAuthService) {
	cfg := config.Config{
		App: config.App{
			Timeout: 5 * time.Second,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mockService.NewMockAuthService(ctrl)
	authHandler := handler.NewAuthHandler(authService)

	server := &Server{
		authHandler: authHandler,
	}

	server.setupRouter(cfg.App)

	return server, authService
}
