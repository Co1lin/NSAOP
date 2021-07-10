package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	ID         uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Comment    string    `gorm:"type:varchar(30)"`
	Detail     string    `gorm:"type:varchar(600)"`
	NCESiteID  string    `gorm:"type:varchar(100)"`
	PayType    string    `gorm:"type:varchar(5)"`
	Status     string    `gorm:"index;type:varchar(8)"`
	LocationID uint      // foreign key
	Require    int       // 0~7 mask, 100 for private, 010 for client, 001 for test
	Users      []*User   `gorm:"many2many:user_service"`
	CreateAt   time.Time `gorm:"index"`
	PassAt     time.Time ``
	OnAt       time.Time ``
	Stamp      uint      ``
	Msg        string    ``
}

func (s *Service) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID, err = uuid.NewRandom()
	s.CreateAt = time.Now()
	s.PassAt = time.Unix(0, 0)
	s.OnAt = time.Unix(0, 0)
	return err
}
