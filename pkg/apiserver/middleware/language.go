package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"v1/pkg/apiserver/request"
)

func WithLanguage() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 尝试从url获取lang参数
		lang := c.Query("lang")

		// 尝试从cookie获取lang
		if lang == "" {
			var err error
			lang, err = c.Cookie("lang")
			if !errors.Is(err, http.ErrNoCookie) {
				zap.L().Warn("WithLanguage get cookieLang error", zap.Error(err))
			}
		}

		// 从请求头获取accept-language
		if lang == "" {
			lang = c.GetHeader("Accept-Language")
		}

		ctx := request.WithLanguageToCtx(c.Request.Context(), lang)
		c.Request = c.Request.WithContext(ctx)
	}
}
