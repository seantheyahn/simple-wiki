package main

import (
	"flag"
	"sean/wiki/config"
	"sean/wiki/server"
	"sean/wiki/services"

	_ "github.com/jackc/pgx/stdlib"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "path/to/config.json")
	flag.Parse()

	config.Init(configPath)
	services.Init()
	server.Run()
}
