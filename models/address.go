package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Address   string
	Latitude  float64
	Longitude float64
}
