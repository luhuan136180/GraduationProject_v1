package contract

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"v1/pkg/apiserver/middleware"
	"v1/pkg/client/cache"
	"v1/pkg/contract"
	"v1/pkg/token"
)

func RegisterRouter(group *gin.RouterGroup, tokenManager token.Manager, cacheClient cache.Interface, db *gorm.DB, bolckClient *contract.Client) {
	resumeG := group.Group("/contract")
	handler := newInterviewHandler(contractHandlerOption{
		bolckClient: bolckClient,
		db:          db,
	})

	resumeG.Use(middleware.CheckToken(tokenManager, cacheClient))

	resumeG.GET("/block/list", handler.blockList)   //
	resumeG.GET("/block/value", handler.getContent) //
	resumeG.GET("/block/:block_hash", handler.blockDetail)
}
