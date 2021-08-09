package main

func main() {
	LoadVersion("version.json")
	LoadConfig("config.json")

	Create()
	ServerStart()
}
