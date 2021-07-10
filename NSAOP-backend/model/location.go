package model

import "gorm.io/gorm"

type Location struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Comment   string         `gorm:"type:varchar(30)"`
	Address   string         `gorm:"index;type:varchar(300)"`
	Contact   string         `gorm:"type:varchar(60)"`
	Phone     string         `gorm:"type:varchar(20)"`
	UserID    uint           `` // foreign key
	Services  []*Service     `` // has many
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
