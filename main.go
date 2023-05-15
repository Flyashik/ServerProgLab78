package main

import (
	"ApiService/internal/app/api"
	"ApiService/internal/app/config"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/api.toml", "Path to config")
}

func main() {
	flag.Parse()

	config := config.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	s := api.New(config)

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
