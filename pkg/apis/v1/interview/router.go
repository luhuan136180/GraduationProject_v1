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
	resumeG.POST("/list/accept", handler.interviewListAccept)

	// 改变状态
	resumeG.POST("/change", handler.interviewChangeStatus) // done
	// 详情
	resumeG.POST("/:id/detail", handler.interviewDetail) // done

	// 我的招聘
	resumeG.GET("/myrecruit/list", handler.getMyRecruitList)
	// 新增招聘信息
	resumeG.POST("/recruit", handler.createRecruit)

	// 删除招聘信息
	resumeG.DELETE("/recruit", handler.deleteRecruit)

	resumeG.GET("/recruit/list", handler.getRecruitList)
	resumeG.GET("/firm/detail", handler.firmDetail) // 获取企业信息

	resumeG.GET("/recruit/detial/:id", handler.getRecruitDeatial) // 详情

}
