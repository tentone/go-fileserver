package main

import "time"

// Sequentially increment numeric identification value.
//
// Should be used for strictly private tables.
type NumID struct {
	// Numeric sequential ID for the database
	ID uint `gorm:"primary_key" json:"id,omitempty"`

	// Created date automatically set by GORM when adding new elements
	CreatedAt *time.Time `json:"-"`

	// Updated date automatically set by GORM on update
	UpdatedAt *time.Time `json:"-"`
}

func NewNumID() NumID {
	return NumID{}
}
