package model

import (
	"startkit/starter"
	"time"

	"github.com/jinzhu/gorm"
	uuid "gopub/src/github.com/satori/go.uuid"
)

type (
	ResetPasswordToken struct {
		starter.MysqlModel
		ExternalID     string    `json:"external_id,omitempty" gorm:"not null;column:external_id"`
		UserExternalID string    `json:"user_external_id,omitempty" gorm:"not null;column:user_external_id"`
		Token          string    `json:"token,omitempty" gorm:"not null;column:token"`
		TokenExpiry    time.Time `json:"token_expiry,omitempty" gorm:"column:token_expiry; sql:DEFAULT:TimeZero"`
		IsRevoked      bool      `json:"is_revoked" gorm:"column:is_revoked"`
		RevokedAt      time.Time `json:"revoked_at,omitempty" gorm:"column:revoked_at; sql:DEFAULT:TimeZero"`
	}
)

func (m *ResetPasswordToken) TableName() string {
	return "reset_password_token"
}

func (m *ResetPasswordToken) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}

func (m *ResetPasswordToken) SetResetPasswordToken(rt time.Time, duration int) (*ResetPasswordToken, bool) {
	if rt.Before(m.TokenExpiry.Add(time.Duration(duration) * time.Minute)) {
		return m, false
	}
	m.Token = uuid.NewV4().String()
	m.TokenExpiry = time.Now().Add(time.Duration(duration) * time.Minute)
	return m, true
}

func (m *ResetPasswordToken) AfterCreate(scope *gorm.Scope) (err error) {
	/* generate external_id (external unique id) for reset_password_token */
	return
}
