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

	projectG.POST("", handler.createProject)   // done
	projectG.DELETE("", handler.deleteProject) // done
	projectG.GET("/list", handler.projectList) // done

	projectG.GET("/user/list", handler.getProjects)     // 用户相关列表(我的)
	projectG.GET("/detail", handler.projectDetail)      // 详情 done
	projectG.GET("/changeStatus", handler.changeStatus) // 更改状态 done

	projectG.GET("/choose", handler.chooseProject) // 学生选择 done
	projectG.GET("/audit", handler.auditProject)   // 审核 废弃
	projectG.POST("/upload/file")                  // 提交文件

	// projectG.GET("/")
}
