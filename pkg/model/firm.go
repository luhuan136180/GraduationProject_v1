package model

type Firm struct {
	ID int64 `gorm:"primary_key;AUTO_INCREMENT"`

	FirmHashID string `gorm:"column:firm_hash_id; not null" json:"firm_hash_id"`
	FirmName   string `gorm:"column:firm_name; not null" json:"firm_name"`
	FirmInfo   string `gorm:"column:firm_info; not null" json:"firm_info"`

	CreatedAt  int64  `gorm:"column:created_at; not null; index:idx_created_at" json:"created_at"`
	Creator    string `gorm:"column:creator; not null; type:varchar(32)" json:"creator"`    // teacher
	CreatorUID string `gorm:"not null; index:idx_uid; type:varchar(32)" json:"creator_uid"` // 面试发起人
}
