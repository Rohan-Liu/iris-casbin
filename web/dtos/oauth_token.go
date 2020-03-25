package dtos

import (
	"../../cache"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"time"
)

type OauthToken struct {
	Token     string
	UserId    uint
	Secret    string
	ExpressIn int64
	Revoked   bool
	Name      string
	RoleIds   []string
	RoleName  string
}

type Token struct {
	Token string `json:"access_token"`
}

func (ot *OauthToken) OauthTokenCreate() *Token {
	value, _ := json.Marshal(ot)

	ttl := time.Duration(1000 * 60 * 60 * 100) //
	stringValue := string(value)
	_, err := cache.Set(ot.Token, stringValue, ttl)
	if err != nil {
		color.Red(fmt.Sprintf("NewAdapter 错误: %v", err))
	}
	//database.Db.Create(ot)
	return &Token{ot.Token}
}

func (ot *OauthToken) GetOauthTokenByToken(token string) {
	value := cache.Get(token)

	_ = json.Unmarshal([]byte(value), &ot)
}

func (ot *OauthToken) RemoveOauthTokenByToken() {
	_, _ = cache.Del(ot.Token)
}
