package main

import (
	"fmt"
	"net/http"
	"strconv"
	"encoding/json"
)
type Rabbitmq struct {
	Enabled bool
	User string
	Pass string
	Host string
	Port int
	Nodes []string
}



func (rmq Rabbitmq) checkRabbitmq(w http.ResponseWriter, r *http.Request) {

	client := &http.Client{}
	url := "http://" + rmq.Host + ":" + strconv.Itoa(rmq.Port) + "/api/nodes"
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rmq.User, rmq.Pass)
	res, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		return

	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	type Node struct {
		Name string
	}
	var n []Node
	err = decoder.Decode(&n)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		return
	}

	var found int
	for _, node := range n {
		for _, nodeToCheck := range rmq.Nodes {
			if (node.Name == "rabbit@" + nodeToCheck) {
				found++
				break
			}
		}
	}

	if found == len(rmq.Nodes) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, reply_200)
		return
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		return
	}
}