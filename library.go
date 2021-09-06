package main

import "gorm.io/gorm"

// A library contains resources of a specific type.
//
// Libraries can have tasks associated with them used to process resources.
type Library struct {
	Model

	// Path of the library, unique used to describe the library
	//
	// Path of the library is directly mapped to its folder in the filesystem
	Path string `gorm:"type:varchar(255);unique;column:path" json:"path"`
}

func LibraryMigrate(db *gorm.DB) {
	_ = db.AutoMigrate(&Library{})
}

func NewLibrary(path string) *Library {
	return &Library{
		Path: path,
	}
}
