package model

import "gorm.io/datatypes"

type ConfigKey string

const (
	ConfigKeyLicense ConfigKey = "license"
)

type ConfigLicense struct {
	License string `json:"license"`
}

// Config common config table struct
type Config struct {
	ID        int            `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Key       ConfigKey      `gorm:"index:uniq_key,unique; not null"`
	Value     datatypes.JSON `gorm:"type:json"`
	CreatedAt int64          `gorm:"not null; autoUpdateTime:milli"`
	UpdatedAt int64          `gorm:"not null; autoUpdateTime:milli"`
}

// TableName return table name
func (Config) TableName() string {
	return "configs"
}
