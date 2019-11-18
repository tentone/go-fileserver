package database

import "github.com/jinzhu/gorm"

// A library contains resources of a specific type.
//
// Libraries can have tasks associated with them used to process resources.
type Library struct {
	NumID

	// Path of the library, unique used to describe the library
	//
	// Path of the library is directly mapped to its folder in the filesystem
	Path string	`gorm:"type:varchar(255)unique;column:path" json:"path"`
}


func LibraryMigrate(db *gorm.DB) {
	db.SingularTable(true)
	db.AutoMigrate(&Library{})
}

func NewLibrary() *Library {
	var l = new(Library)

	return l
}