package main

import (
	"github.com/tentone/godonkey/global"
	"github.com/tentone/godonkey/server"
)

func main() {
	var err error

	err = global.LoadVersion("version.json")
	if err != nil {
		// TODO <ADD CODE HERE>
	}

	err = global.LoadConfig("config.json")
	if err != nil {
		// TODO <ADD CODE HERE>
	}

	var s = server.Server{}
	s.Start()
}
