package identification

import (
"github.com/jinzhu/gorm"
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

func (base *UUID) BeforeCreate(scope *gorm.Scope) error {

	return nil
}

func NewUUID() UUID {

	var u = UUID{}

	var ug, _ = uuid.NewV4()
	u.ID = ug.String()

	return u
}
