package model

// 面试
type Interview struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT"`

	Interviewee string `gorm:"column:interviewee; not null"` // 面试者
}

func (Interview) TableName() string {
	return "interviews"
}
