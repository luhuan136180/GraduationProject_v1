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

func GetFileByProjectID(ctx context.Context, db *gorm.DB, projectID int64) ([]model.File, error) {
	files := make([]model.File, 0)
	err := db.WithContext(ctx).Where("project_id = ?", projectID).Find(&files).Error
	return files, err
}
