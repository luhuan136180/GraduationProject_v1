package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
	"v1/pkg/apiserver/request"
	"v1/pkg/dao"
)

var filter map[string]struct{}

func init() {
	filter = make(map[string]struct{})
	filter["/api/v1/system/logs"] = struct{}{}
	filter["/api/v1/common/fs/"] = struct{}{}
	filter["/api/v1/common/healthy"] = struct{}{}
}

func AddAuditLog(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		url := c.Request.URL
		method := c.Request.Method
		startTime := time.Now().UnixMilli()
		ip := c.ClientIP()
		blueprint := getBluePrint(url.Path)

		c.Next()
		if blueprint == "auth" || ApiFilter(url.Path) {
			return
		}

		status := c.Writer.Status()
		username := request.GetUsernameFromCtx(c)
		name := request.GetNameFromCtx(c)
		endTime := time.Now().UnixMilli()
		duration := endTime - startTime

		err := dao.InsertLog(c, db, username, name, blueprint, method, ip, url.Path, status, duration, startTime)
		if err != nil {
			zap.L().Info("save audit log failed", zap.Error(err))
		}

		return
	}
}

func getBluePrint(path string) string {
	slice := strings.Split(path, "/")
	return slice[3]
}

func ApiFilter(url string) bool {
	for path, _ := range filter {
		if strings.Contains(url, path) {
			return true
		}
	}
	return false
}
