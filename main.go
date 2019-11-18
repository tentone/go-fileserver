package main

import (
	"github.com/tentone/godonkey/database"
	"github.com/tentone/godonkey/global"
	"github.com/tentone/godonkey/server"
)

func main() {
	global.StartLogger("server.log")
	global.LoadVersion("version.json")
	global.LoadConfig("config.json")
	database.Create()
	server.ServerStart()
}
