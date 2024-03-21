package model

import "gorm.io/datatypes"

type UserStatus int
type RoleType string

const (
	RoleTypeSuperAdmin   RoleType = "SuperAdmin"
	RoleTypeCollegeAdmin RoleType = "CollegeAdmin" // 学院管理员
	RoleTypeTeacher      RoleType = "Teacher"      // 老师
	RoleTypeStudent      RoleType = "Student"      // 学生
	RoleTypeNormal       RoleType = "User"         // 保留（预计给 企业 使用）

	SystemUsername            = "system"
	SuperAdminUsername        = "admin"
	SuperAdminDefaultPassword = "root"

	UserStatusNormal   UserStatus = 1
	UserStatusDisabled UserStatus = 2
)

type User struct {
	ID        int64      `gorm:"primary_key;AUTO_INCREMENT"`
	UID       string     `gorm:"not null; index:uniq_uid,unique; type:varchar(32)"` // hash_id
	Username  string     `gorm:"column:username; not null; index:uniq_username,unique; type:varchar(32)"`
	Name      string     `gorm:"column:name; not null; type:varchar(64)"`
	Role      RoleType   `gorm:"not null; type:varchar(32); default:0"`
	Password  string     `gorm:"column:password; not null; type:varchar(32)"`
	CreatedAt int64      `gorm:"column:created_at; not null; index:idx_created_at"`
	Creator   string     `gorm:"column:creator; not null; type:varchar(32)"`
	UpdatedAt int64      `gorm:"not null; default:0"`
	Updater   string     `gorm:"column:updater; not null; type:varchar(32)"`
	Status    UserStatus `gorm:"not null"`
	Phone     string     `gorm:"column:phone; type:varchar(32)"`
	Emial     string     `gorm:"column:email; type:varchar(32)"`
}

func (User) TableName() string {
	return "users"
}

type Profession struct {
	ID             int64  `gorm:"primary_key;AUTO_INCREMENT"`
	HashID         string `gorm:"not null; index:uniq_uid; type:varchar(64)"`
	CollegeName    string `gorm:"not null; type:varchar(32)"`
	CollegeInfo    string `gorm:"not null; type:varchar(32)"`
	ProfessionName string `gorm:"not null; type:varchar(32)"`
	ProfessionInfo string `gorm:"not null; type:varchar(32)"`

	CreatedAt int64  `gorm:"column:created_at; not null; index:idx_created_at"`
	Creator   string `gorm:"column:creator; not null; type:varchar(32)"`
	UpdatedAt int64  `gorm:"not null; default:0"`
	Updater   string `gorm:"column:updater; not null; type:varchar(32)"`
}

func (Profession) TableName() string {
	return "professions"
}

type Class struct {
	ID               int64  `gorm:"primary_key;AUTO_INCREMENT"`
	ProfessionHashID string `gorm:"not null; type:varchar(64)"`
	ClassName        string `gorm:"not null"`
	ClassID          int    `gorm:"not null"`

	CreatedAt int64  `gorm:"column:created_at; not null; index:idx_created_at"`
	Creator   string `gorm:"column:creator; not null; type:varchar(32)"`
	UpdatedAt int64  `gorm:"not null; default:0"`
	Updater   string `gorm:"column:updater; not null; type:varchar(32)"`
}

func (Class) TableName() string {
	return "classes"
}

type Company struct {
	ID          int64          `gorm:"primary_key;AUTO_INCREMENT"`
	CompanyName string         `gorm:"not null; type:varchar(64)"`
	CompanyInfo datatypes.JSON `gorm:"type:json"`
}
