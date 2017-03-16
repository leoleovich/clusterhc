package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Monitoring struct {
	Requests200 int
	Requests500 int
	Address     string
	Prefix      string
	lg          *log.Logger
}

func (m *Monitoring) report() {
	for ; ; time.Sleep(60 * time.Second) {
		conn, err := net.DialTimeout("tcp", m.Address, time.Duration(1)*time.Second)
		if err != nil {
			m.lg.Println("Can not connect to " + m.Address + ": " + err.Error())
			continue
		}
		now := time.Now()
		err = conn.SetWriteDeadline(now.Add(time.Duration(1 * time.Second)))
		if err != nil {
			m.lg.Println("Error while setting up write timeout: " + err.Error())
			conn.Close()
			continue
		}
		_, err = conn.Write(
			[]byte(
				strings.Join(
					[]string{fmt.Sprintf("%s.clusterhc.requests200 %d %d", m.Prefix, m.Requests200, now.Unix()),
						fmt.Sprintf("%s.clusterhc.requests500 %d %d\n", m.Prefix, m.Requests500, now.Unix())},
					"\n")))

		if err != nil {
			m.lg.Println("Error duing write to " + m.Address + ": " + err.Error())
			conn.Close()
			continue
		}

		conn.Close()
		m.Requests200, m.Requests500 = 0, 0
	}
}
