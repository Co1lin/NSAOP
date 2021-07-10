package model

import "time"

type Refresh struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	CreateAt time.Time `gorm:"index"`
	Token    string    `gorm:"unique;type:varchar(255)"`
}
