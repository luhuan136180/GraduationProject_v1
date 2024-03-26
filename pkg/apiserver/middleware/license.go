package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"v1/pkg/apiserver/encoding"
	"v1/pkg/dao"
	"v1/pkg/license"
	"v1/pkg/model"
	"v1/pkg/server/errutil"
)

func VerifyLicense(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		conf, err := dao.GetConfigByKey(context.Background(), db, model.ConfigKeyLicense)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				zap.L().Info("license info not found")
				encoding.HandleSuccess(c, errutil.ErrInvalidLicense)
				return
			}

			zap.L().Error("get license info failed", zap.Error(err))
			encoding.HandleError(c, err)
			return
		}

		cl := model.ConfigLicense{}
		if err = json.Unmarshal(conf.Value, &cl); err != nil {
			zap.L().Error("unmarshal license info failed", zap.Error(err))
			encoding.HandleError(c, err)
			return
		}

		l, err := license.GetAuthorityLicense(cl.License)
		if err != nil {
			zap.L().Error("authority license failed", zap.Error(err))
			encoding.HandleError(c, err)
			return
		}

		if !l.Valid {
			encoding.HandleError(c, errutil.ErrInvalidLicense)
			return
		}
	}
}
