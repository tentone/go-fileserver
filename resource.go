package main

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Resource struct {
	gorm.Model

	// Library where the resource belongs
	Library   Library `gorm:"foreignKey:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LibraryID uint    `gorm:"column:library_id"`

	// UUID of the resource
	UUID string `gorm:"type:varchar(36)unique;column:uuid" json:"uuid"`

	// File encoding format (e.g. png, jpeg, mp3, wav, mp4)
	Format string `gorm:"column:format" json:"format"`
}

func ResourceMigrate(db *gorm.DB) {
	var err = db.AutoMigrate(&Resource{})
	if err != nil {
		print("Failed to migrate resource table.")
	}
}

func NewResource(uuid string, format string) *Resource {
	return &Resource{
		UUID:   uuid,
		Format: format,
	}
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

// Get a resource from the api
func ResourceGetAPI(writer http.ResponseWriter, request *http.Request) {
	// Form data
	var variables = mux.Vars(request)
	var library string = variables["library"]
	var uuid string = variables["uuid"]

	var path string = Config.Storage.Path + "/" + strings.ToLower(library) + "/" + strings.ToLower(uuid)

	// Read file
	var err error
	var file []byte
	file, err = ioutil.ReadFile(path)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write file
	_, err = writer.Write(file)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/octet-stream; charset=UTF-8")
	writer.WriteHeader(http.StatusOK)
}

// Upload a new resource to the api.
func ResourceUploadAPI(writer http.ResponseWriter, request *http.Request) {
	var uuid = request.FormValue("uuid")
	var library = request.FormValue("library")
	var format = request.FormValue("format")

	var path string = Config.Storage.Path + "/" + strings.ToLower(library)

	// Check if path exists and create if necessary
	var err error
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			_, _ = writer.Write([]byte("Failed to create directory to store data."))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// File path
	var fpath string = strings.ToLower(path) + "/" + strings.ToLower(uuid) + "." + format

	// Read request data
	request.Body = http.MaxBytesReader(writer, request.Body, Config.FileServer.MaxUploadSize)
	err = request.ParseMultipartForm(Config.FileServer.MaxUploadSize)
	if err != nil {
		_, _ = writer.Write([]byte("Cannot read data from the request form."))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var file multipart.File
	file, _, err = request.FormFile("file")
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create and store file
	var storeFile *os.File
	storeFile, err = os.Create(fpath)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer storeFile.Close()

	_, err = storeFile.Write(data)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusOK)
}
