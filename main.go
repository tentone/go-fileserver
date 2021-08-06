package main

import (
	"github.com/tentone/go-fileserver/database"
	"github.com/tentone/go-fileserver/global"
	"github.com/tentone/go-fileserver/api"
)

func main() {
	global.StartLogger("api.log")
	global.LoadVersion("version.json")
	global.LoadConfig("config.json")

	database.Create()
	api.ServerStart()
}
