package model

type Resume struct {
	ID         int64       `gorm:"primary_key;AUTO_INCREMENT"`
	UserUid    string      `gorm:"not null; type:varchar(32)"`
	ResumeName string      `gorm:"not null; type:varchar(32)"`
	BasicInfo  interface{} `gorm:"not null; type:json; serializer:json"`
	ProjectIDs []int64     `gorm:"not null; type:json; serializer:json"` // 绑定的项目ids

}
