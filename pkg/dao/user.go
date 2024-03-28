package dao

import (
	"context"
	"gorm.io/gorm"
	"strconv"
	"time"
	"v1/pkg/model"
)

// user
func GetUserByUsername(ctx context.Context, db *gorm.DB, username string) (*model.User, error) {
	u := model.User{}
	err := db.WithContext(ctx).First(&u, "username = ?", username).Error
	return &u, err
}

func GetUserByID(ctx context.Context, db *gorm.DB, id string) (bool, *model.User, error) {
	ID, err := strconv.ParseInt(id, 10, 64)
	if ID < 0 {
		return false, nil, nil
	}
	u := model.User{}

	if err != nil {
		return false, nil, err
	}
	err = db.WithContext(ctx).Model(&u).Where("id = ?", ID).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, &u, nil
}

// 非系统模块用 uid
func GetUserByUID(ctx context.Context, db *gorm.DB, uid string) (bool, *model.User, error) {
	var user model.User
	if err := db.WithContext(ctx).Model(&model.User{}).Where("uid = ?", uid).First(&user).Error; err != nil {
		return false, &user, err
	}
	return true, &user, nil
}

func GetUserByAccount(ctx context.Context, db *gorm.DB, account string) (bool, *model.User, error) {
	u, err := GetUserByUsername(ctx, db, account)
	if err == gorm.ErrRecordNotFound {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}
	return true, u, nil
}

func InsertUser(ctx context.Context, db *gorm.DB, userInfo model.User) (*model.User, error) {
	now := time.Now().UnixMilli()

	userInfo.CreatedAt = now
	userInfo.UpdatedAt = now
	if err := db.WithContext(ctx).Create(&userInfo).Error; err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func DeleteUserByID(ctx context.Context, db *gorm.DB, id string) error {
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return err
	}
	if err := db.WithContext(ctx).Where("id = ?", ID).Delete(&model.User{}).Error; err != nil {
		return err
	}

	return nil
}

func GETUserList(ctx context.Context, db *gorm.DB, page, size int, userOption model.UserOption) (int64, []model.User, error) {
	limit := (page - 1) * size
	var users []model.User
	var count int64

	db = db.WithContext(ctx)
	if userOption.UserNameOption != "" {
		db = db.Where("username LIKE ? ", "%"+userOption.UserNameOption+"%")
	}
	if len(userOption.Status) != 0 {
		db = db.Where("status in (?)", userOption.Status)
	}
	if len(userOption.ProfessionHashIDs) != 0 {
		db = db.Where("profession_hash_id in (?)", userOption.ProfessionHashIDs)
	}
	if len(userOption.ClassHashIDs) != 0 {
		db = db.Where("class_hash_id in (?)", userOption.ClassHashIDs)
	}
	if len(userOption.RoleTypes) != 0 {
		db = db.Where("role in (?)", userOption.RoleTypes)
	}

	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return 0, nil, err
	}

	if err := db.Limit(size).Offset(limit).Order("updated_at DESC").Find(&users).Error; err != nil {
		return count, users, err
	}

	return count, users, nil
}

func UpdateUserInfo(ctx context.Context, db *gorm.DB, id string, userInfo model.User) error {
	ID, err := strconv.ParseInt(id, 10, 64)
	now := time.Now().UnixMilli()
	if err != nil {
		return err
	}
	changeInfo := map[string]interface{}{
		"name":  userInfo.Name,
		"phone": userInfo.Phone,
		"email": userInfo.Phone,

		"profession_hash_id": userInfo.ProfessionHashID,
		"class_hash_id":      userInfo.ClassHashID,

		"updated_at": now,
		"updater":    userInfo.Updater,
	}

	if err = db.WithContext(ctx).Model(&model.User{}).Where("id = ?", ID).Updates(changeInfo).Error; err != nil {
		return err
	}

	return nil
}

// college
func GetCollegeByHashID(ctx context.Context, db *gorm.DB, collegeHashID string) (bool, model.College, error) {
	var collegeItem model.College
	err := db.WithContext(ctx).Model(&model.College{}).Where("hash_id = ?", collegeHashID).First(&collegeItem).Error
	if err != nil {
		return false, collegeItem, err
	}
	return true, collegeItem, nil
}

func InsertCollege(ctx context.Context, db *gorm.DB, collegeInfo model.College) (*model.College, error) {
	now := time.Now().UnixMilli()

	collegeInfo.CreatedAt = now
	collegeInfo.UpdatedAt = now
	if err := db.WithContext(ctx).Create(&collegeInfo).Error; err != nil {
		return nil, err
	}

	return &collegeInfo, nil
}

func DeleteCollege(ctx context.Context, db *gorm.DB, hashID string) error {
	return db.WithContext(ctx).Where("hash_id = ?", hashID).Delete(&model.College{}).Error
}

// profession
func GetProfessionByHashID(ctx context.Context, db *gorm.DB, professionHashID string) (bool, model.Profession, error) {
	var professionItem model.Profession
	err := db.WithContext(ctx).Model(&model.Profession{}).Where("hash_id = ?", professionHashID).First(&professionItem).Error
	if err != nil {
		return false, professionItem, err
	}
	return true, professionItem, nil
}

func GetProfessionsByHashIDs(ctx context.Context, db *gorm.DB, professionHashIDs []string) ([]model.Profession, error) {
	professions := make([]model.Profession, 0)
	err := db.WithContext(ctx).Model(&model.Profession{}).Where("hash_id in  (?)", professionHashIDs).Find(&professions).Error
	return professions, err
}

func GetProfessionByProfessionName(ctx context.Context, db *gorm.DB, professionNames []string) ([]model.Profession, error) {
	var professions []model.Profession
	err := db.WithContext(ctx).Model(&model.Profession{}).Where("profession_name in  (?)", professionNames).Find(&professions).Error
	return professions, err
}

func InsertProfession(ctx context.Context, db *gorm.DB, professionInfo model.Profession) (*model.Profession, error) {
	now := time.Now().UnixMilli()

	professionInfo.CreatedAt = now
	professionInfo.UpdatedAt = now
	if err := db.WithContext(ctx).Create(&professionInfo).Error; err != nil {
		return nil, err
	}

	return &professionInfo, nil
}

func DeleteProfession(ctx context.Context, db *gorm.DB, hashID string) error {
	return db.WithContext(ctx).Where("hash_id = ?", hashID).Delete(&model.Profession{}).Error
}

// class
func GetClassByHashID(ctx context.Context, db *gorm.DB, ClassHashID string) (bool, model.Class, error) {
	var ClassItem model.Class
	err := db.WithContext(ctx).Model(&model.Class{}).Where("class_hash_id = ?", ClassHashID).First(&ClassItem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, ClassItem, nil
		}
		return false, ClassItem, err
	}

	return true, ClassItem, nil
}

func GetClassByHashIDs(ctx context.Context, db *gorm.DB, ClassHashIDs []string) ([]model.Class, error) {
	var ClassItem []model.Class
	err := db.WithContext(ctx).Model(&model.Class{}).Where("class_hash_id in (?)", ClassHashIDs).Find(&ClassItem).Error

	return ClassItem, err
}

func InsertClass(ctx context.Context, db *gorm.DB, ClassInfo model.Class) (*model.Class, error) {
	now := time.Now().UnixMilli()

	ClassInfo.CreatedAt = now
	ClassInfo.UpdatedAt = now
	if err := db.WithContext(ctx).Create(&ClassInfo).Error; err != nil {
		return nil, err
	}

	return &ClassInfo, nil
}

func DeleteClass(ctx context.Context, db *gorm.DB, hashID string) error {
	return db.WithContext(ctx).Where("class_hash_id = ?", hashID).Delete(&model.Class{}).Error
}
