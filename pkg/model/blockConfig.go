package model

type BlockConfig struct {
	ID        int64 `gorm:"primary_key;AUTO_INCREMENT"`
	BlockID   string
	BlockInfo string
	// 未操作
	CreatedAt int64 `gorm:"column:created_at; not null; index:idx_created_at"`
}
