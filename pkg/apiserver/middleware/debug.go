//go:build debug

package middleware

import (
	"github.com/gin-gonic/gin"
	"v1/pkg/apiserver/request"
	"v1/pkg/client/cache"
	"v1/pkg/model"
	"v1/pkg/token"
)

func CheckToken(manager token.Manager, cacheClient cache.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := request.WithTokenPayloadToCtx(c.Request.Context(), &token.Payload{
			ID:       1,
			Username: "admin",
			Name:     "admin",
			Role:     model.RoleTypeSuperAdmin,
		})
		c.Request = c.Request.WithContext(ctx)
	}
}
