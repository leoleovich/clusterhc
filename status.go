package main

import (
	"fmt"
	"net/http"
	"time"
)

type Status struct {
	PartOfCluster bool
	Timestamp     time.Time
	Expired       int
	monitoring    *Monitoring
}

const reply_200 = "Cluster Node is up"
const reply_500 = "Cluster Node is down"

func (s *Status) get(w http.ResponseWriter, r *http.Request) {
	intervalAgo := time.Now().Add(-time.Duration(s.Expired) * time.Second)
	defer r.Body.Close()

	// We need to check if result is fresh
	if s.Timestamp.After(intervalAgo) {
		if s.PartOfCluster {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, reply_200)
			s.monitoring.Requests200++
			return
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintln(w, reply_500)
			s.monitoring.Requests500++
			return
		}
	} else {
		// Think about monitoring of service
		//g.lg.Println("Looks like result of check is outdated")
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		s.monitoring.Requests500++
		return
	}
}
