package source

func main() {
	StartLogger("api.log")
	LoadVersion("version.json")
	LoadConfig("config.json")

	Create()
	ServerStart()
}
