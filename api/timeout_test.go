package api

import (
	apiHelper "basic-trade/api/helper"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestTimeout(t *testing.T) {
	testCases := []struct {
		name          string
		setUpResponse func(resChan chan apiHelper.ResponseData, timeout time.Duration)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Success",
			setUpResponse: func(resChan chan apiHelper.ResponseData, timeout time.Duration) {
				resChan <- apiHelper.ResponseData{
					StatusCode: http.StatusOK,
				}
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "TimedOut",
			setUpResponse: func(resChan chan apiHelper.ResponseData, timeout time.Duration) {
				time.Sleep(timeout + time.Second)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusGatewayTimeout, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			// timeout duration must be less than 30s to avoid timeout from goroutine
			server := newTimeoutTestServer()
			timeout := 2 * time.Second

			url := "/test"
			server.router.GET(
				url,
				Timeout(timeout),
				func(ctx *gin.Context) {
					apiHelper.ResponseHandler(ctx, func(c context.Context, resChan chan apiHelper.ResponseData) {
						tc.setUpResponse(resChan, timeout)
					})
				},
			)

			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func newTimeoutTestServer() *Server {
	server := &Server{}
	server.router = gin.Default()

	return server
}
