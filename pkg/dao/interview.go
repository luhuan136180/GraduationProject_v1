package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"v1/pkg/model"
)

func InsertInterview(ctx context.Context, db *gorm.DB, interview model.Interview) (*model.Interview, error) {
	interview.CreatedAt = time.Now().UnixMilli()
	interview.Status = model.InterviewStatusPost

	if err := db.WithContext(ctx).Create(&interview).Error; err != nil {
		return nil, err
	}
	return &interview, nil
}

func GetInterviewByID(ctx context.Context, db *gorm.DB, id int64) (model.Interview, error) {
	var interview model.Interview

	if err := db.WithContext(ctx).Where("id = ?", id).First(&interview).Error; err != nil {
		return interview, err
	}
	return interview, nil
}

func DeleteInterviewByID(ctx context.Context, db *gorm.DB, id int64) error {
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&model.Interview{}).Error; err != nil {
		return err
	}

	return nil
}
