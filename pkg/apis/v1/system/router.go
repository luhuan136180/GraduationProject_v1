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

	systemG.DELETE("/users", handler.deleteUser) // done
	systemG.POST("/users", handler.createUser)
	systemG.GET("/users/list", handler.getUserList)          // 用户列表 done
	systemG.GET("/user/:id/detail", handler.getUserDetail)   // 用户详情 done
	systemG.PATCH("/users", handler.editUserInfo)            // 编辑用户信息 done
	systemG.PUT("/users/password", handler.changeUserPwd)    // 废弃
	systemG.PUT("/users/:id/password", handler.resetUserPWD) // 管理员重置密码 done

	// college
	systemG.POST("/colleges", handler.createCollege) //
	systemG.DELETE("/colleges", handler.deleteCollege)
	systemG.GET("/colleges/tree", handler.getCollegeTree) // tree done

	// profession
	systemG.POST("/professions", handler.createProfession)
	systemG.DELETE("/professions", handler.deleteProfession)
	systemG.GET("/profession/tree", handler.getProfessionTree) // tree done

	// class
	systemG.POST("/classes", handler.createClass)
	systemG.DELETE("/classes", handler.deleteClass)
	systemG.GET("/:profession_hash_id/class/tree", handler.getClassTree) // tree done

}
