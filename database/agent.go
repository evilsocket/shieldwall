package database

import (
	"errors"
	"fmt"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/shieldwall/firewall"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const MinAgentNameLength = 3
const MaxAgentRules = 10

type Agent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"index" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `json:"name"`
	UserID    uint           `gorm:"index" json:"-"`
	Rules     datatypes.JSON `sql:"type:jsonb" json:"rules"`
	Token     string         `gorm:"index" json:"token"`
	Address   string         `json:"address"`
	UserAgent string         `json:"user_agent"`
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

func validateAgentData(name string, rules []*firewall.Rule) error {
	if name = str.Trim(name); len(name) < MinAgentNameLength {
		return fmt.Errorf("agent name must be at least %d characters long", MinAgentNameLength)
	}

	if len(rules) > MaxAgentRules {
		return fmt.Errorf("max %d rules per agent are allowed", MaxAgentRules)
	}

	for _, rule := range rules {
		if err := rule.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func RegisterAgent(user *User, name string, rules []*firewall.Rule) (*Agent, error) {
	if err := validateAgentData(name, rules); err != nil {
		return nil, err
	}

	var found Agent
	err := db.Where("name=?", name).Where("user_id=?", user.ID).First(&found).Error
	if err == nil {
		return nil, fmt.Errorf("agent name already used")
	}

	for i := range rules {
		rules[i].CreatedAt = time.Now()
	}

	newAgent := Agent{
		UserID: user.ID,
		Name:   name,
		Token:  makeRandomToken(),
		Rules:  ToJSONB(rules),
	}

	if err = db.Create(&newAgent).Error; err != nil {
		return nil, fmt.Errorf("error creating new agent: %v", err)
	}

	return &newAgent, nil
}

func UpdateAgent(agent *Agent, name string, rules []*firewall.Rule) error {
	if err := validateAgentData(name, rules); err != nil {
		return err
	}

	var found Agent
	err := db.Where("name=?", name).Where("id != ?", agent.ID).First(&found).Error
	if err == nil {
		return fmt.Errorf("agent name already used")
	}

	for i, rule := range rules {
		if rule.TTL > 0 && rule.CreatedAt.IsZero() {
			rules[i].CreatedAt = time.Now()
		}
	}

	agent.Name = name
	agent.Rules = ToJSONB(rules)
	agent.UpdatedAt = time.Now()

	return db.Save(agent).Error
}
