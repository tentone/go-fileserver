module github.com/tentone/go-fileserver

go 1.16

replace github.com/tentone/go-fileserver => ./

require (
	github.com/google/logger v1.0.1
	github.com/gorilla/mux v1.7.4
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.12

)
