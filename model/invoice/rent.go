package invoice

import (
	"time"
)

type Rent struct {
	InvoiceNumber    string    `json:"invoice_number,omitempty" gorm:"not null;unique;column:invoice_number"`
	Currency         string    `json:"currency,omitempty" gorm:"not null;column:currency"`
	Charge           float64   `json:"charge,omitempty" gorm:"not null;column:charge"`
	Deposit          float64   `json:"deposit,omitempty" gorm:"not null;column:deposit"`
	DepositRemitTime time.Time `json:"deposit_remit_time,omitempty" gorm:"column:deposit_remit_time; sql:DEFAULT:TimeZero"`
	AdditionalCharge float64   `json:"additional_charge,omitempty" gorm:"not null;column:additional_charge"`
}
