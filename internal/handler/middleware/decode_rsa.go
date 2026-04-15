package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"asona/config"
	"asona/internal/constants"
	"asona/internal/handler/response"
	"asona/internal/pkg/logger"
	"asona/internal/pkg/rsa"
)

// BodyEncrypt is the encrypted request body structure.
type BodyEncrypt struct {
	Data []byte `json:"data"`
}

// RSAAuthMiddleware decodes the request body using the RSA private key.
func RSAAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cfg := config.GetConfig()

		// Skip decryption in dev environment
		if cfg.AppEnv == "dev" {
			c.Next()
			return
		}

		rawBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.ERROR.Printf("[RSAAuthMiddleware] failed to read request body: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, response.NewResponse(
				constants.InvalidRequestParams.Code,
				constants.InvalidRequestParams.Message,
				nil,
			))
			return
		}

		if strings.TrimSpace(string(rawBody)) == "" {
			c.Request.Body = io.NopCloser(strings.NewReader(""))
			c.Next()
			return
		}

		var body BodyEncrypt
		if err := json.Unmarshal(rawBody, &body); err != nil {
			logger.ERROR.Printf("[RSAAuthMiddleware] invalid request body: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, response.NewResponse(
				constants.InvalidRequestParams.Code,
				constants.InvalidRequestParams.Message,
				nil,
			))
			return
		}

		if len(body.Data) == 0 {
			logger.ERROR.Printf("[RSAAuthMiddleware] encrypted body data is empty")
			c.AbortWithStatusJSON(http.StatusBadRequest, response.NewResponse(
				constants.InvalidRequestParams.Code,
				constants.InvalidRequestParams.Message,
				nil,
			))
			return
		}

		if rsa.GlobalRSAKeyPair == nil {
			logger.ERROR.Printf("[RSAAuthMiddleware] RSA key pair not initialized")
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.NewResponse(
				constants.RSANotInitialized.Code,
				constants.RSANotInitialized.Message,
				nil,
			))
			return
		}

		decodedBody, err := rsa.GlobalRSAKeyPair.Decrypt(body.Data)
		if err != nil {
			logger.ERROR.Printf("[RSAAuthMiddleware] failed to decrypt RSA body: %+v", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, response.NewResponse(
				constants.DecryptRSAFail.Code,
				constants.DecryptRSAFail.Message,
				nil,
			))
			return
		}

		c.Request.Body = io.NopCloser(strings.NewReader(string(decodedBody)))
		c.Next()
	}
}
