package main

import (
	"flag"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/seantheyahn/simple-wiki/config"
	"github.com/seantheyahn/simple-wiki/server"
	"github.com/seantheyahn/simple-wiki/services"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "path/to/config.json")
	flag.Parse()

	config.Init(configPath)
	services.Init()
	defer services.Finalize()
	server.Run()
}
