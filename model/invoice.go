package model

import (
	"startkit/starter"

	"github.com/jinzhu/gorm"
	uuid "gopub/src/github.com/satori/go.uuid"
)

type (
	InvoiceRecord struct {
		starter.MysqlModel
		ExternalID      string  `json:"external_id,omitempty" gorm:"not null;column:external_id"`
		PayerID         int     `json:"payer_id,omitempty" gorm:"not null;column:payer_id"`
		PayerExternalID string  `json:"payer_external_id,omitempty" gorm:"not null;column:payer_external_id"`
		InvoiceNo       string  `json:"invoice_no,omitempty" gorm:"not null;column:invoice_no"`
		Price           float64 `json:"price,omitempty" gorm:"not null;column:price"`
		Curency         string  `json:"currency,omitempty" gorm:"not null;column:currency"`
		Unit            string  `json:"unit,omitempty" gorm:"not null;column:unit"` // booking once time unit
		Description     string  `json:"description,omitempty" gorm:"not null;column:description"`
	}
)

func (m *InvoiceRecord) TableName() string {
	return "invoice_record"
}

func (m *InvoiceRecord) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("ExternalID", uuid.NewV4().String())
	return nil
}
