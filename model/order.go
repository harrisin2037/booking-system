package model

import (
	"BookingServer/model/invoice"
	"startkit/starter"
	"time"

	"github.com/jinzhu/gorm"
	uuid "gopub/src/github.com/satori/go.uuid"
)

type (
	OrderRecord struct {
		starter.MysqlModel
		ExternalID               string    `json:"external_id,omitempty" gorm:"not null;column:external_id"`
		CustomerID               int       `json:"customer_id,omitempty" gorm:"not null;column:customer_id"`
		CustomerExternalID       string    `json:"customer_external_id,omitempty" gorm:"not null;column:customer_external_id"`
		PersonInChargeID         int       `json:"person_in_charge_id,omitempty" gorm:"not null;column:person_in_charge_id"`
		PersonInChargeExternalID string    `json:"person_in_charge_external_id,omitempty" gorm:"not null;column:person_in_charge_external_id"`
		EventTime                time.Time `json:"event_time,omitempty" gorm:"not null;column:event_time; sql:DEFAULT:TimeZero"`
		EventDurationInMinute    int       `json:"event_duration_in_minute,omitempty" gorm:"not null;column:event_duration_in_minute"`
		ParticipantNumber        int       `json:"participant_number,omitempty" gorm:"not null;column:participant_number"`
		IsApproved               bool      `json:"is_approved" gorm:"not null;column:is_approved"`
		IsCanceled               bool      `json:"is_canceled" gorm:"not null;column:is_canceled"`
		invoice.Rent
	}
)

func (m *OrderRecord) TableName() string {
	return "order_record"
}

func (m *OrderRecord) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}
