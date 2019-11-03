package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
		eventLogs := make([]EventLog, 0, 10)

		for eventLog := range resc {
			eventLogs = append(eventLogs, eventLog)

			if len(eventLogs) >= 10 {
				valueStrings := make([]string, 0, len(eventLogs))
				valueArgs := make([]interface{}, 0, len(eventLogs)*3)

				for _, eventLog := range eventLogs {
					valueStrings = append(valueStrings, "(?, ?, ?)")
					valueArgs = append(valueArgs, fmt.Sprintf("%s", eventLog.At))
					valueArgs = append(valueArgs, eventLog.Name)
					valueArgs = append(valueArgs, eventLog.Value)
				}

				stmt := fmt.Sprintf("INSERT INTO eventlog(at, name, value) VALUES %s", strings.Join(valueStrings, ","))
				_, e := db.Exec(stmt, valueArgs...)
				if e != nil {
					panic(e.Error())
				}

				eventLogs = make([]EventLog, 0, 10)
			}
		}
	}(resc)

	jst, e := time.LoadLocation("Asia/Tokyo")

	if e != nil {
		panic(e.Error())
	}

	hakaruHandler := func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		value := r.URL.Query().Get("value")

		now := time.Now().In(jst)

		resc <- EventLog{
			At:    now,
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
