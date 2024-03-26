package model

type AuditLog struct {
	ID        int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Username  string `gorm:"column:account; not null; index:account; type:varchar(32)"`
	Name      string `gorm:"column:username; not null; index:username; type:varchar(32)"`
	BluePrint string `gorm:"column:blueprint; not null; type:varchar(32)"`
	Method    string `gorm:"column:method; not null;type:varchar(32)"`
	Duration  int64  `gorm:"column:duration; not null"`
	IP        string `gorm:"column:source_ip; not null; type:varchar(32)"`
	Status    int    `gorm:"column:status; not null"`
	Uri       string `gorm:"column:uri;not null;type:varchar(255)"`
	CreatedAt int64  `gorm:"column:created_at; not null"`
}

func (AuditLog) TableName() string {
	return "api_request_logs"
}
