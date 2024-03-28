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

	projectG.POST("", handler.createProject)
	projectG.DELETE("", handler.deleteProject)
	projectG.GET("/list", handler.projectList)

	projectG.GET("/:uid/list", handler.getProjects) // 用户相关列表
	projectG.GET("/detail", handler.projectDetail)  // 详情
	projectG.GET("/changeStatus", handler.changeStatus)

	projectG.GET("/choose", handler.chooseProject) // 学生选择
	projectG.GET("/audit", handler.auditProject)   // 审核
	projectG.POST("/upload/file")                  // 提交文件

	// projectG.GET("/")
}
