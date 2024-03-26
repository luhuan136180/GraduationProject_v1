package dao

import (
	"context"
	"gorm.io/gorm"
	"v1/pkg/model"
)

func GetAuditLogs(ctx context.Context, db *gorm.DB, page, size int, ip, Account, username string, methodList, statusList []string, createAt, endAt int64) (int64, []model.AuditLog, error) {
	db = db.WithContext(ctx).Model(&model.AuditLog{})

	if ip != "" {
		db = db.Where("source_ip LIKE ?", "%"+ip+"%")
	}
	if Account != "" {
		db = db.Where("account LIKE ?", "%"+Account+"%")
	}
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	if len(methodList) != 0 {
		db = db.Where("method in ?", methodList)
	}
	if len(statusList) != 0 {
		db = db.Where("status in ?", statusList)
	}
	if createAt > 0 {
		db = db.Where("created_at > ?", createAt)
	}
	if endAt > 0 {
		db = db.Where("created_at < ?", endAt)
	}

	var count int64
	var logList []model.AuditLog
	if err := db.Count(&count).Error; err != nil {
		return 0, nil, err
	}
	if size <= 0 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size
	db = db.Limit(size).Offset(offset)

	if err := db.Order("id DESC").Find(&logList).Error; err != nil {
		return count, logList, err
	}

	return count, logList, nil
}

func InsertLog(ctx context.Context, db *gorm.DB, username, name, bluePrint, method, ip, uri string, status int, duration, createAt int64) error {
	log := model.AuditLog{
		Username:  username,
		Name:      name,
		BluePrint: bluePrint,
		Method:    method,
		Duration:  duration,
		IP:        ip,
		Status:    status,
		Uri:       uri,
		CreatedAt: createAt,
	}

	if err := db.WithContext(ctx).Create(&log).Error; err != nil {
		return err
	}
	return nil
}
