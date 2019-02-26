package model

import (
	"database/sql/driver"
	"startkit/starter"
	"time"

	"github.com/danilopolani/gocialite/structs"
	"github.com/jinzhu/gorm"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	uuid "gopub/src/github.com/satori/go.uuid"
)

type (
	RoleType string
	User     struct {
		starter.MysqlModel
		ExternalID                   string    `json:"external_id,omitempty" gorm:"not null;column:external_id"`
		OauthTokenExternalID         string    `json:"oauth_token_external_id,omitempty" gorm:"column:oauth_token_external_id"`
		ResetPasswordTokenExternalID string    `json:"reset_password_token_external_id,omitempty" gorm:"column:reset_password_token_external_id"`
		IsActive                     bool      `json:"is_active" gorm:"column:is_active"`
		ActivatedAt                  time.Time `json:"activated_at,omitempty" gorm:"column:activated_at; sql:DEFAULT:TimeZero"`
		Username                     string    `json:"username,omitempty" gorm:"not null;unique;column:username"`
		Firstname                    string    `json:"firstname,omitempty" gorm:"column:firstname"`
		Lastname                     string    `json:"lastname,omitempty" gorm:"column:lastname"`
		Age                          int       `json:"age,omitempty" gorm:"column:age"`
		Password                     string    `json:"password,omitempty" gorm:"column:password"`
		Email                        string    `json:"email,omitempty" gorm:"not null;unique;column:email"`
		Phone                        int       `json:"phone,omitempty" gorm:"not null;column:phone"`
		Region                       string    `json:"region,omitempty" gorm:"not null;column:region"`
		Role                         RoleType  `json:"role,omitempty" gorm:"column:role; sql:not null;type:ENUM('admin', 'staff', 'boat_owner', 'client'); DEFAULT:'client'"`
		/*
			Role int `json:"role" gorm:"column:role; sql:not null;type:tinyint(2);"`
			OauthProvider int  `json:"oauth_provider" gorm:"column:oauth_provider; sql:not null;type:tinyint(2);"`
		*/
	}
)

const (
	Admin     RoleType = "admin"
	Staff     RoleType = "staff"
	BoatOwner RoleType = "boat_owner"
	Client    RoleType = "client"
)

func (m *User) IsEmpty() bool {
	return m.Email == "" &&
		m.Username == "" &&
		m.Password == "" &&
		m.Age == 0 &&
		m.Phone == 0 &&
		m.Region == "" &&
		m.Firstname == "" &&
		m.Lastname == "" &&
		m.Role == ""
}

func NewUser() *User {
	return &User{
		IsActive: false,
	}
}

func (m *User) TableName() string {
	return "user"
}

func (u RoleType) CorrectTo(to RoleType) RoleType {
	if u != Admin || u != Staff || u != BoatOwner || u != Client {
		return to
	}
	return u
}

func (u *RoleType) Scan(value interface{}) error {
	*u = RoleType(string(value.([]byte)))
	return nil
}

func (u RoleType) Value() (driver.Value, error) {
	return string(u), nil
}

func (m *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}

func (m *User) AfterCreate(scope *gorm.Scope) (err error) {
	/* generate external_id (external unique id) for user */
	return
}

func (m *User) BeforeSave(scope *gorm.Scope) error {
	if m.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		scope.SetColumn("password", password)
	}
	return nil
}

func (m *User) PasswordValidate(password string) bool {
	return password != "" && bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password)) == nil
}

func (m *User) SetOauth(provider string, token *oauth2.Token, user *structs.User) *User {
	m = &User{
		Username:  user.Username,
		Firstname: user.FirstName,
		Lastname:  user.LastName,
		Email:     user.Email,
	}
	return m
}

func (m *User) SetResetPasswordToken(rt time.Time, duration int) (*User, bool) {
	// if rt.Before(m.ResetPasswordTokenExpiry.Add(time.Duration(duration) * time.Minute)) {
	// 	return m, false
	// }
	// m.ResetPasswordToken = uuid.NewV4().String()
	// m.ResetPasswordTokenExpiry = time.Now().Add(time.Duration(duration) * time.Minute)
	return m, true
}
