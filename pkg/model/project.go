package model

import "gorm.io/datatypes"

type ProjectStatus int64

const (
	ProjectStatusAudit   ProjectStatus = 1 // 审核
	ProjectStatusClose   ProjectStatus = 2 // 关闭
	ProjectStatusProceed ProjectStatus = 3 // 进行
	ProjectStatusFinish  ProjectStatus = 4 // 完成
	ProjectStatusPASS    ProjectStatus = 5 // 通过审核

)

type Project struct {
	ID               int64          `gorm:"primary_key;AUTO_INCREMENT"`
	ProjectName      string         `gorm:"not null; type:varchar(32)"`
	ProjectBasicInfo datatypes.JSON `gorm:"type:json"`
	ProjectFile      []byte         `gorm:"column:project_file"`
	Title            string         `gorm:"type:varchar(32)"`
	Status           ProjectStatus  `gorm:"not null default:2"`
	ProfessionHashID string         `gorm:"not null; type:varchar(64)"`

	CreatedAt      int64  `gorm:"column:created_at; not null; index:idx_created_at"`
	Creator        string `gorm:"column:creator; not null; type:varchar(32)"` // teacher
	CreatorUID     string `gorm:"not null; index:idx_uid; type:varchar(32)"`
	AuditUID       string `gorm:"column:audit_uid;not null; type:varchar(32)"` // admin
	Auditor        string `gorm:"column:auditor;not null;type:varchar(32)"`
	Participator   string `gorm:"type:varchar(64)"` // xuesheng
	ParticipatorID string `gorm:"type:varchar(64)"`
}

func (Project) TableName() string {
	return "projects"
}

type DifficultyType string

const (
	DifficultyTypeHard   DifficultyType = "HARD"
	DifficultyTypeEASY   DifficultyType = "EASY"
	DifficultyTypeNORMAL DifficultyType = "NORMAL"
)

type ProjectBasicInfo struct {
	Difficulty  DifficultyType `json:"difficulty"`
	BackGround  string         `json:"back_ground"`
	Requirement string         `json:"requirement"`
	plan        string         `json:"plan"`
}

type ProjectOption struct {
	ProjectName string   `json:"project_name"`
	Title       string   `json:"title"`
	Creator     string   `json:"creator"`
	Auditor     string   `json:"auditor"`
	Status      []string `json:"status"`
	Professions []string `json:"professions"` // profession_hash_ids
}

type ProjectSelectLog struct {
	ID           int64  `gorm:"primary_key;AUTO_INCREMENT"`
	ProjectID    int64  `gorm:"not null; type:varchar(32)"`
	ProjectName  string `gorm:"not null; type:varchar(32)"`
	Applicant    string `gorm:"not null; type:varchar(32)"`
	ApplicantUID string `gorm:"not null; type:varchar(32)"`
}

func (ProjectSelectLog) TableName() string {
	return "projects"
}
