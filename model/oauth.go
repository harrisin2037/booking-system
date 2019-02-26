package model

import (
	"BookingServer/model/oauth"
	"startkit/starter"
	"time"

	"github.com/danilopolani/gocialite/structs"
	"github.com/jinzhu/gorm"
	"golang.org/x/oauth2"
	uuid "gopub/src/github.com/satori/go.uuid"
)

type (
	OauthToken struct {
		starter.MysqlModel
		ExternalID        string    `json:"external_id,omitempty" gorm:"not null;unique;column:external_id"`
		UserExternalID    string    `json:"user_external_id,omitempty" gorm:"not null;unique;column:user_external_id"`
		AccessToken       string    `json:"access_token,omitempty" gorm:"not null;unique;column:access_token"`
		AccessTokenExpiry time.Time `json:"access_token_expiry,omitempty" gorm:"column:access_token_expiry; sql:DEFAULT:TimeZero"`
		oauth.OauthProvider
	}
)

func (m *OauthToken) TableName() string {
	return "oauth_token"
}

func (m *OauthToken) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}

func (m *OauthToken) SetOauth(provider string, token *oauth2.Token, user *structs.User) *OauthToken {
	return m
}
