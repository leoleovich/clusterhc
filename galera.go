package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"net/http"
	"strconv"
)
type Galera struct {
	Enabled bool
	User string
	Pass string
	Host string
	Port int
	Local_state int
	Sst_method string
}


// You do not want to change these queries, to avoid injections
const wsrep_local_state_query = "show global status where variable_name = ?"
const wsrep_sst_method_query = "show global variables where variable_name = ?"



func (g Galera) checkGalera(w http.ResponseWriter, r *http.Request) {
	var wsrep_sst_method string
	var wsrep_local_state int
	var varName string

	db, err := sql.Open("mysql", g.User + ":" + g.Pass + "@tcp(" + g.Host + ":" + strconv.Itoa(g.Port) + ")/")
	if err != nil {
		fmt.Println("Can not connect to database: ", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		return
	}
	defer db.Close()


	err = db.QueryRow(wsrep_local_state_query, "wsrep_local_state").Scan(&varName, &wsrep_local_state)
	if err != nil {
		fmt.Println("Error querying " + wsrep_local_state_query + ": ", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		return
	}

	err = db.QueryRow(wsrep_sst_method_query, "wsrep_sst_method").Scan(&varName, &wsrep_sst_method)
	if err != nil {
		fmt.Println("Error querying " + wsrep_sst_method_query + ": ", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		return
	}

	if wsrep_local_state == g.Local_state && wsrep_sst_method == g.Sst_method {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, reply_200)
		return
	} else {
		fmt.Printf("wsrep_local_state is %s, but should be %d", wsrep_local_state, g.Local_state)
		fmt.Printf("wsrep_sst_method is %s, but should be %d", wsrep_sst_method, g.Sst_method)
		w.WriteHeader(http.StatusServiceUnavailable)
		fmt.Fprintln(w, reply_500)
		return
	}
}