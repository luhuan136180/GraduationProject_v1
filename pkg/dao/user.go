package dao

import (
	"context"
	"gorm.io/gorm"
	"strconv"
	"time"
	"v1/pkg/model"
	"v1/pkg/utils"
)

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

func InsertUser(ctx context.Context, db *gorm.DB, username, password, name, creator, phone, email string, role model.RoleType) (*model.User, error) {
	now := time.Now().UnixMilli()
	user := model.User{
		UID:       utils.NextID(),
		Username:  username,
		Name:      name,
		Password:  utils.MD5Hex(password),
		Role:      role,
		CreatedAt: now,
		Creator:   creator,
		UpdatedAt: now,
		Updater:   creator,
		Status:    model.UserStatusNormal,
		Phone:     phone,
		Emial:     email,
	}

	if err := db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
