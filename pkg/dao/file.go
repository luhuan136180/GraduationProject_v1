package dao

import (
	"context"
	"gorm.io/gorm"
	"v1/pkg/model"
)

func InsterFile(ctx context.Context, db *gorm.DB, file model.File) (*model.File, error) {
	if err := db.WithContext(ctx).Create(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}
