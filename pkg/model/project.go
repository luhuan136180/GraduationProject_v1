package model

import "gorm.io/datatypes"

type projectStatus int64

const (
	ProjectStatusAudit   = 1 // 审核
	ProjectStatusClose   = 2 // 关闭
	ProjectStatusProceed = 3 // 进行
	ProjectStatusFinish  = 4 // 完成
)

type Project struct {
	ID               int64          `gorm:"primary_key;AUTO_INCREMENT"`
	ProjectName      string         `gorm:"not null; type:varchar(32)"`
	ProjectBasicInfo datatypes.JSON `gorm:"type:json"`
	ProjectFile      []byte         `gorm:"column:project_file"`
	Title            string         `gorm:"type；varchar(32)"`
	Status           projectStatus  `gorm:"not null default:2"`

	CreatedAt  int64  `gorm:"column:created_at; not null; index:idx_created_at"`
	Creator    string `gorm:"column:creator; not null; type:varchar(32)"` // student
	CreatorUID string `gorm:"not null; index:idx_uid; type:varchar(32)"`
	AuditUID   string `gorm:"column:audit_uid;not null; type:varchar(32)"` // teacher
	Auditor    string `gorm:"column:auditor;not null;type:varchar(32)"`
}

func (Project) TableName() string {
	return "projects"
}
