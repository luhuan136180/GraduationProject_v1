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

func FindInterviewByOption(ctx context.Context, db *gorm.DB, option model.InterviewOption) (int64, []model.Interview, error) {
	var interviews []model.Interview

	if option.Title != "" {
		db = db.Where("title = ?", option.Title)
	}
	if option.IntervieweeUID != "" {
		db = db.Where("interviewee_uid = ?", option.IntervieweeUID)
	}
	if option.CreatorUID != "" {
		db = db.Where("creator_uid = ?", option.CreatorUID)
	}

	if len(option.Status) != 0 {
		db = db.Where("status in (?)", option.Status)
	}

	var count int64
	err := db.WithContext(ctx).Model(&model.Interview{}).Count(&count).Error
	if err != nil {
		return -1, nil, err
	}

	if option.Size != 0 && option.Page != 0 {
		db = db.Limit(option.Size).Offset((option.Page - 1) * option.Size)
	}

	err = db.WithContext(ctx).Model(&model.Interview{}).Find(&interviews).Error

	return count, interviews, err
}

func UpdateInterviewStatus(ctx context.Context, db *gorm.DB, id int64, status model.InterviewStatus) error {
	changeInfo := map[string]interface{}{
		"status":           status,
		"flag":             false,
		"contract_hash_id": nil,
		"contract_key_id":  nil,
	}

	err := db.WithContext(ctx).Model(&model.Interview{}).Where("id = ?", id).Updates(changeInfo).Error
	return err
}

func FindInterviewsUnContract(ctx context.Context, db *gorm.DB) ([]model.Interview, error) {
	var interviews []model.Interview
	err := db.WithContext(ctx).Where("flag != ?", true).Find(&interviews).Error
	return interviews, err
}
