package model

import (
	"startkit/starter"

	"github.com/jinzhu/gorm"
	uuid "gopub/src/github.com/satori/go.uuid"
)

// record -> payment

type (
	PaymentRecord struct {
		starter.MysqlModel
		ExternalID         string `json:"external_id,omitempty" gorm:"not null;column:external_id"`
		InvoiceID          int    `json:"invoice_id,omitempty" gorm:"not null;unique;column:invoice_id"`
		InvoiceExternalID  string `json:"invoice_external_id,omitempty" gorm:"not null;unique;column:invoice_external_id"`
		PayerID            int    `json:"payer_id,omitempty" gorm:"not null;column:payer_id"`
		PayerExternalID    string `json:"payer_external_id,omitempty" gorm:"not null;column:payer_external_id"`
		ApproverID         int    `json:"approver_id,omitempty" gorm:"not null;column:approver_id"`
		ApproverExternalID string `json:"approver_external_id,omitempty" gorm:"not null;column:approver_external_id"`
		Total              int    `json:"total,omitempty" gorm:"not null;column:total"`
		Currency           string `json:"currency,omitempty" gorm:"not null;column:currency"`
		GatewayType        string `json:"gateway_type,omitempty" gorm:"not null;column:gateway_type"`
		PaymentID          string `json:"payment_id,omitempty" gorm:"not null;column:payment_id"`
		State              string `json:"state,omitempty" gorm:"column:state"`
		Intent             string `json:"intent,omitempty" gorm:"column:intent"`
	}
)

func (m *PaymentRecord) TableName() string {
	return "payment_record"
}

func (m *PaymentRecord) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}
