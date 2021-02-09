package database

import (
	"errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Agent struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	UserID    uint           `gorm:"index"`
	Rules     datatypes.JSON `sql:"type:jsonb"`
	Token     string         `gorm:"index"`
	Address   string
	UserAgent string
}

func FindAgentByToken(token string) (*Agent, error) {
	var found Agent
	if err := db.Where("token=?", token).First(&found).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err == nil {
		return &found, nil
	} else {
		return nil, err
	}
}

func (a *Agent) Save() error {
	return db.Save(a).Error
}
