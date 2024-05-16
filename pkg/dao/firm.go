package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"v1/pkg/model"
)

func GetFirmList(ctx context.Context, db *gorm.DB, firmName string, page, size int) (int64, []*model.Firm, error) {
	var firms []*model.Firm
	var count int64

	db = db.WithContext(ctx).Model(&model.Firm{})
	if firmName != "" {
		db = db.Where("firm_name LIKE ?", firmName)
	}

	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	err = db.Offset((page - 1) * size).Limit(size).Find(&firms).Error
	return count, firms, err
}

func GetFirmTree(ctx context.Context, db *gorm.DB) ([]*model.Firm, error) {
	var firms []*model.Firm

	db = db.WithContext(ctx).Model(&model.Firm{})

	err := db.Find(&firms).Error
	return firms, err
}

func GetFirmByHashID(ctx context.Context, db *gorm.DB, hashID string) (*model.Firm, error) {
	var firm model.Firm
	err := db.WithContext(ctx).Model(&model.Firm{}).Where("firm_hash_id = ?", hashID).First(&firm).Error
	return &firm, err
}

func InsertFirm(ctx context.Context, db *gorm.DB, firmInfo model.Firm) (*model.Firm, error) {
	now := time.Now().UnixMilli()

	firmInfo.CreatedAt = now
	if err := db.WithContext(ctx).Create(&firmInfo).Error; err != nil {
		return nil, err
	}

	return &firmInfo, nil
}
