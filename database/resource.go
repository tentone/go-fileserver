package database

import (
	"github.com/jinzhu/gorm"
)

type Resource struct {

	NumID

	// UUID of the resource stored in the resource server
	UUID string `gorm:"type:varchar(255);unique;column:uuid" json:"uuid"`

	// Format of the resource stored in the resource server
	Format string `gorm:"column:format" json:"format"`
}

func ResourceCreateDB(db *gorm.DB) {
	db.SingularTable(true)
	db.AutoMigrate(&Resource{})
}

func NewResource(uuid string, format string) *Resource {
	var r = new(Resource)

	r.UUID = uuid
	r.Format = format

	return r
}

func (resource *Resource) StoreDB(db *gorm.DB) error {

	var conn = db.Save(resource)

	if conn.Error != nil {
		return conn.Error
	}

	return nil
}

// Store a resource received from the API in the database and return a reference to object.
//
// Checks if there is already a resource with the same UUID and returns it if found.
func StoreResourceDB(db *gorm.DB, res *Resource) (a *Resource, e error) {
	if res != nil {
		var err error
		var resource *Resource

		// Check if picture exists in the database
		resource, err = GetResourceByUuidDB(db, res.UUID)

		// Create new picture
		if err != nil {
			resource = NewResource(res.UUID, res.Format)
			err = resource.StoreDB(db)

			if err != nil {
				return nil, err
			}
		}

		return resource, nil
	} else {
		return nil, nil
	}
}

func GetResourceByIDDB(db *gorm.DB, id uint) (a *Resource, e error) {

	var resource *Resource = new(Resource)

	var conn = db.Where("id = ?", id).First(resource)
	if conn.Error != nil {
		return nil, conn.Error
	}
	return resource, nil
}

// Get resource from the database by its uuid.
func GetResourceByUuidDB(db *gorm.DB, uuid string) (a *Resource, e error) {

	var resource *Resource = new(Resource)

	var conn = db.Where("uuid = ?", uuid).First(resource)
	if conn.Error != nil {
		return nil, conn.Error
	}

	return resource, nil
}
