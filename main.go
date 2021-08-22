package main

func main() {
	LoadVersion("version.json")
	LoadConfig("config.json")

	ConnectDatabase()
	RegistryDatabaseMigrate()

	StartHTTPServer()

	select {}
}
