package main

import (
	"log"
	"net/http"
	"time"

	"database/sql"

	"os"

	_ "github.com/go-sql-driver/mysql"
)

type EventLog struct {
	At    time.Time
	Name  string
	Value string
}

func main() {
	dataSourceName := os.Getenv("HAKARU_DATASOURCENAME")
	if dataSourceName == "" {
		dataSourceName = "root:password@tcp(127.0.0.1:13306)/hakaru"
	}

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	resc := make(chan EventLog)

	go func(resc chan EventLog) {
		eventLogs := make([]EventLog, 10)

		for eventLog := range resc {
			eventLogs = append(eventLogs, eventLog)

			if len(eventLogs) >= 10 {
				_, e := db.Exec("INSERT INTO eventlog(at, name, value) values ?", eventLogs)
				if e != nil {
					panic(e.Error())
				}
				eventLogs = make([]EventLog, 10)
			}
		}
	}(resc)

	hakaruHandler := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		value := r.URL.Query().Get("value")

		resc <- EventLog{
			At:    time.Now(),
			Name:  name,
			Value: value,
		}

		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
	}

	http.HandleFunc("/hakaru", hakaruHandler)
	http.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) })

	// start server
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
