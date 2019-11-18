package database

import (
	"github.com/satori/go.uuid"
	"time"
)

// UUID identification should be used for every public accessible object.
type UUID struct {
	// Unique global identifier of the entry used as the primary key
	ID string `gorm:"type:uniqueidentifier;primary_key;unique;column:id" json:"uuid"`

	// Created date automatically set by GORM when adding new elements
	CreatedAt *time.Time `json:"-"`

	// Updated date automatically set by GORM on update
	UpdatedAt *time.Time `json:"-"`
}

func NewUUID() UUID {
	var u = UUID{}
	u.ID =  uuid.NewV4().String()
	return u
}
