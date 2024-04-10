package model

type Contract struct {
	Flag           bool   `gorm:"column:flag; default:false"` // 是否上链; false:没有;true:上链
	ContractHashID string `gorm:"column:contract_hash_id; type: varchar(256)"`
	ContractKeyID  string `gorm:"column:contract_key_id"`
}
