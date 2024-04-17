package model

import "gorm.io/datatypes"

type UserStatus int
type RoleType string

const (
	RoleTypeSuperAdmin   RoleType = "SuperAdmin"
	RoleTypeCollegeAdmin RoleType = "CollegeAdmin" //  学院管理员
	RoleTypeTeacher      RoleType = "Teacher"      // 老师
	RoleTypeStudent      RoleType = "Student"      // 学生
	RoleTypeNormal       RoleType = "User"         // 保留（预计给 企业 使用）
	RoleTypeFirm         RoleType = "firm"         // 企业

	SystemUsername            = "system"
	SuperAdminUsername        = "admin"
	SuperAdminDefaultPassword = "root"

	UserStatusNormal   UserStatus = 1
	UserStatusDisabled UserStatus = 2
)

type User struct {
	ID               int64    `gorm:"primary_key;AUTO_INCREMENT"`
	UID              string   `gorm:"not null; index:uniq_uid,unique; type:varchar(32)"`                       // hash_id
	Username         string   `gorm:"column:username; not null; index:uniq_username,unique; type:varchar(32)"` // 用户名
	Name             string   `gorm:"column:name; not null; type:varchar(64)"`                                 // 昵称
	Role             RoleType `gorm:"not null; type:varchar(32); default:0"`
	Password         string   `gorm:"column:password; not null; type:varchar(32)"`
	ProfessionHashID string   `gorm:"column:profession_hash_id; not null;type:varchar(64)"`
	ClassHashID      string   `gorm:"column:class_hash_id; type:varchar(64)"`

	CreatedAt int64      `gorm:"column:created_at; not null; index:idx_created_at"`
	Creator   string     `gorm:"column:creator; not null; type:varchar(32)"`
	UpdatedAt int64      `gorm:"not null; default:0"`
	Updater   string     `gorm:"column:updater; not null; type:varchar(32)"`
	Status    UserStatus `gorm:"not null"`
	Phone     string     `gorm:"column:phone; type:varchar(32)"`
	Emial     string     `gorm:"column:email; type:varchar(32)"`
	Head      string     `json:"column:head"` // 头像
}

func (User) TableName() string {
	return "users"
}

type UserOption struct {
	UserNameOption    string     `json:"user_name"`
	ProfessionHashIDs []string   `json:"profession_hash_ids"`
	ClassHashIDs      []string   `json:"class_hash_ids"`
	RoleTypes         []RoleType `json:"role_types"`
	Status            []string   `json:"status"`
}

type College struct {
	ID          int64  `gorm:"primary_key;AUTO_INCREMENT"`
	HashID      string `gorm:"not null; type:varchar(64); index:idx_college_hash_id"` // collegename + professionname
	CollegeName string `gorm:"not null; type:varchar(32); index:idx_college_name"`
	CollegeInfo string `gorm:"not null; type:varchar(32)"`

	CreatedAt int64  `gorm:"column:created_at; not null"`
	Creator   string `gorm:"column:creator; not null; type:varchar(32)"`
	UpdatedAt int64  `gorm:"not null; default:0"`
	Updater   string `gorm:"column:updater; not null; type:varchar(32)"`
}

type Profession struct {
	HashID      string `gorm:"not null; type:varchar(64); index:idx_profession_hash_id"` // collegename + professionname
	CollegeName string `gorm:"not null; type:varchar(32)"`
	// CollegeInfo    string `gorm:"not null; type:varchar(32)"`
	CollegeHashID  string `gorm:"column:college_hash_id"`
	ProfessionName string `gorm:"not null; type:varchar(32); index:idx_profession_name"`
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
	ProfessionHashID string `gorm:"column:profession_hash_id;not null; type:varchar(64)"`
	ClassHashID      string `gorm:"column:class_hash_id;not null; type: varchar(64); index:idx_class_hash_id"`
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

type CollegeInfo struct {
	Info string `json:"info"`
}

type ProfessionInfo struct {
	Info string `json:"info"`
}

type Company struct {
	ID          int64          `gorm:"primary_key;AUTO_INCREMENT"`
	CompanyName string         `gorm:"not null; type:varchar(64)"`
	CompanyInfo datatypes.JSON `gorm:"type:json"`
}
