package model

type File struct {
	ID       int    `gorm:"primary_key;AUTO_INCREMENT"`
	FileName string `gorm:"not null;type:varchar(64)"`
	FilePath string `gorm:"not null;type:varchar(64)"`
}
