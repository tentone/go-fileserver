package main

import (
	"github.com/tentone/godonkey/global"
	"github.com/tentone/godonkey/server"
)

func main() {
	global.StartLogger("server.log")
	global.LoadVersion("version.json")
	global.LoadConfig("config.json")

	var s = server.Server{}
	s.Start()
}
