package main

import (
	"net/http"
	"fmt"
	"github.com/BurntSushi/toml"
)

const reply_200  = "Cluster Node is up"
const reply_500  = "Cluster Node is down"

type Config struct {
	LocalBind string
	Galera Galera
	Rabbitmq Rabbitmq
}

func main() {

	var conf Config
	if _, err := toml.DecodeFile("/etc/clusterhc/clusterhc.toml", &conf); err != nil {
		fmt.Println("Failed to parse config file", err.Error())
	}

	if conf.Galera.Enabled {
		http.HandleFunc("/galera", conf.Galera.checkGalera)
	}

	if conf.Rabbitmq.Enabled {
		http.HandleFunc("/rabbitmq", conf.Rabbitmq.checkRabbitmq)
	}

	err := http.ListenAndServe(conf.LocalBind, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}