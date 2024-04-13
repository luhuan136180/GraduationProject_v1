package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
	"v1/pkg/model"
)

func InsertProject(ctx context.Context, db *gorm.DB, project model.Project) (*model.Project, error) {

	project.CreatedAt = time.Now().UnixMilli()
	project.Status = model.ProjectStatusAudit

	if err := db.WithContext(ctx).Create(&project).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func GetProjectByID(ctx context.Context, db *gorm.DB, id int64) (bool, model.Project, error) {
	var project model.Project
	if err := db.WithContext(ctx).Where("id = ?", id).First(&project).Error; err != nil {
		return false, project, err
	}
	return true, project, nil
}

func UpdateProjectStatus(ctx context.Context, db *gorm.DB, id int64, status model.ProjectStatus) error {
	err := db.WithContext(ctx).Model(&model.Project{}).Where("id = ?", id).Update("status", status).Error

	return err
}

func UpdateProjectParticipator(ctx context.Context, db *gorm.DB, id int64, user model.User) error {
	changeInfo := map[string]interface{}{
		"participator":     user.Username,
		"participator_id":  user.UID,
		"status":           model.ProjectStatusProceed,
		"flag":             false,
		"contract_hash_id": nil,
		"contract_key_id":  nil,
	}
	err := db.WithContext(ctx).Model(&model.Project{}).Where("id = ?", id).Updates(changeInfo).Error

	return err
}

func DeleteProjectByID(ctx context.Context, db *gorm.DB, id int64) error {
	if err := db.WithContext(ctx).Where("id = ?", id).Delete(&model.Project{}).Error; err != nil {
		return err
	}

	return nil
}

// 查询未上链的projects
func FindProjectsUnContract(ctx context.Context, db *gorm.DB) ([]model.Project, error) {
	var projects []model.Project
	err := db.WithContext(ctx).Where("flag != ?", true).Find(&projects).Error
	return projects, err
}

func UpdateProjectFiles(ctx context.Context, db *gorm.DB, id int64, fileID int) error {
	_, project, err := GetProjectByID(ctx, db, id)
	if err != nil {
		return err
	}
	var list []int
	if project.ProjectFile == nil {
		list = make([]int, 0)
	} else {
		list = project.ProjectFile
	}
	list = append(list, fileID)

	changeInfo := map[string]interface{}{
		"project_file":     list,
		"flag":             false,
		"contract_hash_id": nil,
		"contract_key_id":  nil,
	}
	err = db.WithContext(ctx).Model(&model.Project{}).Where("id = ?", id).Updates(changeInfo).Error
	return err
}
