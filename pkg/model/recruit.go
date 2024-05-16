package model

type Recruit struct {
	ID         int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	FirmHashID string `gorm:"column:firm_hash_id; not null" json:"firm_hash_id"`
	FirmName   string `gorm:"column:firm_name" json:"firm_name"`

	JobName      string `gorm:"column:job; not null" json:"job_name"`
	JobIntroduce string `gorm:"column:job_introduce" json:"job_introduce"` // 简介
	JobCondition string `gorm:"column:job_condition" json:"job_condition"` // 条件
	JobSalary    int    `gorm:"column:job_salary" json:"job_salary"`       // 薪资

	CreatedAt  int64  `gorm:"column:created_at; not null; index:idx_created_at" json:"created_at"`
	Creator    string `gorm:"column:creator; not null; type:varchar(32)" json:"creator"`    // teacher
	CreatorUID string `gorm:"not null; index:idx_uid; type:varchar(32)" json:"creator_uid"` // 面试发起人
}
