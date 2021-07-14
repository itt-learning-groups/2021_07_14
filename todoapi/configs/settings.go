package configs

// Settings contains app environment configuration settings
type Settings struct {
	ServerAddr string `envcfg:"SERVER_ADDR" envcfgDefault:""`
	ServerPort string `envcfg:"SERVER_PORT" envcfgDefault:"8080"`
}
