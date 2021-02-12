package api

import (
	"github.com/evilsocket/shieldwall/database"
	"github.com/evilsocket/shieldwall/firewall"
)

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	NewPassword string `json:"password"`
	Use2FA      bool   `json:"use_2fa"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Token   string         `json:"token"`
	User    *database.User `json:"data"`
	Address string         `json:"address"`
}

type Step2Request struct {
	Code string `json:"code"`
}

type AgentCreationRequest struct {
	Name  string           `json:"name"`
	Rules []*firewall.Rule `json:"rules"`
}

type AgentUpdateRequest struct {
	Name  string           `json:"name"`
	Rules []*firewall.Rule `json:"rules"`
}
