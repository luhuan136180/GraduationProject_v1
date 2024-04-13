package model

import "math/big"

type ConfigurationTag string

const (
	ContractTAG ConfigurationTag = "contract"
)

// 配置存储表
type Configuration struct {
	ID      int         `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Tag     string      `gorm:"column:tag; not null; type:varchar(64)"`
	Content interface{} `gorm:"column:content; type:json; serializer:json"`
}

type ContractConfig struct {
	NetworkIP         string   `json:"network_ip"`          // ip
	AccountPrivateKey string   `json:"account_private_key"` // 所用账号私钥
	Address           string   `json:"address"`             // 合约地址
	GasFeeCap         *big.Int `json:"gas_fee_cap"`         // Gas费用上限
	GasLimit          uint64   `json:"gas_limit"`           // Gas费用上限(单次交易) ps：默认设置为0
}
