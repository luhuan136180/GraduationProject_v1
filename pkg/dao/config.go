package dao

import (
	"context"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"v1/pkg/model"
)

func GetConfigByKey(ctx context.Context, db *gorm.DB, key model.ConfigKey) (*model.Config, error) {
	conf := model.Config{}
	err := db.WithContext(ctx).First(&conf, "`key` = ?", key).Error
	return &conf, err
}

func UpsertConfig(ctx context.Context, db *gorm.DB, key model.ConfigKey, val datatypes.JSON) (*model.Config, error) {
	conf := model.Config{
		Key:   key,
		Value: val,
	}

	err := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
	}).Create(&conf).Error

	return &conf, err
}
