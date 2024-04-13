package dao

import (
	"context"
	"gorm.io/gorm"
	"v1/pkg/model"
)

func InitConfiguration(ctx context.Context, db *gorm.DB, configuration model.Configuration) (*model.Configuration, error) {
	if err := db.WithContext(ctx).Create(&configuration).Error; err != nil {
		return nil, err
	}
	return &configuration, nil
}
