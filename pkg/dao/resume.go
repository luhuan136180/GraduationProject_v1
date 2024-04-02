package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"v1/pkg/model"
)

func GetResumeByID(ctx context.Context, db *gorm.DB, id int64) (bool, model.Resume, error) {
	var resume model.Resume
	if err := db.WithContext(ctx).Where("id = ?", id).First(&resume).Error; err != nil {
		return false, resume, err
	}
	return true, resume, nil
}

func GetResumesByUserUid(ctx context.Context, db *gorm.DB, uid string) ([]model.Resume, error) {
	var resumes []model.Resume

	if err := db.WithContext(ctx).Where("user_uid = ?", uid).Find(&resumes).Error; err != nil {
		return resumes, err
	}
	return resumes, nil
}
func FoundResumesByname(ctx context.Context, db *gorm.DB, name string) (bool, model.Resume, error) {
	var resumes model.Resume

	if err := db.WithContext(ctx).Where("resume_name = ?", name).First(&resumes).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, resumes, nil
		}
		return false, resumes, err
	}
	return true, resumes, nil
}

func InsertResume(ctx context.Context, db *gorm.DB, resume model.Resume) (*model.Resume, error) {
	resume.CreatedAt = time.Now().UnixMilli()

	if err := db.WithContext(ctx).Create(&resume).Error; err != nil {
		return nil, err
	}
	return &resume, nil
}

func DeleteResumeByID(ctx context.Context, db *gorm.DB, id int64) error {
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error; err != nil {
		return err
	}

	return nil
}
