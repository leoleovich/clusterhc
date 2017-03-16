package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	LocalBind        string
	Log              string
	MonitoringPrefix string
	MonitoringAddr   string
	Galera           ConfGalera
	Rabbitmq         ConfRabbitmq
}

func main() {

	var conf Config
	if _, err := toml.DecodeFile("/etc/clusterhc/clusterhc.toml", &conf); err != nil {
		fmt.Println("Failed to parse config file", err.Error())
	}
	f, err := os.OpenFile(conf.Log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0660)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	lg := log.New(f, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	monitoring := &Monitoring{0, 0, conf.MonitoringAddr, conf.MonitoringPrefix, lg}
	if conf.MonitoringAddr == "" {
		lg.Println("Monitoring is disabled")
	} else {
		go monitoring.report()
	}

	if conf.Galera.Enabled {
		g := Galera{conf.Galera, Status{false, time.Now(), conf.Galera.Interval * 2, monitoring}, *lg}
		/*
			We do asynchronous checking, that ddos of check will not kill database
		*/
		go g.check()
		http.HandleFunc("/galera", g.Status.get)
	}

	if conf.Rabbitmq.Enabled {
		rmq := Rabbitmq{conf.Rabbitmq, Status{false, time.Now(), conf.Rabbitmq.Interval * 2, monitoring}, *lg}
		/*
			We do asynchronous checking, that ddos of check will not kill database
		*/
		go rmq.check()
		http.HandleFunc("/rabbitmq", rmq.Status.get)
	}

	if http.ListenAndServe(conf.LocalBind, nil) != nil {
		fmt.Printf(err.Error())
	}
}
