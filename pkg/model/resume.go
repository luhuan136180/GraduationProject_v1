package model

// 简历
type Resume struct {
	ID         int64       `gorm:"primary_key;AUTO_INCREMENT"`
	UserUid    string      `gorm:"not null; type:varchar(32)"`
	UserName   string      `gorm:"not null; type；varchar(64)"`
	ResumeName string      `gorm:"not null; type:varchar(32)"`
	BasicInfo  interface{} `gorm:"not null; type:json; serializer:json"`
	ProjectIDs []int64     `gorm:" type:json; serializer:json"` // 绑定的项目ids

	Creator   string `gorm:"column:creator; not null; type:varchar(64)"`
	CreatedAt int64  `gorm:"column:created_at; not null; index:idx_created_at"`
}

func (Resume) TableName() string {
	return "resumes"
}

type ResumeInfo struct {
	Name           string `json:"name"`
	CollegeName    string `json:"collegeName"`
	PorfessionName string `json:"porfessionName"`
	Describe       string `json:"describe"`
	Experience     string `json:"experience"`
}

type ResumeOption struct {
	Title    string `json:"title"` // resume_name
	UserName string `json:"userName"`
	StartAt  string `json:"startAt"`
	EndAt    string `json:"endAt"`
}
