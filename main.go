package main

import (
	"flag"
)

func main() {
	// go run . --config /path/to/config.yml
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	mutill := Mutill{}
	mutill.RegisterService(config)
	mutill.StartAll()
}
