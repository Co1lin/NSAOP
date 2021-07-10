package model

import "time"

type User struct {
	ID        	uint        `gorm:"primaryKey;autoIncrement"`
	Role      	string      `index;gorm:"type:varchar(10);not null;default:'customer'"`
	Username  	string      `gorm:"type:varchar(60);unique"`
	Password  	string      `gorm:"type:varchar(60)"`
	Company   	string      `gorm:"type:varchar(30)"`
	Phone     	string      `gorm:"type:varchar(20)"`
	Email     	string      `gorm:"type:varchar(40)"`
	Locations 	[]*Location ``
	Services  	[]*Service  `gorm:"many2many:user_service;"`
	ResetToken	string		`gorm:"type:varchar(32)"`
	LastOper	time.Time
}
