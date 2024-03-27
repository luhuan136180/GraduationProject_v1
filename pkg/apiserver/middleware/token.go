//go:build !debug

package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
	"time"
	"v1/pkg/apiserver/encoding"
	"v1/pkg/apiserver/request"
	"v1/pkg/client/cache"
	"v1/pkg/server/errutil"
	"v1/pkg/token"
)

// tokenFromHeader tries to retrieve the token string from the
// "Authorization" request header: "Authorization: BEARER T".
func tokenFromHeader(c *gin.Context) string {
	// Get token from authorization header.
	bearer := c.GetHeader("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}

	return bearer
}

// TokenFromCookie tries to retrieve the token string from a cookie named
// "jwt".
func tokenFromCookie(c *gin.Context) string {
	cookieValue, err := c.Cookie("jwt")
	if err != nil {
		return ""
	}
	return cookieValue
}

// tokenFromQuery tries to retrieve the token string from the "jwt" URI
// query parameter.
func tokenFromQuery(c *gin.Context) string {
	// Get token from query param named "jwt".
	return c.Query("jwt")
}

func findTokenVal(c *gin.Context, getTokenFns ...func(c *gin.Context) string) string {
	var tokenStr string

	// Extract token string from the request by calling token find functions in
	// the order they were provided. Further extraction stops if a function
	// returns a non-empty string.
	for _, fn := range getTokenFns {
		tokenStr = fn(c)
		if tokenStr != "" {
			break
		}
	}

	return tokenStr
}

func CheckToken(manager token.Manager, cacheClient cache.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenVal := findTokenVal(c, tokenFromHeader, tokenFromCookie, tokenFromQuery)
		if tokenVal == "" {
			encoding.HandleError(c, errutil.ErrUnauthorized)
			return
		}

		if _, err := cacheClient.Get(context.Background(), "token:"+tokenVal); err != nil {
			zap.L().Error("", zap.Error(err))
			encoding.HandleError(c, errutil.ErrUnauthorized)
			return
		}

		payload, err := manager.Verify(tokenVal)
		if err != nil {
			zap.L().Info("解析token错误", zap.String("tokenVal", tokenVal), zap.Error(err))
			encoding.HandleError(c, errutil.ErrUnauthorized)
			return
		}

		zap.L().Debug("token payload", zap.Any("payload", payload))

		// renew token
		if err = cacheClient.Set(context.Background(), "token:"+tokenVal, payload.Username, time.Minute*30); err != nil {
			zap.L().Error("", zap.Error(err))
			encoding.HandleError(c, errutil.ErrInternalServer)
			return
		}

		ctx := request.WithTokenPayloadToCtx(c.Request.Context(), &payload)
		c.Request = c.Request.WithContext(ctx)
		return
	}
}
