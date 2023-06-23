package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/Stroevik3/OfficeRestaurant/internal/app/customerserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "internal/app/customerserver/configs/customerserver.toml", "path to config file")
}

func main() {
	flag.Parse()

	config := customerserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	if err := customerserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
