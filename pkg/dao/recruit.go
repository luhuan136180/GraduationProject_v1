package dao

import (
	"context"
	"gorm.io/gorm"
	"v1/pkg/model"
)

func GetRecruitList(ctx context.Context, db *gorm.DB, creatorUID string, jobName string, page, size int) (int64, []model.Recruit, error) {
	list := make([]model.Recruit, 0)
	var count int64

	db = db.WithContext(ctx).Model(&model.Recruit{})

	if jobName != "" {
		db = db.Where("job = ?", jobName)
	}

	if creatorUID != "" {
		db = db.Where("creator_uid = ?", creatorUID)
	}

	err := db.Count(&count).Error
	if err != nil {
		return 0, nil, err
	}

	err = db.Limit(size).Offset((page - 1) * size).Find(&list).Error
	if err != nil {
		return 0, nil, err
	}
	return count, list, nil

}

func InsertRecruit(ctx context.Context, db *gorm.DB, data model.Recruit) (*model.Recruit, error) {

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func GetRecruitByID(ctx context.Context, db *gorm.DB, id int) (model.Recruit, error) {
	var data model.Recruit
	err := db.WithContext(ctx).Model(&model.Recruit{}).Where("id = ?", id).Find(&data).Error
	return data, err
}

func DeleteRecruit(ctx context.Context, db *gorm.DB, id int) error {
	return db.WithContext(ctx).Where("id = ?", id).Delete(&model.Recruit{}).Error
}
