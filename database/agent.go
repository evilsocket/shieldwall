package database

import (
	"errors"
	"fmt"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/shieldwall/firewall"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"net"
	"strconv"
	"strings"
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

func validateAgentData(name string, rules []firewall.Rule) error {
	if name = str.Trim(name); len(name) < MinAgentNameLength {
		return fmt.Errorf("agent name must be at least %d characters long", MinAgentNameLength)
	}

	if len(rules) > MaxAgentRules {
		return fmt.Errorf("max %d rules per agent are allowed", MaxAgentRules)
	}

	for _, rule := range rules {
		if net.ParseIP(rule.Address) == nil {
			return fmt.Errorf("%s is not a valid address", rule.Address)
		}

		if rule.TTL < 0 {
			return fmt.Errorf("really? %d", rule.TTL)
		}

		for _, port := range rule.Ports {
			if strings.Index(port, ":") != -1 {
				// parse as range
				if parts := strings.Split(port, ":"); len(parts) != 2 {
					return fmt.Errorf("%s is not a valid port range", port)
				} else if from, err := strconv.ParseInt(parts[0], 10, 32); err != nil {
					return fmt.Errorf("%s is not a valid port", parts[0])
				} else if to, err := strconv.ParseInt(parts[1], 10, 32); err != nil {
					return fmt.Errorf("%s is not a valid port", parts[1])
				} else if to <= from {
					return fmt.Errorf("bad port range, %d is not >= %d", to, from)
				} else if from < 1 || from > 65535 {
					return fmt.Errorf("%d is outside the valid ports range", from)
				} else if to < 1 || to > 65535 {
					return fmt.Errorf("%d is outside the valid ports range", to)
				}
			} else {
				// parse as number
				if p, err := strconv.ParseInt(port, 10, 32); err != nil {
					return fmt.Errorf("%s is not a valid port", port)
				} else if p < 1 || p > 65535 {
					return fmt.Errorf("%d is outside the valid ports range", p)
				}
			}
		}
	}

	return nil
}

func RegisterAgent(user *User, name string, rules []firewall.Rule) (*Agent, error) {
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

func UpdateAgent(agent *Agent, name string, rules []firewall.Rule) error {
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
