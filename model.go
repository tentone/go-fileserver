package main

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUID identification model for tables that use UUID as identifier
type Model struct {
	// Unique global identifier of the entry used as the primary key
	ID uuid.UUID `gorm:"primaryKey;unique;column:id" json:"uuid,omitempty"`

	// Created date automatically set by GORM when adding new elements
	CreatedAt *time.Time `json:"createdAt"`

	// Updated date automatically set by GORM on update
	UpdatedAt *time.Time `json:"updatedAt"`

	// Deleted date automatically set by GORM on delete
	DeletedAt *time.Time `json:"deletedAt"`
}

func (base *Model) BeforeCreate(_ *gorm.DB) error {
	return nil
}

func NewUUID() Model {
	var uuid, _ = uuid.NewUUID()

	return Model{ID: uuid}
}
