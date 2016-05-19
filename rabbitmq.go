package main

import (
	"net/http"
	"strconv"
	"encoding/json"
	"time"
	"log"
)
type ConfRabbitmq struct {
	Enabled bool
	Interval int
	User string
	Pass string
	Host string
	Port int
	Nodes []string
}

type Rabbitmq struct {
	Conf ConfRabbitmq
	Status Status
	lg log.Logger
}

func (rmq * Rabbitmq) check() {

	client := &http.Client{Timeout: time.Duration(rmq.Conf.Interval/2) * time.Second,}
	url := "http://" + rmq.Conf.Host + ":" + strconv.Itoa(rmq.Conf.Port) + "/api/nodes"
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rmq.Conf.User, rmq.Conf.Pass)
	res, err := client.Do(req)
	if err != nil {
		rmq.lg.Println(err.Error())
		rmq.Status.PartOfCluster = false

	} else {
		decoder := json.NewDecoder(res.Body)
		type Node struct {
			Name string
		}
		var n []Node
		err = decoder.Decode(&n)
		if err != nil {
			rmq.lg.Println(err.Error())
			rmq.Status.PartOfCluster = false
		} else {
			var found int
			for _, node := range n {
				for _, nodeToCheck := range rmq.Conf.Nodes {
					if (node.Name == "rabbit@" + nodeToCheck) {
						found++
						break
					}
				}
			}

			if found == len(rmq.Conf.Nodes) {
				rmq.Status.PartOfCluster = true
				rmq.lg.Println("RABBIT IS OK")
			} else {
				rmq.Status.PartOfCluster = false
				rmq.lg.Println("Can not find all nodes in cluster")
			}
		}
		res.Body.Close()
	}
	rmq.Status.Timestamp = time.Now()
}