package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"v1/pkg/apiserver/middleware"
	"v1/pkg/client/cache"
	"v1/pkg/token"
)

func RegisterRouter(group *gin.RouterGroup, tokenManager token.Manager, cacheClient cache.Interface, db *gorm.DB) {
	authG := group.Group("/auth")
	handler := newAuthHandler(authHandlerOption{
		tokenManager: tokenManager,
		db:           db,
		cacheClient:  cacheClient,
	})

	authG.POST("/login", handler.login)
	authG.GET("/captcha", handler.createCaptcha)
	// authG.GET("/license", handler.licenseInfo)      // 获取license信息
	// authG.POST("/license", handler.registerLicense) // 激活license

	authG.Use(middleware.CheckToken(tokenManager, cacheClient))
	authG.POST("/logout", handler.logout)
}
