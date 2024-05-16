package model

type Contract struct {
	Flag           bool   `gorm:"column:flag;not null default:false"` // 是否上链; false:没有;true:上链
	BlockHash      string `gorm:"column:block_hash; type: varchar(256) "`
	ContractHashID string `gorm:"column:contract_hash_id; type: varchar(256)"`
	ContractKeyID  string `gorm:"column:contract_key_id"`
}

type BlockSaveLog struct {
	ID          int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	BlockTXHash string `gorm:"column:block_hash;type:varchar(256)" json:"block_hash"`
	SaveType    string `gorm:"column:save_type" json:"save_type"`
	KeyHash     string `gorm:"column: key_hash" json:"key_hash"`
	Date        string `gorm:"column:date" json:"date"`
}
