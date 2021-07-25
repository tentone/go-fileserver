package database

import (
	"gorm.io/gorm"
)

type Resource struct {
	NumID

	// UUID of the resource
	UUID string `gorm:"type:varchar(36)unique;column:uuid" json:"uuid"`

	// File encoding format
	Format string `gorm:"column:format" json:"format"`
}

func ResourceMigrate(db *gorm.DB) {
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

// Get resources from its numeric ID.
func GetResourceByIDDB(db *gorm.DB, id uint) (a *Resource, e error) {

	var resource *Resource = new(Resource)

	var conn = db.Where("id = ?", id).First(resource)
	if conn.Error != nil {
		return nil, conn.Error
	}
	return resource, nil
}

// Get resource from its UUID.
func GetResourceByUuidDB(db *gorm.DB, uuid string) (a *Resource, e error) {

	var resource *Resource = new(Resource)

	var conn = db.Where("uuid = ?", uuid).First(resource)
	if conn.Error != nil {
		return nil, conn.Error
	}

	return resource, nil
}
