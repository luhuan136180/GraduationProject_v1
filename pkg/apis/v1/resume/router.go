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

	resumeG.POST("", handler.createResume)                 // done
	resumeG.POST("/project/tree", handler.projectTreeList) // done
	resumeG.DELETE("", handler.deleteResume)               // done
	resumeG.POST("/list", handler.resumeList)              // done
	resumeG.POST("/:id/detail", handler.resumeDetail)      // 详情

	resumeG.GET("/:uid/resume/list", handler.resumeListByUid) // 根据传入的uid获取对应用户的简历列表

	// resumeG.POST("/send", handler.SendResume)                            // 发送简历
	// resumeG.GET("/:recruitID/resume/list", handler.getResumeByRecruitID) // 获取招聘接收的简历列表
}
