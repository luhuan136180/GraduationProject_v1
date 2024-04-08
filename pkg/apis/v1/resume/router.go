package resume

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"v1/pkg/apiserver/middleware"
	"v1/pkg/client/cache"
	"v1/pkg/token"
)

func RegisterRouter(group *gin.RouterGroup, tokenManager token.Manager, cacheClient cache.Interface, db *gorm.DB) {
	resumeG := group.Group("/resume")
	handler := newResumeHandler(resumeHandlerOption{
		db: db,
	})

	resumeG.Use(middleware.CheckToken(tokenManager, cacheClient))

	resumeG.POST("", handler.createResume)                // done
	resumeG.GET("/project/tree", handler.projectTreeList) // done
	resumeG.DELETE("", handler.deleteResume)              // done
	resumeG.GET("/list", handler.resumeList)              // done
	resumeG.GET("/:id/detail", handler.resumeDetail)      // 详情

}
