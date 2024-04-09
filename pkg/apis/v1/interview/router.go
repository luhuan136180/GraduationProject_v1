package interview

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"v1/pkg/apiserver/middleware"
	"v1/pkg/client/cache"
	"v1/pkg/token"
)

func RegisterRouter(group *gin.RouterGroup, tokenManager token.Manager, cacheClient cache.Interface, db *gorm.DB) {
	resumeG := group.Group("/interview")
	handler := newInterviewHandler(interviewHandlerOption{
		db: db,
	})

	resumeG.Use(middleware.CheckToken(tokenManager, cacheClient))

	resumeG.POST("", handler.createInterview)    // done
	resumeG.DELETE("", handler.deleteInterview)  // done
	resumeG.POST("/list", handler.interviewList) // done

	// 改变状态
	resumeG.POST("/change", handler.interviewChangeStatus) // done
	// 详情
	resumeG.POST("/:id/detail", handler.interviewDetail) // done
}
