package database

import (
	"errors"
	"fmt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Agent struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"index" json:"updated_at"`
	SeenAt      time.Time      `gorm:"index" json:"seen_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	UserID      uint           `gorm:"index" json:"-"`
	AlertAfter  uint           `gorm:"index" json:"alert_after"`
	AlertPeriod uint           `gorm:"index" json:"alert_period"`
	AlertAt     time.Time      `gorm:"index" json:"alert_at"`
	Name        string         `json:"name"`
	Rules       datatypes.JSON `sql:"type:jsonb" json:"rules"`
	Token       string         `gorm:"index" json:"token"`
	Address     string         `json:"address"`
	UserAgent   string         `json:"user_agent"`
}

func (a *Agent) Save() error {
	return db.Save(a).Error
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

func FindAgentsForAlert() ([]*Agent, error) {
	var agents []*Agent
	query := "SELECT * FROM agents WHERE deleted_at IS NULL AND alert_after > 0 AND (seen_at IS NULL OR EXTRACT(EPOCH FROM current_timestamp-seen_at) >= alert_after)"
	err := db.Raw(query).Find(&agents).Error
	return agents, err
}

func RegisterAgent(fields *AgentWritableFields) (*Agent, error) {
	if err := fields.Validate(); err != nil {
		return nil, err
	}

	for i := range fields.Rules {
		fields.Rules[i].CreatedAt = time.Now()
	}

	newAgent := Agent{
		UserID:      fields.User.ID,
		Name:        fields.Name,
		Token:       makeRandomToken(),
		Rules:       ToJSONB(fields.Rules),
		AlertAfter:  fields.AlertAfter,
		AlertPeriod: fields.AlertPeriod,
	}

	if err := db.Create(&newAgent).Error; err != nil {
		return nil, fmt.Errorf("error creating new agent: %v", err)
	}

	return &newAgent, nil
}

func UpdateAgent(agent *Agent, fields *AgentWritableFields) error {
	// just to be sure
	fields.ID = agent.ID
	if err := fields.Validate(); err != nil {
		return err
	}

	for i, rule := range fields.Rules {
		if rule.TTL > 0 && rule.CreatedAt.IsZero() {
			fields.Rules[i].CreatedAt = time.Now()
		}
	}

	agent.Name = fields.Name
	agent.Rules = ToJSONB(fields.Rules)
	agent.UpdatedAt = time.Now()
	agent.AlertAfter = fields.AlertAfter
	agent.AlertPeriod = fields.AlertPeriod

	return db.Save(agent).Error
}
