package main

import (
	"fmt"
	"net/http"
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
	/*
url := "http://localhost:" + rmq.port + "/api/nodes"
res, err := http.Get(url)
if err != nil {
 */
	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprintln(w, reply_500)
	return
	/*
}
defer res.Body.Close()

query http://localhost:rmq.port/api/nodes
 */
}