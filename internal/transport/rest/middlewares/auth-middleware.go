package middlewares

import (
	"paywise/internal/business/auth/token"
	"paywise/internal/core"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AUTHORIZATION_HEADER = "authorization"
	AUTHORIZATION_TYPE   = "Bearer"
	// to set the payload in the context of the request before passing it to the next handler
	AUTHORIZATION_PAYLOAD_CTX_KEY = "authorization_payload"
)

func Authenticate(tokenProvider token.TokenMaker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// if the authorization header is empty, abort the requst
		if len(ctx.GetHeader(AUTHORIZATION_HEADER)) == 0 {
			appErr := core.NewUnAuthorizedError("authorization header is not proivded, not authorized to perform this request")
			ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{
				"error": appErr,
			})
		}

		// if the authorization header is not empty but the token is not in a valid format, abort the request
		authorizationHeader := strings.Fields(ctx.GetHeader(AUTHORIZATION_HEADER))
		if len(authorizationHeader) < 2 {
			appErr := core.NewUnAuthorizedError("authorization header is in invalid format, not authorized to perform this request")
			ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{
				"error": appErr,
			})
		}

		authorizationTokenType := authorizationHeader[0]
		if authorizationTokenType != AUTHORIZATION_TYPE {
			appErr := core.NewUnAuthorizedError("authorization token type is not supported by the server we only support Bearer tokens, not authorized to perform this request")
			ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{
				"error": appErr,
			})
		}

		// if its in a valid format, verify it using the verify method of the token provider to decrypt the token and return its payload (if the expiration date is still valid)
		tokenPayload, err := tokenProvider.Verify(authorizationHeader[1])
		if err != nil {
			appErr := core.NewUnAuthorizedError("couldn't verify the token, not authorized to perform this request")
			ctx.AbortWithStatusJSON(appErr.StatusCode(), gin.H{
				"error": appErr,
			})
		}

		// set the payload into the context
		ctx.Set(AUTHORIZATION_PAYLOAD_CTX_KEY, tokenPayload)

		// call the next handler
		ctx.Next()
	}
}
