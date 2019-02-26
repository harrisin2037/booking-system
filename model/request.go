package model

import (
	"startkit/starter"

	"github.com/jinzhu/gorm"
	uuid "gopub/src/github.com/satori/go.uuid"
)

type (
	RequestRecord struct {
		starter.MysqlModel
		ExternalID         string `json:"external_id,omitempty" gorm:"not null;column:external_id"`
		CustomerID         int    `json:"customer_id,omitempty" gorm:"not null;column:customer_id"`
		CustomerExternalID string `json:"customer_external_id,omitempty" gorm:"not null;column:customer_external_id"`
	}
)

func (m *RequestRecord) TableName() string {
	return "request_record"
}

func (m *RequestRecord) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}
