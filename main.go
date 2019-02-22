package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"
)

type cmdLine struct {
	test                       bool
	pgHost, pgUser, pgDatabase string
}

var flags cmdLine

func parseFlags() {
	flag.BoolVar(&flags.test, "test", false, "Load the test dataset")
	flag.StringVar(&flags.pgHost, "pg-host", "/run/postgresql", "Postgresql hostname")
	flag.StringVar(&flags.pgUser, "pg-user", "postgres", "Postgresql username")
	flag.StringVar(&flags.pgDatabase, "pg-db", "playground", "Postgresql database")
	flag.Parse()
	log.Debug("Starting with arguments:", "flags", flags)
}

func initLog() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	parseFlags()
	initDb()
	http.HandleFunc("/bc/", bcHandler)
	http.HandleFunc("/cat/", catHandler)
	if err := http.ListenAndServe("0.0.0.0:8012", nil); err != nil {
		log.Warn("ListenAndServe terminated", "error", err)
	}
}
