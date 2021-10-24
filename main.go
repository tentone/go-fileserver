package main

func main() {
	LoadVersion("version.json")
	_ = LoadConfig("config.json")
	_ = ConnectDatabase()
	RegistryDatabaseMigrate()

	StartHTTPServer()

	select {}
}
