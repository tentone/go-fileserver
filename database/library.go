package database

import "gorm.io/gorm"

const BASIC = 0
const FOLDER = 1
const IMAGE = 2

// A library contains resources of a specific type.
//
// Libraries can have tasks associated with them used to process resources.
type Library struct {
	NumID

	// Path of the library, unique used to describe the library
	//
	// Path of the library is directly mapped to its folder in the filesystem
	Path string `gorm:"type:varchar(255);unique;column:path" json:"path"`

	// Type of the library
	Type int `gorm:"type:column:type" json:"type"`
}

func LibraryMigrate(db *gorm.DB) {

	db.SingularTable(true)
	db.AutoMigrate(&Library{})
}

func NewLibrary(path string, _type int) *Library {

	var l = new(Library)
	l.Path = path
	l.Type = _type
	return l
}
