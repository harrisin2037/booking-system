package oauth

import (
	"bytes"
	"database/sql/driver"
	"errors"
)

type (
	JSON              []byte
	OauthProviderType string
	OauthProvider     struct {
		OauthGoogle
		OauthFacebook
		Provider OauthProviderType `json:"provider,omitempty" gorm:"column:provider; sql:not null;type:ENUM('null', 'facebook', 'google')"`
	}
)

const (
	Null     OauthProviderType = "null"
	FaceBook OauthProviderType = "facebook"
	Google   OauthProviderType = "google"
)

func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (m JSON) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

func (m *JSON) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("null point exception")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

func (u *OauthProviderType) Scan(value interface{}) error {
	*u = OauthProviderType(string(value.([]byte)))
	return nil
}

func (u OauthProviderType) Value() (driver.Value, error) {
	return string(u), nil
}
