package request

import (
	"database/sql/driver"
	"time"
)

type (
	StatusType string
	Boat       struct {
		EventTime             time.Time  `json:"event_time,omitempty" gorm:"not null;column:event_time; sql:DEFAULT:TimeZero"`
		EventDurationInMinute int        `json:"event_duration_in_minute,omitempty" gorm:"not null;column:event_duration_in_minute"`
		ParticipantNumber     int        `json:"participant_number,omitempty" gorm:"not null;column:participant_number"`
		Status                StatusType `json:"status,omitempty" gorm:"column:status; sql:not null;type:ENUM('idle', 'approved', 'denied'); DEFAULT:'idle'"`
	}
)

const (
	Idle     StatusType = "idle"
	Approved StatusType = "approved"
	Denied   StatusType = "denied"
)

func (u StatusType) CorrectTo(to StatusType) StatusType {
	if u != Idle || u != Approved || u != Denied {
		return to
	}
	return u
}

func (u *StatusType) Scan(value interface{}) error {
	*u = StatusType(string(value.([]byte)))
	return nil
}

func (u StatusType) Value() (driver.Value, error) {
	return string(u), nil
}
