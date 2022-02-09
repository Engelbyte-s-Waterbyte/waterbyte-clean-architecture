package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username       string `gorm:"uniqueIndex:username_unique"`
	Name           string
	Email          string `gorm:"uniqueIndex:email_unique"`
	ProfilePicture *string
	GoogleID       *string `gorm:"uniqueIndex:google_id_unique"`
	AppleID        *string `gorm:"uniqueIndex:apple_id_unique"`

	// HostedEvents   []Event `gorm:"many2many:user_hosted_events"`
	// AttendedEvents []Event
}
