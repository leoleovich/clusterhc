package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type Galera struct {
	Conf   ConfGalera
	Status Status
	lg     log.Logger
}

type ConfGalera struct {
	Enabled     bool
	Interval    int
	User        string
	Pass        string
	Host        string
	Port        int
	Socket      string
	Local_state int
	Sst_method  string
}

// You do not want to change these queries, to avoid injections
const wsrep_local_state_query = "show global status where variable_name = ?"
const wsrep_sst_method_query = "show global variables where variable_name = ?"

func (c *ConfGalera) getConnectionString() string {
	if c.User == "" {
		c.User = "root"
	}
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 3306
	}
	protocol := fmt.Sprintf("tcp(%s:%d)", c.Host, c.Port)
	if c.Socket != "" {
		protocol = fmt.Sprintf("unix(%s)", c.Socket)
	}

	if c.Pass == "" {
		return fmt.Sprintf("%s@%s/", c.User, protocol)
	} else {
		return fmt.Sprintf("%s:%s@%s/", c.User, c.Pass, protocol)
	}
}

func (g *Galera) check() {
	for ; ; time.Sleep(time.Duration(g.Conf.Interval) * time.Second) {
		var wsrep_sst_method string
		var wsrep_local_state int
		var varName string

		db, err := sql.Open("mysql", fmt.Sprintf("%s?timeout=%ds&readTimeout=%ds", g.Conf.getConnectionString(), g.Conf.Interval/2+1, g.Conf.Interval/2+1))
		if err != nil {
			g.lg.Println("Timeout while connecting to mysql", err.Error())
			g.Status.PartOfCluster = false
			g.Status.Timestamp = time.Now()
			db.Close()
			continue
		}

		err = db.QueryRow(wsrep_local_state_query, "wsrep_local_state").Scan(&varName, &wsrep_local_state)
		if err != nil {
			g.lg.Println("Error querying "+wsrep_local_state_query+": ", err.Error())
			g.Status.PartOfCluster = false

		} else {
			err = db.QueryRow(wsrep_sst_method_query, "wsrep_sst_method").Scan(&varName, &wsrep_sst_method)
			if err != nil {
				g.lg.Println("Error querying "+wsrep_sst_method_query+": ", err.Error())
				g.Status.PartOfCluster = false
			} else {
				if wsrep_local_state == g.Conf.Local_state && wsrep_sst_method == g.Conf.Sst_method {
					g.Status.PartOfCluster = true
				} else {
					g.lg.Printf("wsrep_local_state is %d, but should be %d", wsrep_local_state, g.Conf.Local_state)
					g.lg.Printf("wsrep_sst_method is %s, but should be %s", wsrep_sst_method, g.Conf.Sst_method)
					g.Status.PartOfCluster = false
				}
			}
		}
		g.Status.Timestamp = time.Now()
		db.Close()
	}
}
