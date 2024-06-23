package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type testFormData struct {
	Fieldname string
	Content   []byte
}

func mockRequest(t *testing.T, data []testFormData, url string) *http.Request {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)

	for _, d := range data {
		fw, err := mw.CreateFormFile(d.Fieldname, "")
		require.NoError(t, err)

		n, err := fw.Write(d.Content)
		require.NoError(t, err)
		require.Equal(t, len(d.Content), n)
	}
	err := mw.Close()
	require.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, &body)
	require.NoError(t, err)

	request.Header.Set("Content-Type", mw.FormDataContentType())

	return request
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
