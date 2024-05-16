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
	// 用户个体的API
	systemG.POST("/user/file", handler.uploaduserTouxiang)  // 上传头像
	systemG.POST("/users/changeinfo", handler.editUserInfo) // 编辑用户信息 done

	// yonghuxitong
	// systemG.GET("/users/list", handler.userListTest)
	systemG.DELETE("/users", handler.deleteUser) // done
	systemG.POST("/users", handler.createUser)
	systemG.POST("/users/list", handler.getUserList)         // 用户列表 done
	systemG.POST("/user/:id/detail", handler.getUserDetail)  // 用户详情 done
	systemG.PUT("/users/password", handler.changeUserPwd)    //
	systemG.PUT("/users/:id/password", handler.resetUserPWD) // 管理员重置密码 done

	// college
	systemG.POST("/colleges", handler.createCollege) //
	systemG.DELETE("/colleges", handler.deleteCollege)
	systemG.POST("/colleges/tree", handler.getCollegeTree)   // tree done
	systemG.GET("/college/list", handler.collegeList)        // 获取学院列表
	systemG.GET("/college/detail", handler.getCollegeDetail) // 获取学院详情信息
	// profession
	systemG.POST("/professions", handler.createProfession)
	systemG.DELETE("/professions", handler.deleteProfession)
	systemG.POST("/profession/tree", handler.getProfessionTree) // tree done
	systemG.GET("/profession/list", handler.professionList)
	systemG.GET("/profession/detail", handler.professionDetail)

	// class
	systemG.POST("/classes", handler.createClass)
	systemG.DELETE("/classes", handler.deleteClass)
	systemG.POST("/:profession_hash_id/class/tree", handler.getClassTree) // tree done
	systemG.GET("/class/list", handler.classList)                         // 班级列表

	// 企业管理
	systemG.GET("/firm/list", handler.firmList)
	systemG.GET("/firm", handler.createFirm)
	systemG.DELETE("/firm", handler.deleteFirm)
	systemG.POST("/firm/user", handler.createFirmUser)
	systemG.POST("/firm/user/list", handler.FirmUserList)

	systemG.GET("/firm/detail", handler.firmDetail)
	systemG.GET("/firm/tree", handler.firmTree)
}
