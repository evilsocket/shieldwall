package database

import (
	"fmt"
	"github.com/evilsocket/islazy/str"
	"github.com/evilsocket/shieldwall/firewall"
)

const MinAgentNameLength = 3
const MaxAgentRules = 100

const MinAlertAfter = 30    // 30 seconds
const MaxAlertAfter = 21600 // 6 hours

const MinAlertPeriod = 1800  // 30 minutes
const MaxAlertPeriod = 10800 // 3 hours

type AgentWritableFields struct {
	ID          uint
	Name        string
	Rules       []*firewall.Rule
	User        *User
	AlertAfter  uint
	AlertPeriod uint
}

func (f *AgentWritableFields) Validate() error {
	if f.Name = str.Trim(f.Name); len(f.Name) < MinAgentNameLength {
		return fmt.Errorf("agent name must be at least %d characters long", MinAgentNameLength)
	}

	if len(f.Rules) > MaxAgentRules {
		return fmt.Errorf("max %d rules per agent are allowed", MaxAgentRules)
	}

	if f.AlertAfter != 0 {
		if f.AlertAfter < MinAlertAfter {
			return fmt.Errorf("alert trigger time must be at least %d seconds", MinAlertAfter)
		} else if f.AlertAfter > MaxAlertAfter {
			return fmt.Errorf("alert trigger time must be at most %d seconds", MaxAlertAfter)
		}

		if f.AlertPeriod != 0 {
			if f.AlertPeriod < MinAlertPeriod {
				return fmt.Errorf("alert notification period must be at least %d seconds", MinAlertPeriod)
			} else if f.AlertPeriod > MaxAlertPeriod {
				return fmt.Errorf("alert notification period must be at most %d seconds", MaxAlertPeriod)
			}
		}
	}

	for _, rule := range f.Rules {
		if err := rule.Validate(); err != nil {
			return err
		}
	}

	var found Agent

	if f.ID == 0 {
		err := db.Where("name=?", f.Name).Where("user_id=?", f.User.ID).First(&found).Error
		if err == nil {
			return fmt.Errorf("agent name already used")
		}
	} else {
		err := db.Where("name=?", f.Name).Where("id != ?", f.ID).First(&found).Error
		if err == nil {
			return fmt.Errorf("agent name already used")
		}
	}

	return nil
}
