package models

import "gorm.io/gorm"

type Event struct {
	gorm.Model
	Name   string
	HostID uint
	Host   User `gorm:"foreignkey:HostID;references:ID"`
}
