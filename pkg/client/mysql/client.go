package mysql

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewMysqlClient create a gorm mysql client
func NewMysqlClient(options *Options) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		options.RdbUser, options.RdbPassword, options.RdbHost, options.RdbPort, options.RdbDbname)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		zap.L().Fatal("gorm.Open err", zap.Error(err))
	}
	db.Logger = db.Logger.LogMode(logger.LogLevel(options.RdbLogLevel))

	return db
}
