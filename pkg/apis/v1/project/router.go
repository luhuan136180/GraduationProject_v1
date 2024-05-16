package project

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"v1/pkg/apiserver/middleware"
	"v1/pkg/client/cache"
	"v1/pkg/token"
)

func RegisterRouter(group *gin.RouterGroup, tokenManager token.Manager, cacheClient cache.Interface, db *gorm.DB) {
	projectG := group.Group("/project")
	handler := newProjectHandler(projectHandlerOption{
		db: db,
	})

	projectG.Use(middleware.CheckToken(tokenManager, cacheClient))

	projectG.POST("", handler.createProject)    // done --管理员创建项目
	projectG.DELETE("", handler.deleteProject)  // done
	projectG.POST("/list", handler.projectList) // done

	projectG.POST("/user/list", handler.getProjects)     // 用户相关列表(我的)
	projectG.POST("/detail", handler.projectDetail)      // 详情 done
	projectG.POST("/changeStatus", handler.changeStatus) // 更改状态 done
	projectG.POST("/file/list", handler.FileList)

	// projectG.POST("/choose", handler.chooseProject)       // 学生选择 done
	// projectG.GET("/audit", handler.auditProject)          // 审核 废弃
	projectG.POST("/:id/upload/file", handler.uploadFile) // 提交文件
	projectG.POST("/upload/file", handler.uploadonlyFile) // 单纯上交文件

	// projectG.GET("/")
	projectG.GET("/:uid/project/list", handler.getProjectsByUid) // 通过uid获取项目列表
}
