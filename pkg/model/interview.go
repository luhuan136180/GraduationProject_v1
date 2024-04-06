package model

type InterviewStatus string

const (
	InterviewStatusPost    InterviewStatus = "Post"    // 发起
	InterviewStatusAccept  InterviewStatus = "Accept"  // 接受
	InterviewStatusRefuse  InterviewStatus = "refuse"  // 拒绝
	InterviewStatusProceed InterviewStatus = "Proceed" // 进行
	InterviewStatusFailed  InterviewStatus = "Failed"  // 失败

)

// 面试
type Interview struct {
	ID    int64       `gorm:"primary_key;AUTO_INCREMENT"`
	Ttile string      `gorm:"column:title; not null"`
	Info  interface{} `gorm:"type:json; serializer:json"`

	Interviewee    string          `gorm:"column:interviewee; not null"` // 面试者_name
	IntervieweeUID string          `gorm:"type:varchar(64); not null"`   // 面试者_uid
	Status         InterviewStatus `gorm:"type:varchar(63); not null"`

	CreatedAt  int64  `gorm:"column:created_at; not null; index:idx_created_at"`
	Creator    string `gorm:"column:creator; not null; type:varchar(32)"` // teacher
	CreatorUID string `gorm:"not null; index:idx_uid; type:varchar(32)"`  // 面试发起人
}

func (Interview) TableName() string {
	return "interviews"
}

type InterviewInfo struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Date     string `json:"date"`
	Location string `json:"location"`
	Position string `json:"position"` // 职位

	Creator     string `json:"creator"`      // user_name
	ContactInfo string `json:"contact_info"` // 联系方式

	Interviewee string `json:"interviewee"` // 面试者
}

type InterviewOption struct {
	Title string `json:"title"`

	IntervieweeUID string `json:"interviewee_uid"`
	CreatorUID     string `json:"creator_uid"`

	Page int `json:"page"`
	Size int `json:"size"`
}
