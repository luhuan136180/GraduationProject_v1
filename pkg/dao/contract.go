package dao

import (
	"context"
	"gorm.io/gorm"
	"v1/pkg/model"
)

func InsertBlockSaveLog(ctx context.Context, db *gorm.DB, blockLog model.BlockSaveLog) (*model.BlockSaveLog, error) {
	if err := db.WithContext(ctx).Create(&blockLog).Error; err != nil {
		return nil, err
	}
	return &blockLog, nil
}

func GetBlocks(ctx context.Context, db *gorm.DB, blockHash, saveType string, page, size int) (int64, []model.BlockSaveLog, error) {

	list := make([]model.BlockSaveLog, 0)
	var count int64

	db = db.WithContext(ctx).Model(&model.BlockSaveLog{})

	if blockHash != "" {
		db = db.Where("block_hash = ?", blockHash)
	}

	if saveType != "" {
		db = db.Where("save_type = ?", saveType)
	}

	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	err = db.Offset((page - 1) * size).Limit(size).Find(&list).Error
	return count, list, err
}

func GetBlockByBlockHash(ctx context.Context, db *gorm.DB, hash string) (model.BlockSaveLog, error) {
	var log model.BlockSaveLog

	err := db.WithContext(ctx).Model(&model.BlockSaveLog{}).Where("block_hash = ?", hash).Find(&log).Error
	return log, err
}
