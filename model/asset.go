package model

import (
	"BookingServer/model/asset"
	"startkit/starter"

	"github.com/jinzhu/gorm"
	uuid "gopub/src/github.com/satori/go.uuid"
)

type (
	Asset struct {
		starter.MysqlModel
		ExternalID      string `json:"external_id,omitempty" gorm:"not null;unique;column:external_id"`
		Name            string `json:"name" gorm:"not null;column:name"`
		IsActive        bool   `json:"is_active" gorm:"column:is_active"`
		OwnerID         int    `json:"owner_id,omitempty" gorm:"not null;column:owner_id"`
		OwnerExternalID string `json:"owner_external_id,omitempty" gorm:"not null;column:owner_external_id"`
		asset.Boat
	}
)

func (m *Asset) TableName() string {
	return "asset"
}

func (m *Asset) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}
