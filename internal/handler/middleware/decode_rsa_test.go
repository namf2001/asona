package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"asona/config"
	"asona/internal/constants"
	"asona/internal/pkg/rsa"
)

func setupConfig(t *testing.T, appEnv string) {
	t.Helper()

	t.Setenv("APP_ENV", appEnv)
	t.Setenv("DB_HOST", "localhost")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_NAME", "asona")
	t.Setenv("DB_SSL_MODE", "disable")
	t.Setenv("MAIL_SMTP_USER", "mail@example.com")
	t.Setenv("MAIL_SMTP_PASSWORD", "password")
	t.Setenv("MAIL_EMAIL_FROM", "mail@example.com")
	t.Setenv("AWS_S3_ACCESS_KEY", "access")
	t.Setenv("AWS_S3_SECRET_KEY", "secret")
	t.Setenv("AWS_S3_BUCKET_NAME", "bucket")

	config.Init(appEnv)
}

func makeTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(RSAAuthMiddleware())
	r.POST("/", func(c *gin.Context) {
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		c.String(http.StatusOK, string(payload))
	})
	return r
}

func TestRSAAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("dev environment skips decrypt", func(t *testing.T) {
		setupConfig(t, "dev")
		originalKeyPair := rsa.GlobalRSAKeyPair
		rsa.GlobalRSAKeyPair = nil
		t.Cleanup(func() { rsa.GlobalRSAKeyPair = originalKeyPair })

		router := makeTestRouter()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"name":"test"}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"name":"test"}`, rec.Body.String())
	})

	t.Run("non-dev allows empty body", func(t *testing.T) {
		setupConfig(t, "staging")
		originalKeyPair := rsa.GlobalRSAKeyPair
		rsa.GlobalRSAKeyPair = nil
		t.Cleanup(func() { rsa.GlobalRSAKeyPair = originalKeyPair })

		router := makeTestRouter()
		req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "", rec.Body.String())
	})

	t.Run("non-dev rejects invalid json body", func(t *testing.T) {
		setupConfig(t, "staging")
		originalKeyPair := rsa.GlobalRSAKeyPair
		rsa.GlobalRSAKeyPair = nil
		t.Cleanup(func() { rsa.GlobalRSAKeyPair = originalKeyPair })

		router := makeTestRouter()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("not-json"))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), constants.InvalidRequestParams.Code)
	})

	t.Run("non-dev rejects json body without encrypted data", func(t *testing.T) {
		setupConfig(t, "staging")
		originalKeyPair := rsa.GlobalRSAKeyPair
		rsa.GlobalRSAKeyPair = nil
		t.Cleanup(func() { rsa.GlobalRSAKeyPair = originalKeyPair })

		router := makeTestRouter()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), constants.InvalidRequestParams.Code)
	})

	t.Run("non-dev returns 500 when rsa keypair missing", func(t *testing.T) {
		setupConfig(t, "staging")
		originalKeyPair := rsa.GlobalRSAKeyPair
		rsa.GlobalRSAKeyPair = nil
		t.Cleanup(func() { rsa.GlobalRSAKeyPair = originalKeyPair })

		router := makeTestRouter()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"data":"AQID"}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusInternalServerError, rec.Code)
		require.Contains(t, rec.Body.String(), constants.RSANotInitialized.Code)
	})

	t.Run("non-dev returns bad request when decrypt fails", func(t *testing.T) {
		setupConfig(t, "staging")
		keyPair := rsa.NewKeyPair()
		require.NoError(t, keyPair.GenerateRSAKeyPair(2048))

		originalKeyPair := rsa.GlobalRSAKeyPair
		rsa.GlobalRSAKeyPair = keyPair
		t.Cleanup(func() { rsa.GlobalRSAKeyPair = originalKeyPair })

		router := makeTestRouter()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"data":"AQID"}`))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), constants.DecryptRSAFail.Code)
	})

	t.Run("non-dev decrypts and forwards body", func(t *testing.T) {
		setupConfig(t, "staging")
		keyPair := rsa.NewKeyPair()
		require.NoError(t, keyPair.GenerateRSAKeyPair(2048))

		originalKeyPair := rsa.GlobalRSAKeyPair
		rsa.GlobalRSAKeyPair = keyPair
		t.Cleanup(func() { rsa.GlobalRSAKeyPair = originalKeyPair })

		plainBody := []byte(`{"email":"user@example.com"}`)
		encryptedBody, err := keyPair.Encrypt(plainBody)
		require.NoError(t, err)

		payload, err := json.Marshal(BodyEncrypt{Data: encryptedBody})
		require.NoError(t, err)

		router := makeTestRouter()
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, string(plainBody), rec.Body.String())
	})
}
