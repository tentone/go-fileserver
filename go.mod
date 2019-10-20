module github.com/tentone/godonkey

go 1.12

replace github.com/tentone/godonkey => ./

require github.com/buaazp/fasthttprouter v0.1.1

require github.com/google/logger v1.0.1

require (
	github.com/gorilla/mux v1.7.3
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/valyala/fasthttp v1.5.0
)
