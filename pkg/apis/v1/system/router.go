package system

import (
	"gorm.io/gorm"
	"v1/pkg/apiserver/middleware"
	"v1/pkg/client/cache"
	"v1/pkg/token"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(group *gin.RouterGroup, tokenManager token.Manager, cacheClient cache.Interface, db *gorm.DB) {
	systemG := group.Group("/system")
	handler := newSystemHandler(systemHandlerOption{
		db: db,
	})

	systemG.Use(middleware.CheckToken(tokenManager, cacheClient))

	systemG.DELETE("/users", handler.deleteUser)
	systemG.POST("/users", handler.createUser)
	systemG.GET("/users", handler.getUserList)
	systemG.PATCH("/users", handler.editUserInfo)
	systemG.PUT("/users/password", handler.changeUserPwd)
	systemG.PUT("/users/:uid/password", handler.resetUserPWD)

	// college
	systemG.POST("/colleges", handler.createCollege)
	systemG.DELETE("/colleges", handler.deleteCollege)

	// profession
	systemG.POST("/professions", handler.createProfession)
	systemG.DELETE("/professions", handler.deleteProfession)

	// class
	systemG.POST("/classes", handler.createClass)
	systemG.DELETE("/classes", handler.deleteClass)

}
