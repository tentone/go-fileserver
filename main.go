package main

import (
	"github.com/tentone/go-fileserver/database"
	"github.com/tentone/go-fileserver/global"
	"github.com/tentone/go-fileserver/server"
)

func main() {
	global.StartLogger("server.log")
	global.LoadVersion("version.json")
	global.LoadConfig("config.json")

	database.Create()
	server.ServerStart()
}
